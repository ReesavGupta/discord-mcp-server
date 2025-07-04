.PHONY: build test clean run install deps lint fmt

# Variables
BINARY_NAME=discord-mcp-server
BUILD_DIR=bin
MAIN_PATH=cmd/server/main.go

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

# Run tests
test:
	@echo "Running tests..."
	go test -v ./tests/...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./tests/...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

# Run the application
run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BUILD_DIR)/$(BINARY_NAME)

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod tidy
	go mod download

# Lint code
lint:
	@echo "Running linter..."
	golangci-lint run

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Install the binary
install: build
	@echo "Installing $(BINARY_NAME)..."
	sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/

# Development setup
dev-setup:
	@echo "Setting up development environment..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(MAKE) deps

# Docker build
docker-build:
	@echo "Building Docker image..."
	docker build -t $(BINARY_NAME) .

# Create release
release: clean
	@echo "Creating release builds..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)
