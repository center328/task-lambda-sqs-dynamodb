.PHONY: clean tidy build deploy deployAll

deployAll: clean tidy build deploy

clean:
	rm -rf ./bin ./vendor go.sum ./.serverless

tidy:
	# To update and prune the dependencies
	go mod tidy

build:
	export GO111MODULE=on
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/gateway src/api/gateway.go
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/handler src/handler/processor.go

deploy:
	sls deploy --verbose
