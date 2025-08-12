.PHONY: build test lint run

VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

LDFLAGS := -X github.com/josiarod/multik/internal/cli.version=$(VERSION)

build:
	go build -ldflags "$(LDFLAGS)" -o bin/multik ./cmd/multik

test:
	go test ./... -race -count=1

lint:
	golangci-lint run || true

run:
	go run -ldflags "$(LDFLAGS)" ./cmd/multik
