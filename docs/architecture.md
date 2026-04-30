# Architecture

A small e-commerce platform stitched together so the test suites have something
real to exercise. The shape is deliberately polyglot — Go services for the
backend, C++ for the pricing engine, React + TypeScript for the web UI. The
test layers wrap that with as many stacks as a serious QA org would actually use.

## Service boundaries

| Service              | Responsibility                                  | Stack |
|----------------------|-------------------------------------------------|-------|
| api-gateway          | Edge: routing, CORS, rate limiting              | Go    |
| auth-service         | Register, login, JWT issuance and verification  | Go + Postgres |
| catalog-service      | Products, search, stock reserve/release         | Go + Postgres |
| order-service        | Cart, checkout, idempotency, admin operations   | Go + Postgres + Redis + Kafka |
| payment-service      | Mock payment provider with test tokens          | Go    |
| notification-service | Kafka consumer for order events                 | Go + Kafka |
| pricing-engine-cpp   | Coupons, tax, totals (called from order-service)| C++   |
| web                  | React UI                                        | TS + Vite |

## Data stores

- **Postgres**: source of truth for users, products, orders.
- **Redis**: ephemeral cart state per user.
- **Kafka (Redpanda)**: order events for the notification consumer.

## Why it looks like this

The goal isn't to ship an e-commerce platform. The goal is to give each test
suite real surface area:

- **API tests** want JSON endpoints with state, auth, and edge cases.
- **E2E tests** want a real UI calling real APIs.
- **Performance tests** want a checkout path with caches, queues, and
  downstream calls — anything that would actually move under load.
- **Contract tests** want an OpenAPI spec and an implementation drifting from
  it occasionally.
- **Security tests** want a public surface with auth, rate limits, and headers.
- **Native tests** want native code; the C++ pricing engine is genuinely
  invoked by order-service in production-style deployments and unit-tested in
  isolation by GoogleTest.

## Operational concerns

- Health endpoints on every service (`/healthz`).
- All services log in the same line-oriented format with status, path, latency.
- Prometheus scrapes are stubbed; dashboards are placeholders. Real
  observability is out of scope for this iteration.
