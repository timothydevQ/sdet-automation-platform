#!/usr/bin/env bash
set -euo pipefail
HOST=${PG_HOST:-localhost}
PORT=${PG_PORT:-5432}

echo "waiting for postgres at $HOST:$PORT..."
for i in {1..30}; do
  if docker exec -i $(docker compose ps -q postgres) pg_isready -U sdet >/dev/null 2>&1; then
    break
  fi
  sleep 1
done

echo "applying seed.sql..."
docker exec -i $(docker compose ps -q postgres) \
  psql -U sdet -d sdet < test-data/seed.sql

echo "seeded."
