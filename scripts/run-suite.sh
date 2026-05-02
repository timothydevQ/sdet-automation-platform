#!/usr/bin/env bash
set -euo pipefail
SUITE=${1:?usage: $0 <suite> [profile]}
PROFILE=${2:-smoke}

case "$SUITE" in
  api)
    cd tests/api-python
    if [ "$PROFILE" = smoke ]; then
      python -m pytest -m smoke --junitxml=../../artifacts/pytest-junit.xml
    else
      python -m pytest --junitxml=../../artifacts/pytest-junit.xml
    fi
    ;;
  e2e)
    cd tests/e2e-playwright-ts
    npm install --silent
    npx playwright install --with-deps chromium
    if [ "$PROFILE" = smoke ]; then
      npx playwright test --grep @smoke --project=chromium
    else
      npx playwright test
    fi
    ;;
  java)
    cd tests/selenium-java && mvn -B test
    ;;
  dotnet)
    cd tests/dotnet-nunit && dotnet test --logger "junit;LogFilePath=../../artifacts/nunit-junit.xml"
    ;;
  ruby)
    cd tests/ruby-rspec && bundle install --quiet && bundle exec rspec --format RspecJunitFormatter --out ../../artifacts/rspec-junit.xml
    ;;
  go)
    cd tests/go-integration && go test -v -json ./... > ../../artifacts/go-test.json
    ;;
  cpp)
    cd apps/pricing-engine-cpp
    cmake -S . -B build -DBUILD_TESTING=ON
    cmake --build build
    ctest --test-dir build --output-on-failure
    ;;
  *)
    echo "unknown suite: $SUITE" >&2
    exit 2
    ;;
esac
