## Debug 

dlv debug cmd/main.go 

## Extract texts for translation

gotext -srclang=en-GB update -out="catalog.go" -lang="en-GB,pl-PL,ru-RU,be-BY" creeston/lists/internal/handlers

## TODO
1. Review js code and template logic
2. Archive (or delete) old wishlists, which were not updated in a while