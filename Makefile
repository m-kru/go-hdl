PROJECT_NAME=hdl

.PHONY: default
default: build

.PHONY: help
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
.PHONY: all
all: lint fmt build

.PHONY: build
build:
	go build -v -o $(PROJECT_NAME) ./cmd/hdl


# Quality targets
.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: lint
lint:
	golangci-lint run


# Test targets
.PHONY: test
test:
	go test ./...

.PHONY: test-gen
test-gen:
	@./scripts/test-gen.sh

.PHONY: test-vet
test-vet:
	@./scripts/test-vet.sh

.PHONY: test-all
test-all: test test-gen test-vet


# Installation targets
.PHONY: install
install:
	cp $(PROJECT_NAME) /usr/bin

.PHONY: uninstall
uninstall:
	rm /usr/bin/$(PROJECT_NAME)
