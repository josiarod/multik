.PHONY: build test lint run

build:
   go build -o bin/multik ./cmd/multik

test:
   go test ./... -race -count=1

lint:
   golangci-lint run || true

run:
   go run ./cmd/multik
