package main

import "database/sql"

func migrate(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id BIGSERIAL PRIMARY KEY,
			email TEXT UNIQUE NOT NULL,
			password_hash BYTEA NOT NULL,
			role TEXT NOT NULL DEFAULT 'customer',
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
		CREATE INDEX IF NOT EXISTS users_email_idx ON users (email);
	`)
	return err
}
