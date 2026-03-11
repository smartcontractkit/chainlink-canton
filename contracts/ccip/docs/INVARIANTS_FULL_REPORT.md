# Canton CCIP Invariants — Full Audit Report

**Audit Date:** 2026-03-10
**Codebase:** `/contracts/ccip` (Daml / Canton)
**Classification per:** `invariants/CONTEXT.md`

### Classification Rules

- **PASS**: Invariant is implemented and enforced in the codebase.
- **FAIL**: Feature is missing, validation is missing, or implementation is partial. Further invariants depending on a missing feature are also FAIL.
- **N/A**: Invariant describes a mechanism specific to another chain family (e.g., EVM resolver pattern, SVM executor args, EOAs) that has no Canton equivalent *and* no Canton requirement.

Canton is a **net-new CCIP v2 chain**. There is no legacy V1 support needed. "IPoolV2" = the Canton pool interface. "V2 receiver" = the Canton receiver. Canton has **instant finality on the sending side** — it may always assume `blockConfirmations = 0` when sending and does not need to validate finality with downstream on send. Destination-side finality invariants still apply (Canton OffRamp receives messages from other chains that may use FTF). The **zero address sentinel** for "use default CCVs" must be supported even without a native zero address.

---

## 1. CCV Invariants (`CCV_INVARIANTS.md`)

### 1.1 Configuration


| ID        | Invariant                                                               | Grade    | Notes                                                                                                                                           |
| --------- | ----------------------------------------------------------------------- | -------- | ----------------------------------------------------------------------------------------------------------------------------------------------- |
| INV-CFG-1 | At least one CCV must exist across `defaultCCVs` and `laneMandatedCCVs` | **FAIL** | `GlobalConfig.UpdateDestChainConfig` (line 80) and `UpdateSourceChainConfig` (line 88) do not validate CCV set sizes. Tests use empty CCV sets. |
| INV-CFG-2 | No zero-value address in `defaultCCVs` or `laneMandatedCCVs`            | **FAIL** | No sentinel value defined; no validation at config time.                                                                                        |
| INV-CFG-3 | No duplicates within or across `defaultCCVs` and `laneMandatedCCVs`     | **FAIL** | `UpdateDestChainConfig` has no duplicate check. `UpdateSourceChainConfig` checks `onRampAddresses` uniqueness (line 92-93) but not CCV sets.    |
| INV-CFG-4 | CCV sets validated at config time                                       | **FAIL** | No validation in either `UpdateDestChainConfig` or `UpdateSourceChainConfig` for CCV sets.                                                      |
| INV-CFG-5 | `defaultCCVs` must be non-empty (OffRamp side)                          | **FAIL** | No check that `defaultCCVs` is non-empty.                                                                                                       |
| INV-CFG-6 | No zero-value address in OffRamp CCV sets                               | **FAIL** | Same as INV-CFG-2; no sentinel defined.                                                                                                         |
| INV-CFG-7 | No duplicates enforced in OffRamp CCVs                                  | **FAIL** | Dedup at use time (`OffRamp.daml:249`), not at config time.                                                                                     |
| INV-CFG-8 | Pool `outboundCCVs`/`thresholdOutboundCCVs` no duplicates               | **N/A**  | Canton pools have `outboundCCVs`/`inboundCCVs` only; no threshold sets (AdvancedPoolHooks is EVM-specific).                                     |
| INV-CFG-9 | Pool `inboundCCVs`/`thresholdInboundCCVs` no duplicates                 | **N/A**  | Same.                                                                                                                                           |


### 1.2 Source Side (OnRamp)


| ID         | Invariant                                                               | Grade    | Notes                                                                                                                    |
| ---------- | ----------------------------------------------------------------------- | -------- | ------------------------------------------------------------------------------------------------------------------------ |
| INV-SRC-1  | User-provided CCVs must not contain duplicates                          | **FAIL** | `senderRequiredCCVs` passed through without duplicate check. Dedup happens at merge time but invalid input not rejected. |
| INV-SRC-2  | Zero-value address as "include defaults" placeholder                    | **FAIL** | No sentinel value supported. Context says it must be.                                                                    |
| INV-SRC-3  | Legacy extraArgs get defaults automatically                             | **N/A**  | Canton is net-new; no legacy extraArgs path.                                                                             |
| INV-SRC-4  | Empty user CCVs → defaults for data msgs; token-only → no defaults      | **PASS** | `OnRamp.daml:134-139` branches correctly on `isTokenOnlyTransfer`.                                                       |
| INV-SRC-5  | Defaults used when user provides no CCVs or uses zero-value placeholder | **FAIL** | Defaults apply when empty, but zero-value placeholder not supported.                                                     |
| INV-SRC-6  | Defaults NOT applied for pure token-only unless explicitly requested    | **PASS** | Token-only with empty `senderRequiredCCVs` → `[]` for user/default CCVs.                                                 |
| INV-SRC-7  | Expanding zero-value placeholder skips already-present defaults         | **FAIL** | No zero-value placeholder mechanism.                                                                                     |
| INV-SRC-8  | Lane-mandated CCVs always added                                         | **PASS** | `OnRamp.daml:139` always appends `laneMandatedCCVs`.                                                                     |
| INV-SRC-9  | Lane-mandated already in user list are skipped (dedup)                  | **PASS** | `dedup` applied to merged list.                                                                                          |
| INV-SRC-10 | User entry with custom `ccvArgs` takes precedence                       | **FAIL** | Canton has no per-CCV `ccvArgs`. V2 ExtraArgs requires this.                                                             |
| INV-SRC-11 | Pool without V2 falls back to `defaultCCVs`                             | **N/A**  | Canton is V2-only; all pools implement the V2 interface.                                                                 |
| INV-SRC-12 | Pool with empty `getRequiredCCVs` falls back to `defaultCCVs`           | **FAIL** | If pool returns `[]`, no pool CCVs added and no explicit fallback to defaults.                                           |
| INV-SRC-13 | Pool zero-value address replaced with `defaultCCVs`                     | **FAIL** | No zero-value address sentinel.                                                                                          |
| INV-SRC-14 | Pool-required CCVs already in merged list are skipped                   | **PASS** | `dedup (requiredCCVs ++ poolCCVs)` in `SendingMessageV1.daml:235`.                                                       |
| INV-SRC-15 | User's custom `ccvArgs` preserved for pool-required CCV                 | **FAIL** | No per-CCV args.                                                                                                         |
| INV-SRC-16 | Final CCV list order: user/default → lane-mandated → pool               | **PASS** | Order matches: `(userOrDefaults ++ laneMandatedCCVs)` then `++ poolCCVs`.                                                |
| INV-SRC-17 | Final merged CCV list has no duplicates                                 | **PASS** | `dedup` applied.                                                                                                         |
| INV-SRC-18 | Each CCV's `getFee` called for per-CCV fees                             | **PASS** | `AddCCVFee` called per CCV.                                                                                              |
| INV-SRC-19 | Each CCV generates receipt with fee details                             | **PASS** | `assembleReceipts` builds per-CCV receipts.                                                                              |
| INV-SRC-20 | `ccvAndExecutorHash = keccak256([addressLength][ccv1..ccvN][executor])` | **PASS** | `MessageCodecV1.computeCCVAndExecutorHash` matches.                                                                      |
| INV-SRC-21 | `ccvAndExecutorHash` is informational, not verified on destination      | **PASS** | OffRamp does not verify `ccvAndExecutorHash`.                                                                            |
| INV-SRC-22 | Token-only with no user CCVs → only pool + lane-mandated                | **PASS** | Correct behavior in `OnRamp.daml:134-139`.                                                                               |


