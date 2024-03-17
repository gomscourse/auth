#!/bin/sh

while ! nc -z postgres 5432; do
  >&2 echo "PostgreSQL недоступен - ожидание..."
  sleep 2
done

./auth_server

