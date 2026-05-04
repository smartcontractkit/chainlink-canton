.PHONY: build
build: ## Build the binary.
	go build -o bin/canton-party-ceremony .

.PHONY: test
test: ## Run all unit tests.
	go test -v -count=1 -timeout=5m ./...

.PHONY: test-integration
test-integration: test-integration-kms test-integration-non-kms ## Run all integration tests.

.PHONY: test-integration-kms
test-integration-kms: ## Run KMS-backed integration tests.
	cd integration-tests && go test -v -count=1 -timeout 20m ./ceremony -run '^TestCeremonies$$'

.PHONY: test-integration-non-kms
test-integration-non-kms: ## Run non-KMS integration tests.
	cd integration-tests && go test -v -count=1 -timeout 20m ./ceremony -run '^TestCeremoniesWithoutKMS$$'

.PHONY: lint
lint: ## Run golangci-lint.
	golangci-lint run

.PHONY: lint-all
lint-all: ## Run golangci-lint on all modules.
	golangci-lint run
	cd integration-tests && golangci-lint run

.PHONY: lint-fix-all
lint-fix-all: ## Run golangci-lint with auto-fix on all modules.
	golangci-lint run --fix
	cd integration-tests && golangci-lint run --fix

.PHONY: clean
clean: ## Remove build artifacts.
	rm -rf bin/

.PHONY: help
help: ## Display this help.
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
