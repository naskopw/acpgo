.PHONY: all test lint clean

all: lint test

test:
	go test -v -race -count=1 ./...

lint:
	golangci-lint run

clean:
	go clean -testcache
