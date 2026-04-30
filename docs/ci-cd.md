# CI/CD

## Three workflow shapes

**PR (`ci.yml`)**: lint, build, smoke matrix across the four loudest suites
(api, e2e, java, dotnet). `fail-fast: false` so we see all failures, not just
the first. ~8 minutes target.

**Nightly (`nightly-regression.yml`)**: full suites across every stack, then
the analysis job ingests results into the platform store and produces
`flaky-tests.md` and `release-readiness.md`.

**On-demand**: `api-tests`, `e2e-tests`, `performance-tests`, `security-tests`
can each be triggered manually or on weekly schedules.

## Artifacts

Every job uploads `artifacts/`. The directory is the canonical drop zone — no
test should write its junit XML anywhere else. CI never tries to interpret
trace files; that's the platform layer's job.

## What blocks merge

- Lint or build failure on any service.
- Any smoke test failure (PR matrix).
- Schemathesis contract violation.

## What does NOT block merge

- Performance regressions (reported, not enforced on PR — enforced on main).
- Accessibility violations below `critical`.
- Any test marked `flaky` that hasn't been quarantined.

## Why distinguish

Strict gates that fail constantly get bypassed. Soft gates that report
truthfully and surface in the readiness report stay informative.
