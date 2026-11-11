# GlobalTokenSenderReceiver

A single, Chainlink-operated Daml contract (`ccip-global-token-sender-receiver-v2`) used for **all**
DA Registry token transfers in both directions. It removes the need to deploy per-user orchestrator
contracts: every user keeps their own `PerPartyRouter` (one-time `CreateRouter`, no code deployment)
and drives this contract with `caller` = their own party.

See the [client interfaces](../client/README.md) package for the `ICCIPSender` / `ICCIPReceiver`
interfaces this contract implements.

## Design

- **Signatory** `ccipOwner` (Chainlink). All choices are `nonconsuming` and controlled by a per-call
  `caller : Party`.
- **Baked-in receiver policy** (`requiredCCVs`, `optionalCCVs`, `optionalThreshold`,
  `receiverFinalityConfig`) is set at deployment. It is inert for token-only transfers — the token
  pool is the only counterparty accepting risk in a token transfer and decides CCV and finality
  requirements — but the fields remain because the router's `PrepareExecute` expects them.
- **Pure token transfers only.** `Send` and `GetFee` abort on non-empty `message.payload`; `Execute`
  aborts on message-only and payload-bearing messages (mirroring the OffRamp's definition of a pure
  token transfer).
- **No receipts.** `Execute` returns a `ReceiverExecuteResult` (message id, message, execution
  state, the pool's `ReleaseOrMintResult`, choice context); no contract is created.
- **Self-contained orchestration.** The choice bodies mirror the (to-be-deprecated) `CCIPSender` /
  `CCIPReceiver` orchestrators with `owner := caller` / `receiverParty` bindings. This package
  intentionally does **not** depend on the sender/receiver packages.
- **Double execution is replay-safe** via the router's per-party `executedMessages` tracking; a
  shared contract cannot redirect funds because the OffRamp asserts the message's encoded receiver
  against `receiverParty` / `tokenReceiverParty`.

## Deployment

Deployed through the MCMS-governed `CCIPFactory` via the `DeployGlobalTokenSenderReceiver` choice
(tracked as a singleton in the factory's `deployedContracts`). The direct (non-proposal) deploy path
in `deployment/sequences/deploy_chain_contracts.go` deploys it as well.

## Explicit disclosure (EDS)

The contract's disclosure blob is **static** (created once, never consumed):

- `GET /ccip/v1/global/tokenSenderReceiver` — dedicated lookup for the global contract's disclosure

## Security & trust

Unlike other ccip-owned contracts, this contract orchestrates on behalf of the user and exercises
choices on third-party contracts (most notably the token pool via `ITokenPool`). Since the
signatory is `ccipOwner`, our validators must have all downstream DARs vetted — including the DAR of
every token pool (and its transitive dependencies) used through this contract. Registry-issued
tokens are vetted by us; third-party token packages are out of scope.
