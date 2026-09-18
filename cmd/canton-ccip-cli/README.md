# CCIP Canton <-> EVM CLI

A small Cobra-based CLI that demonstrates sending and executing CCIP messages between Ethereum and
Canton, on either mainnet or testnet (Sepolia / Canton testnet).

## Build / Run

The CLI lives at `cmd/canton-ccip-cli`. All paths below are relative to the **repository root**.

By default the CLI looks for `./config.yaml` in the current working directory. Use `--config`
(shorthand `-c`) to point at a different file.

```bash
# Uses ./config.yaml from the current directory
go run ./cmd/canton-ccip-cli <command>

# Or pass an explicit path
go run ./cmd/canton-ccip-cli --config ./cmd/canton-ccip-cli/config.example.yaml <command>
```

The `--network` flag (shorthand `-n`) selects a static profile baked into the binary. It defaults to
`mainnet`; pass `--network testnet` to target Sepolia / Canton testnet.

```bash
go run ./cmd/canton-ccip-cli --network testnet <command>
```

Both flags are long flags, so they need two dashes (`--network`) or the single-letter shorthand
(`-n`). A single-dash `-network` is parsed as a shorthand cluster and will not work.

To build a binary instead of using `go run`:

```bash
go build -o canton-ccip-cli ./cmd/canton-ccip-cli
./canton-ccip-cli --help
```

## Config file

See [`config.example.yaml`](./config.example.yaml).

Only these fields are required: `canton.participantGRPCLedgerAPIURL`, `canton.validatorAPIURL`,
`canton.userID`, `canton.partyID` and `evm.rpcURL`. Everything else is optional —
`evm.privateKeyHex` can be omitted when signing with a Ledger (see below), and the explorer URL
fields can be left empty.

> **Note:** the top-level `network:` key **overrides** the `--network` flag. If it is set, the flag
> is ignored. Omit `network:` from your config if you want to select the network on the command
> line.

```yaml
# Optional. When set, this wins over --network.
network: mainnet

canton:
  # One of: authorizationCode (default), clientCredentials, static.
  authType: "authorizationCode"
  authServerURL: "https://auth.example.com"
  authClientID: "my-client-id"
  # Only for authType: clientCredentials.
  authClientSecret: ""
  # Only for authType: static — a pre-issued JWT, used as-is without validation.
  authJWT: ""
  participantGRPCLedgerAPIURL: "participant.example.com:5001"
  validatorAPIURL: "https://validator.example.com"
  userID: "my-user"
  partyID: "myParty::1220..."
  # Optional Party ID -> EDS URL mappings, merged over the built-in network profile.
  edsURLs:
    tokenPoolOwner:1220...: "https://token-pool-eds.example.com"
  # Optional Party ID -> Token Standard API URL mappings, merged over the built-in network profile.
  tokenStandardURLs:
    tokenOwner:1220...: "https://token-standard-eds.example.com"

evm:
  rpcURL: "https://eth-sepolia.example.com"
  # Optional; omit to sign with a Ledger via --ledger.
  privateKeyHex: "0xabc..."

ccip_explorer_url: ""
evm_explorer_url: ""
canton_explorer_url: ""
```

## Commands

