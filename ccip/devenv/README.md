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
| `CANTON_LOAD_TOKEN_QUALIFIER` | pool qualifier string (e.g. `TEST` on staging) | _(auto)_ | override token lane when multiple lanes exist |

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

Separate test from message-only load: `TestCanton2EVM_TokenLoad`. Discovers the BM 2.0 ↔ LR 2.0 lane from the env file (override with `CANTON_LOAD_TOKEN_QUALIFIER`). Pre-mints Canton fee + transfer holdings, runs WASP, then asserts EVM receiver token balance delta.

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

Separate test from message-only load: `TestEVM2Canton_TokenLoad`. Discovers token lane from deployment, logs EVM sender balance vs estimated transfer need (devenv pre-funds sender), runs WASP, logs Canton holdings post-run.

```bash
make run-evm2canton-token-load
```

Equivalent:

```bash
cd ccip/devenv/tests/load && go test -timeout 20m -v -count 1 -run '^TestEVM2Canton_TokenLoad$'
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
