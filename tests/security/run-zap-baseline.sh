#!/usr/bin/env bash
set -euo pipefail
TARGET=${TARGET:-http://localhost:8080}
ARTIFACTS=${ARTIFACTS:-$(pwd)/../../artifacts}
mkdir -p "$ARTIFACTS"

docker run --rm \
  --network host \
  -v "$(pwd)":/zap/wrk \
  -v "$ARTIFACTS":/artifacts \
  ghcr.io/zaproxy/zaproxy:stable \
  zap-baseline.py \
    -t "$TARGET" \
    -c zap-baseline.conf \
    -r /artifacts/zap-report.html \
    -J /artifacts/zap-report.json \
    || EXIT=$?

# Baseline scan: warnings allowed, fails only on FAIL-level rules.
exit ${EXIT:-0}
