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

	mux.Handle("/auth/", stripAndProxy("/a