### 1.3 Destination Side (OffRamp)


| ID         | Invariant                                                        | Grade    | Notes                                                                                                       |
| ---------- | ---------------------------------------------------------------- | -------- | ----------------------------------------------------------------------------------------------------------- |
| INV-DST-1  | Receiver V2 → `getCCVsAndMinBlockDepth` called                   | **FAIL** | `CCIPReceiver` has `requiredCCVs` but no `minBlockDepth`, no `optionalCCVs`, no `optionalThreshold`.        |
| INV-DST-2  | Non-V2 receiver → defaults used                                  | **N/A**  | Canton is V2-only; all receivers implement V2.                                                              |
| INV-DST-3  | `requiredCCVs` must not contain duplicates                       | **FAIL** | `UpdateRequiredCCVs` does not reject duplicates.                                                            |
| INV-DST-4  | `optionalCCVs` must not contain duplicates                       | **FAIL** | No optional CCVs implemented.                                                                               |
| INV-DST-5  | `optionalThreshold` must not exceed `optionalCCVs.length`        | **FAIL** | No optional CCV mechanism.                                                                                  |
| INV-DST-6  | Empty required + threshold==0 → fallback to defaults             | **FAIL** | Fallback logic exists (`OffRamp.daml:74-76`) but uses different mechanism — no zero-value address sentinel. |
| INV-DST-7  | Zero-value address in required CCVs → "include defaults"         | **FAIL** | No zero-value address sentinel.                                                                             |
| INV-DST-8  | Pool V2 with empty `getRequiredCCVs` → fallback to `defaultCCVs` | **FAIL** | Empty pool CCVs are just `[]`; no explicit fallback to defaults.                                            |
| INV-DST-9  | Pool V2 with empty array → fallback to defaults                  | **N/A**  | Same as DST-8; counted once.                                                                                |
| INV-DST-10 | Pool zero-value address → include `defaultCCVs`                  | **FAIL** | No sentinel.                                                                                                |
| INV-DST-11 | Lane-mandated CCVs always included                               | **PASS** | `OffRamp.daml:248` includes `sourceConfig.laneMandatedCCVs`.                                                |
| INV-DST-12 | Default CCVs when zero-value appears                             | **N/A**  | No zero-value address; defaults apply when both receiver and pool lack CCVs (`OffRamp.daml:241-244`).       |
| INV-DST-13 | Defaults added at most once                                      | **PASS** | Single `resolvedDefaultCCVs` assignment.                                                                    |
| INV-DST-14 | Token-only with tokens → skip receiver CCVs and defaults         | **PASS** | `resolvedDefaultCCVs = []` and `resolvedReceiverRequiredCCVs = []` for token-only.                          |
| INV-DST-15 | Token-only without tokens (no-op) → use defaults                 | **PASS** | Non-token messages are not token-only; defaults can apply.                                                  |
| INV-DST-16 | Final required CCV list deduplicated                             | **PASS** | `dedup` in `OffRamp.daml:249`.                                                                              |
| INV-DST-17 | CCV in both required and optional → removed from optional        | **N/A**  | No optional CCVs.                                                                                           |
| INV-DST-18 | All required CCVs must be present                                | **PASS** | `Foldable.forA_` asserts each required CCV is present (`OffRamp.daml:258-260`).                             |
| INV-DST-19 | At least `optionalThreshold` of optional CCVs present            | **N/A**  | No optional CCVs.                                                                                           |
| INV-DST-20 | Extra CCVs beyond required/optional are ignored                  | **PASS** | Only required CCVs checked.                                                                                 |
| INV-DST-21 | Each CCV has `verifyMessage` called                              | **PASS** | `AddCCVVerification` called per CCV via `CCIPReceiver.Execute`.                                             |
| INV-DST-22 | `ccvs.length == verifierResults.length`                          | **N/A**  | Canton uses `ccvVerifications` list; no separate `verifierResults` array.                                   |


### 1.4 CCV Address Stability


| ID         | Invariant                                           | Grade    | Notes                                                                           |
| ---------- | --------------------------------------------------- | -------- | ------------------------------------------------------------------------------- |
| INV-ADDR-1 | CCV address stable across verifier upgrades         | **PASS** | CCVs identified by `instanceId@owner`; interface-based design supports upgrade. |
| INV-ADDR-2 | CCV address resolves to concrete verifier per chain | **FAIL** | No resolver mechanism; CCV must be explicitly configured per chain.             |
| INV-ADDR-3 | Resolution transparent to caller                    | **FAIL** | No resolution layer; OnRamp/OffRamp interact with CCV directly.                 |
| INV-RES-1  | EVM: CCV addresses are resolver instances           | **N/A**  | EVM-specific.                                                                   |
| INV-RES-2  | `VersionedVerifierResolver` maps by chain/version   | **N/A**  | EVM-specific.                                                                   |
| INV-RES-3  | Zero-value outbound resolution → cannot send        | **N/A**  | EVM-specific.                                                                   |
| INV-RES-4  | Zero-value inbound resolution → cannot execute      | **N/A**  | EVM-specific.                                                                   |


### 1.5 Cross-Cutting


| ID       | Invariant                                                      | Grade    | Notes                                                                                      |
| -------- | -------------------------------------------------------------- | -------- | ------------------------------------------------------------------------------------------ |
| INV-CC-1 | Every message verified by at least one CCV                     | **FAIL** | Config allows empty CCV sets; token-only with empty pool CCVs can have zero required CCVs. |
| INV-CC-2 | Source-side CCV list is informational for offchain             | **PASS** | Correct.                                                                                   |
| INV-CC-3 | Source and destination CCV lists independently configured      | **PASS** | `DestChainConfig` vs `SourceChainConfig`.                                                  |
| INV-CC-4 | Same verifier may have different addresses on different chains | **PASS** | CCV addresses are chain/participant-specific.                                              |


---

## 2. Finality Invariants (`FINALITY_INVARIANTS.md`)

### 2.1 Encoding and Representation


| ID            | Invariant                                    | Grade    | Notes                                                             |
| ------------- | -------------------------------------------- | -------- | ----------------------------------------------------------------- |
| INV-FIN-ENC-1 | `blockConfirmations == 0` → default finality | **PASS** | `OnRamp.daml:146` uses `fromOptional 0 blockConfirmations`.       |
| INV-FIN-ENC-2 | `blockConfirmations != 0` → FTF              | **PASS** | `SendingMessageV1.daml:286` sets `finality = blockConfirmations`. |
| INV-FIN-ENC-3 | Legacy extraArgs default to `0`              | **N/A**  | Canton is net-new; no legacy extraArgs.                           |


