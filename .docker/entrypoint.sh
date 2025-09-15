#!/bin/sh
set -e

echo "Waiting for database... (Handled by Docker healthcheck)"

echo "Running database migrations..."
/app/migrate up

if [ "$APP_ENV" = "development" ]; then
    echo "Development environment detected. Running database seeder..."
    /app/seeder
else
    echo "Production or other environment detected. Skipping seeder."
fi

echo "Starting application..."
/app/server