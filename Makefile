.PHONY: compile-contracts
compile-contracts:
	@echo "Compiling contracts..."
	dpm clean --all --multi-package-path ./contracts/
	go run ./contracts/cmd/compile -root ./contracts -artifacts ./contracts/dars

.PHONY: generate-bindings
generate-bindings:
	@echo "Generating contract bindings..."
	sh ./scripts/generate-bindings.sh

.PHONY: contracts
contracts: compile-contracts generate-bindings

.PHONY: gomodtidy
gomodtidy: ## Run go mod tidy on all modules.
	go run github.com/jmank88/gomods@v0.1.7 tidy
