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

# Docker tasks
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

# Release module (core, nats, redis) based on a provided version tag
.PHONY: release-core release-nats release-redis
release-core:
	@echo "Releasing core module with tag $(VERSION)..."
	cd core && git tag core/$(VERSION) && git push origin core/$(VERSION)

release-nats:
	@echo "Releasing nats module with tag $(VERSION)..."
	cd nats && git tag nats/$(VERSION) && git push origin nats/$(VERSION)

release-redis:
	@echo "Releasing redis module with tag $(VERSION)..."
	cd redis && git tag redis/$(VERSION) && git push origin redis/$(VERSION)

# Fetch the latest tag for each module
.PHONY: latest-core-tag latest-nats-tag latest-redis-tag
latest-core-tag:
	@echo "Fetching latest core tag..."
	git fetch --tags
	git describe --tags $(shell git rev-list --tags --max-count=1 -- core)

latest-nats-tag:
	@echo "Fetching latest nats tag..."
	git fetch --tags
	git describe --tags $(shell git rev-list --tags --max-count=1 -- nats)

latest-redis-tag:
	@echo "Fetching latest redis tag..."
	git fetch --tags
	git describe --tags $(shell git rev-list --tags --max-count=1 -- redis)

# Help information
.PHONY: help
help:
	@echo "Usage:"
	@echo "                 - Build the application"
	@echo "dev              - Run app in development mode"
	@echo "run              - Run the application"
	@echo "test             - Run tests"
	@echo "clean            - Clean the build files"
	@echo "fmt              - Format Go files"
	@echo "deps             - Install dependencies"
	@echo "tools            - Install/update required tools"
	@echo "lint             - Run static analysis"
	@echo "start-containers - Start containers"
	@echo "stop-containers  - Stop containers"
	@echo "start-redis      - Start Redis container"
	@echo "stop-redis       - Stop Redis container"
	@echo "start-nats       - Start NATS container"
	@echo "stop-nats        - Stop NATS container"
	@echo "release-core VERSION=vX.X.X  - Release core module with the specified tag"
	@echo "release-nats VERSION=vX.X.X  - Release nats module with the specified tag"
	@echo "release-redis VERSION=vX.X.X - Release redis module with the specified tag"
	@echo "latest-core-tag               - Fetch latest core module tag"
	@echo "latest-nats-tag               - Fetch latest nats module tag"
	@echo "latest-redis-tag              - Fetch latest redis module tag"
