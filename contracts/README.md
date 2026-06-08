# CCIP Canton Contracts

This document describes the Daml smart contract architecture for Chainlink CCIP (Cross-Chain Interoperability Protocol)
on Canton.

## Contract artifacts & releases

DARs live under `contracts/dars/` (`current/` for dev, `v1_0_0/` etc. for frozen releases). Go bindings live under `bindings/generated/latest/` for all in-repo code; frozen releases snapshot DARs only (not versioned binding trees).

For day-to-day builds (`make contracts`) and **how to cut and migrate to a new release** (e.g. 1.1.0 — freeze, `ReleaseDir`, import updates, git tag), see **[bindings/README.md](../bindings/README.md)**.

## Overview

CCIP Canton enables cross-chain messaging between Canton and other blockchain networks (EVM, Solana, Aptos, SUI, TVM).
Unlike EVM-based CCIP which supports arbitrary execution via `ccipReceive()` callbacks, Canton CCIP provides **arbitrary
messaging** - users receive verified message data and decide how to process it.

### Key Differences from EVM CCIP

| Aspect              | EVM CCIP                          | Canton CCIP                               |
|---------------------|-----------------------------------|-------------------------------------------|
| Execution model     | Arbitrary execution via callbacks | Arbitrary messaging (user processes data) |
| Receiver            | Contract address                  | PartyId encoded as bytes                  |
| Message retrieval   | Automatic callback                | User fetches from off-chain storage       |
| State model         | Global shared state               | Per-party isolated state                  |
| Contract visibility | Public                            | Explicit disclosure required              |

## Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         CONTRACT ARCHITECTURE                                │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  SHARED CONTRACTS (fetched via explicit disclosure)                          │
│  ══════════════════════════════════════════════════                          │
│                                                                              │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │ GlobalConfig │  │   OnRamp     │  │   OffRamp    │  │  CCVRegistry │     │
│  │              │  │              │  │              │  │              │     │
│  │ chain config │  │ send logic   │  │ receive logic│  │ ticket issuer│     │
│  │ lane configs │  │ msg encoding │  │ verification │  │              │     │
│  └──────────────┘  └──────────────┘  └──────────────┘  └──────────────┘     │
│                                                                              │
│  ┌──────────────────────┐  ┌──────────────────────┐                         │
│  │ TokenAdminRegistry   │  │  CommitteeVerifier   │                         │
│  │                      │  │  (CCV)               │                         │
│  │ token pool registry  │  │ signature validation │                         │
│  │ ticket authority     │  │                      │                         │
│  └──────────────────────┘  └──────────────────────┘                         │
│                                                                              │
│  PER-PARTY CONTRACTS (user is stakeholder)                                   │
│  ═════════════════════════════════════════                                   │
│                                                                              │
│  ┌────────────────────────────────────────────────────────────────────┐     │
│  │ PerPartyRouter (one per user)                                       │     │
│  │                                                                     │     │
│  │ - outboundSequenceNumbers : Map destChain -> seqNum                │     │
│  │ - executionStates : Map messageHash -> state                       │     │
│  │ - receiverRequiredCCVs : [CCVId]                                   │     │
│  └────────────────────────────────────────────────────────────────────┘     │
│                                                                              │
│  TICKETS (ephemeral, created and consumed during flows)                      │
│  ══════════════════════════════════════════════════════                      │
│                                                                              │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │TokenSendTicket│ │  CCVTicket   │  │CCVVerifyTicket│ │TokenReceive- │     │
│  │              │  │              │  │              │  │    Ticket    │     │
│  │ outbound     │  │ outbound     │  │ inbound      │  │ inbound      │     │
│  │ token lock   │  │ attestation  │  │ verification │  │ token release│     │
│  └──────────────┘  └──────────────┘  └──────────────┘  └──────────────┘     │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

## Package Structure

