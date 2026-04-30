# Performance tests (k6)

```
k6 run --summary-export=../../artifacts/k6-summary.json checkout-load.js
k6 run catalog-spike.js
```

Thresholds enforce SLOs in CI (p95 latency, error rate). Failure of a threshold
fails the workflow.
