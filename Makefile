.PHONY: compile-contracts
compile-contracts:
	@echo "Compiling contracts..."
	dpm clean --all --multi-package-path ./contracts/
	go run ./contracts/cmd/compile -root ./contracts -artifacts ./contracts/dars

.PHONY: generate-bindings
generate-bindings:
	@echo "Generating contract bindings..."
	sh ./scripts/generate_bindings.sh

.PHONY: contracts
contracts: compile-contracts generate-bindings