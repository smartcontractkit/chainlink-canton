# Canton CCIP Invariants — Report Card

**Audit Date:** 2026-03-10
**Codebase:** `/contracts/ccip` (Daml / Canton)
**Classification per:** `invariants/CONTEXT.md` — missing features = FAIL, partial = FAIL, N/A only for truly EVM/SVM/Sui-specific details.

---

## Overall Summary

| Category                     | PASS | FAIL | N/A | Total |
|------------------------------|------|------|-----|-------|
| CCV Configuration            |   0  |   7  |  2  |    9  |
| CCV Source Side               |  10  |   8  |  4  |   22  |
| CCV Destination Side          |   8  |   8  |  7  |   23  |
| CCV Address Stability         |   1  |   2  |  4  |    7  |
| CCV Cross-Cutting             |   3  |   1  |  0  |    4  |
| Finality Encoding              |   2  |   0  |  1  |    3  |
| Finality Source Side           |   5  |   2  |  2  |    9  |
| Finality Destination Side      |   0  |   7  |  3  |   10  |
| Finality Pool                  |   2  |   7  |  6  |   15  |
| Finality Executor              |   4  |   0  |  0  |    4  |
| Finality CCV                  |   0  |   2  |  0  |    2  |
| Finality Re-org / Safety       |   2  |   0  |  0  |    2  |
| Finality Opt-in                |   0  |   2  |  1  |    3  |
| Finality Cross-Cutting         |   1  |   1  |  0  |    2  |
| Encoding MessageV1             |   8  |   2  |  1  |   11  |
| Encoding ExtraArgs             |   0  |   1  |  8  |    9  |
| Encoding Pool Remote           |   2  |   2  |  0  |    4  |
| Lane Config Source             |   2  |   2  |  3  |    7  |
| Lane Config Dest               |   3  |   1  |  0  |    4  |
| Lane Pause / Disable           |   3  |   0  |  0  |    3  |
| Lane Inflight                  |   3  |   0  |  0  |    3  |
| Lane Router                   |   0  |   0  |  5  |    5  |
| Lane OnRamp Validation         |   1  |   2  |  0  |    3  |
| Lane RMN / Curse               |   5  |   0  |  0  |    5  |
| Lane Reentrancy                |   2  |   0  |  0  |    2  |
| Lane Upgradability             |   4  |   0  |  0  |    4  |
| Lane Access Control            |   4  |   2  |  0  |    6  |
| Token Pool Types               |   2  |   1  |  0  |    3  |
| Token Pool Access Control      |   1  |   2  |  0  |    3  |
| Token Pool V1/V2               |   0  |   3  |  0  |    3  |
| Token Pool Fee Deduction       |   0  |   4  |  0  |    4  |
| Token Pool Decimals            |   2  |   3  |  0  |    5  |
| Token Pool destBytesOverhead   |   1  |   2  |  0  |    3  |
| Token Pool Rate Limiting       |   8  |   1  |  1  |   10  |
| Token Pool Config              |   3  |   0  |  0  |    3  |
| Token Pool Source Validation   |   1  |   0  |  0  |    1  |
| Fee Structure                  |   2  |   2  |  0  |    4  |
| Fee Pool                       |   1  |   2  |  0  |    3  |
| Fee Executor                   |   0  |   3  |  0  |    3  |
| Fee Network                    |   2  |   0  |  0  |    2  |
| Fee LINK Discount              |   0  |   2  |  0  |    2  |
| Fee Prices                     |   2  |   1  |  0  |    3  |
| Fee Distribution               |   3  |   1  |  0  |    4  |
| Msg Lifecycle Sequencing       |   4  |   0  |  0  |    4  |
| Msg Lifecycle Identity         |   2  |   0  |  0  |    2  |
| Msg Lifecycle Source Flow      |   2  |   1  |  1  |    4  |
| Msg Lifecycle Dest Flow        |   3  |   0  |  0  |    3  |
| Msg Lifecycle Execution        |   4  |   3  |  0  |    7  |
| Msg Lifecycle Receiver         |   1  |   1  |  0  |    2  |
| Msg Lifecycle Token Recv       |   3  |   1  |  0  |    4  |
| Msg Lifecycle Token-Only       |   3  |   2  |  0  |    5  |
| Msg Lifecycle No-Exec          |   0  |   3  |  0  |    3  |
| **TOTALS**                     | **124** | **89** | **48** | **261** |

**Pass rate (excl. N/A): 124 / 213 = 58.2%**

---

## Grade by Invariant File

| File | PASS | FAIL | N/A | Pass Rate (excl. N/A) |
|------|------|------|-----|-----------------------|
| CCV_INVARIANTS.md | 22 | 26 | 17 | 45.8% |
| FINALITY_INVARIANTS.md | 16 | 21 | 13 | 43.2% |
| ENCODING_INVARIANTS.md | 10 | 5 | 9 | 66.7% |
| LANE_INVARIANTS.md | 27 | 7 | 8 | 79.4% |
| TOKEN_POOL_INVARIANTS.md | 18 | 16 | 1 | 52.9% |
| FEE_INVARIANTS.md | 10 | 11 | 0 | 47.6% |
| MESSAGE_LIFECYCLE_INVARIANTS.md | 22 | 11 | 0 | 66.7% |

---

## Key Gaps (Prioritized)

### Critical (Security / Correctness)

| # | Invariants | Issue |
|---|-----------|-------|
| 1 | INV-CFG-1..7 | No CCV config validation: empty sets allowed, no duplicate/sentinel checks |
| 2 | INV-CC-1 | No guarantee every message is verified by at least one CCV |
| 3 | INV-ONRVAL-2/3 | OffRamp does not validate `offRampAddress` or `destChainSelector` |
| 4 | INV-POOL-4/5 | Pool does not verify caller is OnRamp/OffRamp |
| 5 | INV-ENC-1 | Uses `CantonExtraArgsV1` instead of `GenericExtraArgsV3` wire format |
| 6 | INV-MSG-2/3 | Decoder does not validate version byte or reject trailing bytes |

### High (Missing V2 Features)

| # | Invariants | Issue |
|---|-----------|-------|
| 7 | INV-POOL-7..13 | Missing V2 pool features: `feeBps`, finality-aware fees |
| 8 | INV-FIN-DST-* | Destination-side finality not implemented: no `minBlockDepth`, no pool finality on inbound |
| 9 | INV-NOEXEC-1..3 | No-execution address not implemented |
| 10 | INV-FEE-3 | No max USD cents per message cap |
| 11 | INV-SRC-2, INV-DST-7 | Zero-value address sentinel for "use defaults" not supported |

### Medium (Feature Completeness)

| # | Invariants | Issue |
|---|-----------|-------|
| 12 | INV-FEE-8..10 | No gas-based execution cost component |
| 13 | INV-FEE-13/14 | No LINK-specific discount multiplier |
| 14 | INV-EXEC-1 | Execution requires `controller partyOwner`, not permissionless |
| 15 | INV-LCFG-2/9 | No chain selector validation (non-zero, not self) |
| 16 | INV-RECV-1 | No `ccipReceive` callback; uses ticket model |
| 17 | INV-DST-3..5 | Receiver CCV duplicate validation and optional CCV quorum missing |
