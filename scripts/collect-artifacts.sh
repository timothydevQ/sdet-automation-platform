#!/usr/bin/env bash
set -euo pipefail
ART=${ARTIFACTS_DIR:-artifacts}
mkdir -p "$ART"

# Pytest
cp tests/api-python/pytest-junit.xml "$ART/" 2>/dev/null || true
# Playwright
cp -r tests/e2e-playwright-ts/playwright-report "$ART/" 2>/dev/null || true
# Surefire
cp -r tests/selenium-java/target/surefire-reports "$ART/surefire-reports" 2>/dev/null || true
# NUnit
cp tests/dotnet-nunit/TestResults/*.trx "$ART/" 2>/dev/null || true
# RSpec
cp tests/ruby-rspec/rspec-junit.xml "$ART/" 2>/dev/null || true

echo "collected -> $ART"
ls -la "$ART"
