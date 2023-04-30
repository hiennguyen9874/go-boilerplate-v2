#!/bin/sh

export POSTGRES_HOST=${POSTGRES_HOST:-db}
export POSTGRES_PORT=${POSTGRES_PORT:-5432}
export POSTGRES_USER=${POSTGRES_USER:-postgres}
export POSTGRES_PASSWORD=${POSTGRES_PASSWORD:-postgres}
export POSTGRES_DBNAME=${POSTGRES_DBNAME:-db}

# Run migrations
echo "migrate..."
atlas migrate apply \
  --dir "file://ent/migrate/migrations" \
  --url "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DBNAME}?sslmode=disable"

# Create initial data in DB
echo "init data..."
go run main.go initdata
