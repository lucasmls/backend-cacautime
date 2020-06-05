#!/bin/sh

# Exit on any error
# set -e

until psql -v ON_ERROR_STOP=1 -Atx "$POSTGRES_CONNECTION_STRING"; do
  >&2 echo "Postgres is unavailable - sleeping"
  sleep 2
done

printf "Starting application...\n"
go run cmd/server/main.go