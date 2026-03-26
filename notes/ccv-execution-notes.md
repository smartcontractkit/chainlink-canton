# CCV Execution Notes

Date: 2026-03-26

## Canton inbound execution

- `GlobalConfig.sourceChainConfigs[sourceChainSelector]` stores Canton-side CCV policy for messages coming from that source chain.
- `defaultCCVs` are fallback CCVs, not mandatory CCVs.
- `laneMandatedCCVs` are the mandatory CCVs.
- `receiverRequiredCCVs` lets the receiver choose its own required CCV set for execution.
- `receiverOptionalCCVs` with `receiverOptionalThreshold` means "require any N of these optional CCVs".

Effective inbound behavior on Canton:

- If `receiverRequiredCCVs = []`, execution falls back to `defaultCCVs`.
- If `receiverRequiredCCVs` contains explicit CCVs, those are used instead of `defaultCCVs`.
- If `receiverRequiredCCVs` contains the `useDefaultCCVs` sentinel, defaults are included plus any extra explicit CCVs.
- `laneMandatedCCVs` are always added on top.

Examples:

- `receiverRequiredCCVs = [attackerCCV]`, `laneMandatedCCVs = []`, `defaultCCVs = [officialCCV]`
  Result: required CCVs are `[attackerCCV]`.

- `receiverRequiredCCVs = []`, `laneMandatedCCVs = []`, `defaultCCVs = [officialCCV]`
  Result: required CCVs are `[officialCCV]`.

- `receiverRequiredCCVs = [attackerCCV]`, `laneMandatedCCVs = [officialCCV]`
  Result: required CCVs are `[attackerCCV, officialCCV]`.

## Why the audit PoC is not a security finding as written

The PoC showed:

- `defaultCCVs = [officialCCVRawAddress]`
- `laneMandatedCCVs = []`
- receiver supplied `receiverRequiredCCVs = [attackerCCVRawAddress]`
- execution succeeded with only the attacker-controlled CCV

This is consistent with the current intended Canton design:

- defaults are fallback only
- receivers are allowed to opt out of defaults and choose their own CCV set
- mandatory source-lane authentication must live in `laneMandatedCCVs`, not `defaultCCVs`

So this is best understood as a design clarification / docs ambiguity, not a vuln, unless product intended `defaultCCVs` to be mandatory.

## Receiver trust model on Canton

Per discussion, Canton treats the receiver as trusted configuration authority for arbitrary messaging.

- The executing party may be acting on behalf of the receiver contract owner.
- The receiver contract is expected to forward its configured CCV requirements into the router/offramp path.
- Because of that, allowing receiver-level CCV override is considered intentional.

## EVM comparison

EVM has similar fallback-vs-mandatory semantics, but the enforcement boundary is different.

- Receiver can return CCV requirements via `IAny2EVMMessageReceiverV2.getCCVsAndMinBlockConfirmations(...)`.
- If the receiver returns no CCVs, EVM falls back to lane defaults.
- `laneMandatedCCVs` are always enforced.
- OffRamp directly calls verifier contracts and runs `verifyMessage(...)`.

Key difference:

- EVM OffRamp enforces verifier execution directly.
- Canton trusts the selected CCV party to add a verification after doing its verification flow.

## Source-side CCV selection

On the source side, the verifier set is chosen during send from:

- user-specified CCVs in extra args
- fallback defaults if none specified
- lane-mandated CCVs
- pool-required CCVs

Then the final verifier set is fixed into source-side message state.

## EVM to Canton

For `evm -> canton`:

- EVM extra args contain EVM CCV addresses.
- Canton `GlobalConfig.sourceChainConfigs[sourceChainSelector]` contains Canton CCV raw addresses.
- There is no direct onchain mapping like `evmCCV -> cantonCCV`.
- Alignment is by verifier policy / proof compatibility, not by identical cross-chain addresses.

## Canton to EVM

For `canton -> evm`:

- Canton send flow uses Canton CCVs.
- EVM destination uses EVM-side verifier contracts/config for inbound verification.

## Signature / attestation alignment

The destination CCV checks verifier results against its own configured signer set for that source chain.

For the CommitteeVerifier-style Canton flow:

- signatures are checked against the Canton CCV's configured signer keys
- keyed by `sourceChainSelector`
- over the expected message hash / version-tagged hash

So the destination verifier must know the signer set for the verifier scheme, not the literal source-chain CCV address.
