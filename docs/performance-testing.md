# Performance testing

## What is measured

`tests/performance/checkout-load.js`: end-to-end checkout under a ramping load
profile. Covers register → cart → checkout. Records latency per step plus
overall throughput.

`tests/performance/catalog-spike.js`: rapid spike from 5 to 200 RPS on the
catalog endpoint. Used to verify the gateway and rate-limiter survive bursts.

## SLOs (enforced as k6 thresholds)

| Metric                                  | Threshold |
|-----------------------------------------|-----------|
| `http_req_failed`                       | < 2%      |
| `checkout_latency_ms` p95               | < 800 ms  |
| `checkout_latency_ms` p99               | < 1500 ms |
| `checkout_errors`                       | < 5%      |
| `http_req_duration` p95 (catalog spike) | < 500 ms  |

A failed threshold fails the workflow. Output is exported to
`artifacts/k6-summary.json` and read by the report ingestor into
`performance_results`.

## What's intentionally simple

- Single-region, single-VU pool. No long-soak yet.
- No data warmup. The first iterations are part of the measurement, which
  reflects cold-start behaviour.
- No comparison against a baseline budget file. Future work.

## Future

- Lock thresholds against historical p95 instead of static numbers.
- Run on-merge against main with budget alerts.
- Add a long-soak (60 min) once we have a separate environment to avoid
  blocking CI.
