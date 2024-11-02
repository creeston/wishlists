## Debug 

dlv debug cmd/main.go 

## Extract texts for translation

gotext -srclang=en-GB update -out="catalog.go" -lang="en-GB,pl-PL,ru-RU,be-BY" creeston/lists/internal/handlers

## To do

1. Create Dockerfile

### When becomes a problem

1. Archive (or delete) old wishlists, which were not updated in a while