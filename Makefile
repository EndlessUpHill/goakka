# Variables
APP_NAME := goakka
GO_FILES := $(shell find . -type f -name '*.go')

# Default target
.PHONY: all
all: build

# Build the application
.PHONY: build
build:
	@echo "Building the application..."
	go build -o bin/$(APP_NAME) ./main

# Run the application
.PHONY: run
run:
	@echo "Running the application..."
	./bin/$(APP_NAME)

.PHONY: dev
dev:
	@echo "Running the application in development mode..."
	go run main/main.go

# Test the application
.PHONY: test
test:
	@echo "Running tests..."
	go test ./... -v

# Clean the build and generated files
.PHONY: clean
clean:
	@echo "Cleaning up..."
	rm -rf bin/
	go clean

# Format Go files
.PHONY: fmt
fmt:
	@echo "Formatting Go files..."
	go fmt ./...

# Install dependencies
.PHONY: deps
deps:
	@echo "Tidying dependencies..."
	go mod tidy

# Install/update required Go tools
.PHONY: tools
tools:
	@echo "Installing/updating required tools..."
	go install golang.org/x/tools/cmd/goimports@latest

# Run static analysis
.PHONY: lint
lint:
	@echo "Running linter..."
	golangci-lint run ./...

start-containers:
	@echo "Starting containers..."
	docker-compose up -f docker/docker-compose.yml -d

stop-containers:
	@echo "Stopping containers..."
	docker-compose down -f docker/docker-compose.yml

start-redis:
	@echo "Starting Redis container..."
	docker-compose -f docker/docker-compose.yml up -d redis

stop-redis:
	@echo "Stopping Redis container..."
	docker-compose -f docker/docker-compose.yml down redis

start-nats:
	@echo "Starting NATS container..."
	docker-compose -f docker/docker-compose.yml up -d nats

stop-nats:
	@echo "Stopping NATS container..."
	docker-compose -f docker/docker-compose.yml down nats



# Help information
.PHONY: help
help:
	@echo "Usage:"
	@echo "  make                  - Build the application"
	@echo "  make dec              - Run app in development mode"
	@echo "  make run              - Run the application"
	@echo "  make test             - Run tests"
	@echo "  make clean            - Clean the build files"
	@echo "  make fmt              - Format Go files"
	@echo "  make deps             - Install dependencies"
	@echo "  make tools            - Install/update required tools"
	@echo "  make lint             - Run static analysis"
	@echo "  make start-containers - Start containers"
	@echo "  make stop-containers  - Stop containers"
	@echo "  make start-redis      - Start Redis container"
	@echo "  make stop-redis       - Stop Redis container"
	@echo "  make start-nats       - Start NATS container"
	@echo "  make stop-nats        - Stop NATS container"
	@echo "  make help             - Show this help message"

