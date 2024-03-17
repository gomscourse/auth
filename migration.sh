#!/bin/sh
source .env

while ! nc -z postgres 5432; do
  >&2 echo "PostgreSQL недоступен для миграции - ожидание..."
  sleep 2
done

export MIGRATION_DSN="host=postgres port=5432 dbname=$PG_DATABASE_NAME user=$PG_USER password=$PG_PASSWORD sslmode=disable"

goose -dir "${MIGRATION_DIR}" postgres "${MIGRATION_DSN}" up -v