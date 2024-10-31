## Debug 

dlv debug cmd/main.go 

## Extract texts for translation

gotext -srclang=en-GB update -out="catalog.go" -lang="en-GB,pl-PL,ru-RU,be-BY" creeston/lists/internal/handlers

## TODO
2. Review labels and translations
3. Review js code and template logic
4. Limit number of lists that user (from IP) can create daily
5. Archive (or delete) old lists, which were not updated in a while