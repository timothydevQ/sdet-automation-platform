package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type cartItem struct {
	SKU  string `json:"sku"`
	Qty  int    `json:"qty"`
	Name string `json:"name"`
	Price int64 `json:"price_cents"`
}

type order struct {
	ID         int64      `json:"id"`
	UserID     int64      `json:"user_id"`
	Status     string     `json:"status"`
	TotalCents int64      `json:"total_cents"`
	Items      []cartItem `json:"items"`
	CreatedAt  time.Time  `json:"created_at"`
}

func cartKey(uid int64) string {
	return fmt.Sprintf("cart:%d", uid)
}

func (d *deps) addToCart(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value(ctxUserID).(int64)
	var req struct {
		SKU string `json:"sku"`
		Qty int    `json:"qty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad json", 400)
		return
	}
	if req.Qty <= 0 {
		http.Error(w, "qty must be positive", 400)
		return
	}
	p, err := d.fetchProductBySKU(r.Context(), req.SKU)
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}
	item := cartItem{SKU: req.SKU, Qty: req.Qty, Name: p.Name, Price: p.Price}
	b, _ := json.Marshal(item)
	if err := d.rdb.HSet(r.Context(), cartKey(uid), req.SKU, b).Err(); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	d.rdb.Expire(r.Context(), cartKey(uid), 24*time.Hour)
	writeJSON(w, 200, item)
}

func (d *deps) getCart(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value(ctxUserID).(int64)
	items, err := d.loadCart(r.Context(), uid)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	writeJSON(w, 200, items)
}

func (d *deps) removeFromCart(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value(ctxUserID).(int64)
	sku := r.PathValue("sku")
	d.rdb.HDel(r.Context(), cartKey(uid), sku)
	w.WriteHeader(204)
}

func (d *deps) loadCart(ctx context.Context, uid int64) ([]cartItem, error) {
	raw, err := d.rdb.HGetAll(ctx, cartKey(uid)).Result()
	if err != nil {
		return nil, err
	}
	out := make([]cartItem, 0, len(raw))
	for _, v := range raw {
		var item cartItem
		if err := json.Unmarshal([]byte(v), &item); err == nil {
			out = append(out, item)
		}
	}
	return out, nil
}

func (d *deps) checkout(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value(ctxUserID).(int64)
	idemKey := r.Header.Get("Idempotency-Key")

	var req struct {
		Coupon string `json:"coupon"`
		Card   string `json:"card_token"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	if idemKey != "" {
		if existing, err := d.lookupIdempotent(r.Context(), idemKey); err == nil && existing != 0 {
			o, err := d.loadOrder(r.Context(), existing)
			if err == nil {
				writeJSON(w, 200, o)
				return
			}
		}
	}

	items, err := d.loadCart(r.Context(), uid)
	if err != nil || len(items) == 0 {
		http.Error(w, "cart empty", 400)
		return
	}

	var subtotal int64
	for _, it := range items {
		subtotal += it.Price * int64(it.Qty)
	}
	total := applyDiscount(subtotal, req.Coupon)

	for _, it := range items {
		if err := d.reserve(r.Context(), it.SKU, it.Qty); err != nil {
			http.Error(w, fmt.Sprintf("reserve %s: %v", it.SKU, err), 409)
			return
		}
	}

	pay, err := d.charge(r.Context(), req.Card, total)
	if err != nil || !pay.Approved {
		for _, it := range items {
			d.release(r.Context(), it.SKU, it.Qty)
		}
		msg := "payment declined"
		if err != nil {
			msg = err.Error()
		}
		http.Error(w, msg, 402)
		return
	}

	tx, err := d.db.Begin()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	var oid int64
	err = tx.QueryRow(`
		INSERT INTO orders (user_id, status, total_cents, payment_ref)
		VALUES ($1, 'paid', $2, $3) RETURNING id
	`, uid, total, pay.Ref).Scan(&oid)
	if err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), 500)
		return
	}
	for _, it := range items {
		if _, err := tx.Exec(`
			INSERT INTO order_items (order_id, sku, name, qty, price_cents)
			VALUES ($1, $2, $3, $4, $5)
		`, oid, it.SKU, it.Name, it.Qty, it.Price); err != nil {
			tx.Rollback()
			http.Error(w, err.Error(), 500)
			return
		}
	}
	if err := tx.Commit(); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if idemKey != "" {
		d.storeIdempotent(r.Context(), idemKey, oid)
	}
	d.rdb.Del(r.Context(), cartKey(uid))

	o := order{ID: oid, UserID: uid, Status: "paid", TotalCents: total, Items: items, CreatedAt: time.Now()}
	d.publisher.Publish(r.Context(), "order.created", o)
	// BUG (intentional): order.created published twice on retry path.
	// Caught by test_order_event_idempotency (consumer-side).
	writeJSON(w, 201, o)
}

