package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/segmentio/kafka-go"
)

type event struct {
	Type      string         `json:"type"`
	Payload   map[string]any `json:"payload"`
	Timestamp time.Time      `json:"timestamp"`
}

func main() {
	brokers := os.Getenv("KAFKA_BROKERS")
	topic := os.Getenv("KAFKA_TOPIC")
	if brokers == "" || topic == "" {
		log.Fatal("KAFKA_BROKERS and KAFKA_TOPIC required")
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: strings.Split(brokers, ","),
		Topic:   topic,
		GroupID: "notification-service",
		MinBytes: 1,
		MaxBytes: 10e6,
	})
	defer r.Close()

	ctx, cancel := context.WithCancel(context.Background())
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	dedupe := newSeen(10000)
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		log.Printf("consuming %s", topic)
		for {
			m, err := r.ReadMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				log.Printf("read: %v", err)
				continue
			}
			var e event
			if err := json.Unmarshal(m.Value, &e); err != nil {
				log.Printf("unmarshal: %v", err)
				continue
			}
			id, _ := e.Payload["id"].(float64)
			key := e.Type + ":" + jsonNum(id)
			if dedupe.Has(key) {
				log.Printf("DROP duplicate %s", key)
				continue
			}
			dedupe.Add(key)
			handle(e)
		}
	}()

	<-stop
	cancel()
	wg.Wait()
}

func handle(e event) {
	switch e.Type {
	case "order.created":
		log.Printf("[email] order created: %v", e.Payload["id"])
	case "order.refunded":
		log.Printf("[email] order refunded: %v", e.Payload["id"])
	default:
		log.Printf("unknown event: %s", e.Type)
	}
}

type seen struct {
	mu   sync.Mutex
	max  int
	keys map[string]time.Time
}

func newSeen(max int) *seen {
	return &seen{max: max, keys: make(map[string]time.Time)}
}

func (s *seen) Has(k string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.keys[k]
	return ok
}

func (s *seen) Add(k string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.keys[k] = time.Now()
	if len(s.keys) > s.max {
		oldest := time.Now()
		var ok string
		for kk, t := range s.keys {
			if t.Before(oldest) {
				oldest = t
				ok = kk
			}
		}
		delete(s.keys, ok)
	}
}

func jsonNum(f float64) string {
	b, _ := json.Marshal(f)
	return string(b)
}
