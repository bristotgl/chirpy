#!/bin/bash

URL="http://localhost:8080/api/validate_chirp"
CONTENT_TYPE="Content-Type: application/json"

post() {
  echo ""
  echo "> $1"

  response=$(curl -s \
    -w "\n%{http_code}" \
    -X POST "$URL" \
    -H "$CONTENT_TYPE" \
    -d "$(jq -n --arg body "$2" '{body: $body}')")

  body=$(echo "$response" | head -n -1)
  status=$(echo "$response" | tail -n 1)

  echo "$body" | jq
  echo "Status: $status"
}

post "Success" "This is my first chirp!"

post "Error empty" "   "

post "Error too long" "We need to talk about panettone. How can this bread be so soft and tasty? Honestly, it's a crime it's only sold on Christmas eve! If I had the money I'd hire a bakery to deliver me fresh panettones every week! Wouldn't that be marvelous?? Who else's with me?"