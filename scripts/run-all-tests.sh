#!/usr/bin/env bash
set -euo pipefail
./scripts/run-smoke.sh
./scripts/run-regression.sh
make test-performance || true
make test-security || true
make analyze-results || true
make flaky-report || true
