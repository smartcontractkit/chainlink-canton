# E2E test suite

The **E2E** suite validates end‑to‑end message‑flow behaviour of the CCIP system (Canton ↔ EVM). It exercises the full stack, including contract deployment, message sending, and execution confirmation.

## Running the suite

All E2E tests share the common environment selector (`-ccip-env` / `CCIP_ENV`). The supported environments are the same as described in the top‑level README:

| Value                | Config file                     | Remote |
|----------------------|---------------------------------|--------|
| `devenv` (default)   | `env-canton-evm-out.toml`       | no     |
| `prod-testnet`       | `env-prod-testnet-out.toml`      | yes    |

### Local `devenv`

```bash
cd ccip/devenv/tests/e2e
go test -v -count=1 -run 'TestCanton2EVM_Basic/EOA'
```

### Remote `prod-testnet`

```bash
CCIP_ENV=prod-testnet \
  CANTON_GRPC_URL=... CANTON_PARTY_ID=... CANTON_AUTH_CLIENT_ID=... CANTON_AUTH_CLIENT_SECRET=... \
  PRIVATE_KEY=... \
  go test -timeout 8m -v -count=1 -ccip-env=prod-testnet \
  -run 'TestEVM2Canton_Basic/message|TestCanton2EVM_Basic/EOA'
```

**Prod‑testnet prerequisites**

* Canton party wallet funded with at least **50 Amulet** units (message fee).
* Auth environment variables (`CANTON_GRPC_URL`, `CANTON_PARTY_ID`, `CANTON_AUTH_*`). See the [Prod testnet connection smoke test](../README.md#prod-testnet-connection-smoke-test) for details.
* Sepolia gas via `PRIVATE_KEY` (EVM sender/receiver).

### Optional instance ID overrides

The following variables can be set to reuse existing router/sender/receiver contracts on the ledger:

| Variable                     | Default |
|------------------------------|---------|
| `CANTON_ROUTER_INSTANCE_ID` | `test-router` |
| `CANTON_SENDER_INSTANCE_ID` | `e2e-ccipsender` |
| `CANTON_RECEIVER_INSTANCE_ID` | `e2e-receiver` |

### Token‑E2E tests (prod‑testnet only)

Both **EVM → Canton** and **Canton → EVM** token transfers are supported on `prod-testnet`. They use the same command line as above; the test runner will automatically exercise the token lanes defined in `tests/token_transfer_config.toml`.

---

For a one‑step build‑run‑E2E shortcut, see the top‑level README.