Flags shared by several commands (`--fee-token`, `--executor`, `--finality`, `--gas-limit`,
`--ledger`, holding selection, `--wait`) are described once under
[Common flags](#common-flags).

### EVM → Canton

| Command                                       | Description                                                                                  |
|-----------------------------------------------|----------------------------------------------------------------------------------------------|
| `evm send-message`                            | Send a message from EVM to Canton with no token transfer.                                    |
| `evm send-token --amount <wei>`               | Send a token transfer (LINK by default) from EVM to Canton.                                  |
| `canton execute --message-id <0xhash>`        | Execute on Canton a message that was sent from EVM.                                          |

Command-specific flags:

* `evm send-message`: `--receiver-party <party>` (defaults to own party),
  `--payload <text>` (default `Hello, Canton!`).
* `evm send-token`: `--receiver-party <party>` (defaults to own party), `--amount <wei>`
  (**required**; supports exponents, e.g. `1e18`), `--token <address>` (defaults to LINK).
* `canton execute`: `--message-id <0xhash>` (**required**). Must use the same `--finality` value as
  the send if the send used faster-than-finality.

Because no executor currently supports Canton as a destination, the EVM-side send commands always
use the `noExecution` tag in extraArgs.

#### Example
`go run ./cmd/canton-ccip-cli  --config  ./cmd/canton-ccip-cli/config.test.yaml --network testnet evm send-message --receiver-party u_5da1d5aca0c7::1220c250c23c55120f7c758bccc5cbc739629015ab921594e1c29656981f985bffa7 --payload "test-message-evm-to-canton" --fee-token native --finality 1`

### Canton → EVM

| Command                                        | Description                                                                                 |
|------------------------------------------------|-----------------------------------------------------------------------------------------------|
| `canton send-message`                          | Send a message-only CCIP message from Canton to EVM.                                        |
| `canton send-token --amount <decimal>`         | Send a token transfer (LINK by default) from Canton to EVM.                                 |
| `evm execute --message-id <0xhash>`            | Execute on EVM a message that was sent from Canton with `--executor none`.                  |

Command-specific flags:

* `canton send-message`: `--receiver <0xhex>` (defaults to the profile's `CCIPReceiver` contract),
  `--payload <text>` (default `Hello, EVM from Canton!`). Default `--fee-token` is `link`.
* `canton send-token`: `--receiver <0xhex>` (defaults to own address), `--amount <decimal>`
  (**required**, e.g. `0.12345`, `1e-2`), `--token <id>@<admin>` (defaults to LINK),
  `--payload <text>` (optional, empty by default). Default `--fee-token` is **`native`**, unlike
  `canton send-message`. When paying the fee in LINK, two separate input holdings must be provided —
  one for the fee (`--fee-input`) and one for the token transfer (`--token-input`).
* `evm execute`: `--message-id <0xhash>` (**required**).

#### Example
`go run ./cmd/canton-ccip-cli  --config  ./cmd/canton-ccip-cli/config.test.yaml --network testnet canton send-message --receiver 0x15a8a0831a7FEdda5CFabF66452b92a980d204A1 --payload "test-message-canton-to-evm" --executor default --fee-token native`

### Canton utilities

| Command                                           | Description                                                                                                          |
|---------------------------------------------------|------------------------------------------------------------------------------------------------------------------------|
| `canton list-events [--event {sent\|executed}]`   | List active `CCIPMessageSent` (default) or `ExecutionStateChanged` contracts visible to the configured party.        |
| `canton list-holdings [--cid]`                    | List all holdings for the configured party. `--cid` also prints each holding's Contract ID.                          |
| `canton list-transfer-instructions`               | List all `TransferInstruction` contracts for the configured party.                                                   |
| `canton create-transfer --amount <decimal>`       | Create an outgoing `TransferInstruction`. `--receiver <party>` defaults to own party; `--token` defaults to `link`.   |
| `canton accept-transfer --contract-id <cid>`      | Accept an incoming `TransferInstruction` by contract ID. `--token` defaults to `link`.                               |
| `canton sync-receiver-ccv --required-ccv <addr>`  | Deploy or update the `CCIPReceiver` required CCVs. Run once before inbound load/e2e on prod.                         |

`--required-ccv` takes a Canton committee verifier raw address in `instanceId@owner` form. Note that
`sync-receiver-ccv` defaults `--finality` to `1`, not `finality`.

### External party onboarding (Ledger)

These commands are grouped under `canton external-party` and support the external-party onboarding
flow using a Ledger device and the Canton admin APIs.

| Command                                                                                             | Description                                                                                                                                                     |
|-----------------------------------------------------------------------------------------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `canton external-party get-fingerprint [--derivation-path <path-or-index>]`                         | Reads the Ledger key, prints its fingerprint (`ldg::<fingerprint>`), and prints the DER public key to share with the operator.                                  |
| `canton external-party prepare-topology [--synchronizer-id <id>] <party-hint> <der-public-key-hex>` | Calls `GenerateExternalPartyTopology` and prints a topology JSON payload containing `partyId`, `publicKeyFingerprint`, `multiHash`, and `topologyTransactions`. |
| `canton external-party sign-topology [--derivation-path <path-or-index>] '<topology-json>'`         | Verifies the topology payload against the Ledger key and signs the topology multihash. Prints the same JSON with `signature` added.                             |
| `canton external-party submit-topology [--synchronizer-id <id>] '<signed-topology-json>'`           | Calls `AllocateExternalParty` with the signed topology payload and verifies the party was created.                                                              |

`--derivation-path` defaults to `m/44'/6767'/0'/0'/0'`. You can also pass just an index (for example
`42`), which maps to `m/44'/6767'/0'/0'/42'`. `--synchronizer-id` defaults to the global
synchronizer.

## Common flags

| Flag                          | Applies to                                                                                              | Notes                                                                                                                     |
|-------------------------------|---------------------------------------------------------------------------------------------------------|-----------------------------------------------------------------------------------------------------------------------------|
| `--fee-token`                 | `evm send-message`, `evm send-token`, `canton send-message`, `canton send-token`                        | `link\|native`; `canton send-message` also accepts an `<InstrumentId>`. Defaults to `link`, except `canton send-token` which defaults to `native`. |
| `--executor {default\|none}`  | `canton send-message`, `canton send-token`                                                              | Defaults to `default`. `none` skips automatic execution on the destination chain; settle manually with `evm execute`.     |
| `--finality`                  | `evm send-message`, `evm send-token`, `canton execute`, `canton sync-receiver-ccv`                      | `finality` (full), `safe`, or a block depth `1`–`65535`. See below.                                                       |
| `--gas-limit <int>`           | `canton send-message`, `canton send-token`                                                              | Gas limit for EVM execution. Defaults to `-1`, meaning `50000`.                                                           |
| `--fee-input <cid>`           | `canton send-message`, `canton send-token`                                                              | Repeatable. Holding(s) to use as input for the fee payment. If unspecified, all current holdings are used.                |
| `--token-input <cid>`         | `canton send-token`                                                                                     | Repeatable. Holding(s) to use as input for the token transfer. If unspecified, all current holdings are used.             |
| `--input <cid>`               | `canton create-transfer`                                                                                | Repeatable. Holding(s) to use as input for the transfer. If unspecified, all current holdings are used.                   |
| `--wait <duration>`           | `evm execute`, `canton execute`                                                                         | Go duration syntax (e.g. `15m`, `30s`). Max time to wait for verifier results. Default `15m`.                             |
| `--ledger <path-or-index>`    | all signing commands on both sides                                                                      | See [Ledger signing](#ledger-signing).                                                                                    |

**Finality (`--finality`):** defaults to `finality` (wait for full Sepolia finalization, slowest).
Use `1` for one block confirmation (~12s), `safe` for the Ethereum safe head, or any depth
`1`–`65535`. When using faster finality, pass the same value to `canton execute`. The CLI selects
(or creates) a `CCIPReceiver` whose finality config matches `--finality`, sets `requiredCCVs` from
the indexer attestation (`verifier_dest_address`), and can keep separate receivers for full finality
and faster-than-finality on the same party.

## Ledger signing

Both the EVM and Canton sides support Ledger signing via `--ledger`, but they use **different
derivation schemes**.

**EVM (`evm send-message`, `evm send-token`, `evm execute`):** transactions are signed with a
connected Ledger device (with the Ethereum app open) instead of the configured private key,
including the ERC20 allowance approval. The flag accepts either a full path like `m/44'/60'/0'/0/0`
or a bare index like `42`, which replaces the **hardened account component**, e.g.
`m/44'/60'/42'/0/0`. The sender address is read from the device, so `evm.privateKeyHex` can be
omitted from the config when using a Ledger.

**Canton (`canton send-message`, `canton send-token`, `canton execute`, `canton create-transfer`,
`canton accept-transfer`, `canton sync-receiver-ccv`):** enables interactive Ledger signing. The
flag accepts either a full path like `m/44'/6767'/0'/0'/0'` or a bare index like `42`, which
replaces the **last component**, e.g. `m/44'/6767'/0'/0'/42'`.

## Notes

* Run `go run ./cmd/canton-ccip-cli <command> --help` for the authoritative flag list.
