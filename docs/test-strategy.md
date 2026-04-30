# Test strategy

## Pyramid

```
              ▲
              │   E2E (Playwright)        ← business-critical UI flows only
              │
              │   API & Contract           ← bulk of regression
              │
              │   Unit (Go, C++)           ← every service / library
              ▼
```

## Suite ownership

Each suite has a clear remit. Suites that overlap exist for cross-stack
demonstration, not duplication.

| Suite                      | Owns                                             |
|----------------------------|--------------------------------------------------|
| Go unit                    | Per-service business logic                       |
| C++ GoogleTest             | Pricing engine math: coupons, tax, totals        |
| Python pytest              | Core API contract, all happy + negative paths    |
| Playwright (TS)            | UI flows that matter to a user                   |
| Selenium (Java)            | Compatibility smoke; demonstrates WebDriver use  |
| NUnit (.NET)               | API integration in a third stack                 |
| RSpec (Ruby + Capybara)    | Cross-language smoke                             |
| Go integration             | Concurrency, idempotency, cross-service paths    |
| Schemathesis (contract)    | OpenAPI fuzz                                     |
| k6 (performance)           | Latency budget, throughput                       |
| ZAP baseline (security)    | OWASP Top 10 surface checks                      |
| axe-core (accessibility)   | WCAG 2 AA on key pages                           |

## Markers / tags

`smoke`, `regression`, `auth`, `orders`, `catalog`, `admin`, `payments`,
`negative`, `flaky`. PRs run `smoke`; nightly runs everything.

## Seeded bug scenarios

Each one is wired to at least one test that should fail until the bug is fixed.

| Bug | Catching test |
|-----|---------------|
| Admin JWT 60s clock-skew tolerance | `tests/api-python/tests/test_jwt_expiration.py::test_admin_jwt_skew_bug` |
| Rate-limit bypass via X-Forwarded-For rotation | `tests/api-python/tests/test_rate_limit.py::test_rate_limit_bypass_via_xff` |
| LIKE pattern doesn't escape `%` and `_` | `tests/api-python/tests/test_search.py::test_search_special_chars` |
| Non-atomic stock decrement under concurrency | `tests/api-python/tests/test_inventory.py::test_concurrent_checkout_inventory` |
| BULK20 coupon off-by-one rounding | `tests/api-python/tests/test_pricing_rounding.py::test_bulk_discount_rounding` and `tests/cpp-gtest/test_coupons.cpp::Coupons_Bulk20OnlyAboveThreshold` |
| `order.created` event published twice on retry | `tests/go-integration/integration_test.go::TestIdempotency` (consumer-side dedupe in notification-service) |
| Admin role check is case-sensitive | `tests/api-python/tests/test_rate_limit.py::test_admin_authorization_case` |

## Flake policy

A test is `flaky` when its score (failure rate, retry-pass rate, duration
variance, recent streak) crosses 0.3, and `quarantined` at 0.7. Quarantined
tests are skipped in PR runs, run on main, and surfaced in the readiness
report. Never silently skipped.

## Release readiness

Generated nightly. A release is `NO-GO` if any `broken` test exists or there
are more than 5 `flaky` tests. Output goes to `artifacts/release-readiness.md`.
