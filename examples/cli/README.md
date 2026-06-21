# CCIP Canton <-> EVM Demo CLI

A small Cobra-based CLI that demonstrates sending and executing CCIP messages
between Ethereum Mainnet and Canton Mainnet.

## Build / Run

By default, the CLI looks for `./config.yaml` in the current working directory.
Use `--config <path>` to point at a different file.

```bash
# Uses ./config.yaml from the current directory
go run ./examples/cli <command>

# Or pass an explicit path
go run ./examples/cli --config ./examples/cli/config.example.yaml <command>
```

The `--network` flag selects a static profile baked into the binary. It
defaults to `mainnet`; pass `--network testnet` to target Sepolia / Canton
testnet.

```bash
go run ./examples/cli --network testnet <command>
```

## Config file

See [`config.example.yaml`](./config.example.yaml). All fields are required.

```yaml
canton:
  authServerURL: "https://auth.example.com"
  authClientID: "my-client-id"
  participantGRPCLedgerAPIURL: "participant.example.com:5001"
  validatorAPIURL: "https://validator.example.com"
  userID: "my-user"
  partyID: "myParty::1220..."

evm:
  rpcURL: "https://eth-sepolia.example.com"
  privateKeyHex: "0xabc..."
```

## Commands

### EVM → Canton

| Command                                                                                 | Description                                                                           |
|-----------------------------------------------------------------------------------------|---------------------------------------------------------------------------------------|
| `evm send-message --receiver-party <p> [--payload <text>] [--fee-token {link\|native}]` | Send a message from EVM to Canton with no token transfer.                             |
| `evm send-token   --receiver-party <p> --amount <wei> [--fee-token {link\|native}]`     | Send a token-only transfer (LINK) from EVM to Canton.                                 |
| `evm execute      --message-id <0xhash> [--wait <duration>]`                            | Execute on EVM a message that was sent from Canton with the `noneExecution` executor. |

Because no executor currently supports Canton as a destination, the EVM-side
send commands always use the `noExecution` tag in extraArgs.

### Canton → EVM

| Command                                                                                                                                  | Description                                                                                                                                                                                        |
|------------------------------------------------------------------------------------------------------------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `canton send-message --receiver <0xhex> [--payload <text>] [--executor {default\|none}] [--fee-token {link\|native}]`                    | Send a message-only CCIP message from Canton to EVM.                                                                                                                                               |
| `canton send-token   --receiver <0xhex> --amount <decimal> [--payload <text>] [--executor {default\|none}] [--fee-token {link\|native}]` | Send a LINK token transfer CCIP message from Canton to EVM. Note: when using LINK as the fee token, two sepoarate input holdings must be provided, one for the fee and one for the token transfer. |
| `canton execute --message-id <0xhash> [--wait <duration>]`                                                                               | Execute on Canton a message sent from EVM.                                                                                                                                                         |
| `canton list-events --event {sent\|executed}`                                                                                            | List active `CCIPMessageSent` or `ExecutionStateChanged` contracts visible to the configured party.                                                                                                |
| `canton list-holdings [--cid]`                                                                                                           | List all holdings for the configured party. Specify `--cid` to include the holding's Contract ID in the output.                                                                                    |
| `canton create-transfer --amount <decimal> [--token {link\|native}] [--receiver <party>] [--input <contractId>]`                         | Initiates a transfer of an given amount of an asset to another party/itself. Defaults to LINK, use `--token native` to switch to Amulet.                                                           |

`--executor none` skips automatic execution on the destination chain; you can
then settle the message manually with `evm execute --message-id ...`.

## Notes

* `--wait` accepts Go duration syntax (e.g. `15m`, `30s`). Default is `15m`.

