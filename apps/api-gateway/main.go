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

type config struct {
	Auth    string
	Catalog string
	Order   string
}

func loadConfig() config {
	return config{
		Auth:    must("AUTH_URL"),
		Catalog: must("CATALOG_URL"),
		Order:   must("ORDER_URL"),
	}
}

func must(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("%s required", k)
	}
	return v
}

func main() {
	cfg := loadConfig()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()

	// Auth: strip /auth prefix before forwarding
	mux.Handle("/auth/", stripAndProxy("/auth", cfg.Auth))

	// Catalog: strip /catalog prefix before forwarding
	mux.Handle("/catalog/", stripAndProxy("/catalog", cfg.Catalog))

	// Order service: forward full path (order-service owns these prefixes)
	mux.Handle("/cart", proxyTo(cfg.Order))
	mux.Handle("/cart/", proxyTo(cfg.Order))
	mux.Handle("/checkout", proxyTo(cfg.Order))
	mux.Handle("/orders/", proxyTo(cfg.Order))
	mux.Handle("/admin/", proxyTo(cfg.Order))

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	rl := newRateLimiter(100, time.Minute)
	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           rl.middleware(withCORS(mux)),
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Printf("api-gateway listening on :%s", port)
	log.Fatal(srv.ListenAndServe())
}

// stripAndProxy removes the prefix from the path before forwarding.
func stripAndProxy(prefix, target string) http.Handler {
	u, err := url.Parse(target)
	if err != nil {
		log.Fatalf("bad target %s: %v", target, err)
	}
	rp := httputil.NewSingleHostReverseProxy(u)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.TrimPrefix(r.URL.Path, prefix)
		if r.URL.Path == "" {
			r.URL.Path = "/"
		}
		rp.ServeHTTP(w, r)
	})
}

// proxyTo forwards the request unchanged.
func proxyTo(target string) http.Handler {
	u, err := url.Parse(target)
	if err != nil {
		log.Fatalf("bad target %s: %v", target, err)
	}
	return httputil.NewSingleHostReverseProxy(u)
}

func withCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Idempotency-Key")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		h.ServeHTTP(w, r)
	})
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
