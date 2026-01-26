# Canton Integration Tests

This directory contains integrations tests for the Canton network.

## Running the Tests

To simply run an integration test, run the corresponding Go test. For example, to run the CCIP integration tests:

```bash
go test -v ./canton/ccip/...
```

## Using a local Canton environment

By default, the tests will use CLDF/CTF to spin up a local Canton network with a variable number of participants.
For local testing, it can be quicker to use an existing Canton network instead of spinning up a new one each time.
To do this, set the `PARTICIPANT_INPUT` environment variable to a `.toml` file containing the configuration for each participant.
Example configuration files include `local-docker-compose.toml` and `local-docker-compose-direct.toml` in this directory.

To run the CCIP integration tests with the `local-docker-compose.toml` configuration file:

```bash
PARTICIPANT_INPUT=../local-docker-compose-direct.toml go test ./ccip -v -count 1
```

```toml
selector=8706591216959472610 # The chain selector for the Canton network

scan_api_url="http://scan.localhost:8080" # The Scan API URL to use
registry_api_url="http://scan.localhost:8080" # The Token Registry API URL to use

[[participants]]
name="Participant 1" # A human-readable name for the participant
jwt="ey...jbJ0" # JWT token to authenticate with this participant
username="ledger-api-user" # The username for this participant, the JWT must have a subject claim matching this username
json_ledger_api_url="http://participant1.json-ledger-api.localhost:8080" # The JSON Ledger API URL to use
grpc_ledger_api_url="participant1.grpc-ledger-api.localhost:8080" # The gRPC Ledger API URL to use
admin_api_url="participant1.admin-api.localhost:8080" # The gRPC Admin API URL to use
validator_api_url="http://participant1.validator-api.localhost:8080/api/validator" # The Validator API URL to use

... # Other participants
```

Tests require a variable number of participants depending on the functionality being tested.
If the provided configuration does not contain enough participants, the test will fail.