#!/usr/bin/env bash
set -euo pipefail
mkdir -p artifacts

echo "==> API smoke"
cd tests/api-python && python -m pytest -m smoke --junitxml=../../artifacts/pytest-smoke.xml -q
cd - >/dev/null

echo "==> Playwright smoke"
cd tests/e2e-playwright-ts && npx playwright test --grep @smoke --project=chromium
cd - >/dev/null

echo "==> Java smoke"
cd tests/selenium-java && mvn -B test -Dgroups=smoke
cd - >/dev/null

echo "smoke complete."
