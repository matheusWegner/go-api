build:
    @go build -o bin/go-api

run: build 
    @./bin/go-api

test:
    @go test -v ./...