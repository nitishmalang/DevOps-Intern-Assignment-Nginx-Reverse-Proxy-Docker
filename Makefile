# YubiKey Setup Tool Makefile
.PHONY: build test clean install help run-check run-setup-dry cross-compile

# Variables
BINARY_NAME=yubikey-setup
GO_MODULE=yubikey-setup
VERSION?=1.0.0
LDFLAGS=-ldflags "-X yubikey-setup/cmd.Version=$(VERSION)"

# Default target
all: test build

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	go build $(LDFLAGS) -o $(BINARY_NAME) .

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -cover ./...

# Run benchmarks
benchmark:
	@echo "Running benchmarks..."
	go test -bench=. ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME)-*
	go clean

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Install the binary
install: build
	@echo "Installing $(BINARY_NAME) to /usr/local/bin/"
	sudo cp $(BINARY_NAME) /usr/local/bin/

# Uninstall the binary
uninstall:
	@echo "Removing $(BINARY_NAME) from /usr/local/bin/"
	sudo rm -f /usr/local/bin/$(BINARY_NAME)

# Cross-compile for different platforms
cross-compile: clean
	@echo "Cross-compiling for multiple platforms..."
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME)-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BINARY_NAME)-linux-arm64 .
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BINARY_NAME)-darwin-arm64 .

# Development commands
run-check: build
	@echo "Running system check..."
	./$(BINARY_NAME) check

run-setup-dry: build
	@echo "Running setup in dry-run mode..."
	./$(BINARY_NAME) setup --dry-run --skip-pin

run-help: build
	@echo "Showing help..."
	./$(BINARY_NAME) --help

# Lint the code
lint:
	@echo "Running linter..."
	golangci-lint run

# Format the code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Vet the code
vet:
	@echo "Vetting code..."
	go vet ./...

# Security check
security:
	@echo "Running security checks..."
	gosec ./...

# Full quality check
quality: fmt vet test lint security
	@echo "All quality checks passed!"

# Development setup
dev-setup:
	@echo "Setting up development environment..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest

# Show version
version:
	@echo "Version: $(VERSION)"

# Help
help:
	@echo "Available targets:"
	@echo "  build           - Build the binary"
	@echo "  test            - Run tests"
	@echo "  test-coverage   - Run tests with coverage"
	@echo "  benchmark       - Run benchmarks"
	@echo "  clean           - Clean build artifacts"
	@echo "  deps            - Install dependencies"
	@echo "  install         - Install binary to /usr/local/bin"
	@echo "  uninstall       - Remove binary from /usr/local/bin"
	@echo "  cross-compile   - Build for multiple platforms"
	@echo "  run-check       - Run system check"
	@echo "  run-setup-dry   - Run setup in dry-run mode"
	@echo "  run-help        - Show application help"
	@echo "  lint            - Run linter"
	@echo "  fmt             - Format code"
	@echo "  vet             - Vet code"
	@echo "  security        - Run security checks"
	@echo "  quality         - Run all quality checks"
	@echo "  dev-setup       - Set up development environment"
	@echo "  version         - Show version"
	@echo "  help            - Show this help"