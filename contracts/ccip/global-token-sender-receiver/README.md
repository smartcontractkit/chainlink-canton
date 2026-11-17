# GlobalTokenSenderReceiver

A single, Chainlink-operated Daml contract (`ccip-global-token-sender-receiver-v1`) used for **all**
DA Registry token transfers in both directions. It removes the need to deploy per-user orchestrator
contracts: every user keeps their own `PerPartyRouter` (one-time `CreateRouter`, no code deployment)
and drives this contract with `caller` = their own party.

The contract can only be used for **pure token transfers** — no arbitrary messages, no payload, no
receiver callbacks. See [Pure token transfers only](#pure-token-transfers-only).

See the [client interfaces](../client/README.md) package for the `ICCIPSender` / `ICCIPReceiver`
interfaces this contract implements.

## Pure token transfers only

On Canton, a message receiver is a **party**, not a contract: a message addressed to a party has to
be executed by that party against a contract that reacts to the message content. Since this contract
is shared between all users, it is deliberately a **blank receiver** — it contains no application
logic and only ever moves tokens. Messages that carry anything beyond a token transfer have no logic
to react to them and are therefore rejected.

A **pure token transfer** is a message that consists solely of a token transfer, mirroring the
OnRamp / OffRamp definition of a token-only transfer:

- **Outgoing** (`GetRequiredCCVsForSend`, `GetFee`, `Send`): `message.tokenTransfer` present,
  `message.payload` empty, and `extraArgs.gasLimit == 0` — the field that becomes
  `ccipReceiveGasLimit` on the wire.
- **Incoming** (`GetRequiredCCVsForExecute`, `Execute`): decoded `message.tokenTransfer` present,
  `message.messageData` empty, and `message.ccipReceiveGasLimit == 0`.

Implications:

- **No arbitrary messages and no programmable receivers.** Message-only messages (no token
  transfer), payload-bearing messages, and messages with a non-zero callback gas limit are rejected
  by the choices above. Users that need custom receiver logic must deploy their own orchestrator 
  (e.g. the per-user `CCIPSender` / `CCIPReceiver`).
- **The baked-in receiver policy is inert.** The token pool is the only counterparty accepting risk
  in a pure token transfer and decides CCV and finality requirements (see
  [Design](#design) for the fields involved).

## Design

- **Signatory** `ccipOwner` (Chainlink). All choices are `nonconsuming` and controlled by a per-call
  `caller : Party`.
- **Baked-in receiver policy** (`requiredCCVs`, `optionalCCVs`, `optionalThreshold`,
  `receiverFinalityConfig`) is set at deployment. It is inert for token-only transfers — the token
  pool is the only counterparty accepting risk in a token transfer and decides CCV and finality
  requirements — but the fields remain because the router's `PrepareExecute` expects them.
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
