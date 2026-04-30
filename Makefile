SHELL := /bin/bash
ARTIFACTS := $(CURDIR)/artifacts

.PHONY: up down restart seed reset logs ps \
	test-all test-smoke test-regression \
	test-api test-e2e test-java test-dotnet test-ruby test-go test-cpp \
	test-performance test-security test-accessibility test-contract \
	analyze-results flaky-report release-readiness

up:
	docker compose up -d --build

down:
	docker compose down -v

restart: down up

seed:
	./scripts/seed-db.sh

reset:
	./scripts/reset-env.sh

logs:
	docker compose logs -f

ps:
	docker compose ps

test-smoke:
	./scripts/run-smoke.sh

test-regression:
	./scripts/run-regression.sh

test-all: test-api test-e2e test-java test-dotnet test-ruby test-go test-cpp

test-api:
	mkdir -p $(ARTIFACTS)
	cd tests/api-python && python -m pytest -m "smoke or regression" \
		--junitxml=$(ARTIFACTS)/pytest-junit.xml \
		--html=$(ARTIFACTS)/pytest-report.html --self-contained-html

test-e2e:
	mkdir -p $(ARTIFACTS)
	cd tests/e2e-playwright-ts && npm ci && npx playwright test \
		--reporter=html,junit \
		&& cp -r playwright-report $(ARTIFACTS)/playwright-report || true

test-java:
	cd tests/selenium-java && mvn -B test

test-dotnet:
	cd tests/dotnet-nunit && dotnet test --logger "trx;LogFileName=nunit-results.trx" \
		--logger "junit;LogFilePath=$(ARTIFACTS)/nunit-junit.xml"

test-ruby:
	cd tests/ruby-rspec && bundle install --quiet && bundle exec rspec \
		--format documentation --format RspecJunitFormatter \
		--out $(ARTIFACTS)/rspec-junit.xml

test-go:
	mkdir -p $(ARTIFACTS)
	cd tests/go-integration && go test -v -json ./... > $(ARTIFACTS)/go-test.json

test-cpp:
	cd apps/pricing-engine-cpp && cmake -S . -B build -DBUILD_TESTING=ON \
		&& cmake --build build && ctest --test-dir build --output-on-failure

test-contract:
	cd tests/contract && python -m pytest --junitxml=$(ARTIFACTS)/contract-junit.xml

test-performance:
	cd tests/performance && k6 run --summary-export=$(ARTIFACTS)/k6-summary.json checkout-load.js

test-security:
	./tests/security/run-zap-baseline.sh

test-accessibility:
	cd tests/e2e-playwright-ts && npx playwright test specs/accessibility.spec.ts

analyze-results:
	cd platform/report_ingestor && python -m report_ingestor $(ARTIFACTS)

flaky-report:
	cd platform/flaky_detector && python -m flaky_detector --out $(ARTIFACTS)/flaky-tests.md

release-readiness:
	cd platform/flaky_detector && python -m flaky_detector --release-check --out $(ARTIFACTS)/release-readiness.md
