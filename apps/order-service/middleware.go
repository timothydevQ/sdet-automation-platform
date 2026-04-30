package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

type ctxKey string

const (
	ctxUserID ctxKey = "user_id"
	ctxRole   ctxKey = "role"
)

func auth(d *deps, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tok := r.Header.Get("Authorization")
		if !strings.HasPrefix(tok, "Bearer ") {
			http.Error(w, "missing token", 401)
			return
		}
		req, _ := http.NewRequestWithContext(r.Context(), "GET", d.authURL+"/verify", nil)
		req.Header.Set("Authorization", tok)
		resp, err := httpClient.Do(req)
		if err != nil {
			http.Error(w, "auth unavailable", 503)
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			body, _ := io.ReadAll(resp.Body)
			http.Error(w, string(body), resp.StatusCode)
			return
		}
		var v struct {
			UserID int64  `json:"user_id"`
			Role   string `json:"role"`
		}
		buf, _ := io.ReadAll(resp.Body)
		json.NewDecoder(bytes.NewReader(buf)).Decode(&v)
		ctx := context.WithValue(r.Context(), ctxUserID, v.UserID)
		ctx = context.WithValue(ctx, ctxRole, v.Role)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func requireRole(d *deps, role string, h http.Handler) http.Handler {
	return auth(d, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got, _ := r.Context().Value(ctxRole).(string)
		if got != role {
			// BUG (intentional): "admin" check is case-sensitive but tokens
			// can be issued with role "Admin" via certain registration paths.
			// Caught by test_admin_authorization_case.
			http.Error(w, "forbidden", 403)
			return
		}
		h.ServeHTTP(w, r)
	}))
}
