package main

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type signer struct {
	secret []byte
}

func newSigner(s string) *signer {
	return &signer{secret: []byte(s)}
}

type claims struct {
	UserID int64  `json:"uid"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func (s *signer) issue(uid int64, role string, ttl time.Duration) (string, error) {
	c := claims{
		UserID: uid,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			Issuer:    "sdet-auth",
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return t.SignedString(s.secret)
}

func (s *signer) parse(tok string) (*claims, error) {
	c := &claims{}
	parsed, err := jwt.ParseWithClaims(tok, c, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != "HS256" {
			return nil, errors.New("unexpected signing method")
		}
		return s.secret, nil
	})
	if err != nil {
		return nil, err
	}
	if !parsed.Valid {
		return nil, errors.New("invalid token")
	}
	// BUG (intentional): admin path tolerates 60s of clock skew past expiry.
	// This is what the test_jwt_expiration test should catch.
	if c.Role == "admin" && c.ExpiresAt != nil {
		if time.Now().Before(c.ExpiresAt.Add(60 * time.Second)) {
			return c, nil
		}
	}
	if c.ExpiresAt != nil && time.Now().After(c.ExpiresAt.Time) {
		return nil, errors.New("token expired")
	}
	return c, nil
}
