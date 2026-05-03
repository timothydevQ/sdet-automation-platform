package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

func mustEnv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("%s required", k)
	}
	return v
}

func makeProxy(target string) *httputil.ReverseProxy {
	u, err := url.Parse(target)
	if err != nil {
		log.Fatalf("bad target %s: %v", target, err)
	}
	return httputil.NewSingleHostReverseProxy(u)
}

func main() {
	authURL := mustEnv("AUTH_URL")
	catalogURL := mustEnv("CATALOG_URL")
	orderURL := mustEnv("ORDER_URL")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	authProxy := makeProxy(authURL)
	catalogProxy := makeProxy(catalogURL)
	orderProxy := makeProxy(orderURL)

	rl := newRateLimiter(100, time.Minute)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Idempotency-Key")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		path := r.URL.Path

		switch {
		case path == "/healthz":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ok"))

		case strings.HasPrefix(path, "/auth/"):
			r.URL.Path = strings.TrimPrefix(path, "/auth")
			if r.URL.Path == "" {
				r.URL.Path = "/"
			}
			authProxy.ServeHTTP(w, r)

		case strings.HasPrefix(path, "/catalog/"):
			r.URL.Path = strings.TrimPrefix(path, "/catalog")
			if r.URL.Path == "" {
				r.URL.Path = "/"
			}
			catalogProxy.ServeHTTP(w, r)

		default:
			// cart, checkout, orders, admin — all go to order-service unchanged
			orderProxy.ServeHTTP(w, r)
		}
	})

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           rl.middleware(handler),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("api-gateway listening on :%s", port)
	log.Fatal(srv.ListenAndServe())
}

type rateLimiter struct {
	mu      sync.Mutex
	max     int
	window  time.Duration
	buckets map[string][]time.Time
}

func newRateLimiter(max int, window time.Duration) *rateLimiter {
	return &rateLimiter{max: max, window: window, buckets: make(map[string][]time.Time)}
}

func (r *rateLimiter) middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		key := clientKey(req)
		r.mu.Lock()
		now := time.Now()
		cutoff := now.Add(-r.window)
		hits := r.buckets[key]
		fresh := hits[:0]
		for _, t := range hits {
			if t.After(cutoff) {
				fresh = append(fresh, t)
			}
		}
		if len(fresh) >= r.max {
			r.buckets[key] = fresh
			r.mu.Unlock()
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		r.buckets[key] = append(fresh, now)
		r.mu.Unlock()
		h.ServeHTTP(w, req)
	})
}

func clientKey(r *http.Request) string {
	// BUG (intentional): trusts X-Forwarded-For unconditionally, allowing
	// rate-limit bypass by rotating that header. Caught by test_rate_limit_bypass_via_xff.
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}
	return r.RemoteAddr
}
