#!/bin/bash

URL="http://localhost:8080/api/users"
CONTENT_TYPE="Content-Type: application/json"

post() {
  echo ""
  echo "> $1"

  response=$(curl -s \
    -w "\n%{http_code}" \
    -X POST "$URL" \
    -H "$CONTENT_TYPE" \
    -d "$(jq -n --arg email "$2" '{email: $email}')")

  body=$(echo "$response" | head -n -1)
  status=$(echo "$response" | tail -n 1)

  echo "$body" | jq
  echo "Status: $status"
}

post "Success" "anakin.skywalker@empire.com"
