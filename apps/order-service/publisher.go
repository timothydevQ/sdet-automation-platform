package main

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

type publisher struct {
	w *kafka.Writer
}

func newPublisher(brokers, topic string) *publisher {
	return &publisher{
		w: &kafka.Writer{
			Addr:         kafka.TCP(strings.Split(brokers, ",")...),
			Topic:        topic,
			Balancer:     &kafka.LeastBytes{},
			RequiredAcks: kafka.RequireOne,
			BatchTimeout: 50 * time.Millisecond,
		},
	}
}

func (p *publisher) Publish(ctx context.Context, eventType string, payload any) {
	body, _ := json.Marshal(map[string]any{
		"type":      eventType,
		"payload":   payload,
		"timestamp": time.Now().UTC(),
	})
	if err := p.w.WriteMessages(ctx, kafka.Message{Value: body}); err != nil {
		log.Printf("publish: %v", err)
	}
}

func (p *publisher) Close() error {
	return p.w.Close()
}
