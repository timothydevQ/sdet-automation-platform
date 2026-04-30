# Flaky test process

## Detection

`platform/flaky-detector` reads from the ingested test history and computes a
score per test:

```
score = 0.5 * failure_rate
      + 0.3 * pass_after_retry_rate
      + 0.1 * min(1.0, duration_cv)
      + 0.1 * min(1.0, recent_failure_streak / 5.0)
```

## Classifications

| Class       | Threshold                                  | Treatment in CI                     |
|-------------|--------------------------------------------|-------------------------------------|
| stable      | score < 0.3 and failure rate < 90%         | Runs everywhere                     |
| slow        | high duration variance, low failure rate   | Runs everywhere; warning only       |
| flaky       | 0.3 <= score < 0.7                         | Runs everywhere; surfaced in report |
| broken      | failure rate >= 90%                        | Blocks release                      |
| quarantined | score >= 0.7                               | Skipped on PR, runs on main         |

## Workflow

1. Detector runs nightly after the regression matrix.
2. Output goes to `artifacts/flaky-tests.md` and as a comment on the nightly
   run summary.
3. Tests classified as `quarantined` are added to a quarantine list (not
   silently — the change is a PR with the score and link to the failing run).
4. Owners triage within 5 business days. After triage: fix, delete, or move to
   the long-term skip list with a Jira link.

## What we do not do

- Auto-retry on PR runs without surfacing it. Retries hide flake; they don't
  fix it. Pytest reruns are limited to `regression` runs only.
- Silently `xfail` tests. If a test is wrong, fix it or delete it.
