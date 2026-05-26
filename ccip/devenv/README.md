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

### CI (on demand)

The **CCIP Canton Load Tests** workflow (`ccip-load-tests.yml`) can be triggered manually from GitHub Actions
(`workflow_dispatch`). It reuses the same devenv setup as the CCIP E2E workflow. Inputs:

| Input | Default | Maps to |
|---|---|---|
| `message_rate` | `1/1s` | `CANTON_LOAD_MESSAGE_RATE` |
| `load_duration` | `90s` | `CANTON_LOAD_DURATION` |
| `test_timeout` | `20m` | `go test -timeout` |
| `canton_ref` | workflow ref | chainlink-canton checkout |

chainlink-ccv is pinned in `.github/actions/setup-ccip-devenv` (same as CCIP E2E).

## Shortcut

If you want to build docker images, spin up a new env, and run the test in a single command, the following is useful:

```bash
# From repo root
make build-run-e2e-tests
```
