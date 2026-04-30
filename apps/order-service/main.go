package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

type deps struct {
	db          *sql.DB
	rdb         *redis.Client
	catalogURL  string
	paymentURL  string
	authURL     string
	publisher   *publisher
}

func main() {
	cfg := struct {
		dsn, redis, kafka, catalog, payment, auth, port string
	}{
		dsn:     mustEnv("DB_DSN"),
		redis:   mustEnv("REDIS_ADDR"),
		kafka:   mustEnv("KAFKA_BROKERS"),
		catalog: mustEnv("CATALOG_URL"),
		payment: mustEnv("PAYMENT_URL"),
		auth:    mustEnv("AUTH_URL"),
		port:    envOr("PORT", "8084"),
	}

	db, err := sql.Open("postgres", cfg.dsn)
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < 30; i++ {
		if err := db.Ping(); err == nil {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	if err := migrate(db); err != nil {
		log.Fatal(err)
	}

	rdb := redis.NewClient(&redis.Options{Addr: cfg.redis})
	pub := newPublisher(cfg.kafka, "orders.events")
	defer pub.Close()

	d := &deps{
		db:         db,
		rdb:        rdb,
		catalogURL: cfg.catalog,
		paymentURL: cfg.payment,
		authURL:    cfg.auth,
		publisher:  pub,
	}

	mux := http.NewServeMux()
	mux.Handle("POST /cart/items", auth(d, http.HandlerFunc(d.addToCart)))
	mux.Handle("GET /cart", auth(d, http.HandlerFunc(d.getCart)))
	mux.Handle("DELETE /cart/items/{sku}", auth(d, http.HandlerFunc(d.removeFromCart)))
	mux.Handle("POST /checkout", auth(d, http.HandlerFunc(d.checkout)))
	mux.Handle("GET /orders/{id}", auth(d, http.HandlerFunc(d.getOrder)))
	mux.Handle("GET /admin/orders", requireRole(d, "admin", http.HandlerFunc(d.listOrders)))
	mux.Handle("POST /admin/orders/{id}/refund", requireRole(d, "admin", http.HandlerFunc(d.refundOrder)))
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	log.Printf("order-service listening on :%s", cfg.port)
	log.Fatal(http.ListenAndServe(":"+cfg.port, mux))
}

func mustEnv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("%s required", k)
	}
	return v
}

func envOr(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

var _ = context.Background
