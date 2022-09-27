PROJECT_NAME=hdl

default: build

help:
	@echo "Build targets:"
	@echo "  all      Run lint fmt build."
	@echo "  build    Build binary."
	@echo "  default  Run build."
	@echo "Quality targets:"
	@echo "  fmt   Format files with go fmt."
	@echo "  lint  Lint files with golangci-lint."
	@echo "Test targets:"
	@echo "  test      Run go test."
	@echo "  test-gen  Run gen command tests."
	@echo "  test-vet  Run vet command tests."
	@echo "  test-all  Run all tests."
	@echo "Other targets:"
	@echo "  help  Print help message."


# Build targets
all: lint fmt build

build:
	go build -v -o $(PROJECT_NAME) ./cmd/hdl


# Quality targets
fmt:
	go fmt ./...

lint:
	golangci-lint run


# Test targets
test:
	go test ./...

test-gen:
	@./scripts/test-gen.sh
	
test-vet:
	@./scripts/test-vet.sh

test-all: test test-gen test-vet


# Installation targets
install:
	cp $(PROJECT_NAME) /usr/bin

uninstall:
	rm /usr/bin/$(PROJECT_NAME)

.PHONY: all build test test-vet test-all test-gen install uninstall