### 2.2 Source Side (OnRamp) — Canton instant-finality exemption applies


| ID            | Invariant                                                                                 | Grade    | Notes                                                                                                                                                                                        |
| ------------- | ----------------------------------------------------------------------------------------- | -------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| INV-FIN-SRC-1 | `blockConfirmations` flows through OnRamp to downstream                                   | **PASS** | Flows to `message.finality`, pool `GetRequiredCCVs` (has `finality` param), executor. Canton sending side always sends 0, so flow is nominal.                                                |
| INV-FIN-SRC-2 | OnRamp does not validate `blockConfirmations`                                             | **PASS** | Passes through without validation. Correct for instant-finality chain.                                                                                                                       |
| INV-FIN-SRC-3 | Each CCV receives `blockConfirmations` in `getFee`                                        | **FAIL** | CCV `CalculateFee` receives `sendingMessageCid`, not `blockConfirmations` as an explicit parameter. Canton exempt from sending-side finality validation, but interface should still pass it. |
| INV-FIN-SRC-4 | V2 pools receive `blockConfirmationsRequested` in `lockOrBurn`/`getFee`/`getRequiredCCVs` | **FAIL** | `GetRequiredCCVs` has `finality : Int` but `LockOrBurn`/`CalculateFee` do not receive it explicitly.                                                                                         |
| INV-FIN-SRC-5 | V1 pools cannot support FTF                                                               | **N/A**  | No V1 pools.                                                                                                                                                                                 |
| INV-FIN-SRC-6 | V1 pool + non-empty `tokenArgs` → revert                                                  | **N/A**  | Same.                                                                                                                                                                                        |
| INV-FIN-SRC-7 | Executor receives `requestedBlockDepth` in `getFee`                                       | **PASS** | `TestExecutor.daml:51` reads `sm.blockConfirmations`.                                                                                                                                        |
| INV-FIN-SRC-8 | Executor rejects FTF below `minBlockConfirmations`                                        | **PASS** | `TestExecutor.daml:50-51`.                                                                                                                                                                   |
| INV-FIN-SRC-9 | Finality (`0`) always accepted by executor                                                | **PASS** | Same logic; `0` bypasses minimum.                                                                                                                                                            |


### 2.3 Destination Side (OffRamp) — Must handle inbound FTF


| ID             | Invariant                                                         | Grade    | Notes                                                                        |
| -------------- | ----------------------------------------------------------------- | -------- | ---------------------------------------------------------------------------- |
| INV-FIN-DST-1  | OffRamp queries receiver for `minBlockDepth`                      | **FAIL** | No `minBlockDepth` in receiver interface.                                    |
| INV-FIN-DST-2  | `minBlockDepth == 0` → reject FTF messages                        | **FAIL** | No receiver finality check.                                                  |
| INV-FIN-DST-3  | `minBlockDepth > 0` → require `message.finality >= minBlockDepth` | **FAIL** | Same.                                                                        |
| INV-FIN-DST-4  | Non-V2 receiver → `minBlockDepth` defaults to 0                   | **N/A**  | Canton is V2-only.                                                           |
| INV-FIN-DST-5  | Non-programmable accounts cannot receive FTF                      | **N/A**  | No EOAs in Canton/Daml.                                                      |
| INV-FIN-DST-6  | Token-only: `message.finality` passed to pool `releaseOrMint`     | **FAIL** | `TokenReceiveTicket` has no `finality`; pool does not receive it on inbound. |
| INV-FIN-DST-7  | Token-only: receiver's `minBlockDepth` not checked                | **N/A**  | No receiver finality check exists.                                           |
| INV-FIN-DST-8  | OffRamp passes `message.finality` to V2 pool `releaseOrMint`      | **FAIL** | Pool does not receive finality on inbound.                                   |
| INV-FIN-DST-9  | V1 pools receive no finality parameter                            | **N/A**  | No V1 pools.                                                                 |
| INV-FIN-DST-10 | OffRamp passes `message.finality` to pool CCV resolution          | **FAIL** | Pool does not receive finality for inbound CCV resolution.                   |


### 2.4 Pool Finality


| ID              | Invariant                                                | Grade    | Notes                                                                                |
| --------------- | -------------------------------------------------------- | -------- | ------------------------------------------------------------------------------------ |
| INV-FIN-POOL-1  | `WAIT_FOR_FINALITY = 0`                                  | **PASS** | `0` used for finality throughout.                                                    |
| INV-FIN-POOL-2  | Pool's min block depth config                            | **FAIL** | `ChainPoolConfig` has no `minBlockDepth`.                                            |
| INV-FIN-POOL-3  | Pool owner configures min block depth                    | **FAIL** | Not implemented.                                                                     |
| INV-FIN-POOL-4  | Pool reverts if FTF requested but not enabled            | **N/A**  | Canton sending side exempt (instant finality).                                       |
| INV-FIN-POOL-5  | Pool reverts if FTF below configured minimum             | **N/A**  | Same.                                                                                |
| INV-FIN-POOL-6  | V1 `lockOrBurn` passes `WAIT_FOR_FINALITY`               | **N/A**  | No V1 pools.                                                                         |
| INV-FIN-POOL-7  | `releaseOrMint` same finality validation as `lockOrBurn` | **FAIL** | No finality validation on inbound.                                                   |
| INV-FIN-POOL-8  | Separate rate limit buckets for default vs FTF           | **FAIL** | `RateLimitMode_CustomFinality` defined but unused; single bucket per chain.          |
| INV-FIN-POOL-9  | FTF → custom bucket consumed                             | **FAIL** | Only `RateLimitMode_DefaultFinality` used.                                           |
| INV-FIN-POOL-10 | Custom bucket not configured → fallback to default       | **N/A**  | No custom bucket implemented.                                                        |
| INV-FIN-POOL-11 | Default finality → default bucket consumed               | **PASS** | Rate limiter always uses default bucket (which is correct when only default exists). |
| INV-FIN-POOL-12 | Rate limiter config per-remote-chain                     | **N/A**  | Rate limiters are per-chain but no finality split.                                   |
| INV-FIN-POOL-13 | Separate fee fields for default vs custom finality       | **FAIL** | Single `feeUSDCents`; no finality-based fee split.                                   |
| INV-FIN-POOL-14 | `customBlockConfirmationsTransferFeeBps` < 10000         | **N/A**  | No `feeBps` in Canton.                                                               |
| INV-FIN-POOL-15 | `getRequiredCCVs` receives `blockConfirmationsRequested` | **PASS** | `TokenPool_GetRequiredCCVs` has `finality : Int` parameter.                          |


### 2.5 Executor Finality


