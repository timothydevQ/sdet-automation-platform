#!/usr/bin/env bash
set -euo pipefail
SUITE=${1:?usage: $0 <suite> [profile]}
PROFILE=${2:-smoke}
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
ARTIFACTS="$ROOT/artifacts"
mkdir -p "$ARTIFACTS"

case "$SUITE" in
  api)
    cd "$ROOT/tests/api-python"
    pip install -q -r requirements.txt
    if [ "$PROFILE" = smoke ]; then
      python -m pytest -m smoke --junitxml="$ARTIFACTS/pytest-junit.xml" -v
    else
      python -m pytest --junitxml="$ARTIFACTS/pytest-junit.xml" -v
    fi
    ;;
  e2e)
    cd "$ROOT/tests/e2e-playwright-ts"
    npm install --silent
    npx playwright install --with-deps chromium
    if [ "$PROFILE" = smoke ]; then
      npx playwright test --grep @smoke --project=chromium
    else
      npx playwright test
    fi
    ;;
  java)
    cd "$ROOT/tests/selenium-java"
    mvn -B test -Dgroups=smoke
    ;;
  dotnet)
    cd "$ROOT/tests/dotnet-nunit"
    dotnet test --logger "junit;LogFilePath=$ARTIFACTS/nunit-junit.xml"
    ;;
  ruby)
    cd "$ROOT/tests/ruby-rspec"
    bundle install --quiet
    bundle exec rspec --format RspecJunitFormatter --out "$ARTIFACTS/rspec-junit.xml"
    ;;
  go)
    cd "$ROOT/tests/go-integration"
    go test -v -json ./... > "$ARTIFACTS/go-test.json"
    ;;
  cpp)
    cd "$ROOT/apps/pricing-engine-cpp"
    cmake -S . -B build -DBUILD_TESTING=ON
    cmake --build build
    ctest --test-dir build --output-on-failure
    ;;
  *)
    echo "unknown suite: $SUITE" >&2
    exit 2
    ;;
esac
