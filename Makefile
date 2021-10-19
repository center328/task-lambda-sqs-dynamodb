
.EXPORT_ALL_VARIABLES:
STAGE_NAME ?= dev

.PHONY: build clean deploy gomodgen tidy

build: gomodgen
	export GO111MODULE=on
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/gateway src/api/gateway.go
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/handler src/handler/processor.go

clean:
	rm -rf ./bin ./vendor go.sum

deploy: clean build
	sls deploy --verbose

tidy:
	# To update and prune the dependencies
	go mod tidy
