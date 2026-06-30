BINARY=acp

.PHONY: all build test lint clean

all: lint test build

build:
	go build -o $(BINARY) .

test:
	go test -v -race -count=1 ./...

lint:
	golangci-lint run

clean:
	rm -f $(BINARY)
