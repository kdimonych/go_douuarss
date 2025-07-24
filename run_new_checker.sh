#!/bin/bash

#load environment variables
if [ -f .env ]; then
  #export $(grep -v '^#' .env | xargs)
  set -a
  source .env
  set +a
fi

DATABASE_URL="postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable"
MIGRATIONS_DIR=${MIGRATIONS_PATH}

echo "DATABASE_URL=" $DATABASE_URL
echo "MIGRATIONS_DIR=" $MIGRATIONS_DIR

docker-compose -f ./docker-compouse.yaml up -d postgres
# Wait for Postgres to be ready
echo "Waiting for Postgres to be ready..."
until docker-compose -f ./docker-compouse.yaml exec postgres pg_isready -U "$POSTGRES_USER" >/dev/null 2>&1; do
  sleep 1
done
echo "Postgres is ready!"

scripts/build_and_run.sh ./services/news_checker/
#stop postgres container
docker-compose -f docker-compouse.yaml down postgres
#docker network prune -f
