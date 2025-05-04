#!/usr/bin/env sh

set -e

if [[ "$APP_TYPE" =~ ^api$ ]]; then
  echo "Starting the api"
  exec "/app/order-service-api"
elif [[ "$APP_TYPE" =~ ^worker$ ]]; then
  echo "Starting the worker"
  exec "/app/order-service-worker"
fi