```
contracts/ccip/
├── common/              # Shared types and utilities
│   └── daml/CCIP/
│       ├── Common.daml           # Chain family selectors
│       ├── Internal.daml         # MessageExecutionState enum
│       ├── Client.daml           # Message types for users
│       ├── CCVId.daml            # CCV identifier utilities
│       ├── GlobalConfig.daml     # Lane configuration
│       ├── MessageCodecV1.daml   # Cross-chain message encoding
│       ├── Tickets.daml          # All ticket templates
│       └── Interfaces/
│           ├── CrossChainVerifier.daml
│           └── Any2CantonMessageReceiver.daml
├── perpartyrouter/      # Per-party routing/state contract
├── onramp/              # Outbound message handling
├── offramp/             # Inbound message handling
├── ccipsender/          # Sender entrypoint that calls verifiers + OnRamp
├── ccvs/                # CommitteeVerifier implementation
├── tokenAdminRegistry/  # Token pool registry and ticket authority
├── feequoter/           # Fee calculation
├── ccipreceiver/        # Example receiver implementation
└── test/                # Test contracts and mocks
```

## Core Contracts

### PerPartyRouter

Per-party routing/state contract used by sender/receiver flows. Each party gets their own router to avoid state
contention.

```daml
template PerPartyRouter
    with
        ccipOwner : Party
        partyOwner : Party
        environmentId : Text
        outboundSequenceNumbers : Map (Numeric 0) (Numeric 0)  -- destChain -> seqNum
        executionStates : Map BytesHex MessageExecutionState   -- messageHash -> state
        receiverRequiredCCVs : [CCVId]                         -- CCVs required for inbound
```

**Key choices:**

- `CCIPSend` - Send a cross-chain message
- `Execute` - Execute an inbound message
- `GetRequiredCCVsForSend` - Query required CCVs for sending
- `GetRequiredCCVsForExecute` - Query required CCVs for receiving

### OnRamp

Contains business logic for outbound messages. Called by PerPartyRouter.

```daml
template OnRamp
    with
        ccipOwner : Party
        environmentId : Text
```

**Key choice:**

- `CCIPSendFromRouter` - Validates tickets, encodes message, returns encoded result

### OffRamp

Contains business logic for inbound messages. Called by PerPartyRouter.

```daml
template OffRamp
    with
        ccipOwner : Party
        environmentId : Text
```

**Key choice:**

- `ExecuteFromRouter` - Decodes message, validates CCVs, issues TokenReceiveTicket

### GlobalConfig

Shared configuration for lane settings.

```daml
template GlobalConfig
    with
        ccipOwner : Party
        environmentId : Text
        chainSelector : Numeric 0              -- This chain's selector
        onRampAddress : BytesHex               -- This chain's OnRamp address
        destChainConfigs : Map (Numeric 0) DestChainConfig
        sourceChainConfigs : Map (Numeric 0) SourceChainConfig

data DestChainConfig = DestChainConfig
    with
        isEnabled : Bool
        offRampAddress : BytesHex
        laneMandatedCCVs : [CCVId]             -- Required for this lane
        defaultCCVs : [CCVId]                  -- Fallback if none specified

data SourceChainConfig = SourceChainConfig
    with
        isEnabled : Bool
        onRampAddress : BytesHex
        laneMandatedCCVs : [CCVId]
        defaultCCVs : [CCVId]
```

### TokenAdminRegistry

Central authority for token-related tickets. Manages which pools are authorized for which tokens.

```daml
template TokenAdminRegistry
    with
        owner : Party
        tokenConfigs : Map InstrumentId TokenConfig

data TokenConfig = TokenConfig
    with
        admin : Optional Party
        pendingAdmin : Optional Party
        tokenPoolOwner : Optional Party
        requiredCCVs : [CCVId]                 -- Pool's required CCVs
```

**Key choices:**

- `SetPool` - Registers or clears the pool for an instrument
- `ProposeAdministrator` / `AcceptAdminRole` / `TransferAdminRole` - Two-step token admin management
- `ConsumeReceiveTicket` - Called by pool to release tokens against an inbound ticket
- `SetOutboundPoolCCVs` / `SetInboundPoolCCVs` - Records pool-required CCVs on outbound and inbound message state
- `AddTokenSendFee` / `AddTokenSend` / `FinalizeExecute` - Finalizes token-side message accounting and execution state

### CommitteeVerifier

