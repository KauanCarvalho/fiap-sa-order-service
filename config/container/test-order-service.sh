#!/bin/bash

# Check for required tools
command -v uuidgen >/dev/null 2>&1 || { echo >&2 "uuidgen is required but not installed. Aborting."; exit 1; }
command -v jq >/dev/null 2>&1 || { echo >&2 "jq is required but not installed. Aborting."; exit 1; }

base_url=$1
shift
skus=("$@")

if [ -z "$base_url" ]; then
  echo "Usage: $0 <base_url> <sku1> <sku2> ..."
  exit 1
fi

if [ ${#skus[@]} -eq 0 ]; then
  echo "Error: You must provide at least one SKU."
  echo "Usage: $0 <base_url> <sku1> <sku2> ..."
  exit 1
fi

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m'

clients=()
cpfs=()
orders=()

generate_random_cpf() {
  tr -dc '0-9' </dev/urandom | head -c11
}

generate_random_quantity() {
  echo $(( RANDOM % 5 + 1 ))
}

print_pretty_body() {
  body="$1"
  if [ -z "$body" ]; then
    echo "no body"
  elif echo "$body" | jq . >/dev/null 2>&1; then
    echo "$body" | jq .
  else
    echo "$body"
  fi
}

print_status_code() {
  code="$1"
  if [ "$code" -ge 200 ] && [ "$code" -lt 300 ]; then
    echo -e "${GREEN}$code${NC}"
  elif [ "$code" -ge 400 ] && [ "$code" -lt 500 ]; then
    echo -e "${BLUE}$code${NC}"
  else
    echo -e "${RED}$code${NC}"
  fi
}

# Healthcheck
echo "Checking healthcheck..."
response=$(curl -s -w "\n%{http_code}" "$base_url/healthcheck")
body=$(echo "$response" | sed '$d')
status_code=$(echo "$response" | tail -n1)

print_status_code "$status_code"
print_pretty_body "$body"
echo "---------------------------------------------"

if [ "$status_code" -lt 200 ] || [ "$status_code" -ge 300 ]; then
  echo "Healthcheck failed. Aborting."
  exit 1
fi

# Create 2 clients
for i in {1..2}; do
  uuid=$(uuidgen)
  cpf=$(generate_random_cpf)

  echo "Creating client #$i"
  response=$(curl -s -w "\n%{http_code}" -X POST "$base_url/api/v1/clients" \
    -H "Content-Type: application/json" \
    -d "{
      \"name\": \"name $uuid\",
      \"cpf\": \"$cpf\"
    }")

  body=$(echo "$response" | sed '$d')
  status_code=$(echo "$response" | tail -n1)

  client_id=$(echo "$body" | jq -r '.id // empty')
  clients+=("$client_id")
  cpfs+=("$cpf")

  echo -n "Status: "
  print_status_code "$status_code"
  echo "Response body:"
  print_pretty_body "$body"
  echo "---------------------------------------------"
done

# Create 5 orders for each client
for client_id in "${clients[@]}"; do
  for i in {1..5}; do
    sku_index=$(( RANDOM % ${#skus[@]} ))
    quantity=$(generate_random_quantity)
    sku="${skus[$sku_index]}"

    echo "Creating order #$i for client_id: $client_id (sku: $sku, qty: $quantity)"
    response=$(curl -s -w "\n%{http_code}" -X POST "$base_url/api/v1/checkout" \
      -H "Content-Type: application/json" \
      -d "{
        \"client_id\": $client_id,
        \"items\": [
          {
            \"sku\": \"$sku\",
            \"quantity\": $quantity
          }
        ]
      }")

    body=$(echo "$response" | sed '$d')
    status_code=$(echo "$response" | tail -n1)

    order_id=$(echo "$body" | jq -r '.id // empty')
    if [ -n "$order_id" ]; then
      orders+=("$order_id")
    fi

    echo -n "Status: "
    print_status_code "$status_code"
    echo "Response body:"
    print_pretty_body "$body"
    echo "---------------------------------------------"
  done
done

# PATCH last 2 orders to /ready
for order_id in "${orders[@]: -2}"; do
  echo "PATCH /ready for order_id: $order_id"
  response=$(curl -s -w "\n%{http_code}" -X PATCH "$base_url/api/v1/admin/orders/$order_id/ready")
  body=$(echo "$response" | sed '$d')
  status_code=$(echo "$response" | tail -n1)

  echo -n "Status: "
  print_status_code "$status_code"
  echo "Response body:"
  print_pretty_body "$body"
  echo "---------------------------------------------"
done

# PATCH first 2 orders to /delivered
for order_id in "${orders[@]:0:2}"; do
  echo "PATCH /delivered for order_id: $order_id"
  response=$(curl -s -w "\n%{http_code}" -X PATCH "$base_url/api/v1/admin/orders/$order_id/delivered")
  body=$(echo "$response" | sed '$d')
  status_code=$(echo "$response" | tail -n1)

  echo -n "Status: "
  print_status_code "$status_code"
  echo "Response body:"
  print_pretty_body "$body"
  echo "---------------------------------------------"
done

# GET /admin/orders
echo "GET /api/v1/admin/orders"
response=$(curl -s -w "\n%{http_code}" "$base_url/api/v1/admin/orders")
body=$(echo "$response" | sed '$d')
status_code=$(echo "$response" | tail -n1)

echo -n "Status: "
print_status_code "$status_code"
echo "Response body:"
print_pretty_body "$body"
echo "---------------------------------------------"
