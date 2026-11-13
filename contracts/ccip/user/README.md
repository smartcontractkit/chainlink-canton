# ccip-client

Minimal Daml package (`ccip-user-v1`) containing the **stable client interfaces** for driving CCIP
sender/receiver orchestrator contracts, plus the canonical input/result types those interfaces
reference.

## Purpose

The SDK (and any other client) needs to exercise sender/receiver functionality on *conforming
contracts* without depending on a specific implementation. This package defines that contract
surface:

- **`CCIP.UserV1.Sender`** — the `ICCIPSender` interface (`ICCIPSender_GetRequiredCCVs`,
  `ICCIPSender_GetFee`, `ICCIPSender_Send`) with the `SenderFeeTokenInput`, `SenderTokenTransferInput`,
  `SenderExecutorInput`, `SenderCCVSendInput` input types and the `SenderGetFeeResult` return type.
- **`CCIP.UserV1.Receiver`** — the `ICCIPReceiver` interface (`ICCIPReceiver_GetRequiredCCVs`,
  `ICCIPReceiver_Execute`) with the `ReceiverCCVInput`, `ReceiverTokenTransferInput` input types and
  the `ReceiverExecuteResult` return type.

Implementations:

- `CCIP.GlobalTokenSenderReceiver.GlobalTokenSenderReceiver` (see
  [global-token-sender-receiver](../global-token-sender-receiver/README.md)) — the single
  Chainlink-operated orchestrator.
- `CCIP.CCIPSender` / `CCIP.CCIPReceiver` (per-user orchestrators) — implemented via interface
  instances that convert between their local types and the canonical types defined here.
  These packages are scheduled for deprecation.

## Design notes

- **Why not `client-api`?** Daml packages containing interfaces are in practice un-upgradeable.
  `client-api` already has a released lineage; adding interfaces there would freeze it. A dedicated
  minimal package contains the blast radius.
- **ChoiceContexts everywhere.** Every interface choice takes a
  `Splice.Api.Token.MetadataV1.ChoiceContext`, and `ReceiverExecuteResult` carries one. Interfaces
  themselves are un-upgradeable; the context allows arbitrary data to be passed around without
  touching the interface.
- **`GetFee` returns `GetFeeResult`**, combining both the protocol-level CCIP fee, and
  the quote returned by the token pool, if a token transfer is present in the message.
- **`Execute` creates no contracts.** The choice returns everything the caller needs (message,
  execution state, the token pool's `ReleaseOrMintResult`) as `ReceiverExecuteResult`.

## Dependencies

Mirrors `client-api` (splice token APIs, `chainlink-api`, `ccip-codec-v2`, `ccip-api-v2`) plus
`ccip-extension-api-v2` (token pool / executor / verifier interfaces) and `ccip-client-api-v2`
(`PerPartyRouter`, message and result types). No dependency on any orchestrator implementation.

## Testing

Implementability is proven by the real implementations and the unit tests in
[`ccip/test`](../test) (`GlobalTokenSenderReceiverTest.daml`).