A CCV implementation using committee-based ECDSA signature verification.

```daml
template CommitteeVerifier
    with
        owner : Party
        ccipOwner : Party
        versionTag : BytesHex          -- e.g., "e9a05a20"
        threshold : Int                -- Minimum signatures required
        signers : [BytesHex]           -- Authorized signer public keys
```

## Ticket Pattern

Canton's privacy model means CCIP infrastructure cannot directly call user contracts. The **ticket pattern** solves
this:

1. User calls their contract (e.g., TokenPool) to perform an action
2. Contract issues a ticket proving the action occurred
3. User passes ticket to CCIP contracts
4. CCIP validates and consumes the ticket

### Ticket Types

| Ticket               | Issued By                        | Consumed By                  | Purpose                                       |
|----------------------|----------------------------------|------------------------------|-----------------------------------------------|
| `TokenSendTicket`    | TokenPool via TokenAdminRegistry | OnRamp                       | Proves tokens were locked                     |
| `VerifierData`       | CCV (`ForwardToVerifier`)        | OnRamp / message pipeline    | Attestation data attached to outbound message |
| `CCVVerification`    | CCV (`VerifyMessage`)            | OffRamp / execution pipeline | Proves inbound message verification           |
| `TokenReceiveTicket` | OffRamp via TokenAdminRegistry   | TokenPool                    | Authorizes token release                      |

## CCV Identification

Cross-Chain Verifiers are identified by a `CCVId` combining:

- **Version tag** (4 bytes hex): Identifies the CCV type/version
- **Party**: Identifies the CCV operator

```
CCVId = "<versionTag>@<partyId>"
Example: "e9a05a20@participant1::ccv-owner"
```

The version tag is embedded in the first 4 bytes of `verifierBlob`/`ccvData`, preventing CCVs from lying about their
type.

## Message Encoding

Messages are encoded to match EVM's `MessageV1Codec.sol` exactly for cross-chain `messageId` compatibility:

```
messageId = keccak256(encodedMessage)
```

### MessageV1 Structure

```daml
data MessageV1 = MessageV1
    with
        -- Static section (69 bytes)
        sourceChainSelector : Numeric 0   -- uint64
        destChainSelector : Numeric 0     -- uint64
        sequenceNumber : Numeric 0        -- uint64
        executionGasLimit : Int           -- uint32
        ccipReceiveGasLimit : Int         -- uint32
        finality : BytesHex               -- bytes4 finality config
        ccvAndExecutorHash : BytesHex     -- bytes32
        -- Variable section
        onRampAddress : BytesHex
        offRampAddress : BytesHex
        sender : BytesHex                 -- PartyId encoded
        receiver : BytesHex               -- PartyId encoded
        destBlob : BytesHex
        tokenTransfer : Optional TokenTransferV1
        messageData : BytesHex
```

## Explicit Disclosure

Canton participants maintain private contract stores. Users need **explicit disclosure** to access shared contracts
they're not stakeholders of.

### How It Works

1. CCIP operator queries contracts with `IncludeCreatedEventBlob: true`
2. Contract blobs are distributed via Web2 (disclosure server)
3. Users attach `DisclosedContract` to their command submissions
4. Participants validate disclosed contracts against their hash

### Contracts Needed by Operation

**CCIPSend:**

- `PerPartyRouter` (stakeholder)
- `OnRamp` (disclosed)
- `GlobalConfig` (disclosed)
- `TokenAdminRegistry` (disclosed)
- User's tickets (stakeholder)

**Execute:**

- `PerPartyRouter` (stakeholder)
- `OffRamp` (disclosed)
- `GlobalConfig` (disclosed)
- `TokenAdminRegistry` (disclosed)
- User's tickets (stakeholder)

---

## Send Flow (Outbound)

Sending a cross-chain message from Canton to another chain.