| ID             | Invariant                                          | Grade    | Notes                      |
| -------------- | -------------------------------------------------- | -------- | -------------------------- |
| INV-FIN-EXEC-1 | Configurable `minBlockConfirmations`               | **PASS** | `TestExecutor.daml:20`.    |
| INV-FIN-EXEC-2 | Rejects FTF below `minBlockConfirmations`          | **PASS** | `TestExecutor.daml:50-51`. |
| INV-FIN-EXEC-3 | `minBlockConfirmations` is floor for FTF depth     | **PASS** | Same logic.                |
| INV-FIN-EXEC-4 | CCV allowlist and max CCVs independent of finality | **PASS** | `TestExecutor.daml:52-56`. |


### 2.6 CCV Finality


| ID            | Invariant                                         | Grade    | Notes                                                                                        |
| ------------- | ------------------------------------------------- | -------- | -------------------------------------------------------------------------------------------- |
| INV-FIN-CCV-1 | CCVs receive `blockConfirmations` in `getFee`     | **FAIL** | CCV `CalculateFee` receives `sendingMessageCid`; `blockConfirmations` not passed explicitly. |
| INV-FIN-CCV-2 | CCVs can reject `blockConfirmations` by reverting | **FAIL** | CCV can abort generally, but has no access to `blockConfirmations` as a parameter to check.  |


### 2.7 Re-org and Safety


| ID              | Invariant                                                   | Grade    | Notes                                        |
| --------------- | ----------------------------------------------------------- | -------- | -------------------------------------------- |
| INV-FIN-REORG-1 | `messageNumber` unique per-lane only for finalized messages | **PASS** | Protocol-level; Canton has instant finality. |
| INV-FIN-REORG-2 | FTF shifts re-org risk to receiver/pool/integrators         | **PASS** | Conceptual; FTF is opt-in.                   |


### 2.8 Opt-in Requirements


| ID              | Invariant                                           | Grade    | Notes                                                                 |
| --------------- | --------------------------------------------------- | -------- | --------------------------------------------------------------------- |
| INV-FIN-OPTIN-1 | Any layer rejecting FTF blocks the message          | **FAIL** | Executor can reject; pool and CCV do not enforce finality on inbound. |
| INV-FIN-OPTIN-2 | All layers default to rejecting FTF                 | **FAIL** | Receiver and pool lack opt-in/rejection mechanism.                    |
| INV-FIN-OPTIN-3 | V1 pools and non-V2 receivers implicitly reject FTF | **N/A**  | Canton is V2-only.                                                    |


### 2.9 Cross-Cutting


| ID           | Invariant                                              | Grade    | Notes                                                                       |
| ------------ | ------------------------------------------------------ | -------- | --------------------------------------------------------------------------- |
| INV-FIN-CC-1 | `finality` is message-level                            | **PASS** | `MessageV1.finality` is a single field.                                     |
| INV-FIN-CC-2 | Finality affects fees, rate limiting, CCV requirements | **FAIL** | Only executor fees affected; pool and CCV do not differentiate by finality. |


---

## 3. Encoding Invariants (`ENCODING_INVARIANTS.md`)

### 3.1 MessageV1 Encoding


| ID         | Invariant                                                              | Grade    | Notes                                                                                                                                                                                    |
| ---------- | ---------------------------------------------------------------------- | -------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| INV-MSG-1  | Chain-agnostic wire format; `messageId = keccak256(encodedMessageV1)`  | **PASS** | `MessageCodecV1.daml` matches EVM wire format; `messageId` uses `keccak256`.                                                                                                             |
| INV-MSG-2  | Version byte is `1`; decoding rejects other versions                   | **FAIL** | Version decoded (`decodeMessageV1` line 664-666) but not validated against `1`.                                                                                                          |
| INV-MSG-3  | Decoding strict: entire byte array consumed exactly                    | **FAIL** | No check that all bytes are consumed; trailing bytes accepted.                                                                                                                           |
| INV-MSG-4  | Max 1 token transfer                                                   | **PASS** | `tokenTransfer : Optional TokenTransferV1`.                                                                                                                                              |
| INV-MSG-5  | `MESSAGE_V1_BASE_SIZE = 77`; reject input < 77 bytes                   | **PASS** | Static section is 67 bytes in Canton (no separate `executionGasLimit`/`ccipReceiveGasLimit` as uint32 fields; they are present but at same offsets). `extractBytes` rejects short input. |
| INV-MSG-6  | Variable-length fields: `uint8` for addresses, `uint16` for data blobs | **PASS** | Correct length prefix conventions used.                                                                                                                                                  |
| INV-MSG-7  | Addresses use minimal native byte length; EVM source abi-encoded       | **PASS** | Canton uses 32-byte addresses consistently.                                                                                                                                              |
| INV-MSG-8  | Destination addresses validated against `addressBytesLength`           | **N/A**  | Canton uses fixed 32-byte party addresses; no per-chain address length.                                                                                                                  |
| INV-MSG-9  | Token transfer encoding follows same address conventions               | **PASS** | Same length prefix conventions.                                                                                                                                                          |
| INV-MSG-10 | `extraData` carries pool-specific data                                 | **PASS** | `LockReleaseTokenPool` sets `extraData = destPoolData`.                                                                                                                                  |
| INV-MSG-11 | `tokenTransfer` length-prefixed with `uint16`; bytes must match        | **PASS** | `decodeMessageV1` uses `tokenLen` for extraction.                                                                                                                                        |


### 3.2 ExtraArgs Encoding


| ID        | Invariant                                                           | Grade    | Notes                                                                                                                                            |
| --------- | ------------------------------------------------------------------- | -------- | ------------------------------------------------------------------------------------------------------------------------------------------------ |
| INV-ENC-1 | ExtraArgsV3 tag `0xa69dd4aa`; all chains should use the same format | **FAIL** | Canton uses `CantonExtraArgsV1` (Daml native types), not `GenericExtraArgsV3` wire format. Invariant explicitly says "not a chain-specific one". |
| INV-ENC-2 | CCV count as `uint8` (max 255)                                      | **N/A**  | Depends on INV-ENC-1; Canton uses Daml lists.                                                                                                    |
| INV-ENC-3 | `ccvs` and `ccvArgs` same length                                    | **N/A**  | Same.                                                                                                                                            |
| INV-ENC-4 | `GENERIC_EXTRA_ARGS_V3_BASE_SIZE = 17`                              | **N/A**  | Same.                                                                                                                                            |
| INV-ENC-5 | `addressLength == 0` encodes zero-value address                     | **N/A**  | Same.                                                                                                                                            |
| INV-ENC-6 | Strict decoding for ExtraArgsV3                                     | **N/A**  | Same.                                                                                                                                            |
| INV-ENC-7 | SVM executor args                                                   | **N/A**  | SVM-specific.                                                                                                                                    |
| INV-ENC-8 | Sui executor args                                                   | **N/A**  | Sui-specific.                                                                                                                                    |
| INV-ENC-9 | Strict decoding for executor args                                   | **N/A**  | SVM/Sui-specific.                                                                                                                                |


