# CCIP Canton <-> EVM Demo CLI

A small Cobra-based CLI that demonstrates sending and executing CCIP messages between Ethereum Mainnet and Canton
Mainnet.

## Build / Run

By default, the CLI looks for `./config.yaml` in the current working directory. Use `--config <path>` to point at a
different file.

```bash
# Uses ./config.yaml from the current directory
go run ./examples/cli <command>

# Or pass an explicit path
go run ./examples/cli --config ./examples/cli/config.example.yaml <command>
```

The `--network` flag selects a static profile baked into the binary. It defaults to `mainnet`; pass `--network testnet`
to target Sepolia / Canton testnet.

```bash
go run ./examples/cli --network testnet <command>
```

## Config file

See [`config.example.yaml`](./config.example.yaml). All `canton.*` and `evm.*` fields are required; the optional
explorer URL fields can be left empty.

```yaml
canton:
  authType: "authorizationCode"
  authServerURL: "https://auth.example.com"
  authClientID: "my-client-id"
  participantGRPCLedgerAPIURL: "participant.example.com:5001"
  validatorAPIURL: "https://validator.example.com"
  userID: "my-user"
  partyID: "myParty::1220..."

evm:
  rpcURL: "https://eth-sepolia.example.com"
  privateKeyHex: "0xabc..."

ccip_explorer_url: ""
evm_explorer_url: ""
canton_explorer_url: ""
```

## Commands

### EVM → Canton

| Command                                                                                                                  | Description                                                                                                  |
|--------------------------------------------------------------------------------------------------------------------------|--------------------------------------------------------------------------------------------------------------|
| `evm send-message --receiver-party <p> [--payload <text>] [--fee-token {link\|native}] [--finality {finality\|safe\|N}]` | Send a message from EVM to Canton with no token transfer.                                                    |
| `evm send-token   --receiver-party <p> --amount <wei> [--fee-token {link\|native}] [--finality {finality\|safe\|N}]`     | Send a token-only transfer (LINK) from EVM to Canton.                                                        |
| `evm execute      --message-id <0xhash> [--wait <duration>]`                                                             | Execute on EVM a message that was sent from Canton with the `noneExecution` executor.                        |
| `canton execute   --message-id <0xhash> [--finality {finality\|safe\|N}] [--wait <duration>]`                            | Execute on Canton a message sent from EVM. Must use matching `--finality` if send used faster-than-finality. |

Because no executor currently supports Canton as a destination, the EVM-side send commands always use the `noExecution`
tag in extraArgs.

**Finality (`--finality`):** defaults to `finality` (wait for full Sepolia finalization, slowest). Use `1` for one block
confirmation (~12s), `safe` for Ethereum safe head, or any depth `1`–`65535`. When using faster finality, pass the same
value to `canton execute`. The CLI selects (or creates) a
`CCIPReceiver` whose finality config matches `--finality`, sets `requiredCCVs`
from the indexer attestation (`verifier_dest_address`), and can keep separate receivers for full finality and
faster-than-finality on the same party.

### Canton → EVM

| Command                                                                                                                                  | Description                                                                                                                                                                                       |
|------------------------------------------------------------------------------------------------------------------------------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `canton send-message --receiver <0xhex> [--payload <text>] [--executor {default\|none}] [--fee-token {link\|native}]`                    | Send a message-only CCIP message from Canton to EVM.                                                                                                                                              |
| `canton send-token   --receiver <0xhex> --amount <decimal> [--payload <text>] [--executor {default\|none}] [--fee-token {link\|native}]` | Send a LINK token transfer CCIP message from Canton to EVM. Note: when using LINK as the fee token, two separate input holdings must be provided, one for the fee and one for the token transfer. |
| `canton execute --message-id <0xhash> [--wait <duration>]`                                                                               | Execute on Canton a message sent from EVM (see EVM → Canton section for `--finality`).                                                                                                            |
| `canton list-events --event {sent\|executed}`                                                                                            | List active `CCIPMessageSent` or `ExecutionStateChanged` contracts visible to the configured party.                                                                                               |
| `canton list-holdings [--cid]`                                                                                                           | List all holdings for the configured party. Specify `--cid` to include the holding's Contract ID in the output.                                                                                   |
| `canton create-transfer --amount <decimal> [--token {link\|native}] [--receiver <party>] [--input <contractId>]`                         | Initiates a transfer of an given amount of an asset to another party/itself. Defaults to LINK, use `--token native` to switch to Amulet.                                                          |
| `canton accept-transfer --contract-id <contractId> [--token {link\|native}]`                                                             | Accept an incoming `TransferInstruction` by contract ID. Defaults to LINK; use `--token native` for Amulet.                                                                                       |

`--executor none` skips automatic execution on the destination chain; you can then settle the message manually with
`evm execute --message-id ...`.

### External party onboarding (Ledger)

These commands are grouped under `canton external-party` and support the external-party onboarding flow using a Ledger
device and the Canton admin APIs.

| Command                                                                                             | Description                                                                                                                                                     |
|-----------------------------------------------------------------------------------------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `canton external-party get-fingerprint [--derivation-path <path-or-index>]`                         | Reads the Ledger key, prints its fingerprint (`ldg::<fingerprint>`), and prints the DER public key to share with the operator.                                  |
| `canton external-party prepare-topology [--synchronizer-id <id>] <party-hint> <der-public-key-hex>` | Calls `GenerateExternalPartyTopology` and prints a topology JSON payload containing `partyId`, `publicKeyFingerprint`, `multiHash`, and `topologyTransactions`. |
| `canton external-party sign-topology [--derivation-path <path-or-index>] '<topology-json>'`         | Verifies the topology payload against the Ledger key and signs the topology multihash. Prints the same JSON with `signature` added.                             |
| `canton external-party submit-topology [--synchronizer-id <id>] '<signed-topology-json>'`           | Calls `AllocateExternalParty` with the signed topology payload and verifies the party was created.                                                              |

`--derivation-path` defaults to `m/44'/6767'/0'/0'/0'`. You can also pass just an index (for example `42`), which maps
to `m/44'/6767'/0'/0'/42'`.

## Notes

* `--wait` accepts Go duration syntax (e.g. `15m`, `30s`). Default is `15m`.