Note: some legacy diagrams below still use `CCVTicket` naming. Current verifier integration is direct via verifier
interface choices.

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              SEND FLOW                                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  1. PREPARE TOKENS (if token transfer)                                      │
│  ═══════════════════════════════════════                                    │
│                                                                             │
│     User ──► TokenPool.Lock()                                               │
│                   │                                                         │
│                   ▼                                                         │
│              TokenAdminRegistry_IssueSendTicket                             │
│                   │                                                         │
│                   ▼                                                         │
│              TokenSendTicket (created)                                      │
│                                                                             │
│  2. GET CCV ATTESTATION                                                     │
│  ══════════════════════                                                     │
│                                                                             │
│     User ──► ForwardToVerifier                            │
│                   │                                                         │
│                   ▼                                                         │
│              VerifierData appended to SendingMessageV1                      │
│                                                                             │
│  3. SEND MESSAGE                                                            │
│  ═══════════════                                                            │
│                                                                             │
│     User ──► CCIPSender.Send(                                               │
│                  context,                                                   │
│                  routerCid,                                                 │
│                  destChainSelector,                                         │
│                  receiver,                                                  │
│                  payload,                                                   │
│                  tokenTransfer,                                             │
│                  ccvSendInputs                                              │
│              )                                                              │
│                   │                                                         │
│                   ▼                                                         │
│     ┌─────────────────────────────────────────────────────────────┐         │
│     │ CCIPSender                                                  │         │
│     │   1. Prepare send via PerPartyRouter                        │         │
│     │   2. Calculate/finalize fee                                 │         │
│     │   3. Forward to verifier interfaces                         │         │
│     │   4. Finalize send via PerPartyRouter.CCIPSend              │         │
│     └─────────────────────────────────────────────────────────────┘         │
│                   │                                                         │
│                   ▼                                                         │
│     ┌─────────────────────────────────────────────────────────────┐         │
│     │ OnRamp.CCIPSendFromRouter                                   │         │
│     │   1. Validate router environment matches                    │         │
│     │   2. Validate GlobalConfig                                  │         │
│     │   3. Get dest chain config                                  │         │
│     │   4. Process TokenSendTicket → TokenTransferV1              │         │
│     │   5. Get required CCVs (lane + pool + defaults)             │         │
│     │   6. Validate required CCV inputs                           │         │
│     │   7. Build final message                                    │         │
│     │   8. Build MessageV1                                        │         │
│     │   9. Encode message                                         │         │
│     │  10. Compute messageId = keccak256(encodedMessage)          │         │
│     └─────────────────────────────────────────────────────────────┘         │
│                   │                                                         │
│                   ▼                                                         │
│     ┌─────────────────────────────────────────────────────────────┐         │
│     │ PerPartyRouter (continued)                                  │         │
│     │   4. Update outboundSequenceNumbers                         │         │
│     │   5. Create CCIPMessageSent event contract                  │         │
│     │   6. Return CCIPSendResult                                  │         │
│     └─────────────────────────────────────────────────────────────┘         │
│                                                                             │
│  4. OUTPUT                                                                  │
│  ═════════                                                                  │
│                                                                             │
│     CCIPMessageSent contract contains:                                      │
│       - destChainSelector                                                   │
│       - sequenceNumber                                                      │
│       - messageId                                                           │
│       - encodedMessage                                                      │
│       - verifierBlobs (for off-chain DON to pick up)                        │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Send Flow Code Example

```daml
-- Step 1: Lock tokens and get TokenSendTicket (if transferring tokens)
tokenSendTicket <- exercise tokenPoolCid TokenPool_LockOrBurn with
    sender = userParty
    amount = 100.0
    destTokenAddress = destTokenAddr
    tokenReceiver = receiverPartyBytes
    ...

-- Step 2/3: Send the message
result <- exercise senderCid Send with
    context = context
    routerCid = routerCid
    destChainSelector = 1.0  -- e.g., Ethereum mainnet
    receiver = receiverPartyBytes
    payload = userPayload
    tokenTransfer = Some tokenTransferInput
    ccvSendInputs = ccvSendInputs
    ...
```

---

## Receive Flow (Inbound)