### 3.3 Pool Remote Config Encoding


| ID             | Invariant                                                  | Grade    | Notes                                                                                                                                               |
| -------------- | ---------------------------------------------------------- | -------- | --------------------------------------------------------------------------------------------------------------------------------------------------- |
| INV-POOL-ENC-1 | `remotePoolAddresses` stored as raw bytes                  | **PASS** | `ChainPoolConfig.remotePools : [BytesHex]`.                                                                                                         |
| INV-POOL-ENC-2 | `remoteTokenAddress` must be non-empty                     | **FAIL** | No explicit non-empty validation.                                                                                                                   |
| INV-POOL-ENC-3 | Multiple remote pool addresses per chain                   | **PASS** | `remotePools : [BytesHex]` — list supports multiple.                                                                                                |
| INV-POOL-ENC-4 | Remote pool addresses must be non-empty; format must match | **FAIL** | No non-empty validation on config. Pool checks `sourcePoolAddress elem remotePools` at use time but doesn't validate address format at config time. |


---

## 4. Lane Configuration Invariants (`LANE_INVARIANTS.md`)

### 4.1 Destination Chain Config (Source Side)


| ID         | Invariant                                                   | Grade    | Notes                                                                                                                                                  |
| ---------- | ----------------------------------------------------------- | -------- | ------------------------------------------------------------------------------------------------------------------------------------------------------ |
| INV-LCFG-1 | OnRamp stores per-lane config                               | **PASS** | `DestChainConfig` has `offRampAddress`, CCV sets, network fees. Missing EVM-specific: `addressBytesLength`, `baseExecutionGasCost`, `defaultExecutor`. |
| INV-LCFG-2 | `destChainSelector` must be non-zero and not equal to local | **FAIL** | `UpdateDestChainConfig` does not validate.                                                                                                             |
| INV-LCFG-3 | `addressBytesLength` must be non-zero                       | **N/A**  | Canton uses fixed 32-byte addresses.                                                                                                                   |
| INV-LCFG-4 | `baseExecutionGasCost` must be non-zero                     | **N/A**  | No gas model in Canton.                                                                                                                                |
| INV-LCFG-5 | Default executor must be configured                         | **N/A**  | Executor chosen per message via `executorCid` in `CantonExtraArgsV1`.                                                                                  |
| INV-LCFG-6 | Dest OffRamp address length must equal `addressBytesLength` | **PASS** | Canton stores `offRampAddress : BytesHex` directly.                                                                                                    |
| INV-LCFG-7 | Lane config updates restricted to owner                     | **PASS** | `controller ccipOwner`.                                                                                                                                |


### 4.2 Source Chain Config (Destination Side)


| ID          | Invariant                                       | Grade    | Notes                                                                               |
| ----------- | ----------------------------------------------- | -------- | ----------------------------------------------------------------------------------- |
| INV-LCFG-8  | OffRamp stores per-lane config                  | **PASS** | `SourceChainConfig` has `isEnabled`, `onRampAddresses`, CCV sets.                   |
| INV-LCFG-9  | `sourceChainSelector` must be non-zero          | **FAIL** | No validation.                                                                      |
| INV-LCFG-10 | Source chain config updates restricted to owner | **PASS** | `controller ccipOwner`.                                                             |
| INV-LCFG-11 | OnRamp addresses in allowlist must be non-empty | **PASS** | **Fixed.** `GlobalConfig.daml:90-91` validates `not (null config.onRampAddresses)`. |


### 4.3 Lane Pause and Disable


| ID          | Invariant                                    | Grade    | Notes                                                      |
| ----------- | -------------------------------------------- | -------- | ---------------------------------------------------------- |
| INV-PAUSE-1 | Mechanism to pause lane on source side       | **PASS** | `DestChainConfig.isEnabled` checked; disabled lanes abort. |
| INV-PAUSE-2 | Paused lane → fee quoting also reverts       | **PASS** | Disabled lanes abort before fee quoting.                   |
| INV-PAUSE-3 | `isEnabled = false` on OffRamp disables lane | **PASS** | `OffRamp.daml:213` checks `c.isEnabled`.                   |


### 4.4 Inflight Message Handling


| ID             | Invariant                                               | Grade    | Notes                                                                                                     |
| -------------- | ------------------------------------------------------- | -------- | --------------------------------------------------------------------------------------------------------- |
| INV-INFLIGHT-1 | Removing OnRamp from allowlist blocks inflight messages | **PASS** | **Fixed.** OffRamp now validates `message.onRampAddress elem c.onRampAddresses` (`OffRamp.daml:214-215`). |
| INV-INFLIGHT-2 | `isEnabled = false` blocks inflight messages            | **PASS** | `isEnabled` checked at execution time.                                                                    |
| INV-INFLIGHT-3 | Removed remote pool fails token validation              | **PASS** | Pool validates `sourcePoolAddress` against `chainPoolConfigs`.                                            |


### 4.5 Router Configuration


| ID        | Invariant                                   | Grade   | Notes                                                                                            |
| --------- | ------------------------------------------- | ------- | ------------------------------------------------------------------------------------------------ |
| INV-RTR-1 | Router maps each dest to exactly one OnRamp | **N/A** | "Not all chain families use a router." Canton uses `PerPartyRouter` with different architecture. |
| INV-RTR-2 | Multiple OffRamps per source chain          | **N/A** | Same.                                                                                            |
| INV-RTR-3 | Zero-value OnRamp disables sending          | **N/A** | Canton uses `isEnabled` flag instead.                                                            |
| INV-RTR-4 | Only registered OffRamps trigger delivery   | **N/A** | Different architecture.                                                                          |
| INV-RTR-5 | Router config updates restricted to owner   | **N/A** | Same.                                                                                            |


### 4.6 OnRamp Validation on Destination


| ID           | Invariant                                                   | Grade    | Notes                                                                                    |
| ------------ | ----------------------------------------------------------- | -------- | ---------------------------------------------------------------------------------------- |
| INV-ONRVAL-1 | OffRamp validates `message.onRampAddress`                   | **PASS** | **Fixed.** `OffRamp.daml:214-215` checks `message.onRampAddress elem c.onRampAddresses`. |
| INV-ONRVAL-2 | OffRamp validates `message.offRampAddress` matches itself   | **FAIL** | No validation of `message.offRampAddress`.                                               |
| INV-ONRVAL-3 | OffRamp validates `message.destChainSelector` matches local | **FAIL** | No validation of `message.destChainSelector` against `globalConfig.chainSelector`.       |


### 4.7 Risk Management (Curse)


| ID        | Invariant                                          | Grade    | Notes                                                     |
| --------- | -------------------------------------------------- | -------- | --------------------------------------------------------- |
| INV-RMN-1 | Per-chain curse blocks operations                  | **PASS** | `IsCursedForChain` used; OnRamp, OffRamp, pool all check. |
| INV-RMN-2 | Curse on dest chain blocks sending                 | **PASS** | OnRamp checks curse.                                      |
| INV-RMN-3 | Curse on source chain blocks execution             | **PASS** | OffRamp checks curse (`OffRamp.daml:122-125, 192-195`).   |
| INV-RMN-4 | Curse blocks pool `lockOrBurn` and `releaseOrMint` | **PASS** | Pool checks curse for both directions.                    |
| INV-RMN-5 | Curse independent of lane config                   | **PASS** | `RMNRemote` is independent.                               |


