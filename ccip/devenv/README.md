# devenv

The `devenv` package contains the needed implementations needed to properly spin up a CCIP
environment that includes Canton.

## Prerequisites

Make sure to checkout [chainlink-ccv](https://github.com/smartcontractkit/chainlink-ccv/) in a parallel directory,
i.e. so that chainlink-canton and chainlink-ccv have the same parent directory. This is needed so that we can build
needed docker images.

## Spin up an environment

The quickest way to run the EVM/Canton environment is from the root of the repo:

```bash
# From chainlink-canton root
make start-devenv
```

This will build all the necessary docker images from chainlink-ccv and chainlink-canton, and spin up an environment with
both EVM and Canton chains.

## Run a test

After the environment is spun up, you can run a test like:

```bash
# From repo root
make run-e2e-tests
```

## Load tests

Load tests live in `ccip/devenv/tests/load`. They use [WASP](https://pkg.go.dev/github.com/smartcontractkit/chainlink-testing-framework/wasp) and run sequentially (RPS=1) because Canton holdings are single-flight.

### Canton → EVM load (requires devenv)

Sequential Canton→EVM messages round-robined across every EVM destination in the env file. Requires `make start-devenv` so `ccip/devenv/env-canton-evm-out.toml` exists.

Schedule is configured via env vars (defaults are `1/1s` for 90s):

| Env var | Form | Default | Meaning |
|---|---|---|---|
| `CANTON_LOAD_MESSAGE_RATE` | `<int>/<duration>` (e.g. `1/1s`, `1/20s`, `10/5m`) | `1/1s` | rate per rate-limit window |
| `CANTON_LOAD_DURATION` | Go duration (e.g. `90s`, `10m`, `1h`) | `90s` | total runtime |

Example — 1 message every 20 seconds for 10 minutes:

```bash
CANTON_LOAD_MESSAGE_RATE=1/20s CANTON_LOAD_DURATION=10m make run-canton2evm-load
```

```bash
# From repo root
make run-canton2evm-load
```

Equivalent:

```bash
cd ccip/devenv/tests/load && go test -timeout 15m -v -count 1 -run '^TestCanton2EVM_Load$'
```

If the out file is missing the test skips with a hint.

### Canton → EVM token load (requires devenv)

Separate test from message-only load: `TestCanton2EVM_TokenLoad`. Resolves the token lane declared in [`token_transfer_config.toml`](./tests/token_transfer_config.toml) (see [Token lane configuration](#token-lane-configuration)) against the source chain's `GetTokenTransferConfigs`, validating every destination has the lane. Pre-mints Canton fee + transfer holdings, runs WASP, then asserts EVM receiver token balance delta.

```bash
make run-canton2evm-token-load
```

Equivalent:

```bash
cd ccip/devenv/tests/load && go test -timeout 20m -v -count 1 -run '^TestCanton2EVM_TokenLoad$'
```

### EVM → Canton load (requires devenv)

Sequential EVM→Canton messages against the Canton destination in the env file. Uses the same schedule env vars as Canton→EVM (`CANTON_LOAD_MESSAGE_RATE`, `CANTON_LOAD_DURATION`). EVM accounts are pre-funded by devenv; no Canton pre-mint.

Example — 1 message every 20 seconds for 10 minutes:

```bash
CANTON_LOAD_MESSAGE_RATE=1/20s CANTON_LOAD_DURATION=10m make run-evm2canton-load
```

```bash
# From repo root
make run-evm2canton-load
```

Equivalent:

```bash
cd ccip/devenv/tests/load && go test -timeout 15m -v -count 1 -run '^TestEVM2Canton_Load$'
```

### EVM → Canton token load (requires devenv)

Separate test from message-only load: `TestEVM2Canton_TokenLoad`. Resolves the token lane declared in [`token_transfer_config.toml`](./tests/token_transfer_config.toml) (see [Token lane configuration](#token-lane-configuration)), logs EVM sender balance vs estimated transfer need (devenv pre-funds sender), runs WASP, logs Canton holdings post-run.

```bash
make run-evm2canton-token-load
```

Equivalent:

```bash
cd ccip/devenv/tests/load && go test -timeout 20m -v -count 1 -run '^TestEVM2Canton_TokenLoad$'
```

### Token lane configuration

Token transfer tests (both e2e and load) declare the token lane to send in [`tests/token_transfer_config.toml`](./tests/token_transfer_config.toml). The test resolves the declared **token pool identity** against the source chain's `GetTokenTransferConfigs`, then validates that every destination chain has that lane configured before running. The committed file holds the devenv defaults.

Override the file path with `CANTON_TOKEN_TEST_CONFIG`:

```bash
CANTON_TOKEN_TEST_CONFIG=/path/to/custom.toml make run-canton2evm-token-load
```

The file has one block per direction (`[evm_to_canton]`, `[canton_to_evm]`):

| Key | Required | Meaning |
|---|---|---|
| `pool_type` | yes | token pool contract type on the source chain (e.g. `BurnMintTokenPool`, `LockReleaseTokenPool`) |
| `pool_version` | yes | semantic version of the token pool (e.g. `2.0.0`) |
| `pool_qualifier` | yes | datastore qualifier identifying the exact pool |
| `transfer_amount` | no | per-message token amount (string integer); falls back to a per-direction default |
| `execution_gas_limit` | no | per-message execution gas limit; falls back to a per-direction default |
| `finality_config` | no | per-message finality config; falls back to a per-direction default |

Token **identity** (`pool_type` / `pool_version` / `pool_qualifier`) is always required and is never defaulted; only the numeric send params have code-level fallbacks. If the qualifier matches zero or multiple pools, or a destination lacks the lane, the test fails fast listing the available pool refs / selectors.

```toml
[evm_to_canton]
pool_type      = "BurnMintTokenPool"
pool_version   = "2.0.0"
pool_qualifier = "TEST (BurnMintTokenPool 2.0.0 [default], LockReleaseTokenPool 2.0.0 [default])::BurnMintTokenPool 2.0.0 [default]"
transfer_amount     = "100000000000"
execution_gas_limit = 200000
finality_config     = 0

[canton_to_evm]
pool_type      = "LockReleaseTokenPool"
pool_version   = "2.0.0"
pool_qualifier = "TEST (BurnMintTokenPool 2.0.0 [default], LockReleaseTokenPool 2.0.0 [default])::LockReleaseTokenPool 2.0.0 [default]"
transfer_amount     = "1000"
execution_gas_limit = 500000
finality_config     = 1
```

### CI (on demand)

The **CCIP Canton Load Tests** workflow (`ccip-load-tests.yml`) can be triggered manually from GitHub Actions
(`workflow_dispatch`). It reuses the same devenv setup as the CCIP E2E workflow and runs one of the load tests depending on the `direction` input. Inputs:

| Input | Default | Maps to |
|---|---|---|
| `direction` | `canton2evm` | `canton2evm` → `TestCanton2EVM_Load`; `evm2canton` → `TestEVM2Canton_Load`; `canton2evm-token` → `TestCanton2EVM_TokenLoad`; `evm2canton-token` → `TestEVM2Canton_TokenLoad` |
| `message_rate` | `1/1s` | `CANTON_LOAD_MESSAGE_RATE` |
| `load_duration` | `90s` | `CANTON_LOAD_DURATION` |
| `test_timeout` | `40m` | `go test -timeout` |
| `canton_ref` | workflow ref | chainlink-canton checkout |

chainlink-ccv is pinned in `.github/actions/setup-ccip-devenv` (same as CCIP E2E).

## Shortcut

If you want to build docker images, spin up a new env, and run the test in a single command, the following is useful:

```bash
# From repo root
make build-run-e2e-tests
```
