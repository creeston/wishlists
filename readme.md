[![npm](https://img.shields.io/badge/demo-online-008000.svg)](https://wishlist-app.proudbush-84ce1133.polandcentral.azurecontainerapps.io/)

# Wishlist

A simple app that was created with purpose of learning Go with HTMX and Alpine.js.

## App overview 

App allows users to create wishlists and share them with others. Other users can see the wishlist and mark items as "reserved". Once item is reserved, it is not possible to reserve it again by other user, or delete / modify it by the owner.

### Technologies used

- Go
- HTMX
- Alpine.js
- SQLite
- Pico CSS

## Development

### Pre-requisites

- Go 1.23
- Air
- Docker (optional)
- dlv (optional for debugging)

### Running with Docker

`docker compose up -d`

### Running in terminal

`air`

## Debug 

dlv debug cmd/main.go 

### Extracting texts for translation

gotext -srclang=en-GB update -out="catalog.go" -lang="en-GB,pl-PL,ru-RU,be-BY" creeston/lists/internal/handlers