### 4.8 Reentrancy Protection


| ID          | Invariant                                     | Grade    | Notes                                                                                  |
| ----------- | --------------------------------------------- | -------- | -------------------------------------------------------------------------------------- |
| INV-REENT-1 | OnRamp send protected by reentrancy guard     | **PASS** | Daml is single-threaded with atomic transactions; reentrancy is inherently impossible. |
| INV-REENT-2 | OffRamp execute protected by reentrancy guard | **PASS** | Same.                                                                                  |


### 4.9 Upgradability


| ID        | Invariant                                    | Grade    | Notes                                                           |
| --------- | -------------------------------------------- | -------- | --------------------------------------------------------------- |
| INV-UPG-1 | CCV upgradeable behind stable address        | **PASS** | Interface-based design supports upgrades.                       |
| INV-UPG-2 | Pool supports multiple remote pool addresses | **PASS** | `remotePools : [BytesHex]`.                                     |
| INV-UPG-3 | Executor upgradeable via config update       | **PASS** | Executor chosen per message; new messages use updated executor. |
| INV-UPG-4 | Upgrades require no user action or downtime  | **PASS** | Consistent with Daml upgrade patterns.                          |


### 4.10 Access Control


| ID                                      | Invariant      | Grade    | Notes                                                    |
| --------------------------------------- | -------------- | -------- | -------------------------------------------------------- |
| INV-AC (OnRamp lane config)             | Contract owner | **PASS** | `controller ccipOwner`.                                  |
| INV-AC (OffRamp source chain config)    | Contract owner | **PASS** | `controller ccipOwner`.                                  |
| INV-AC (Router ramp updates)            | Contract owner | **PASS** | `PerPartyRouter` controlled by `partyOwner`/`ccipOwner`. |
| INV-AC (Token pool remote chain config) | Contract owner | **PASS** | `controller poolOwner`.                                  |
| INV-AC (Message sending)                | Permissionless | **FAIL** | `controller partyOwner` — not permissionless.            |
| INV-AC (Message execution)              | Permissionless | **FAIL** | `controller partyOwner` — not permissionless.            |


---

## 5. Token Pool Invariants (`TOKEN_POOL_INVARIANTS.md`)

### 5.1 Pool Types


| ID         | Invariant                               | Grade    | Notes                                                           |
| ---------- | --------------------------------------- | -------- | --------------------------------------------------------------- |
| INV-POOL-1 | All pool types implement same interface | **PASS** | `LockReleaseTokenPool` implements `ITokenPool`.                 |
| INV-POOL-2 | Lock-release pools use lockbox          | **PASS** | Tokens move via `TransferFactory` to/from `poolOwner`.          |
| INV-POOL-3 | Burn-mint pools destroy/create tokens   | **FAIL** | No burn-mint pool implementation. Required for V2 completeness. |


### 5.2 Access Control


| ID         | Invariant                             | Grade    | Notes                                                                                                                 |
| ---------- | ------------------------------------- | -------- | --------------------------------------------------------------------------------------------------------------------- |
| INV-POOL-4 | Only OnRamp may call `lockOrBurn`     | **FAIL** | `LockOrBurn` uses `controller caller, poolOwner`. No check that caller is OnRamp.                                     |
| INV-POOL-5 | Only OffRamp may call `releaseOrMint` | **FAIL** | `ReleaseFromTicket` checks `caller == ticket.tokenReceiver`. Security via ticket ownership, not OffRamp verification. |
| INV-POOL-6 | Remote chain must be configured       | **PASS** | Aborts if config missing.                                                                                             |


### 5.3 V1 vs V2 Interface


| ID         | Invariant                                                                          | Grade    | Notes                                                                                                                      |
| ---------- | ---------------------------------------------------------------------------------- | -------- | -------------------------------------------------------------------------------------------------------------------------- |
| INV-POOL-7 | V1 pools: lock/burn without finality                                               | **FAIL** | Canton is V2-only but pool `LockOrBurn` does not take `blockConfirmationsRequested` as a parameter. Missing V2 feature.    |
| INV-POOL-8 | V2 pools: `getFee`, `getRequiredCCVs`, finality-aware `lockOrBurn`/`releaseOrMint` | **FAIL** | `CalculateFee` and `LockOrBurn` do not receive `blockConfirmationsRequested`. `getTokenTransferFeeConfig` not implemented. |
| INV-POOL-9 | V1 pools cannot support FTF                                                        | **FAIL** | Canton is V2-only but lacks FTF support in pool.                                                                           |


### 5.4 Fee Deduction (V2)


| ID          | Invariant                                    | Grade    | Notes                                                                      |
| ----------- | -------------------------------------------- | -------- | -------------------------------------------------------------------------- |
| INV-POOL-10 | V2 pools deduct proportional fee (`feeBps`)  | **FAIL** | Pool uses flat `feeUSDCents`; no proportional deduction from token amount. |
| INV-POOL-11 | `feeBps` depends on finality                 | **FAIL** | No `feeBps`.                                                               |
| INV-POOL-12 | `feeBps` must be < 10000                     | **FAIL** | No `feeBps`.                                                               |
| INV-POOL-13 | Fee config disabled → `feeBps` defaults to 0 | **FAIL** | No `feeBps`.                                                               |


### 5.5 Decimal Handling


| ID          | Invariant                                                  | Grade    | Notes                                                      |
| ----------- | ---------------------------------------------------------- | -------- | ---------------------------------------------------------- |
| INV-POOL-14 | Pools encode local token decimals into `sourcePoolData`    | **PASS** | `encodeUint256 (intToNumeric decimals)` in `LockOrBurn`.   |
| INV-POOL-15 | Empty `sourcePoolData` → same decimals assumed             | **FAIL** | Destination pool not implemented; should handle this case. |
| INV-POOL-16 | `sourcePoolData` exactly 32 bytes encoding `uint8` decimal | **PASS** | `encodeUint256` produces 32 bytes.                         |
| INV-POOL-17 | Higher→lower decimals: rounds down                         | **FAIL** | No inbound decimal conversion in pool.                     |
| INV-POOL-18 | Lower→higher decimals: overflow reverts                    | **FAIL** | Same — no inbound decimal conversion.                      |


### 5.6 destBytesOverhead


| ID          | Invariant                                             | Grade    | Notes                                             |
| ----------- | ----------------------------------------------------- | -------- | ------------------------------------------------- |
| INV-POOL-19 | `destBytesOverhead` data budget for fee calculation   | **PASS** | Field exists in `PoolFeeConfig`.                  |
| INV-POOL-20 | Minimum `destBytesOverhead` is 32 bytes               | **FAIL** | No validation; tests use `destBytesOverhead = 0`. |
| INV-POOL-21 | `destPoolData` exceeding `destBytesOverhead` → revert | **FAIL** | No size validation.                               |


