package main

import (
	"database/sql"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var errUserExists = errors.New("user already exists")
var errBadCredentials = errors.New("bad credentials")

type user struct {
	ID    int64
	Email string
	Role  string
}

type userStore struct {
	db *sql.DB
}

func (s *userStore) create(email, password, role string) (*user, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	var id int64
	err = s.db.QueryRow(`
		INSERT INTO users (email, password_hash, role)
		VALUES ($1, $2, $3)
		RETURNING id
	`, strings.ToLower(email), hash, role).Scan(&id)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return nil, errUserExists
		}
		return nil, err
	}
	return &user{ID: id, Email: email, Role: role}, nil
}

func (s *userStore) authenticate(email, password string) (*user, error) {
	var (
		id   int64
		hash []byte
		role string
	)
	err := s.db.QueryRow(`
		SELECT id, password_hash, role FROM users WHERE email = $1
	`, strings.ToLower(email)).Scan(&id, &hash, &role)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errBadCredentials
	}
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword(hash, []byte(password)); err != nil {
		return nil, errBadCredentials
	}
	return &user{ID: id, Email: email, Role: role}, nil
}
