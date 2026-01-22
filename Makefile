SHELL := /bin/bash

# Localnet docker compose settings
COMPOSE_DIR  ?= compose/localnet
COMPOSE_FILE ?= $(COMPOSE_DIR)/compose.yaml
CONTRACTS_DIR ?= contracts

# Docker Compose v2 (preferred). Override if needed, e.g. `make DC="docker-compose" up`
DC ?= docker compose

# Common compose invocation (project-directory makes relative paths resolve correctly)
DC_CMD := $(DC) -f $(COMPOSE_FILE) --project-directory $(COMPOSE_DIR)

.PHONY: help up down stop start restart ps logs pull build config console clean contracts-build contracts-clean compile-contracts generate-bindings contracts

help:
	@echo "Localnet docker compose targets:"
	@echo "  make up        - Start localnet (detached)"
	@echo "  make down      - Stop and remove containers/networks"
	@echo "  make stop      - Stop containers"
	@echo "  make start     - Start existing containers"
	@echo "  make restart   - Restart containers"
	@echo "  make ps        - Show container status"
	@echo "  make logs      - Follow logs"
	@echo "  make pull      - Pull images"
	@echo "  make build     - Build images (console, etc.)"
	@echo "  make config    - Render compose config"
	@echo "  make console   - Run console (profile: console)"
	@echo "  make clean     - Down + remove volumes (DANGEROUS)"

up:
	$(DC_CMD) up -d

down:
	$(DC_CMD) down

stop:
	$(DC_CMD) stop

start:
	$(DC_CMD) start

restart:
	$(DC_CMD) restart

ps:
	$(DC_CMD) ps

logs:
	$(DC_CMD) logs -f --tail=200

pull:
	$(DC_CMD) pull

config:
	$(DC_CMD) config

# compose/localnet/compose.yaml declares `console` under `profiles: [console]`
console:
	$(DC_CMD) --profile console run --rm console


build:
	make contracts
	$(DC_CMD) build

clean:
	$(DC_CMD) down -v --remove-orphans
	make contracts-clean

compile-contracts:
	@echo "Compiling contracts..."
	dpm clean --all --multi-package-path ./contracts/
	go run ./contracts/cmd/compile -root ./contracts -artifacts ./contracts/dars

generate-bindings:
	@echo "Generating contract bindings..."
	sh ./scripts/generate-bindings.sh
  
contracts: compile-contracts generate-bindings
