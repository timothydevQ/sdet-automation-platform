package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

type product struct {
	ID        int64   `json:"id"`
	SKU       string  `json:"sku"`
	Name      string  `json:"name"`
	Price     int64   `json:"price_cents"`
	Stock     int     `json:"stock"`
	Category  string  `json:"category"`
}

func main() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("DB_DSN required")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	for i := 0; i < 30; i++ {
		if err := db.Ping(); err == nil {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	if err := migrate(db); err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /products", listProducts(db))
	mux.HandleFunc("GET /products/{id}", getProduct(db))
	mux.HandleFunc("POST /products/{id}/reserve", reserveStock(db))
	mux.HandleFunc("POST /products/{id}/release", releaseStock(db))
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	log.Printf("catalog-service listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func listProducts(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")
		var rows *sql.Rows
		var err error
		if q != "" {
			// BUG (intentional): doesn't escape % and _ in LIKE pattern.
			// Search breaks for special chars; caught by test_search_special_chars.
			rows, err = db.Query(`SELECT id, sku, name, price_cents, stock, category
				FROM products WHERE name ILIKE $1 OR sku ILIKE $1`, "%"+q+"%")
		} else {
			rows, err = db.Query(`SELECT id, sku, name, price_cents, stock, category FROM products`)
		}
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer rows.Close()
		out := []product{}
		for rows.Next() {
			var p product
			if err := rows.Scan(&p.ID, &p.SKU, &p.Name, &p.Price, &p.Stock, &p.Category); err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			out = append(out, p)
		}
		writeJSON(w, 200, out)
	}
}

func getProduct(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			http.Error(w, "bad id", 400)
			return
		}
		var p product
		err = db.QueryRow(`SELECT id, sku, name, price_cents, stock, category
			FROM products WHERE id = $1`, id).
			Scan(&p.ID, &p.SKU, &p.Name, &p.Price, &p.Stock, &p.Category)
		if err == sql.ErrNoRows {
			http.Error(w, "not found", 404)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		writeJSON(w, 200, p)
	}
}

type qtyReq struct {
	Qty int `json:"qty"`
}

func reserveStock(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
		var req qtyReq
		json.NewDecoder(r.Body).Decode(&req)
		// BUG (intentional): non-atomic decrement under concurrency.
		// Two parallel reservations can both pass the read check.
		// Caught by test_concurrent_checkout_inventory.
		var stock int
		if err := db.QueryRow(`SELECT stock FROM products WHERE id = $1`, id).Scan(&stock); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		if stock < req.Qty {
			http.Error(w, "out of stock", 409)
			return
		}
		if _, err := db.Exec(`UPDATE products SET stock = stock - $1 WHERE id = $2`, req.Qty, id); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		writeJSON(w, 200, map[string]any{"ok": true, "remaining": stock - req.Qty})
	}
}

func releaseStock(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.ParseInt(r.PathValue("id"), 10, 64)
		var req qtyReq
		json.NewDecoder(r.Body).Decode(&req)
		if _, err := db.Exec(`UPDATE products SET stock = stock + $1 WHERE id = $2`, req.Qty, id); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		writeJSON(w, 200, map[string]any{"ok": true})
	}
}

func migrate(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS products (
			id BIGSERIAL PRIMARY KEY,
			sku TEXT UNIQUE NOT NULL,
			name TEXT NOT NULL,
			price_cents BIGINT NOT NULL CHECK (price_cents >= 0),
			stock INT NOT NULL DEFAULT 0,
			category TEXT NOT NULL DEFAULT 'general',
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
		CREATE INDEX IF NOT EXISTS products_name_idx ON products (name);
		CREATE INDEX IF NOT EXISTS products_category_idx ON products (category);
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

var _ = strings.Contains
