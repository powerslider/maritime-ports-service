#!/bin/bash

echo -e ">>> Creating a new non existing port with ID NEWPORT...\n"

curl -X 'POST' \
  'http://0.0.0.0:8080/api/v1/ports' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
    "id": "NEWPORT",
    "name": "Newest Port",
    "coordinates": [
      123.321,
      43.34
    ],
    "city": "Some City",
    "country": "Some Country",
    "alias": [],
    "regions": [],
    "timezone": "My/Timezone",
    "unlocs": [
      "NEWPORT"
    ]
}'

echo -e "\n\n>>> Querying that newly created port via its ID NEWPORT...\n"

curl -X 'GET' \
  'http://0.0.0.0:8080/api/v1/ports/NEWPORT' \
  -H 'accept: application/json'

echo -e "\n\n>>> Updating existing port with ID AEKLF by changing its city and country properties...\n"

curl -X 'POST' \
  'http://0.0.0.0:8080/api/v1/ports' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
    "id": "AEKLF",
    "city": "London",
    "country": "United Kingdom"
}'

echo -e "\n\n>>> Querying that newly modified port via its ID AEKLF to verify the changed properties...\n"

curl -X 'GET' \
  'http://0.0.0.0:8080/api/v1/ports/AEKLF' \
  -H 'accept: application/json'
