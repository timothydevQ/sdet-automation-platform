package main

import (
	"context"
	"database/sql"
)

func (d *deps) lookupIdempotent(ctx context.Context, key string) (int64, error) {
	var id int64
	err := d.db.QueryRowContext(ctx, `SELECT order_id FROM idempotency_keys WHERE key = $1`, key).Scan(&id)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return id, err
}

func (d *deps) storeIdempotent(ctx context.Context, key string, orderID int64) {
	d.db.ExecContext(ctx, `
		INSERT INTO idempotency_keys (key, order_id) VALUES ($1, $2)
		ON CONFLICT (key) DO NOTHING
	`, key, orderID)
}
