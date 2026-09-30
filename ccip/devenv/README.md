# devenv

The `devenv` package provides helpers and bindings to spin up a CCIP environment that includes Canton. It is used by three independent test suites: **E2E**, **Integration**, and **Load**.

## Prerequisites

* Checkout the companion repository **[chainlink-ccv](https://github.com/smartcontractkit/chainlink-ccv/)** in a sibling directory so that `chainlink-canton` and `chainlink-ccv` share the same parent. This is required to build the Docker images used by the tests.

## Spin up a local environment

```bash
# From the repository root
make start-devenv
```

This builds the Docker images from `chainlink-ccv` and `chainlink-canton` and starts a local EVM + Canton network.

## Test suites

The repository contains three test suites. Each suite has its own README with detailed usage instructions:

* **[E2E tests](tests/e2e/README.md)** – end‑to‑end message‑flow tests.
* **[Integration tests](tests/integration/README.md)** – connectivity and smoke‑test checks.
* **[Load tests](tests/load/README.md)** – performance / stress tests using WASP.

All suites share a common environment selector (`-ccip-env` / `CCIP_ENV`). The supported environments are:

| Value                | Config file                     | Remote |
|----------------------|---------------------------------|--------|
| `devenv` (default)   | `env-canton-evm-out.toml`       | no     |
| `prod-testnet`       | `env-prod-testnet-out.toml`      | yes    |
| `mainnet` (future)   | –                               | –      |

Use the flag `-ccip-env=<env>` or the environment variable `CCIP_ENV` (the flag takes precedence) to choose the target.

## Ledger bindings (`-tags=prodledger`)

When targeting production contracts, build the tests with the tag `-tags=prodledger` so they use the older `bindings/generated/v1_0_0` layout. The default builds use `bindings/generated/latest`.

| Target               | Build tag          | Bindings                |
|----------------------|--------------------|------------------------|
| `devenv` (default)   | _(none)_           | `bindings/generated/latest` |
| `prod-testnet` / `mainnet` | `-tags=prodledger` | `bindings/generated/v1_0_0` |

Implementation lives in [`ccip/devenv/ledgertarget/`](./ledgertarget/).

## Shortcut

To build the Docker images, start the environment, and run the E2E suite in one step:

```bash
# From the repository root
make build-run-e2e-tests
```

For detailed usage of each test suite, follow the links above.
