---
name: tooling-api-changesets
description: Author or refactor deployment sequences and changesets so they conform to the CCIP Deployment Tooling API. Use when adding or changing code under deployment/sequences or deployment/changesets, or when an operation needs to be lifted into a workflow entry point.
---

# Tooling API Changesets

Use sequences and changesets to preserve the Tooling API hierarchy around deployment operations.

## Workflow

1. Read `.cursor/rules/references/tooling-api-authoring.md`.
2. Open the upstream `chainlink-ccip` docs linked there for architecture, changesets, types, and adapter expectations.
3. If the task touches token pools or token-transfer setup, also follow the token-pool guidance in `.cursor/rules/references/tooling-api-authoring.md`.
4. Keep sequences single-chain, serializable, and focused on composing operations plus collecting `OnChainOutput`.
5. Keep environment access, datastore reads/writes, and entry-point validation in the changeset layer.
6. Do not bury orchestration or datastore behavior inside leaf operations when the behavior belongs in a sequence or changeset.
7. Update tests and any nearby docs when you introduce a new reusable deployment workflow.

## Local Examples

- `deployment/sequences/deploy_chain_contracts.go`
- `deployment/sequences/register_token_pool.go`
- `deployment/sequences/token_pools.go`
- `deployment/changesets/deploy_chain_contracts.go`
- `deployment/changesets/deploy_token_pool.go`
