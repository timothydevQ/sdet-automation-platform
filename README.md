# SDET Automation Platform

A polyglot test automation platform built around a microservice e-commerce application. Demonstrates API testing, end-to-end UI testing, contract testing, performance testing, accessibility checks, native code testing, and a small platform layer for report ingestion and flaky-test analytics.

## Why this exists

Most test repos show a single suite in a single language. This one shows the shape of a real QA org: multiple test layers, multiple stacks, CI matrix execution, artifact collection, and analytics on top of the results.

## Architecture

```
                     ┌──────────────┐
                     │ web (React)  │
                     └──────┬───────┘
                            │
                     ┌──────▼───────┐
                     │ api-gateway  │
                     └──┬─────┬─────┘
              ┌─────────┘     └─────────┐
        ┌─────▼─────┐             ┌─────▼─────┐
        │   auth    │             │  catalog  │
        └───────────┘             └─────┬─────┘
                                        │
                              ┌─────────▼────────┐
                              │      order       │◄──┐
                              └────┬────────┬────┘   │
                              ┌────▼───┐ ┌──▼─────┐  │
                              │payment │ │pricing │──┘
                              └────────┘ │ (C++)  │
                                         └────────┘
                              ┌─────────────────┐
                              │ notifications   │  (Kafka consumer)
                              └─────────────────┘
```

## Test layers

| Layer         | Stack                        | Location                    |
|---------------|------------------------------|-----------------------------|
| Unit          | Go testing, GoogleTest       | `apps/*`, `tests/cpp-gtest` |
| API           | pytest + httpx               | `tests/api-python`          |
| E2E UI        | Playwright + TypeScript      | `tests/e2e-playwright-ts`   |
| Compat UI     | Selenium + Java + JUnit 5    | `tests/selenium-java`       |
| .NET          | NUnit + Playwright .NET      | `tests/dotnet-nunit`        |
| Smoke         | RSpec + Capybara             | `tests/ruby-rspec`          |
| Integration   | Go testing                   | `tests/go-integration`      |
| Native        | GoogleTest                   | `tests/cpp-gtest`           |
| Contract      | OpenAPI + schemathesis       | `tests/contract`            |
| Performance   | k6                           | `tests/performance`         |
| Security      | OWASP ZAP baseline           | `tests/security`            |
| Accessibility | axe-core + Playwright        | `tests/accessibility`       |

## Languages and tools

| Language    | Role                                           |
|-------------|------------------------------------------------|
| Go          | Backend services, integration tests            |
| Python      | API tests, report ingestor, flaky detector     |
| TypeScript  | Playwright E2E, web frontend                   |
| Java        | Selenium compat suite                          |
| C#          | NUnit API/UI tests                             |
| Ruby        | RSpec smoke suite                              |
| C++         | Pricing engine + GoogleTest                    |
| SQL         | Seeds, schema, assertions                      |

## Running locally

```sh
make up           # docker compose up -d --build
make seed         # seed postgres
make test-smoke   # fast smoke across stacks
make test-api     # full pytest API suite
make test-e2e     # full Playwright suite
make test-all     # everything
```

## CI/CD

Three workflow shapes:

- **PR**: lint, unit, API smoke, Playwright smoke, Selenium smoke, NUnit smoke. Fails fast on hot paths, runs everything else in matrix with `fail-fast: false`.
- **Nightly**: full regression across every suite, all browsers, performance, security, accessibility, flaky-test analysis.
- **Release**: contract tests, full regression, performance budget check, release-readiness report.

Artifacts uploaded on every run: JUnit XML, Playwright HTML report, traces, screenshots, videos, k6 JSON summary, ZAP report, and the platform's generated `flaky-tests.md` and `release-readiness.md`.

## Platform layer

`platform/report-ingestor` parses JUnit XML, Playwright JSON, NUnit XML, k6 summaries, and Go test JSON into a normalized schema. `platform/flaky-detector` reads from that store and computes per-test flake scores using failure rate, retry-pass rate, duration variance, and recent failure streak. Outputs:

- `flaky-tests.md` — ranked list with classification (stable / slow / flaky / broken / quarantined)
- `release-readiness.md` — go/no-go summary based on suite health
- `ci-summary.md` — markdown summary for the GitHub Actions step output

## Seeded bug scenarios

The application has a handful of intentionally broken paths to give the test suites something real to catch:

- Duplicate checkout creates two orders if idempotency key is reused incorrectly
- Expired JWTs accepted on a specific admin endpoint (off-by-one on exp check)
- Inventory can go negative under concurrent checkout
- Discount rounding off-by-one cent on certain bulk thresholds
- Rate limit bypass via `X-Forwarded-For` rotation
- Order confirmation event published twice on retry

Each is documented in `docs/test-strategy.md` with the catching test.

## Future work

- Mobile testing (Appium)
- Service virtualization for payment provider
- Mutation testing
- Visual regression
- Multi-region deploy + chaos
