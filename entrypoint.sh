#!/bin/sh
set -e

export GOOSE_DRIVER="${GOOSE_DRIVER:-postgres}"
export GOOSE_DBSTRING="postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT:-5432}/${POSTGRES_DB}?sslmode=disable"
export GOOSE_MIGRATION_DIR="${GOOSE_MIGRATION_DIR:-./migrations}"

# depends_on не ждёт готовности Postgres, поэтому ретраим миграции.
attempts=30
until goose up; do
  attempts=$((attempts - 1))
  if [ "$attempts" -le 0 ]; then
    echo "migrations failed: database is unavailable" >&2
    exit 1
  fi
  echo "waiting for database... ($attempts attempts left)"
  sleep 2
done

exec ./app
