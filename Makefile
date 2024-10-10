# Variables
APP_NAME := goakka
SUBMODULES := core redis

.PHONY: all core nats redis

.PHONY: core
core:
	$(MAKE) -C core
core-test:
	$(MAKE) -C core test

nats:
	$(MAKE) -C nats
nats-test:
	$(MAKE) -C nats test

redis:
	$(MAKE) -C redis
redis-test:
	$(MAKE) -C redis test

# Install/update required Go tools
.PHONY: tools
tools:
	@echo "Installing/updating required tools..."
	go install golang.org/x/tools/cmd/goimports@latest


tidy:
	@for dir in $(SUBMODULES); do \
    	echo "Running go mod tidy in $$dir..."; \
    	cd $$dir && go mod tidy; \
	done

# Run go install in each submodule
install:
	@for dir in $(SUBMODULES); do \
    	echo "Running go install in $$dir..."; \
    	cd $$dir && go install; \
	done

# Run go install in each submodule
lint: 
	@echo "Running linter..."
	@for dir in $(SUBMODULES); do \
    	echo "Running go install in $$dir..."; \
    	cd $$dir && go install; \
	done

test:
	@> go-test-report.txt  # Clear existing report file
	@for dir in $(SUBMODULES); do \
		echo "Running go test in $$dir..."; \
		(cd $$dir && go test ./... -v) >> go-test-report.txt 2>&1; \
	done

fmt:
	@for dir in $(SUBMODULES); do \
		echo "Running go fmt in $$dir..."; \
		cd $$dir && go fmt ./...; \
	done

.PHONY: cover
cover:
	@mkdir -p cover
	@for dir in $(SUBMODULES); do \
		echo "Generating coverage for $$dir..."; \
		(cd $$dir && go test -coverprofile=coverage.out ./... && mv coverage.out ../cover/$$dir-coverage.out && go tool cover -html=../cover/$$dir-coverage.out -o ../cover/$$dir-coverage.html); \
	done

.PHONY: coverage
coverage:
	@mkdir -p coverage
	@for dir in $(SUBMODULES); do \
		echo "Generating coverage for $$dir..."; \
		(cd $$dir && go test -coverprofile=coverage.out ./... && mv coverage.out ../cover/$$dir-coverage.out && go tool cover -html=../cover/$$dir-coverage.out -o ../cover/$$dir-coverage.html); \
	done

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


ci-build: tools tidy install fmt lint test