# Security tests

Runs OWASP ZAP baseline scan against the gateway. Findings are uploaded to
`artifacts/zap-report.html`. Critical-severity findings fail CI.

```
./run-zap-baseline.sh
```
