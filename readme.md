# Wish!

[![azure](https://img.shields.io/badge/demo-online-008000.svg)](https://wishlist-app.proudbush-84ce1133.polandcentral.azurecontainerapps.io/)

## Overview

Wish! — A wishlist creation app built with Go as the backend and HTMX and Alpine.js for the frontend.

 - Create and share wishlists with others.
 - Users can view your wishlist and mark items as "reserved."
 - Once an item is reserved, it cannot be reserved by another user, nor can it be deleted or modified by the owner.

### Technologies Used

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

`dlv debug cmd/main.go `

### Extracting texts for translation

`gotext -srclang=en-GB update -out="catalog.go" -lang="en-GB,pl-PL,ru-RU,be-BY" creeston/lists/internal/handlers`