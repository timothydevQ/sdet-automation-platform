#!/usr/bin/env bash
set -euo pipefail
mkdir -p artifacts

make test-api
make test-e2e
make test-java
make test-dotnet
make test-ruby
make test-go
make test-cpp
