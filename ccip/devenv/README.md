# devenv

The `devenv` package contains the needed implementations needed to properly spin up a CCIP
environment that includes Canton.

## Spin up an environment

```bash
# From repo root
cd ccip/devenv
go run cmd/ccv/main.go up env-canton-evm.toml
```

## Run a test

After the environment is spun up, you can run a test like:

```bash
# From repo root
cd ccip/devenv/tests/e2e
go test -v -run "TestEVM2Canton_Basic" .
```
