package integration_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"
)

var apiBase = func() string {
	if v := os.Getenv("API_BASE"); v != "" {
		return v
	}
	return "http://localhost:8080"
}()

func registerCustomer(t *testing.T) (string, string) {
	t.Helper()
	email := fmt.Sprintf("u%d@test.local", time.Now().UnixNano())
	body, _ := json.Marshal(map[string]string{"email": email, "password": "Hunter22!"})
	resp, err := http.Post(apiBase+"/auth/register", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("register: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("register status %d: %s", resp.StatusCode, b)
	}
	var out struct{ Token string }
	json.NewDecoder(resp.Body).Decode(&out)
	return email, out.Token
}

func TestHealthEndpoints(t *testing.T) {
	resp, err := http.Get(apiBase + "/healthz")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Errorf("gateway healthz: %d", resp.StatusCode)
	}
}

func TestRegisterAndCheckout(t *testing.T) {
	_, token := registerCustomer(t)
	addCart(t, token, "SKU-001", 1)
	resp := checkout(t, token, "", "tok_test_visa", "")
	if resp.StatusCode != 201 {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("checkout: %d %s", resp.StatusCode, b)
	}
}

func TestIdempotency(t *testing.T) {
	_, token := registerCustomer(t)
	addCart(t, token, "SKU-001", 1)
	key := fmt.Sprintf("idem-%d", time.Now().UnixNano())
	r1 := checkout(t, token, "", "tok_test_visa", key)
	if r1.StatusCode != 201 {
		t.Fatalf("first checkout: %d", r1.StatusCode)
	}
	addCart(t, token, "SKU-001", 1)
	r2 := checkout(t, token, "", "tok_test_visa", key)
	if r2.StatusCode != 200 {
		t.Errorf("idempotent replay should be 200, got %d", r2.StatusCode)
	}
}

func TestConcurrentCheckoutSameUser(t *testing.T) {
	_, token := registerCustomer(t)
	addCart(t, token, "SKU-001", 1)
	var wg sync.WaitGroup
	results := make([]int, 5)
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			r := checkout(t, token, "", "tok_test_visa", "")
			results[i] = r.StatusCode
		}(i)
	}
	wg.Wait()
	t.Logf("results: %v", results)
}

func addCart(t *testing.T, token, sku string, qty int) {
	body, _ := json.Marshal(map[string]any{"sku": sku, "qty": qty})
	req, _ := http.NewRequest("POST", apiBase+"/cart/items", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("addCart: %v", err)
	}
	resp.Body.Close()
}

func checkout(t *testing.T, token, coupon, card, idem string) *http.Response {
	body, _ := json.Marshal(map[string]any{"coupon": coupon, "card_token": card})
	req, _ := http.NewRequest("POST", apiBase+"/checkout", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	if idem != "" {
		req.Header.Set("Idempotency-Key", idem)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("checkout: %v", err)
	}
	return resp
}
