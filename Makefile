
.EXPORT_ALL_VARIABLES:
STAGE_NAME ?= dev

.PHONY: build clean deploy gomodgen

build: gomodgen
	export GO111MODULE=on
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/gateway src/api/gateway.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/handler src/handler/processor.go

clean:
	rm -rf ./bin ./vendor go.sum

deploy: clean build
	sls deploy --verbose
