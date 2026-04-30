# Security testing

## Layers

**Static (today): none.** Future: `gosec`, `bandit`, `npm audit`,
`dotnet list package --vulnerable`.

**Dynamic baseline:** OWASP ZAP baseline against the gateway. Findings are
classified `WARN`, `IGNORE`, `FAIL` in `tests/security/zap-baseline.conf`.
Only `FAIL` blocks CI.

**API fuzz:** schemathesis (in `tests/contract/`) drives random inputs through
each endpoint and checks the response matches the schema. This catches input
validation gaps without anyone writing per-endpoint negative tests.

## What ZAP catches in this codebase

- Missing security headers (CSP, X-Content-Type-Options, X-Frame-Options).
- Cookies without `Secure` / `HttpOnly` flags.
- CORS that allows everything (currently it does, intentionally — see ADR).

## Out of scope

- Authentication bypass attempts beyond what schemathesis surfaces.
- Penetration testing.
- Supply-chain SCA.

## Why baseline rather than full scan

Full ZAP scans take 30+ minutes and produce too much noise to gate on. The
baseline scan is fast and catches the obvious. Anything deeper belongs in a
dedicated security review, not a CI check.
