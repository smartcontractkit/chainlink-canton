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

.PHONY: update-contract-version
update-contract-version: ## Update contract version and rebuild. Usage: make update-contract-version OLD=0.0.1 NEW=1.0.0
	@./contracts/scripts/update-version.sh $(OLD) $(NEW)
	@$(MAKE) contracts

.PHONY: go-generate
go-generate:
	go generate ./...

.PHONY: gomodtidy
gomodtidy: ## Run go mod tidy on all modules.
	go run github.com/jmank88/gomods@v0.1.7 tidy

.PHONY: test-daml-contracts
test-daml-contracts:
	cd contracts && dpm build --all
	go run ./contracts/cmd/test --root ./contracts

# GolangCI-Lint targets

.PHONY: golangci-lint-main golangci-lint-integration-tests
golangci-lint-main: ## Run golangci-lint on the main module.
	golangci-lint run
golangci-lint-integration-tests: ## Run golangci-lint on the integration-tests module.
	cd integration-tests && golangci-lint run

.PHONY: golangci-lint-all
golangci-lint-all: golangci-lint-main golangci-lint-integration-tests ## Run golangci-lint on all modules.

.PHONY: golangci-lint-fix-main golangci-lint-fix-integration-tests
golangci-lint-fix-main: ## Run golangci-lint --fix on the main module.
	golangci-lint run --fix
golangci-lint-fix-integration-tests: ## Run golangci-lint --fix on the integration-tests module.
	cd integration-tests && golangci-lint run --fix

.PHONY: golangci-lint-fix-all
golangci-lint-fix-all: golangci-lint-fix-main golangci-lint-fix-integration-tests ## Run golangci-lint --fix on all modules.

## Run all fix targets.
## Compiles contracts, generates bindings, runs all go generates, runs go mod tidy, and runs golangci-lint --fix on all modules.
.PHONY: fix-all
fix-all: contracts go-generate gomodtidy golangci-lint-fix-all

.PHONY: build-committeeverifier
build-committeeverifier:
	docker build -t committeeverifier-canton:latest -f ccip/committee_verifier.Dockerfile .

.PHONY: build-eds
build-eds:
	docker build -t canton-eds:latest -f eds/eds.Dockerfile .

## Assuming chainlink-ccv is checked out in ../chainlink-ccv.
.PHONY: build-ccv-images
build-ccv-images:
	cd ../chainlink-ccv/build/devenv && just build-docker

.PHONY: start-devenv
start-devenv: build-ccv-images build-committeeverifier
	cd ccip/devenv && go run cmd/ccv/main.go down && go run cmd/ccv/main.go up env-canton-evm.toml

.PHONY: run-e2e-tests
run-e2e-tests:
	cd ccip/devenv/tests/e2e && go test -timeout 5m -v -count 1 -run TestEVM2Canton_Basic && go test -timeout 5m -v -count 1 -run TestCantonSourceReader

.PHONY: build-run-e2e-tests
build-run-e2e-tests: start-devenv run-e2e-tests
