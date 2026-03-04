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

## Shortcut

If you want to build docker images, spin up a new env, and run the test in a single command, the following is useful:

```bash
# From repo root
make build-run-e2e-tests
```
