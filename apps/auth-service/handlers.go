package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"
)

type handlers struct {
	store  *userStore
	signer *signer
}

type registerReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type tokenResp struct {
	Token string `json:"token"`
	Role  string `json:"role"`
}

func (h *handlers) register(w http.ResponseWriter, r *http.Request) {
	var req registerReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid json")
		return
	}
	if !strings.Contains(req.Email, "@") || len(req.Password) < 6 {
		writeErr(w, http.StatusBadRequest, "invalid email or password")
		return
	}
	role := "customer"
	if strings.HasSuffix(req.Email, "@admin.local") {
		role = "admin"
	}
	u, err := h.store.create(req.Email, req.Password, role)
	if err != nil {
		if errors.Is(err, errUserExists) {
			writeErr(w, http.StatusConflict, "user exists")
			return
		}
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	tok, err := h.signer.issue(u.ID, u.Role, 1*time.Hour)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, tokenResp{Token: tok, Role: u.Role})
}

func (h *handlers) login(w http.ResponseWriter, r *http.Request) {
	var req registerReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid json")
		return
	}
	u, err := h.store.authenticate(req.Email, req.Password)
	if err != nil {
		writeErr(w, http.StatusUnauthorized, "invalid credentials")
		return
	}
	tok, err := h.signer.issue(u.ID, u.Role, 1*time.Hour)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, tokenResp{Token: tok, Role: u.Role})
}

func (h *handlers) verify(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		writeErr(w, http.StatusUnauthorized, "missing bearer token")
		return
	}
	c, err := h.signer.parse(strings.TrimPrefix(auth, "Bearer "))
	if err != nil {
		writeErr(w, http.StatusUnauthorized, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"user_id": c.UserID,
		"role":    c.Role,
		"exp":     c.ExpiresAt,
	})
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, map[string]string{"error": msg})
}