func (d *deps) getOrder(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value(ctxUserID).(int64)
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	o, err := d.loadOrder(r.Context(), id)
	if err != nil {
		http.Error(w, "not found", 404)
		return
	}
	if o.UserID != uid {
		http.Error(w, "forbidden", 403)
		return
	}
	writeJSON(w, 200, o)
}

func (d *deps) listOrders(w http.ResponseWriter, r *http.Request) {
	rows, err := d.db.Query(`SELECT id, user_id, status, total_cents, created_at
		FROM orders ORDER BY created_at DESC LIMIT 200`)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()
	out := []order{}
	for rows.Next() {
		var o order
		rows.Scan(&o.ID, &o.UserID, &o.Status, &o.TotalCents, &o.CreatedAt)
		out = append(out, o)
	}
	writeJSON(w, 200, out)
}

func (d *deps) refundOrder(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if _, err := d.db.Exec(`UPDATE orders SET status='refunded' WHERE id=$1`, id); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	writeJSON(w, 200, map[string]any{"ok": true, "id": id})
}

func (d *deps) loadOrder(ctx context.Context, id int64) (*order, error) {
	var o order
	err := d.db.QueryRowContext(ctx, `SELECT id, user_id, status, total_cents, created_at
		FROM orders WHERE id = $1`, id).Scan(&o.ID, &o.UserID, &o.Status, &o.TotalCents, &o.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, errors.New("not found")
	}
	if err != nil {
		return nil, err
	}
	rows, _ := d.db.QueryContext(ctx, `SELECT sku, name, qty, price_cents
		FROM order_items WHERE order_id=$1`, id)
	defer rows.Close()
	for rows.Next() {
		var it cartItem
		rows.Scan(&it.SKU, &it.Name, &it.Qty, &it.Price)
		o.Items = append(o.Items, it)
	}
	return &o, nil
}

func applyDiscount(subtotal int64, coupon string) int64 {
	switch coupon {
	case "":
		return subtotal
	case "WELCOME10":
		return subtotal - subtotal/10
	case "BULK20":
		// BUG (intentional): rounds wrong direction at 5-cent boundaries
		// when subtotal*20 is not divisible by 100. Off-by-one cent.
		return subtotal - (subtotal*20)/100
	default:
		return subtotal
	}
}

func migrate(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS orders (
			id BIGSERIAL PRIMARY KEY,
			user_id BIGINT NOT NULL,
			status TEXT NOT NULL,
			total_cents BIGINT NOT NULL,
			payment_ref TEXT,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
		CREATE TABLE IF NOT EXISTS order_items (
			id BIGSERIAL PRIMARY KEY,
			order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
			sku TEXT NOT NULL,
			name TEXT NOT NULL,
			qty INT NOT NULL,
			price_cents BIGINT NOT NULL
		);
		CREATE TABLE IF NOT EXISTS idempotency_keys (
			key TEXT PRIMARY KEY,
			order_id BIGINT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
		CREATE INDEX IF NOT EXISTS orders_user_idx ON orders(user_id);
	`)
	return err
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if v != nil {
		json.NewEncoder(w).Encode(v)
	}
}
