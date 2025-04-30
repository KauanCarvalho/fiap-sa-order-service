#!/usr/bin/env sh

set -e

DATABASE_URL="mysql://${DB_USER}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}?charset=utf8mb4&parseTime=true"

echo "DATABASE_URL: ${DATABASE_URL}"

until nc -z -v -w30 $DB_HOST 3306; do
  >&2 echo "MySQL is unavailable - sleeping"
  sleep 1
done

# Install migrate.
go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Setup database.
migrate -verbose -database ${DATABASE_URL} -path ./migrations up

# Run tests.
go test ./... -coverprofile=coverage.out -cover -p 1

# Run coverrage check.
go tool go-test-coverage --config=./.testcoverage.yml
