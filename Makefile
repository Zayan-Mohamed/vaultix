.PHONY: build clean test install help

# Build variables
BINARY_NAME=vaultix
BUILD_DIR=.
GO=go

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	$(GO) build -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Build for multiple platforms
build-all:
	@echo "Building for multiple platforms..."
	GOOS=linux GOARCH=amd64 $(GO) build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 .
	GOOS=darwin GOARCH=amd64 $(GO) build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 $(GO) build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 .
	GOOS=windows GOARCH=amd64 $(GO) build -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe .
	@echo "Multi-platform build complete"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f $(BUILD_DIR)/$(BINARY_NAME)
	rm -f $(BUILD_DIR)/$(BINARY_NAME)-*
	@echo "Clean complete"

# Install to system (use install script for full setup)
install: build
	@echo "Installing $(BINARY_NAME)..."
	@if [ -f "install.sh" ]; then \
		./install.sh; \
	else \
		echo "Note: install.sh not found, using basic install"; \
		sudo mv $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/; \
		echo "Installation complete"; \
	fi

# Run tests
test:
	@echo "Running tests..."
	$(GO) test -v ./...

# Format code
fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...

# Run linter (requires golangci-lint)
lint:
	@echo "Running linter..."
	golangci-lint run

# Show help
help:
	@echo "Available targets:"
	@echo "  build      - Build the binary for current platform"
	@echo "  build-all  - Build for Linux, macOS, and Windows"
	@echo "  clean      - Remove build artifacts"
	@echo "  install    - Install to /usr/local/bin (requires sudo)"
	@echo "  test       - Run tests"
	@echo "  fmt        - Format code"
	@echo "  lint       - Run linter"
	@echo "  help       - Show this help message"
