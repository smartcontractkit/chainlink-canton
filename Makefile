.PHONY: compile-contracts
compile-contracts:
	@echo "Compiling contracts..."
	dpm clean --all --multi-package-path ./contracts/
	go run ./contracts/cmd/compile -root ./contracts -artifacts ./contracts/dars

.PHONY: generate-bindings
generate-bindings:
	go run ./contracts/cmd/bindings

.PHONY: contracts
contracts: compile-contracts generate-bindings

.PHONY: gomodtidy
gomodtidy: ## Run go mod tidy on all modules.
	go run github.com/jmank88/gomods@v0.1.7 tidy

.PHONY: test-daml-contracts
test-daml-contracts:
	cd contracts && dpm build --all
	go run ./contracts/cmd/test --root ./contracts --color

# GolangCI-Lint targets

.PHONY: golangci-lint-main golangci-lint-eds golangci-lint-integration-tests
golangci-lint-main: ## Run golangci-lint on the main module.
	golangci-lint run
golangci-lint-eds: ## Run golangci-lint on the eds module.
	cd eds && golangci-lint run
golangci-lint-integration-tests: ## Run golangci-lint on the integration-tests module.
	cd integration-tests && golangci-lint run

.PHONY: golangci-lint-all
golangci-lint-all: golangci-lint-main golangci-lint-eds golangci-lint-integration-tests ## Run golangci-lint on all modules.

.PHONY: golangci-lint-fix-main golangci-lint-fix-eds golangci-lint-fix-integration-tests
golangci-lint-fix-main: ## Run golangci-lint --fix on the main module.
	golangci-lint run --fix
golangci-lint-fix-eds: ## Run golangci-lint --fix on the eds module.
	cd eds && golangci-lint run --fix
golangci-lint-fix-integration-tests: ## Run golangci-lint --fix on the integration-tests module.
	cd integration-tests && golangci-lint run --fix

.PHONY: golangci-lint-fix-all
golangci-lint-fix-all: golangci-lint-fix-main golangci-lint-fix-eds golangci-lint-fix-integration-tests ## Run golangci-lint --fix on all modules.

## Run all fix targets.
## Compiles contracts, generates bindings, runs go mod tidy, and runs golangci-lint --fix on all modules.
.PHONY: fix-all
fix-all: contracts gomodtidy golangci-lint-fix-all

.PHONY: build-committeeverifier
build-committeeverifier:
	docker build -t committeeverifier-canton:latest -f ccip/committee_verifier.Dockerfile .
