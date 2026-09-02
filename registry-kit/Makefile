.PHONY: test
test: ## Run all unit tests.
	go test -v -count=1 -timeout=5m ./...

.PHONY: lint
lint: ## Run golangci-lint.
	golangci-lint run

.PHONY: clean
clean: ## Remove build artifacts.
	rm -rf bin/
