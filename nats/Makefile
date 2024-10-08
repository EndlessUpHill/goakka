## Make file for nats

# Default target
.PHONY: test
test: ## Run tests
	@echo "core :: Running tests..."
	go test ./... -v

.PHONY: build
build: ## Build the binary
	@echo "core :: Building binary..."
	go build -o ashgo

.PHONY: fmt
fmt: ## Format the code
	@echo "core :: Formatting code..."
	go fmt ./...

# Run static analysis
.PHONY: lint
lint:
	@echo "Running linter..."
	golangci-lint run ./...