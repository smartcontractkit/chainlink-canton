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

### E2E environment selection (`-ccip-env` / `CCIP_ENV`)

Message e2e tests run against **local devenv** by default or **Canton TestNet + Sepolia** when prod-testnet is selected. Set the environment by name (not TOML path):

| Value | Config file | Remote |
|---|---|---|
| `devenv` (default) | [`env-canton-evm-out.toml`](./env-canton-evm-out.toml) | no |
| `prod-testnet` | [`env-prod-testnet-out.toml`](./env-prod-testnet-out.toml) | yes |

Use the `-ccip-env` flag or `CCIP_ENV` env var (flag wins if both are set):

```bash
# Local devenv (default)
cd ccip/devenv/tests/e2e && go test -v -run 'TestCanton2EVM_Basic/EOA' -count=1

# Prod testnet
CCIP_ENV=prod-testnet \
  CANTON_GRPC_URL=... CANTON_PARTY_ID=... CANTON_AUTH_*=... \
  PRIVATE_KEY=... \
  go test -timeout 8m -v -count=1 -ccip-env=prod-testnet \
  -run 'TestEVM2Canton_Basic/message|TestCanton2EVM_Basic/EOA'
```

**Prod prerequisites**

- Canton party wallet funded with at least **50 Amulet units** (message fee)
- Canton auth env vars (`CANTON_GRPC_URL`, `CANTON_PARTY_ID`, `CANTON_AUTH_*`) — see [Prod testnet connection smoke test](#prod-testnet-connection-smoke-test)
- Sepolia gas via `PRIVATE_KEY` (EVM sender / receiver for prod runs)

**Optional instance ID overrides** (defaults: `test-router`, `e2e-ccipsender`, `e2e-receiver`):

| Env var | Default |
|---|---|
| `CANTON_ROUTER_INSTANCE_ID` | `test-router` |
| `CANTON_SENDER_INSTANCE_ID` | `e2e-ccipsender` |
| `CANTON_RECEIVER_INSTANCE_ID` | `e2e-receiver` |

Token e2e subtests are skipped on prod-testnet. A second prod run reuses existing router/sender/receiver contracts on ledger when instance IDs match.

## Load tests

Load tests live in `ccip/devenv/tests/load`. They use [WASP](https://pkg.go.dev/github.com/smartcontractkit/chainlink-testing-framework/wasp) and run sequentially (RPS=1) because Canton holdings are single-flight.

### Canton → EVM load

Sequential Canton→EVM messages round-robined across every EVM destination in the env file.

**Devenv** (requires `make start-devenv` so `ccip/devenv/env-canton-evm-out.toml` exists): pre-mints fee holdings and calls `SetupSend` once before WASP starts. Full send + EVM exec confirmation per message.

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

**Prod-testnet** (Canton TestNet + Sepolia): send-only message load — Canton send + confirm send, no `ConfirmExecOnDest` on EVM (executor not available on prod). Set `CANTON_LOAD_SKIP_EXEC_CONFIRM=true`. Verify delivery via indexer/CCIP ops; the test does not assert EVM execution.

Prerequisites: `CANTON_GRPC_URL`, `CANTON_PARTY_ID`, `CANTON_AUTH_*`, `PRIVATE_KEY` (EVM message receiver wallet), and a pre-funded Canton party (~50 Amulet per message at `CantonToEVMFeeAmount`).

```bash
CCIP_ENV=prod-testnet \
CANTON_LOAD_SKIP_EXEC_CONFIRM=true \
CANTON_GRPC_URL=... CANTON_PARTY_ID=... CANTON_AUTH_CLIENT_ID=... CANTON_AUTH_CLIENT_SECRET=... \
PRIVATE_KEY=0x... \
CANTON_LOAD_MESSAGE_RATE=1/10s \
CANTON_LOAD_DURATION=5m \
go test -timeout 30m -v -count=1 -ccip-env=prod-testnet \
  -run '^TestCanton2EVM_Load$' ./ccip/devenv/tests/load/
```

Or from repo root:

```bash
make run-canton2evm-load-prod
```

### Canton → EVM token load (requires devenv)

Separate test from message-only load: `TestCanton2EVM_TokenLoad`. Resolves the token lane declared in [`token_transfer_config.toml`](./tests/token_transfer_config.toml) (see [Token lane configuration](#token-lane-configuration)) against the source chain's `GetTokenTransferConfigs`, validating every destination has the lane. Pre-mints Canton fee + transfer holdings, runs WASP, then asserts EVM receiver token balance delta.

```bash
make run-canton2evm-token-load
```

Equivalent:

```bash
cd ccip/devenv/tests/load && go test -timeout 20m -v -count 1 -run '^TestCanton2EVM_TokenLoad$'
```

### EVM → Canton load

Sequential EVM→Canton messages against the Canton destination in the env file. Uses the same schedule env vars as Canton→EVM (`CANTON_LOAD_MESSAGE_RATE`, `CANTON_LOAD_DURATION`).

**Devenv** (requires running devenv + `env-canton-evm-out.toml`): EVM accounts are pre-funded by devenv; no Canton pre-mint.

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

**Prod-testnet** (Canton TestNet + Sepolia): message-only load with full per-message confirmation (send → receipt on EVM → `ConfirmExecOnDest` on Canton). Use a conservative rate — each iteration is synchronous (~30–60s end-to-end on prod), so WASP RPS=1 is effectively bounded by confirm latency. Budget for Sepolia gas per send plus Canton execution fees. Token load remains devenv-only; Canton→EVM message-only prod load is supported send-only (see above).

Prerequisites: `CANTON_GRPC_URL`, `CANTON_PARTY_ID`, `CANTON_AUTH_*`, `PRIVATE_KEY` (Sepolia sender/receiver wallet), and a pre-funded Canton party.

```bash
CCIP_ENV=prod-testnet \
CANTON_GRPC_URL=... CANTON_PARTY_ID=... CANTON_AUTH_CLIENT_ID=... CANTON_AUTH_CLIENT_SECRET=... \
PRIVATE_KEY=0x... \
CANTON_LOAD_MESSAGE_RATE=1/30s \
CANTON_LOAD_DURATION=5m \
CANTON_CONFIRM_EXEC_TIMEOUT=10m \
go test -timeout 45m -v -count=1 -ccip-env=prod-testnet \
  -run '^TestEVM2Canton_Load$' ./ccip/devenv/tests/load/
```

Or from repo root:

```bash
make run-evm2canton-load-prod
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

Load tests share one composite action (`.github/actions/ccip-load-test`) used by:

- **CCIP Canton Load Tests** (`ccip-load-tests.yml`) — manual `workflow_dispatch`
- **CCIP Load (ChatOps)** (`ccip-load-command.yml`) — PR comment `/ccip-load`

#### Manual workflow (`workflow_dispatch`)

| Input | Default (devenv) | Default (prod-testnet) | Maps to |
|---|---|---|---|
| `ccip_env` | `devenv` | select `prod-testnet` | `-ccip-env` |
| `direction` | `canton2evm` | same | test `-run` regex |
| `message_rate` | `1/1s` | `1/45s` (when left at devenv default) | `CANTON_LOAD_MESSAGE_RATE` |
| `load_duration` | `90s` | `2m` (when left at devenv default) | `CANTON_LOAD_DURATION` |
| `test_timeout` | `40m` | `30m` / `45m` for evm2canton | `go test -timeout` |
| `config_file` | — | `env-prod-testnet.ci.toml` | `CCIP_CONFIG_FILE` |
| `skip_exec_confirm` | — | `true` | `CANTON_LOAD_SKIP_EXEC_CONFIRM` |
| `confirm_exec_timeout` | — | `10m` | `CANTON_CONFIRM_EXEC_TIMEOUT` |
| `canton_ref` | workflow ref | devenv only | chainlink-canton checkout |

**Devenv** spins up Docker via `setup-ccip-devenv` (same as CCIP E2E). **Prod-testnet** hits live Canton TestNet + Sepolia with no local devenv.

#### ChatOps (PR comment)

On an open PR, comment (all named args optional):

```
/ccip-load direction=canton2evm config_file=env-prod-testnet.ci.toml message_rate=1/45s load_duration=2m
```

| Named arg | Default | Notes |
|---|---|---|
| `ccip_env` | `prod-testnet` | |
| `direction` | `canton2evm` | `evm2canton`, `*-token` (token skips on prod) |
| `config_file` | `env-prod-testnet.ci.toml` | basename under `ccip/devenv/`; sets `CCIP_CONFIG_FILE` |
| `message_rate` | `1/45s` | |
| `load_duration` | `2m` | |
| `test_timeout` | `30m` / `45m` for evm2canton | |
| `skip_exec_confirm` | `true` | required on prod for canton→evm |
| `confirm_exec_timeout` | `10m` | |
| `party_id` | ccip-app party (see composite action default) | |
| `grpc_url` | `testnet.cv1.bcy-v.metalhosts.com:443` | |
| `auth_type` | `clientCredentials` | CI only; local dev uses `authorizationCode` |
| `auth_url`, `client_id` | from GitHub secrets | override secrets when set |

Requires write access to the repo (same as `/auto-fix`).

#### Prod-testnet config file

`*-out.toml` files are gitignored (local devenv output). CI uses the committed snapshot [`env-prod-testnet.ci.toml`](./env-prod-testnet.ci.toml) instead. Override with `config_file=` or `CCIP_CONFIG_FILE` when running locally:

```bash
CCIP_CONFIG_FILE=env-prod-testnet.ci.toml CCIP_ENV=prod-testnet go test ...
```

#### GitHub secrets (prod-testnet CI)

Workflows pass these secrets to `.github/actions/ccip-load-test` using the same names (1:1); the composite action maps them to test env vars:

| Secret (input name) | Env var |
|---|---|
| `CANTON_OKTA_AUTHORIZER_TESTNET` | `CANTON_AUTH_URL` |
| `CANTON_OKTA_CLIENT_ID_TESTNET` | `CANTON_CLIENT_ID` |
| `CANTON_OKTA_CLIENT_SECRET_TESTNET` | `CANTON_CLIENT_SECRET` |
| `CCIP_PROD_TESTNET_PRIVATE_KEY` | `PRIVATE_KEY` |

Devenv secrets unchanged: `CCV_IAM_ROLE`, `JD_REGISTRY`, `JD_IMAGE`.

chainlink-ccv is pinned in `.github/actions/setup-ccip-devenv` (devenv only).

## Shortcut

If you want to build docker images, spin up a new env, and run the test in a single command, the following is useful:

```bash
# From repo root
make build-run-e2e-tests
```

## Prod testnet connection smoke test

Minimal Canton-only connectivity check against real testnet infrastructure. Uses [`env-prod-testnet-out.toml`](./env-prod-testnet-out.toml) (Canton TestNet ↔ Sepolia message lane: contract refs, indexer URLs, EDS URL, verifier issuers); auth secrets and party identity come from environment variables.

The test is opt-in: it skips unless `CANTON_GRPC_URL` is set (even though the TOML file includes default URLs). This keeps CI green while allowing manual or workflow-triggered runs.

### Environment variables

| Variable | Required | Notes |
|---|---|---|
| `CANTON_PARTY_ID` | yes | Ledger party to query; skips `GetUser` when set |
| `CANTON_AUTH_URL` | yes | OIDC issuer |
| `CANTON_CLIENT_ID` | yes | OAuth2 client ID |
| `CANTON_AUTH_TYPE` | no | `authorizationCode` (local), `clientCredentials` (CI), `static`, `insecureStatic` |
| `CANTON_USER_ID` | no | Required for `clientCredentials`; optional for `authorizationCode` (extracted from token `sub` after login) |
| `CANTON_CLIENT_SECRET` | CI only | Required when `CANTON_AUTH_TYPE=clientCredentials` |
| `CANTON_JWT` | static only | For `static` / `insecureStatic` auth |
| `CANTON_GRPC_URL` | opt-in signal | Must be set to run the test; overrides TOML gRPC URL |
| `CANTON_VALIDATOR_API_URL` | no | Overrides TOML validator API URL |

**Local (browser Okta login):**

```bash
export CANTON_GRPC_URL='testnet.cv1.bcy-v.metalhosts.com:443'
export CANTON_PARTY_ID='u_0e0328cbbcb7::1220c250c23c55120f7c758bccc5cbc739629015ab921594e1c29656981f985bffa7'
export CANTON_AUTH_URL='https://smartcontract.okta.com/oauth2/austsuml9q2WhPBMM5d7'
export CANTON_CLIENT_ID='0oau1l22b1Jv3dcih5d7'
export CANTON_AUTH_TYPE='authorizationCode'
```

**CI (`clientCredentials`, no browser):**

```bash
export CANTON_GRPC_URL='testnet.cv1.bcy-v.metalhosts.com:443'
export CANTON_PARTY_ID='...'
export CANTON_AUTH_URL='https://smartcontract.okta.com/oauth2/austsuml9q2WhPBMM5d7'
export CANTON_CLIENT_ID='0oau1l22b1Jv3dcih5d7'
export CANTON_AUTH_TYPE='clientCredentials'
export CANTON_CLIENT_SECRET='...'
export CANTON_USER_ID='...'
```

### Run command

```bash
cd ccip/devenv/tests/integration && go test -v -run TestIntegration_CantonProdTestnet_Connection -count=1
```

Use `-ccip-env=prod-testnet` (or `CCIP_ENV=prod-testnet`) so the test loads `env-prod-testnet-out.toml` locally, or `CCIP_CONFIG_FILE=env-prod-testnet.ci.toml` for the committed CI snapshot; default devenv config is `env-canton-evm-out.toml` via `-ccip-env=devenv`.

The test connects via `NewCLDF`, asserts `PartyID` is set, and lists holdings for the party (empty balance is OK).
