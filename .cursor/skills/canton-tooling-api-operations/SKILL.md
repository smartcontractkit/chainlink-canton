---
name: canton-tooling-api-operations
description: Author or refactor Chainlink Canton deployment operations so they conform to the CCIP Deployment Tooling API. Use when adding or changing code under deployment/operations or deployment/utils/operations, especially contract deploy/exercise wrappers.
---

# Canton Tooling API Operations

Write Canton operations as Tooling API leaf nodes, not mini changesets.

## Workflow

1. Read `.cursor/rules/references/tooling-api-authoring.md`.
2. Use the upstream `chainlink-ccip` docs linked there when hierarchy, interfaces, or changeset boundaries are unclear.
3. Keep the operation to one side effect. If the task needs multiple writes, address threading, or datastore work, move that logic into a sequence or changeset.
4. Prefer `deployment/utils/operations/contract.NewDeploy` and `NewExercise` for Canton contract create/exercise flows.
5. Match existing naming, versioning, validation, and output patterns from the local operation packages.
6. Update tests when the operation's validation, defaults, or execution behavior changes.

## Local Examples

- `deployment/operations/ccip/committee_verifier/committee_verifier.go`
- `deployment/operations/mcms/mcms.go`
- `deployment/utils/operations/contract/deploy.go`
- `deployment/utils/operations/contract/exercise.go`
