# Canton Party Ceremony

A CLI tool for managing **decentralized party onboarding ceremonies** on Canton distributed ledgers.


## Prerequisites

- Go 1.25


## Build

```bash
make build
```

This produces `./bin/canton-party-ceremony`.

---

## Configuration

Each invocation is driven by a `participant-config.json` file that identifies the **calling participant**:

```json
{
  "participant_id": "p1",
  "admin_jwt": "not_needed_for_example",
  "admin_host": "not_needed_for_example",
  "admin_port": 0
}
```

| Field            | Description                                                          |
|------------------|----------------------------------------------------------------------|
| `participant_id` | The ID of the participant running this invocation (`p1`, `p2`, …).  |
| `admin_jwt`      | Bearer token for the Canton Admin gRPC API (empty = no auth).       |
| `admin_host`     | Canton Admin API host (defaults to `localhost`).                     |
| `admin_port`     | Canton Admin API port (defaults to `5001`).                         |

The config file path defaults to `./participant-config.json` and can be overridden with `--config <path>`.

---

## Example Ceremony: Local multi-party onboarding ceremony

The following walkthrough runs a 3-participant onboarding with a threshold of 2.

### 1. Build

```bash
make build
```

### 2. Initialise the ceremony (run once, as the coordinator)

Set `participant_id` to `p1` in `participant-config.json`, then run:

```bash
./bin/canton-party-ceremony init onboarding \
  --coordinator p1 \
  --new-namespace-name "decentralized-namespace" \
  --new-party-name "dec-party" \
  --participants p1,p2,p3 \
  --synchronizer-id global \
  --threshold 2
```

This creates a new ceremony directory under `ceremonies/` and prints the generated `<workflow_ID>`, for example:

```
Ceremony initialised: 3c015d2b-e957-4e56-8bb3-2d54b2a82633
```

The coordinator's own step is executed automatically during `init`.

### 3. Collect signatures (run once per additional participant)

Each remaining participant updates `participant_id` in `participant-config.json` to their own ID and runs:

```bash
# As participant p2
# participant-config.json → "participant_id": "p2"
./bin/canton-party-ceremony resume <workflow_ID>

# As participant p3
# participant-config.json → "participant_id": "p3"
./bin/canton-party-ceremony resume <workflow_ID>
```

Repeat until the threshold is reached (`2` signatures in this example).  
The command exits with code `2` and a human-readable message when more signatures are still needed.

Once the threshold is met the final `submit` step is executed automatically and the ceremony is complete.

## Development

```bash
make test      # run unit tests
make lint      # run golangci-lint
make lint-fix  # run golangci-lint with auto-fix
make clean     # remove build artifacts
```