Executing an inbound cross-chain message on Canton.

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                             RECEIVE FLOW                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  1. FETCH MESSAGE (off-chain)                                               │
│  ════════════════════════════                                               │
│                                                                             │
│     User fetches from off-chain storage:                                    │
│       - encodedMessage                                                      │
│       - verifierBlobs (ccvData with signatures)                             │
│                                                                             │
│  2. VERIFY MESSAGE WITH CCVs                                                │
│  ═══════════════════════════                                                │
│                                                                             │
│     User provides ccvInputs to CCIPReceiver.Execute                         │
│     (receiver calls CrossChainVerifier_VerifyMessage internally)            │
│                                                                             │
│     CrossChainVerifier_VerifyMessage(                                       │
│                  executingMessageCid,                                       │
│                  verifierResults,                                           │
│                  caller                                                     │
│              )                                                              │
│                   │                                                         │
│                   ▼                                                         │
│     ┌─────────────────────────────────────────────────────────────┐         │
│     │ CommitteeVerifier                                           │         │
│     │   1. Verify version tag matches                             │         │
│     │   2. Parse signatures from ccvData                          │         │
│     │   3. Compute signedHash = keccak256(versionTag || encoded)  │         │
│     │   4. Verify >= threshold valid signatures                   │         │
│     │   5. Append verification to ExecutingMessageV1              │         │
│     └─────────────────────────────────────────────────────────────┘         │
│                   │                                                         │
│                   ▼                                                         │
│              ExecutingMessageV1 updated                                     │
│                                                                             │
│  3. EXECUTE MESSAGE                                                         │
│  ══════════════════                                                         │
│                                                                             │
│     User ──► CCIPReceiver.Execute(                                          │
│                  context,                                                   │
│                  routerCid,                                                 │
│                  encodedMessage,                                            │
│                  tokenTransfer,                                             │
│                  ccvInputs,                                                 │
│              )                                                              │
│                   │                                                         │
│                   ▼                                                         │
│     ┌─────────────────────────────────────────────────────────────┐         │
│     │ PerPartyRouter                                              │         │
│     │   1. Validate OffRamp (ccipOwner, environmentId)            │         │
│     │   2. Compute messageHash = keccak256(encodedMessage)        │         │
│     │   3. Check execution state (replay prevention)              │         │
│     │   4. Delegate to OffRamp.ExecuteFromRouter                  │         │
│     └─────────────────────────────────────────────────────────────┘         │
│                   │                                                         │
│                   ▼                                                         │
│     ┌─────────────────────────────────────────────────────────────┐         │
│     │ OffRamp.ExecuteFromRouter                                   │         │
│     │   1. Validate router environment matches                    │         │
│     │   2. Decode MessageV1 from encodedMessage                   │         │
│     │   3. Verify receiver matches router party                   │         │
│     │   4. Validate GlobalConfig                                  │         │
│     │   5. Get source chain config                                │         │
│     │   6. Get required CCVs (lane + receiver + pool + defaults)  │         │
│     │   7. Validate CCV verifications                             │         │
│     │   8. Validate all required CCVs                             │         │
│     │   8. Validate all optional CCVs meet optional threshold     │         │
│     │  10. Issue TokenReceiveTicket via TokenAdminRegistry        │         │
│     └─────────────────────────────────────────────────────────────┘         │
│                   │                                                         │
│                   ▼                                                         │
│     ┌─────────────────────────────────────────────────────────────┐         │
│     │ PerPartyRouter (continued)                                  │         │
│     │   5. Update executionStates[messageHash] = SUCCESS          │         │
│     │   6. Create ExecutionStateChanged event contract            │         │
│     │   7. Return ExecuteResult with TokenReceiveTicket           │         │
│     └─────────────────────────────────────────────────────────────┘         │
│                                                                             │
│  4. RELEASE TOKENS (if token transfer)                                      │
│  ═════════════════════════════════════                                      │
│                                                                             │
│     User ──► TokenPool.ReleaseFromTicket(tokenReceiveTicket)                │
│                   │                                                         │
│                   ▼                                                         │
│              ConsumeReceiveTicket                        │
│                   │                                                         │
│                   ▼                                                         │
│              Tokens released to receiver                                    │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Receive Flow Code Example

