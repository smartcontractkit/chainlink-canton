.PHONY: build
build: ## Build the binary.
	go build -o bin/canton-registry-kit .

.PHONY: test
test: ## Run all unit tests.
	go test -v -count=1 -timeout=5m ./...

.PHONY: lint
lint: ## Run golangci-lint.
	golangci-lint run

.PHONY: clean
clean: ## Remove build artifacts.
	rm -rf bin/

.PHONY: help
help: ## Display this help.
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
