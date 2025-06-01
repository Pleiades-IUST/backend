#!/bin/bash
source .env

COMMON_DB_ARGS="-h $DB_HOST -U $DB_USER -p $DB_PORT"

# Create the database unconditionally
echo "Creating database $DB_NAME..."
CREATEDB_CMD="CREATE DATABASE \"$DB_NAME\";"
psql "postgresql://$DB_USER:$DB_PASS@$DB_HOST:$DB_PORT/postgres" -c "$CREATEDB_CMD"

migrate -database "postgresql://$DB_USER:$DB_PASS@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable" -path database/migrations up