### 5.7 Rate Limiting


| ID        | Invariant                                           | Grade    | Notes                                              |
| --------- | --------------------------------------------------- | -------- | -------------------------------------------------- |
| INV-RL-1  | Token bucket: `capacity`, `rate`, `tokens`          | **PASS** | `RateLimiter.daml:46-59`.                          |
| INV-RL-2  | Bucket refills based on elapsed time                | **PASS** | `availableTokensAt`.                               |
| INV-RL-3  | Requested exceeds available → revert                | **PASS** | Enforced.                                          |
| INV-RL-4  | Disabled or zero request → no consumption           | **PASS** | Checked.                                           |
| INV-RL-5  | Separate buckets for outbound and inbound per chain | **PASS** | `outboundRateLimiters` and `inboundRateLimiters`.  |
| INV-RL-6  | FTF may use additional rate limit buckets           | **FAIL** | `RateLimitMode_CustomFinality` defined but unused. |
| INV-RL-7  | Rate limiter config per-remote-chain                | **PASS** | Per-chain config.                                  |
| INV-RL-8  | Enabled: `rate <= capacity`                         | **PASS** | Validated.                                         |
| INV-RL-9  | Disabled: `rate == 0` and `capacity == 0`           | **PASS** | Enforced.                                          |
| INV-RL-10 | Config updated by owner or rate limit admin         | **PASS** | `poolOwner` controls updates.                      |


### 5.8 Pool Configuration


| ID         | Invariant                                                | Grade    | Notes                                               |
| ---------- | -------------------------------------------------------- | -------- | --------------------------------------------------- |
| INV-PCFG-1 | Add/remove remote chains; non-empty `remoteTokenAddress` | **PASS** | `UpdateChainPoolConfigs` in `LockReleaseTokenPool`. |
| INV-PCFG-2 | Remote pool addresses added/removed individually         | **PASS** | `remotePools : [BytesHex]`.                         |
| INV-PCFG-3 | Pool config updates restricted to owner                  | **PASS** | `controller poolOwner`.                             |


### 5.9 Source Pool Validation


| ID         | Invariant                                                     | Grade    | Notes                                          |
| ---------- | ------------------------------------------------------------- | -------- | ---------------------------------------------- |
| INV-PVAL-1 | `sourcePoolAddress` validated against configured remote pools | **PASS** | `sourcePoolAddress elem chainCfg.remotePools`. |


---

## 6. Fee Invariants (`FEE_INVARIANTS.md`)

### 6.1 Fee Structure


| ID        | Invariant                                     | Grade    | Notes                                                                                |
| --------- | --------------------------------------------- | -------- | ------------------------------------------------------------------------------------ |
| INV-FEE-1 | Total fee = CCV + pool + executor + network   | **PASS** | `FeeTokenAmount` sums all components.                                                |
| INV-FEE-2 | All fees in USD cents, converted to fee token | **PASS** | `FinalizeFee` converts via `feeTokenPrice`.                                          |
| INV-FEE-3 | Max USD cents per message cap                 | **FAIL** | No `maxUSDCentsPerMessage` enforcement.                                              |
| INV-FEE-4 | `getFee` and send path compute same fee       | **FAIL** | No separate `getFee` path; fees built incrementally. Consistency cannot be verified. |


### 6.2 Pool Fee


| ID        | Invariant                                                        | Grade    | Notes                                                                                              |
| --------- | ---------------------------------------------------------------- | -------- | -------------------------------------------------------------------------------------------------- |
| INV-FEE-5 | V2 pool with enabled fee config → use pool's `getFee`            | **PASS** | Pool has `CalculateFee` that reports fees.                                                         |
| INV-FEE-6 | V1 pool or disabled config → FeeQuoter fallback                  | **FAIL** | No V1/FeeQuoter fallback logic in OnRamp. Canton is V2-only but should handle disabled fee config. |
| INV-FEE-7 | FeeQuoter uses token-specific config if available, else defaults | **FAIL** | FeeQuoter exists but the OnRamp→pool→FeeQuoter fallback chain is not wired.                        |


### 6.3 Executor Fee


| ID         | Invariant                                          | Grade    | Notes                                                   |
| ---------- | -------------------------------------------------- | -------- | ------------------------------------------------------- |
| INV-FEE-8  | Executor fee = flat fee + gas-based execution cost | **FAIL** | Executor fee is flat only; no gas-based execution cost. |
| INV-FEE-9  | No-execution address → zero fee                    | **FAIL** | No-execution address not implemented.                   |
| INV-FEE-10 | Execution cost from gas estimates                  | **FAIL** | No gas-based execution cost.                            |


### 6.4 Network Fee


| ID         | Invariant                                      | Grade    | Notes                                             |
| ---------- | ---------------------------------------------- | -------- | ------------------------------------------------- |
| INV-FEE-11 | `tokenNetworkFeeUSDCents` for token msgs       | **PASS** | `OnRamp.daml:149-152` branches on token presence. |
| INV-FEE-12 | Exactly one network fee value used per message | **PASS** | One selected.                                     |


### 6.5 LINK Discount


| ID         | Invariant                                      | Grade    | Notes                                                                               |
| ---------- | ---------------------------------------------- | -------- | ----------------------------------------------------------------------------------- |
| INV-FEE-13 | LINK → `linkFeeMultiplierPercent` discount     | **FAIL** | `premiumMultiplier` per fee token exists but no LINK-specific multiplier mechanism. |
| INV-FEE-14 | LINK discount does NOT apply to execution cost | **FAIL** | No separate execution cost to exclude from discount.                                |


### 6.6 Price Requirements


| ID         | Invariant                                          | Grade    | Notes                                                          |
| ---------- | -------------------------------------------------- | -------- | -------------------------------------------------------------- |
| INV-FEE-15 | Fee token price must be non-zero                   | **PASS** | `feeTokenPrice > 0.0` enforced.                                |
| INV-FEE-16 | `usdPerUnitGas` must be set                        | **FAIL** | Gas price exists in FeeQuoter but not used in fee calculation. |
| INV-FEE-17 | Prices valid until overwritten; no staleness check | **PASS** | Stored with timestamp; no staleness check.                     |


### 6.7 Fee Distribution


| ID         | Invariant                                    | Grade    | Notes                                                                               |
| ---------- | -------------------------------------------- | -------- | ----------------------------------------------------------------------------------- |
| INV-FEE-18 | CCV fees transferred to CCV addresses        | **PASS** | CCV receipts built.                                                                 |
| INV-FEE-19 | Executor fee transferred to executor address | **PASS** | Executor receipt built.                                                             |
| INV-FEE-20 | Pool fee transferred to pool (V2 only)       | **PASS** | Pool fee in pool receipt.                                                           |
| INV-FEE-21 | Network fee stays on OnRamp                  | **FAIL** | Network fee in `networkReceipt`; transferred to CCIP owner, not retained by OnRamp. |