```daml
-- Step 1: Fetch message from off-chain (pseudo-code)
-- encodedMessage, ccvDataList <- fetchFromStorage(messageId)

-- Step 2/3: Execute the message
result <- exercise receiverCid Execute with
    context = context
    routerCid = routerCid
    encodedMessage = encodedMessage
    ccvInputs = ccvInputs
    ...

-- Step 4: Release tokens (if message had token transfer)
case result.tokenReceiveTicket of
    Some ticketCid -> do
        exercise tokenPoolCid TokenPool_ReleaseFromTicket with
            ticketCid = ticketCid
            ...
    None -> pure ()

-- Step 5: Process message data as needed
let messageData = result.message.messageData
-- Application-specific logic here
```

---

## Required CCVs

CCVs (Cross-Chain Verifiers) are aggregated from multiple sources:

### For Sending (OnRamp)

```
Required CCVs = lane-mandated + pool + defaults
```

1. **Lane-mandated**: `GlobalConfig.destChainConfigs[dest].laneMandatedCCVs`
2. **Pool CCVs**: `TokenAdminRegistry.tokenConfigs[instrument].requiredCCVs` (if token transfer)
3. **Defaults**: `GlobalConfig.destChainConfigs[dest].defaultCCVs` (if no others specified)

Note: Sender does NOT know receiver's requirements on the destination chain.

### For Receiving (OffRamp)

```
Required CCVs = lane-mandated + receiver + pool + defaults
```

1. **Lane-mandated**: `GlobalConfig.sourceChainConfigs[source].laneMandatedCCVs`
2. **Receiver CCVs**: `PerPartyRouter.receiverRequiredCCVs`
3. **Pool CCVs**: `TokenAdminRegistry.tokenConfigs[instrument].requiredCCVs` (if token transfer)
4. **Defaults**: `GlobalConfig.sourceChainConfigs[source].defaultCCVs` (if no others specified)

---

## Bi-directional Validation

Security is enforced through mutual validation:

```
┌───────────────────┐                    ┌───────────────────┐
│ PerPartyRouter    │                    │ OnRamp/OffRamp    │
│                   │                    │                   │
│ 1. Fetch ramp     │───────────────────►│                   │
│ 2. Validate:      │                    │                   │
│    - ccipOwner    │                    │                   │
│    - environmentId│                    │                   │
│                   │                    │                   │
│ 3. Exercise       │───────────────────►│ 4. Validate:      │
│    choice         │                    │    - environmentId│
│                   │◄───────────────────│ 5. Return result  │
│ 6. Update state   │                    │                   │
└───────────────────┘                    └───────────────────┘
```

---

## Events

CCIP emits event contracts that can be observed by off-chain systems:

### CCIPMessageSent

Created after successful send. Observed by sender and DON nodes.

```daml
template CCIPMessageSent
    with
        ccipOwner : Party
        ccvOwners : [Party]          -- CCV operators that this message has been forwarded to
        sender : Party
        observers : [Party]          -- Includes (optional) additional observers that CCVs have specified
        event : CCIPMessageSentEvent

data CCIPMessageSentEvent = CCIPMessageSentEvent
    with
        destChainSelector : Numeric 0
        sequenceNumber : Numeric 0
        messageId : BytesHex
        encodedMessage : BytesHex
        verifierBlobs : [BytesHex]   -- DON picks these up
        receipts : [CCIP.Tickets.Receipt]
```

### ExecutionStateChanged

Created after successful execute. Observed by receiver.

```daml
template ExecutionStateChanged
    with
        ccipOwner : Party
        ccvOwners : [Party]          -- CCV operators that this message has been verified by
        receiver : Party
        event : ExecutionStateChangedEvent

data ExecutionStateChangedEvent = ExecutionStateChangedEvent
    with
        sourceChainSelector : Numeric 0
        sequenceNumber : Numeric 0
        messageId : BytesHex
        state : MessageExecutionState
        returnData : BytesHex
```

---

## Building

Build all Daml contracts:

```bash
cd contracts && dpm build --all
```

Build a single package:

```bash
cd contracts/ccip/perpartyrouter && dpm build
```

Clean all build artifacts:

```bash
cd contracts && dpm clean --all
```

Built DARs are output to `.daml/dist/` in each package directory.
