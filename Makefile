PROJECT_NAME=hdl

default: build

help:
	@echo "Build targets:"
	@echo "  all      Run fmt vet build."
	@echo "  build    Build binary."
	@echo "  default  Run build."
	@echo "Quality targets:"
	@echo "  fmt       Format files with go fmt."
	@echo "  vet       Examine go sources with go vet."
	@echo "  errcheck  Examine go sources with errcheck."
	@echo "Test targets:"
	@echo "  test      Run go test."
	@echo "  test-all  Run all tests."
	@echo "Other targets:"
	@echo "  help  Print help message."


# Build targets
all: fmt vet build

build:
	go build -v -o $(PROJECT_NAME) .


# Quality targets
fmt:
	go fmt ./...

vet:
	go vet ./...

errcheck:
	errcheck -verbose ./...


# Test targets
test:
	go test ./...

test-all: test


# Installation targets
install:
	cp $(PROJECT_NAME) /usr/bin

uninstall:
	rm /usr/bin/$(PROJECT_NAME)
