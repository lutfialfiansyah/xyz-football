#!/bin/sh
set -e

echo "Running database migrations..."
./migrate up

echo "Starting API server..."
exec ./api
