# Platform walkthrough

A guide to running the platform and understanding what each layer does. Intended for anyone reviewing the repo who wants to see it working, not just read the code.

## 1. Start the stack

```sh
docker compose up -d --build
```

Ten containers come up: postgres, redis, kafka, auth-service, catalog-service, payment-service, order-service, notification-service, api-gateway, web. All services expose health endpoints at `/healthz`.

Open http://localhost:3000 — you should see the shop catalog.

## 2. Register and place an order

1. Go to http://localhost:3000/login
2. Click **Create account**, register with any email + 6+ char password
3. Browse the catalog, add a product to cart
4. Go to Cart → Checkout
5. Card token `tok_test_visa` → succeeds. `tok_decline_card` → 402 declined.

For admin access, register with an email ending in `@admin.local`. The Admin nav link shows the order list with refund controls.

## 3. Run the API smoke suite

```sh
cd tests/api-python
pip install -r requirements.txt
API_BASE=http://localhost:8080 python -m pytest -m smoke -v
```

Expected: login passes, checkout passes, admin authorization passes. Two tests intentionally fail — they assert *desired* behavior against *known bugs*:

- `test_rate_limit_bypass_via_xff` — asserts XFF rotation should be blocked (it isn't yet)
- `test_admin_expired_token_skew_bug` — asserts expired admin tokens should be rejected (they aren't yet)

This is intentional. The tests document what needs fixing, not what currently works.

## 4. Understand the intentional bugs

Every bug is marked `BUG (intentional)` in the source. Each has a test that fails until the bug is fixed.

| Bug | File | Catching test |
|-----|------|---------------|
| Admin JWT 60s clock-skew tolerance | `apps/auth-service/jwt.go` | `test_admin_expired_token_skew_bug` |
| Rate-limit bypass via X-Forwarded-For | `apps/api-gateway/main.go` | `test_rate_limit_bypass_via_xff` |
| LIKE pattern doesn't escape `%` and `_` | `apps/catalog-service/main.go` | `test_search_special_chars` |
| Non-atomic stock decrement | `apps/catalog-service/main.go` | `test_concurrent_checkout_inventory_race` |
| BULK20 off-by-one rounding | `apps/order-service/orders.go` | `test_bulk20_rounding_bug` |
| `order.created` published twice on retry | `apps/order-service/orders.go` | `TestIdempotency` (Go) |
| Admin role check case-sensitive | `apps/order-service/middleware.go` | `test_admin_role_check_is_case_sensitive` |

Fix one, run its test, watch it go green. That's the intended workflow.

## 5. Run Playwright E2E

```sh
cd tests/e2e-playwright-ts
npm install
npx playwright install chromium
API_BASE=http://localhost:8080 WEB_BASE_URL=http://localhost:3000 npx playwright test --project=chromium --headed
```

`--headed` shows the browser. Playwright will register a user, click through the catalog, add to cart, and complete checkout. On failure it captures a trace file you can open with `npx playwright show-trace`.

## 6. Run the C++ pricing tests

```sh
cd apps/pricing-engine-cpp
cmake -S . -B build -DBUILD_TESTING=ON
cmake --build build
ctest --test-dir build --output-on-failure
```

GoogleTest runs 16 cases covering coupons (WELCOME10, BULK20, FLAT5, VIP25), tax rounding, and invalid input handling. All should pass — the C++ engine has no intentional bugs; the rounding bug lives in the Go layer that calls it.

## 7. Run performance tests

Requires k6 (https://k6.io/docs/get-started/installation/).

```sh
cd tests/performance
k6 run --summary-export=../../artifacts/k6-summary.json checkout-load.js
```

Ramps from 1 to 25 virtual users over 2 minutes. Thresholds: p95 checkout latency < 800ms, error rate < 2%. A threshold failure exits non-zero and fails CI.

## 8. CI artifacts

On every CI run, the following are uploaded:

| Artifact | Source |
|----------|--------|
| `pytest-junit.xml` | API test results |
| `playwright-report/` | HTML report with traces and screenshots |
| `playwright-junit.xml` | Playwright JUnit XML |
| `surefire-reports/` | Java/Selenium results |
| `nunit-junit.xml` | .NET NUnit results |
| `rspec-junit.xml` | Ruby RSpec results |
| `go-test.json` | Go integration test JSON |
| `k6-summary.json` | Performance thresholds and metrics |
| `zap-report.html` | OWASP ZAP baseline scan |
| `flaky-tests.md` | Flaky test rankings (nightly only) |
| `release-readiness.md` | Go/No-go verdict (nightly only) |

## 9. Flaky test analytics

After a nightly run the platform layer ingests all results and scores each test:

```
score = 0.5 × failure_rate
      + 0.3 × pass_after_retry_rate
      + 0.1 × duration_coefficient_of_variation
      + 0.1 × min(1.0, recent_failure_streak / 5)
```

Tests scoring above 0.3 are `flaky`, above 0.7 are `quarantined`. The output is `artifacts/flaky-tests.md` — a ranked table with scores, failure rates, and streaks. `artifacts/release-readiness.md` gives a GO/NO-GO verdict: no-go if any test is `broken` (≥90% failure rate) or more than 5 are `flaky`.

## 10. Fix a bug (example walkthrough)

**Bug:** BULK20 coupon rounding off by one cent.

**Test that catches it:** `tests/api-python/tests/test_pricing_rounding.py::test_bulk20_rounding_bug`

**Run it:**
```sh
API_BASE=http://localhost:8080 python -m pytest tests/test_pricing_rounding.py::test_bulk20_rounding_bug -v
```
It fails with an assertion showing the actual vs expected total.

**Find the bug:** `apps/order-service/orders.go` function `applyDiscount()`, the `BULK20` case:
```go
return subtotal - (subtotal*20)/100  // wrong: truncates instead of rounds
```

**Fix it:**
```go
return subtotal - int64(math.Round(float64(subtotal)*0.20))
```

**Rebuild and retest:**
```sh
docker compose up -d --build order-service
API_BASE=http://localhost:8080 python -m pytest tests/test_pricing_rounding.py::test_bulk20_rounding_bug -v
```

Test goes green.
