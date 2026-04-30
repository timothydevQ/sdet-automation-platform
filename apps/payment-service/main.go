package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type chargeReq struct {
	CardToken   string `json:"card_token"`
	AmountCents int64  `json:"amount_cents"`
}

type chargeResp struct {
	Approved bool   `json:"approved"`
	Ref      string `json:"ref"`
	Reason   string `json:"reason,omitempty"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}
	failureRate := 0.05
	if v := os.Getenv("FAILURE_RATE"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			failureRate = f
		}
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /charge", handleCharge(failureRate))
	mux.HandleFunc("POST /refund", handleRefund)
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	log.Printf("payment-service listening on :%s (failure rate %.2f)", port, failureRate)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func handleCharge(failureRate float64) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req chargeReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, 400, chargeResp{Approved: false, Reason: "bad request"})
			return
		}
		if req.AmountCents <= 0 {
			writeJSON(w, 400, chargeResp{Approved: false, Reason: "invalid amount"})
			return
		}
		if strings.HasPrefix(req.CardToken, "tok_decline_") {
			writeJSON(w, 200, chargeResp{Approved: false, Reason: "declined by issuer"})
			return
		}
		if strings.HasPrefix(req.CardToken, "tok_timeout_") {
			time.Sleep(10 * time.Second)
			writeJSON(w, 504, chargeResp{Approved: false, Reason: "timeout"})
			return
		}
		if randFloat() < failureRate {
			writeJSON(w, 200, chargeResp{Approved: false, Reason: "transient failure"})
			return
		}
		writeJSON(w, 200, chargeResp{Approved: true, Ref: newRef()})
	}
}

func handleRefund(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Ref         string `json:"ref"`
		AmountCents int64  `json:"amount_cents"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	writeJSON(w, 200, map[string]any{"refunded": true, "ref": req.Ref})
}

func newRef() string {
	b := make([]byte, 8)
	rand.Read(b)
	return "pay_" + hex.EncodeToString(b)
}

func randFloat() float64 {
	n, _ := rand.Int(rand.Reader, big.NewInt(10000))
	return float64(n.Int64()) / 10000.0
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}
