#!/usr/bin/env sh

set -e

direction=${1:-up}

echo "Running database migrations..."

until nc -z -v -w30 "$DB_HOST" 3306; do
  echo "MySQL is unavailable - sleeping" >&2
  sleep 1
done

database_url="mysql://${DB_USER}:${DB_PASSWORD}@tcp(${DB_HOST}:3306)/${DB_NAME}?charset=utf8mb4&parseTime=true"

case "$direction" in
  up)
    migrate -path /app/migrations -database "${database_url}" -verbose up
    ;;
  down)
    migrate -path /app/migrations -database "${database_url}" -verbose down 1
    ;;
  *)
    echo "invalid migration direction" >&2
    exit 1
esac
