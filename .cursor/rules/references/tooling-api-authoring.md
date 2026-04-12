# Tooling API Authoring

Use the upstream CCIP Deployment Tooling API docs on `main` as the canonical reference:

- Overview: <https://github.com/smartcontractkit/chainlink-ccip/blob/main/deployment/docs/index.md>
- Architecture: <https://github.com/smartcontractkit/chainlink-ccip/blob/main/deployment/docs/architecture.md>
- Changesets: <https://github.com/smartcontractkit/chainlink-ccip/blob/main/deployment/docs/changesets.md>
- Types: <https://github.com/smartcontractkit/chainlink-ccip/blob/main/deployment/docs/types.md>
- Interfaces: <https://github.com/smartcontractkit/chainlink-ccip/blob/main/deployment/docs/interfaces.md>
- Implementing Adapters: <https://github.com/smartcontractkit/chainlink-ccip/blob/main/deployment/docs/implementing-adapters.md>
- MCMS and Utilities: <https://github.com/smartcontractkit/chainlink-ccip/blob/main/deployment/docs/mcms-and-utilities.md>

## Tooling API Hierarchy

- Operations are the leaf layer. Keep them to one side effect and serializable inputs/outputs.
- Sequences compose ordered operations for a single-chain workflow and return `OnChainOutput`.
- Changesets are environment-aware wrappers. They verify preconditions, run operations or sequences, and return datastore updates plus reports.

## Operation Rules

- Put contract-facing operations under `deployment/operations/<domain>/`.
- Prefer existing chain-appropriate operation helpers such as `deployment/utils/operations/contract.NewDeploy` and `NewExercise` when the target chain exposes contract deploy/exercise wrappers.
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

## Token Pool Setup Rules

- Keep token adapters thin. They should return existing sequences and implement interface glue, not own deploy/config workflow logic.
- Put token-pool deployment, token-transfer configuration, and pool rate-limit orchestration in sequences.
- Use changesets as CLDF wrappers around the minimal deploy/config sequence or operation.
- Drive deploys through the chain's standard deploy operations. Do not add custom deploy utility layers for token pools.
- Keep first-time token-transfer setup and later rate-limit replacement separate:
  - deploy pool
  - configure pool for transfers
  - later update pool rate limits
- Build token-transfer configs from actual deployed refs and intended lane combinations, not from ad hoc topology guesses or post-connect hacks.
- When the workflow deploys auxiliary contracts such as rate limiters, return them in `OnChainOutput.Addresses` so later steps can resolve them from datastore.
- Prefer datastore and existing helper outputs over re-deriving state from encoded remote-config blobs.
- Prefer shared helpers in `testhelpers` and EDS helpers over local devenv wrappers when they already provide the needed validator clients, disclosures, or accept-context handling.
- Keep token identity explicit:
  - require token refs and labels when needed
  - avoid hidden fallback derivation of instrument identity or token address
- Keep setup comments tied to the actual failure mode being prevented, not broad architecture narration.

## Good Local Examples

- Deploy/exercise helpers:
  - `deployment/utils/operations/contract/deploy.go`
  - `deployment/utils/operations/contract/exercise.go`
- Contract operation package:
  - `deployment/operations/ccip/committee_verifier/committee_verifier.go`
- Sequence orchestration:
  - `deployment/sequences/deploy_chain_contracts.go`
  - `deployment/sequences/token_pools.go`
- Basic changeset wrapper:
  - `deployment/changesets/deploy_chain_contracts.go`
- Changeset mixing a direct operation with a follow-up sequence:
  - `deployment/changesets/deploy_token_pool.go`
