# Tooling API Authoring

Use the upstream CCIP Deployment Tooling API docs on `main` as the canonical reference:

- Overview: <https://github.com/smartcontractkit/chainlink-ccip/blob/main/deployment/docs/index.md>
- Architecture: <https://github.com/smartcontractkit/chainlink-ccip/blob/main/deployment/docs/architecture.md>
- Changesets: <https://github.com/smartcontractkit/chainlink-ccip/blob/main/deployment/docs/changesets.md>
- Types: <https://github.com/smartcontractkit/chainlink-ccip/blob/main/deployment/docs/types.md>
- Interfaces: <https://github.com/smartcontractkit/chainlink-ccip/blob/main/deployment/docs/interfaces.md>
- Implementing Adapters: <https://github.com/smartcontractkit/chainlink-ccip/blob/main/deployment/docs/implementing-adapters.md>
- MCMS and Utilities: <https://github.com/smartcontractkit/chainlink-ccip/blob/main/deployment/docs/mcms-and-utilities.md>

## Canton Hierarchy

- Operations are the leaf layer. Keep them to one side effect and serializable inputs/outputs.
- Sequences compose ordered operations for a single-chain workflow and return `OnChainOutput`.
- Changesets are environment-aware wrappers. They verify preconditions, run operations or sequences, and return datastore updates plus reports.

## Operation Rules

- Put contract-facing operations under `deployment/operations/<domain>/`.
- Prefer `deployment/utils/operations/contract.NewDeploy` and `NewExercise` for Canton contract deploy/exercise flows.
- Match the existing shape:
  - `ContractType`
  - `Version`
  - stable slash-separated operation `Name`
  - concise `Description`
  - `Validate` for caller-controlled inputs
- Use `Modifier` only for chain-derived defaults or caller injection that must happen inside the operation.
- Keep contract lookup, submission, and single-write execution inside the operation. Do not orchestrate multiple ledger writes there.
- Do not make the operation responsible for cross-chain dispatch, datastore population, or MCMS proposal construction.
- Return framework-native outputs such as `datastore.AddressRef` or a small typed struct, not ad hoc maps.
- Add or update tests when validation, defaults, or execution semantics change.

## Sequence Rules

- Put multi-step workflows under `deployment/sequences/`.
- Keep sequence input serializable and scoped to one chain.
- Use sequences to compose multiple operations, thread intermediate addresses forward, and collect `AddressRef` outputs.
- Wrap failures with enough context to identify the failing operation.

## Changeset Rules

- Put Tooling API entry points under `deployment/changesets/`.
- `VerifyPreconditions` should validate chain presence, participant indices, and obvious environment assumptions.
- `Apply` should create or update the datastore, run the minimal operation or sequence needed, and return `cldf.ChangesetOutput`.
- Keep datastore reads/writes, environment access, and future MCMS wiring in the changeset layer rather than in the leaf operations.

## Good Local Examples

- Deploy/exercise helpers:
  - `deployment/utils/operations/contract/deploy.go`
  - `deployment/utils/operations/contract/exercise.go`
- Contract operation package:
  - `deployment/operations/ccip/committee_verifier/committee_verifier.go`
- Sequence orchestration:
  - `deployment/sequences/deploy_chain_contracts.go`
- Basic changeset wrapper:
  - `deployment/changesets/deploy_chain_contracts.go`
- Changeset mixing a direct operation with a follow-up sequence:
  - `deployment/changesets/deploy_token_pool.go`
