#!/bin/sh

set -e

echo "run db migration"
/app/migrate -path /app/db/migration -database "$DB_SOURCE" -verbose up

echo "running go tests"
go test -v -cover ./... 

echo "start the app"
exec "$@"
