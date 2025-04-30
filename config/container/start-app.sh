#!/usr/bin/env sh

set -e

if [[ "$APP_TYPE" =~ ^api$ ]]; then
  echo "Starting the api"
  exec "/app/order-service-api"
fi
