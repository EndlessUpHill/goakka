## Make file for redis

# Default target
.PHONY: test
test: ## Run tests
	@echo "redis :: Running tests..."
	go test ./... -v

.PHONY: build
build: ## Build the binary
	@echo "redis :: Building binary..."
	go build -o ashgo

.PHONY: fmt
fmt: ## Format the code
	@echo "redis :: Formatting code..."
	go fmt ./...

# Run static analysis
.PHONY: lint
lint:
	@echo "Running linter..."
	golangci-lint run ./...