---

## 7. Message Lifecycle Invariants (`MESSAGE_LIFECYCLE_INVARIANTS.md`)

### 7.1 Message Sequencing


| ID        | Invariant                             | Grade    | Notes                                                   |
| --------- | ------------------------------------- | -------- | ------------------------------------------------------- |
| INV-SEQ-1 | `messageNumber` per-lane              | **PASS** | `outboundSequenceNumbers` keyed by `destChainSelector`. |
| INV-SEQ-2 | Strictly monotonic, first message = 1 | **PASS** | `newSequenceNumber = currentSequenceNumber + 1.0`.      |
| INV-SEQ-3 | `messageNumber = 0` reserved          | **PASS** | First message gets 1.                                   |
| INV-SEQ-4 | Persists across lane config updates   | **PASS** | Stored in router; unaffected by `GlobalConfig`.         |


### 7.2 Message Identity


| ID       | Invariant                                 | Grade    | Notes                                   |
| -------- | ----------------------------------------- | -------- | --------------------------------------- |
| INV-ID-1 | `messageId = keccak256(encodedMessageV1)` | **PASS** | `MessageCodecV1.messageId`.             |
| INV-ID-2 | Execution outcomes keyed by `messageId`   | **PASS** | `executionStates` keyed by `messageId`. |


### 7.3 Source Side Flow


| ID        | Invariant                               | Grade    | Notes                                                                               |
| --------- | --------------------------------------- | -------- | ----------------------------------------------------------------------------------- |
| INV-SRC-1 | Fees computed before `lockOrBurn`       | **FAIL** | `LockOrBurn` happens before `CCIPSend` distributes fees. Fees collected after lock. |
| INV-SRC-2 | `messageId` computed after `lockOrBurn` | **PASS** | Correct; built in `AddExecutorWithFee` after `LockOrBurn`.                          |
| INV-SRC-3 | Fee must not exceed provided amount     | **PASS** | Transfer fails if insufficient.                                                     |
| INV-SRC-4 | Token amount must be non-zero           | **FAIL** | No explicit non-zero check for token amount.                                        |


### 7.4 Destination Side Flow


| ID        | Invariant                             | Grade    | Notes                                                                 |
| --------- | ------------------------------------- | -------- | --------------------------------------------------------------------- |
| INV-DST-1 | CCV verification before token release | **PASS** | CCV verification via `AddCCVVerification` → then `ReleaseFromTicket`. |
| INV-DST-2 | `releaseOrMint` before `ccipReceive`  | **PASS** | Execute → ReleaseFromTicket → CCIPMessageReceived.                    |
| INV-DST-3 | Token-only: `ccipReceive` not called  | **PASS** | Token-only skips receiver callback.                                   |


### 7.5 Execution Semantics


| ID         | Invariant                                              | Grade    | Notes                                                                                                        |
| ---------- | ------------------------------------------------------ | -------- | ------------------------------------------------------------------------------------------------------------ |
| INV-EXEC-1 | Execution is permissionless                            | **FAIL** | `controller partyOwner` — requires specific party.                                                           |
| INV-EXEC-2 | No in-protocol ordering                                | **PASS** | Execution by `messageId`; no ordering.                                                                       |
| INV-EXEC-3 | `messageNumber` for identification only                | **PASS** | Not used for ordering.                                                                                       |
| INV-EXEC-4 | Never-executed message may be executed                 | **PASS** | `UNTOUCHED` messages can proceed.                                                                            |
| INV-EXEC-5 | Failed message may be retried                          | **FAIL** | Canton does not persist failure state; failed executions abort the Daml transaction entirely. No retry path. |
| INV-EXEC-6 | Successfully executed message cannot be re-executed    | **PASS** | SUCCESS messages abort on re-execution.                                                                      |
| INV-EXEC-7 | Retry that still fails → no redundant state transition | **FAIL** | No failure state persistence means no retry mechanism exists.                                                |


### 7.6 Receiver Callback


| ID         | Invariant                                                            | Grade    | Notes                                                                                                                                                                             |
| ---------- | -------------------------------------------------------------------- | -------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| INV-RECV-1 | `ccipReceive` called on `message.receiver`; only configured OffRamps | **FAIL** | No `ccipReceive` callback; Canton uses `CCIPMessageReceived` template creation. The receiver can add custom logic before `create CCIPMessageReceived` but this is not a callback. |
| INV-RECV-2 | `ccipReceive` not called for token-only                              | **PASS** | Token-only skips receiver callback/message creation flow.                                                                                                                         |


### 7.7 Token Receiver vs Message Receiver


| ID       | Invariant                                          | Grade    | Notes                                            |
| -------- | -------------------------------------------------- | -------- | ------------------------------------------------ |
| INV-TR-1 | `message.receiver` gets callback                   | **PASS** | `CCIPReceiver.owner` receives the message.       |
| INV-TR-2 | `tokenReceiver` may differ from `message.receiver` | **PASS** | `tokenReceiver` in `TokenTransferV1` can differ. |
| INV-TR-3 | Empty `tokenReceiver` → use `message.receiver`     | **PASS** | Fallback in `OffRamp.daml:102-106`.              |
| INV-TR-4 | `tokenReceiverAllowed` per-lane config flag        | **FAIL** | No `tokenReceiverAllowed` config.                |


### 7.8 Token-Only Transfer Behavior


| ID       | Invariant                                                | Grade    | Notes                                                                                                                                   |
| -------- | -------------------------------------------------------- | -------- | --------------------------------------------------------------------------------------------------------------------------------------- |
| INV-TO-1 | Token-only: `gasLimit == 0`, `data` empty, token present | **PASS** | `isTokenOnlyTransfer` in `OffRamp.daml:229-232` matches.                                                                                |
| INV-TO-2 | Token-only: `ccipReceive` not called                     | **PASS** | Correct.                                                                                                                                |
| INV-TO-3 | CCV verification and `releaseOrMint` still apply         | **PASS** | Correct.                                                                                                                                |
| INV-TO-4 | Token-only: receiver/default CCVs excluded               | **FAIL** | Logic exists in OffRamp but source side should also exclude defaults for token-only; no zero-value sentinel to explicitly request them. |
| INV-TO-5 | Receiver finality not checked for token-only             | **FAIL** | No receiver finality checks exist at all (see INV-FIN-DST-1).                                                                           |


### 7.9 No-Execution Address


| ID           | Invariant                                             | Grade    | Notes                                                   |
| ------------ | ----------------------------------------------------- | -------- | ------------------------------------------------------- |
| INV-NOEXEC-1 | No-execution address = `address(bytes20(0xeba517d2))` | **FAIL** | Not implemented. Protocol-wide requirement per Context. |
| INV-NOEXEC-2 | No-execution → zero fee, no execution cost            | **FAIL** | Same.                                                   |
| INV-NOEXEC-3 | Tokens still released for no-execution messages       | **FAIL** | Same.                                                   |


