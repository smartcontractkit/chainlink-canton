# Canton Party Ceremony

A CLI tool for governing and managing **decentralized parties** on Canton distributed ledgers.

## Table of Contents

- [Why Party-Ceremony](#why-party-ceremony)
- [Architecture](#architecture)
  - [Operations Framework](#operations-framework)
  - [Async Multi-Actor Model](#async-multi-actor-model)
  - [State Persistence](#state-persistence)
  - [Coordinator Role](#coordinator-role)
  - [Exit Codes](#exit-codes)
- [Production Workflow](#production-workflow)
- [Workflows](#workflows)
  - [Onboarding](#onboarding)
  - [Kick (Remove Participant)](#kick-remove-participant)
  - [Add Participant](#add-participant)
  - [Contract Deploy](#contract-deploy)
  - [Example (Mock)](#example-mock)
- [Configuration](#configuration)
- [CLI Reference](#cli-reference)
  - [init onboarding](#init-onboarding)
  - [init kick](#init-kick)
  - [init add-participant](#init-add-participant)
  - [init contract-deploy](#init-contract-deploy)
  - [resume](#resume)
  - [query-parties](#query-parties)
- [Quick Start](#quick-start)
- [Development](#development)

---

## Why Party-Ceremony

Canton's decentralized party model allows multiple independent participants to
jointly own a namespace through threshold signing. This is the mechanism used
to create the **decentralized party that owns the MCMS contract** — the central
on-ledger governance primitive for multi-operator Canton deployments.

Setting up (and later maintaining) a decentralized party requires a
sequence of Canton topology operations: generating signing keys, publishing
namespace delegations, creating and threshold-signing a
`DecentralizedNamespaceDefinition`, and establishing `PartyToParticipant`
mappings. Each step involves admin and/or ledger gRPC calls to different participant nodes, and
many steps require collecting signatures from a quorum of independent operators
before the workflow can proceed.

**Party-ceremony automates this entire lifecycle.** Its primary use case is
deploying and owning the MCMS contract through the decentralized party. After
that initial deployment, every subsequent workflow exists for maintenance
reasons:

| Lifecycle stage                                        | Workflow          |
| ------------------------------------------------------ | ----------------- |
| **Create** the decentralized party                     | `onboarding`      |
| **Deploy** the MCMS contract (or other DAML contracts) | `contract-deploy` |
| **Add** a new participant to the party                 | `add-participant` |
| **Remove** a compromised or decommissioned participant | `kick`            |

---

## Architecture

### Operations Framework

Party-ceremony is built on the
[chainlink-deployments-framework/operations](https://github.com/smartcontractkit/chainlink-deployments-framework)
package. Every workflow is defined as a **Sequence** composed of individual
**Operations**:

- A **Sequence** is a named, versioned function that chains multiple operations
  together with control flow (loops, gates, validation).
- An **Operation** is a named, versioned function that performs a single
  side-effect (e.g. generate a key, sign a transaction, submit a proposal).

Each operation is **idempotent**: the framework computes a deterministic hash
of `(operation definition, input data)` and stores the result in a **Reporter**.
If the same operation+input has already succeeded, the cached result is returned
immediately without re-executing the side-effect. This is the foundation of the
multi-actor resume pattern.

### Async Multi-Actor Model

A ceremony involves multiple independent operators, each running the CLI at
different times from their own machine. The tool is designed around this
constraint:

1. **Operator A** runs `init` — this creates the ceremony, executes Operator A's
   steps, and persists the results.
2. **Operator B** runs `resume` with the same workflow ID — the framework loads
   Operator A's cached results, skips their steps instantly, and executes
   Operator B's pending steps.
3. This repeats until a threshold of operators have contributed, at which point
   submission and finalization steps execute automatically.

Operations that belong to a different participant are skipped (the operation
detects that the current participant doesn't own the relevant key). Threshold
gates return `ErrThresholdNotMet` until enough actors have contributed, causing
the CLI to exit with code `2` ("come back later").

### State Persistence

Each ceremony is stored in a directory (`<state-dir>/<workflow-id>/`) containing
two files:

| File            | Mutability  | Purpose                                                                                                                                   |
| --------------- | ----------- | ----------------------------------------------------------------------------------------------------------------------------------------- |
| `workflow.json` | Immutable   | Ceremony type + full input parameters. Written once by `init`. Allows `resume` to reconstruct the workflow without re-supplying flags.    |
| `reports.json`  | Append-only | Cached operation results (the idempotency store). Updated after every CLI invocation — even on error — so partial progress is never lost. |

Both files use **atomic writes** (write to temp file, then rename) to prevent
corruption if the process is interrupted.

### Coordinator Role

One participant is designated as the **coordinator**. The coordinator is
responsible for:

- Creating proposals (e.g. the `DecentralizedNamespaceDefinition` proposal)
- Merging collected signatures and submitting the fully-signed transaction
- Preparing and executing contract submissions (for `contract-deploy`)

Non-coordinator participants contribute keys, signatures, and P2P mappings.
Any participant can be the coordinator; the choice is made at `init` time.

### Exit Codes

| Code | Meaning                                                                                  |
| ---- | ---------------------------------------------------------------------------------------- |
| `0`  | Ceremony completed successfully                                                          |
| `1`  | Unrecoverable error                                                                      |
| `2`  | Threshold not yet met — more participants must `resume` before the ceremony can complete |

---

## Production Workflow

In production, the party-ceremony tool is used by multiple independent operators
who each run the binary locally on their own machine. Each operator's binary has
access to their own Canton Admin API (and Ledger API for contract deployments).

Ceremony state is shared through a **dedicated Git repository** where each
workflow is tracked as a **pull request**:

```
┌──────────────────────────────────────────────────────────────────────────┐
│                         Ceremony Git Repository                          │
│                                                                          │
│  PR #42: "Onboard decentralized party for MCMS"                         │
│  ┌────────────────────────────────────────────────────────────────────┐  │
│  │  ceremonies/<workflow-id>/workflow.json   ← immutable input        │  │
│  │  ceremonies/<workflow-id>/reports.json    ← grows with each run    │  │
│  └────────────────────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────────────────────┘
```

**Operator flow:**

1. **Coordinator** runs `init` with the desired workflow type and parameters.
   This creates `workflow.json` and the coordinator's initial `reports.json`.
2. Coordinator **commits** both files and **opens a PR** in the ceremony repo.
3. **Each operator** pulls the PR branch, runs `resume <workflow-id>` with
   their own `participant-config.json` (pointing to their local Admin API),
   then **commits and pushes** the updated `reports.json`.
4. The next operator **pulls** the latest reports and runs `resume` again.
   Previously completed operations are instant (served from cache); only that
   operator's pending steps execute.
5. Once the threshold is met, the final `resume` run automatically submits and
   completes the ceremony. The PR is merged.

The **workflow ID** (a UUID printed by `init`) is the coordination key shared
across all operators. It identifies the ceremony directory and allows `resume`
to load the correct state.

---

## Workflows

| Type                                  | Steps | Purpose                                               |
| ------------------------------------- | ----- | ----------------------------------------------------- |
| [`onboarding`](#onboarding)           | 6     | Create a new decentralized party from scratch         |
| [`kick`](#kick-remove-participant)    | 5     | Remove a participant from an existing party           |
| [`add-participant`](#add-participant) | 7     | Add a new participant to an existing party            |
| [`contract-deploy`](#contract-deploy) | 7     | Deploy a DAML contract (e.g. MCMS) owned by the party |
| [`example`](#example-mock)            | 4     | Mock reference implementation for testing             |

### Onboarding

Creates a new decentralized party by generating keys, establishing namespace
delegations, and building the topology mappings from scratch.

| #   | Operation                      | Actor            | Description                                                                |
| --- | ------------------------------ | ---------------- | -------------------------------------------------------------------------- |
| 1   | `CreateMemberKeyOp`            | Each participant | Generate a namespace signing key and a DAML (protocol) signing key         |
| 2   | `ProposeNamespaceDelegationOp` | Each participant | Publish own namespace delegation to the synchronizer                       |
| 3   | `CreateDNSProposalOp`          | Coordinator      | Create the `DecentralizedNamespaceDefinition` proposal with all owner keys |
| 4   | `SignDNSProposalOp`            | Each signer      | Add participant's signature to the DNS proposal                            |
| 5   | `SubmitDNSOp`                  | Coordinator      | Merge all collected signatures and submit the fully-signed DNS             |
| 6   | `CreateAndSubmitP2POp`         | Each participant | Authorize the `PartyToParticipant` mapping independently                   |

**Threshold gates:**

- Step 3 requires all participants to have completed step 1 (all keys must exist).
- Step 5 requires `threshold` signatures to have been collected in step 4.
- Step 6 requires the DNS to be confirmed on the synchronizer.

### Kick (Remove Participant)

Removes a participant from an existing decentralized party by updating the
`DecentralizedNamespaceDefinition` (remove the kicked owner) and the
`PartyToParticipant` mapping (remove the kicked hosting participant).

| #   | Operation                 | Actor                  | Description                                                |
| --- | ------------------------- | ---------------------- | ---------------------------------------------------------- |
| 1   | `ReadCurrentStateOp`      | Any                    | Read current DNS and P2P topology state                    |
| 2   | `CreateKickDNSProposalOp` | Coordinator            | Create updated DNS proposal with the kicked owner removed  |
| 3   | `SignKickDNSProposalOp`   | Remaining participants | Sign the updated DNS proposal                              |
| 4   | `SubmitKickDNSOp`         | Coordinator            | Merge signatures and submit the updated DNS                |
| 5   | `ProposeKickP2POp`        | Remaining participants | Propose updated P2P mapping without the kicked participant |

**Authorization rules:**

- Canton requires `threshold-of-current-owners` signatures for serial > 1 DNS
  updates. The kicked participant is still a current owner until the update is
  confirmed and **can sign the DNS proposal** in step 3.
- The kicked participant is **excluded from P2P proposals** in step 5 — only
  remaining participants update the hosting mapping.
- At least 2 owners must remain after the kick.

### Add Participant

Adds a new participant to an existing decentralized party. The reverse of the
kick workflow.

| #   | Operation                | Actor                | Description                                               |
| --- | ------------------------ | -------------------- | --------------------------------------------------------- |
| 1   | `GenerateNewMemberKeyOp` | New participant      | Generate namespace + DAML signing keys                    |
| 2   | `ProposeNewNSDOp`        | New participant      | Publish namespace delegation to the synchronizer          |
| 3   | `ReadCurrentStateOp`     | Any                  | Read current DNS and P2P topology state                   |
| 4   | `CreateAddDNSProposalOp` | Coordinator          | Create updated DNS proposal with the new owner added      |
| 5   | `SignAddDNSProposalOp`   | Existing members     | Sign the updated DNS proposal                             |
| 6   | `SubmitAddDNSOp`         | Coordinator          | Merge signatures and submit the updated DNS               |
| 7   | `ProposeAddP2POp`        | All (existing + new) | Propose updated P2P mapping including the new participant |

**Authorization rules:**

- Only existing members sign the DNS update (the new participant is not yet an
  owner and cannot sign).
- The new participant receives `CONFIRMATION` permission in the P2P mapping.

### Contract Deploy

Deploys a DAML contract owned by the decentralized party using the Canton
Ledger API's `InteractiveSubmissionService`. This is the workflow used to
deploy the MCMS contract.

| #   | Operation             | Actor            | Description                                                                        |
| --- | --------------------- | ---------------- | ---------------------------------------------------------------------------------- |
| 1   | `VerifyPartyOp`       | Any              | Verify the party is visible on the Ledger API                                      |
| 2   | `FetchParticipantsOp` | Any              | Fetch hosting participant UIDs from topology                                       |
| 3   | `UploadDarsOp`        | Each participant | Upload DAML package archives (DARs) via Admin `PackageService`                     |
| 4   | `PrepareSubmissionOp` | Coordinator      | Prepare the contract creation via `InteractiveSubmissionService.PrepareSubmission` |
| 5   | `SignSubmissionOp`    | Each participant | Sign the prepared transaction hash with their DAML signing key                     |
| 6   | `ExecuteSubmissionOp` | Coordinator      | Aggregate all signatures and call `ExecuteSubmission`                              |
| 7   | `VerifyContractOp`    | Any              | Query the Active Contract Set and verify the contract was created                  |

**Threshold gates:**

- Step 4 requires all participants to have uploaded DARs in step 3.
- Step 6 requires all participants to have signed in step 5.

**Known contracts:** The CLI has built-in profiles for known contract packages.
When a known package name is used with `--packages`, the template module,
entity, and args file are auto-populated:

| Package name | Module      | Entity | Args file        |
| ------------ | ----------- | ------ | ---------------- |
| `mcms`       | `MCMS.Main` | `MCMS` | `mcms-args.json` |

DARs are loaded from the `dars/` directory by default (relative to the working
directory), under versioned subdirectories such as `dars/current/` (dev) or
`dars/v1_0_0/` (pinned release). Each participant must have the required DAR
files available locally. When the repo bumps its Canton release, update
`releaseDir` in `ceremony/ops/ledger/deps.go` to match `contracts.ReleaseDir`
(see [bindings/README.md](../bindings/README.md)).

### Example (Mock)

A mock-backed reference implementation of the onboarding ceremony that uses an
in-memory `MockCantonClient` instead of real gRPC calls. Useful for development
and testing without Canton infrastructure. Follows the same 4-step pattern
(init → propose → sign → submit) with deterministic, SHA256-derived outputs.

---

## Configuration

Each operator creates a `participant-config.json` that identifies their
participant and points to their Canton APIs:

```json
{
  "participant_id": "p1",
  "admin_host": "localhost",
  "admin_port": 5001,
  "admin_jwt": "",
  "ledger_host": "localhost",
  "ledger_port": 5002,
  "ledger_jwt": "",
  "kms_namespace_key_id": "",
  "kms_protocol_key_id": ""
}
```

| Field                  | Required | Default     | Description                                                                                                  |
| ---------------------- | -------- | ----------- | ------------------------------------------------------------------------------------------------------------ |
| `participant_id`       | Yes      | —           | Identifies this operator in the ceremony (`p1`, `p2`, …)                                                     |
| `admin_host`           | No       | `localhost` | Canton Admin gRPC API host                                                                                   |
| `admin_port`           | No       | `5001`      | Canton Admin gRPC API port                                                                                   |
| `admin_jwt`            | No       | (empty)     | Bearer token for Admin API authentication                                                                    |
| `ledger_host`          | No       | `localhost` | Canton Ledger gRPC API host (used by `contract-deploy`)                                                      |
| `ledger_port`          | No       | `5002`      | Canton Ledger gRPC API port                                                                                  |
| `ledger_jwt`           | No       | (empty)     | Bearer token for Ledger API authentication                                                                   |
| `kms_namespace_key_id` | No       | (empty)     | Local external KMS key ID (e.g. AWS KMS ARN) for NAMESPACE signing key registration                         |
| `kms_protocol_key_id`  | No       | (empty)     | Local external KMS key ID (e.g. AWS KMS ARN) for PROTOCOL/DAML signing key registration and tx signing      |

### KMS Key Registration

By default, the ceremony generates new signing keys via Canton's internal vault.
For production deployments where Canton participants run in environments with
restricted IAM permissions (e.g. AWS), you can **pre-create** signing keys in
an external KMS and register them instead:

1. Use Terraform (or other infrastructure tooling) to create two KMS keys per
   participant: one for NAMESPACE signing, one for PROTOCOL (DAML) signing.
2. Grant the Canton participant's IAM role permissions to **use** those keys
   (e.g. `kms:Sign`, `kms:GetPublicKey`) but **not** to create new keys.
3. Each operator populates their own `participant-config.json` with the external
   key identifiers (e.g. AWS KMS ARNs). These values stay local and are injected
   on every `init` or `resume`.

When these fields are set, the ceremony calls Canton's
`VaultService.RegisterKmsSigningKey` instead of `GenerateSigningKey`.
Onboarding and add-participant require both namespace and protocol key IDs.
Key rotation only requires the key ID for the key type being rotated.
Contract deploy uses `kms_protocol_key_id` to sign prepared transaction hashes
through AWS KMS; when it is empty, it signs with the participant vault.

**Example with AWS KMS keys:**

```json
{
  "participant_id": "cv1",
  "admin_host": "cv1.devnet.canton",
  "admin_port": 5001,
  "admin_jwt": "eyJhbGciOi...",
  "kms_namespace_key_id": "arn:aws:kms:us-east-1:123456789012:key/a1b2c3d4-...",
  "kms_protocol_key_id": "arn:aws:kms:us-east-1:123456789012:key/e5f6g7h8-..."
}
```

The config file path defaults to `./participant-config.json` and can be
overridden with `--config <path>`.

---

## CLI Reference

### init onboarding

Create a new decentralized party onboarding ceremony.

```bash
canton-party-ceremony init onboarding \
  --new-namespace-name <name> \
  --new-party-name <prefix> \
  --participants <p1,p2,p3> \
  --synchronizer-id <id> \
  [--coordinator <id>] \
  [--threshold <n>] \
  [--config <path>] \
  [--state-dir <dir>]
```

| Flag                   | Required | Default                   | Description                                                       |
| ---------------------- | -------- | ------------------------- | ----------------------------------------------------------------- |
| `--new-namespace-name` | Yes      | —                         | Unique ceremony identifier                                        |
| `--new-party-name`     | Yes      | —                         | Party ID prefix used to derive the final party identifier         |
| `--participants`       | Yes      | —                         | Comma-separated list of participant IDs                           |
| `--synchronizer-id`    | Yes      | —                         | Canton synchronizer ID                                            |
| `--coordinator`        | No       | —                         | Coordinator participant ID                                        |
| `--threshold`          | No       | `0`                       | Minimum signatures required. `0` = strict majority `floor(n/2)+1` |
| `--config`             | No       | `participant-config.json` | Path to participant config JSON                                   |
| `--state-dir`          | No       | `ceremonies`              | Root directory for ceremony state                                 |

### init kick

Remove a participant from an existing decentralized party.

```bash
canton-party-ceremony init kick \
  --decentralized-party-id <prefix::namespace> \
  --kicked-participant-id <PAR::name::fp> \
  --kicked-namespace-fingerprint <1220...> \
  --remaining-participants <PAR::a::fp1,PAR::b::fp2> \
  --synchronizer-id <id> \
  [--new-threshold <n>] \
  [--config <path>] \
  [--state-dir <dir>]
```

| Flag                             | Required | Default                   | Description                                                 |
| -------------------------------- | -------- | ------------------------- | ----------------------------------------------------------- |
| `--decentralized-party-id`       | Yes      | —                         | Full party ID (`<prefix>::<namespace>`)                     |
| `--kicked-participant-id`        | Yes      | —                         | Canton UID of the participant to remove                     |
| `--kicked-namespace-fingerprint` | Yes      | —                         | Namespace fingerprint of the kicked participant             |
| `--remaining-participants`       | Yes      | —                         | Comma-separated Canton UIDs of remaining participants (≥ 2) |
| `--synchronizer-id`              | Yes      | —                         | Canton synchronizer ID                                      |
| `--new-threshold`                | No       | `0`                       | Post-kick threshold. `0` = strict majority of remaining     |
| `--config`                       | No       | `participant-config.json` | Path to participant config JSON                             |
| `--state-dir`                    | No       | `ceremonies`              | Root directory for ceremony state                           |

Use [`query-parties`](#query-parties) to obtain the identifiers needed for
these flags.

### init add-participant

Add a new participant to an existing decentralized party.

```bash
canton-party-ceremony init add-participant \
  --decentralized-party-id <prefix::namespace> \
  --new-participant-id <PAR::newnode::fp> \
  --namespace-name <label> \
  --synchronizer-id <id> \
  [--new-threshold <n>] \
  [--config <path>] \
  [--state-dir <dir>]
```

| Flag                       | Required | Default                   | Description                                    |
| -------------------------- | -------- | ------------------------- | ---------------------------------------------- |
| `--decentralized-party-id` | Yes      | —                         | Full party ID (`<prefix>::<namespace>`)        |
| `--new-participant-id`     | Yes      | —                         | Canton UID of the participant to add           |
| `--namespace-name`         | Yes      | —                         | Label for the new participant's key generation |
| `--synchronizer-id`        | Yes      | —                         | Canton synchronizer ID                         |
| `--new-threshold`          | No       | `0`                       | Post-addition threshold. `0` = keep current    |
| `--config`                 | No       | `participant-config.json` | Path to participant config JSON                |
| `--state-dir`              | No       | `ceremonies`              | Root directory for ceremony state              |

### init contract-deploy

Deploy a DAML contract owned by the decentralized party.

```bash
canton-party-ceremony init contract-deploy \
  --decentralized-party-id <prefix::namespace> \
  --synchronizer-id <id> \
  --packages <name:version,...> \
  [--template-module <module>] \
  [--template-entity <entity>] \
  [--contract-args-file <path>] \
  [--config <path>] \
  [--state-dir <dir>]
```

| Flag                       | Required | Default                        | Description                                                |
| -------------------------- | -------- | ------------------------------ | ---------------------------------------------------------- |
| `--decentralized-party-id` | Yes      | —                              | Full party ID (`<prefix>::<namespace>`)                    |
| `--synchronizer-id`        | Yes      | —                              | Canton synchronizer ID                                     |
| `--packages`               | Yes      | —                              | Comma-separated `name:version` pairs (e.g. `mcms:current`) |
| `--template-module`        | No       | Auto-filled for known packages | Fully-qualified DAML module name (e.g. `MCMS.Main`)        |
| `--template-entity`        | No       | Auto-filled for known packages | DAML template entity name (e.g. `MCMS`)                    |
| `--contract-args-file`     | No       | Auto-filled for known packages | Path to JSON file with contract creation arguments         |
| `--config`                 | No       | `participant-config.json`      | Path to participant config JSON                            |
| `--state-dir`              | No       | `ceremonies`                   | Root directory for ceremony state                          |

For [known contracts](#contract-deploy), `--template-module`, `--template-entity`,
and `--contract-args-file` are auto-populated. For custom contracts, all three
must be provided explicitly.

### resume

Resume an existing ceremony from its persisted state.

```bash
canton-party-ceremony resume <ceremony-id> \
  [--config <path>] \
  [--state-dir <dir>]
```

| Flag          | Required | Default                   | Description                       |
| ------------- | -------- | ------------------------- | --------------------------------- |
| `--config`    | No       | `participant-config.json` | Path to participant config JSON   |
| `--state-dir` | No       | `ceremonies`              | Root directory for ceremony state |

`resume` loads `workflow.json` to determine the ceremony type and input, then
loads `reports.json` to seed the reporter cache. All previously completed
operations are served from cache; only the current participant's pending
operations execute.

**Exit behavior:**

- Exit `0` — ceremony complete.
- Exit `2` — threshold not met. Print a message and wait for more operators.
- Exit `1` — unrecoverable error.

### query-parties

List all decentralized parties visible to the local participant. The output
provides the identifiers required for the `kick` and `add-participant`
ceremonies.

```bash
canton-party-ceremony query-parties \
  --synchronizer-id <id> \
  [--config <path>]
```

| Flag                | Required | Default                   | Description                     |
| ------------------- | -------- | ------------------------- | ------------------------------- |
| `--synchronizer-id` | Yes      | —                         | Canton synchronizer ID to query |
| `--config`          | No       | `participant-config.json` | Path to participant config JSON |

**Example output:**

```
Namespace:           1220abcdef...
  Threshold:         2-of-3
  Serial:            1
  Owners:
    - fingerprint:   1220aaa...
    - fingerprint:   1220bbb...
    - fingerprint:   1220ccc...
  Party ID:          prefix::1220abcdef...
  P2P Threshold:     2
  P2P Serial:        2
  Hosting Participants:
    - uid: PAR::p1::1220...  permission: CONFIRMATION
    - uid: PAR::p2::1220...  permission: CONFIRMATION
    - uid: PAR::p3::1220...  permission: CONFIRMATION
---
```

---

## Quick Start

### Prerequisites

- Go 1.25+

### Build

```bash
make build
```

This produces `./bin/canton-party-ceremony`.

### Local 3-participant onboarding (using `example` mock)

This walkthrough uses the `example` workflow (mock Canton client) so no running
Canton nodes are required.

**1. Initialize the ceremony as the coordinator (`p1`):**

Create `participant-config.json`:

```json
{
  "participant_id": "p1",
  "admin_host": "not_needed_for_example",
  "admin_port": 0
}
```

```bash
./bin/canton-party-ceremony init example \
  --new-namespace-name "decentralized-namespace" \
  --new-party-name "dec-party" \
  --participants p1,p2,p3 \
  --synchronizer-id global \
  --threshold 2
```

This prints the workflow ID: `Ceremony initialised: 3c015d2b-...`

**2. Resume as `p2`:**

Update `participant-config.json` to `"participant_id": "p2"`, then:

```bash
./bin/canton-party-ceremony resume <workflow-id>
```

If the threshold (2) is now met, the ceremony completes. Otherwise exit code `2`
is returned — continue with the next participant.

**3. Resume as `p3` (if needed):**

```bash
# participant-config.json → "participant_id": "p3"
./bin/canton-party-ceremony resume <workflow-id>
```

### Real onboarding (with Canton nodes)

For a real ceremony against live Canton nodes, use `init onboarding` instead of
`init example` and configure each operator's `participant-config.json` with
their actual Admin API endpoint and JWT.

---

## Development

```bash
make build             # Build the binary
make test              # Run unit tests
make test-integration  # Run integration tests (requires Canton)
make lint              # Run golangci-lint
make lint-all          # Lint all modules (including integration-tests)
make lint-fix-all      # Lint with auto-fix across all modules
make clean             # Remove build artifacts
```
