package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type productInfo struct {
	ID    int64  `json:"id"`
	SKU   string `json:"sku"`
	Name  string `json:"name"`
	Price int64  `json:"price_cents"`
	Stock int    `json:"stock"`
}

var httpClient = &http.Client{Timeout: 5 * time.Second}

func (d *deps) fetchProductBySKU(ctx context.Context, sku string) (*productInfo, error) {
	url := fmt.Sprintf("%s/products?q=%s", d.catalogURL, sku)
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("catalog status %d", resp.StatusCode)
	}
	var ps []productInfo
	if err := json.NewDecoder(resp.Body).Decode(&ps); err != nil {
		return nil, err
	}
	for _, p := range ps {
		if p.SKU == sku {
			return &p, nil
		}
	}
	return nil, fmt.Errorf("sku %s not found", sku)
}

func (d *deps) reserve(ctx context.Context, sku string, qty int) error {
	p, err := d.fetchProductBySKU(ctx, sku)
	if err != nil {
		return err
	}
	body, _ := json.Marshal(map[string]int{"qty": qty})
	url := fmt.Sprintf("%s/products/%d/reserve", d.catalogURL, p.ID)
	req, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("reserve status %d", resp.StatusCode)
	}
	return nil
}

func (d *deps) release(ctx context.Context, sku string, qty int) error {
	p, err := d.fetchProductBySKU(ctx, sku)
	if err != nil {
		return err
	}
	body, _ := json.Marshal(map[string]int{"qty": qty})
	url := fmt.Sprintf("%s/products/%d/release", d.catalogURL, p.ID)
	req, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := httpClient.Do(req)
	if resp != nil {
		resp.Body.Close()
	}
	return nil
}

type paymentResp struct {
	Approved bool   `json:"approved"`
	Ref      string `json:"ref"`
}

func (d *deps) charge(ctx context.Context, card string, amountCents int64) (*paymentResp, error) {
	body, _ := json.Marshal(map[string]any{"card_token": card, "amount_cents": amountCents})
	req, _ := http.NewRequestWithContext(ctx, "POST", d.paymentURL+"/charge", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var pr paymentResp
	if err := json.NewDecoder(resp.Body).Decode(&pr); err != nil {
		return nil, err
	}
	return &pr, nil
}
