#!/bin/bash

# Run application migrations
if [[ $DATABASE_SCHEMA != "" && $DATABASE_MIGRATIONS_PATH != "" ]]; then
  echo "Searching for migrations..."
  uri="$DATABASE_SCHEMA://$POSTGRES_USER:$POSTGRES_PASSWORD@$DATABASE_HOST:$DATABASE_PORT/$POSTGRES_NAME?sslmode=$DATABASE_SSL_MODE"
  migrate -path $DATABASE_MIGRATIONS_PATH -database $uri up
fi

# Run application binary
./main
