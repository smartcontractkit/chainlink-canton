# Scripts

## Archive Active Canton Contracts

Use [`archive_active_canton_contracts/main.go`](./archive_active_canton_contracts/main.go) to archive active Canton contracts from ACS before removing old DARs.

This helper:
- connects directly to the Canton ledger gRPC API with JWT auth
- queries ACS for the templates you specify
- archives every matching active contract for the given party
- archives duplicate and repeated deployments too, not just the latest instance

### Dry Run

Always inspect first:

```sh
cd /Users/sish/Desktop/chainlink-canton
export ONCHAIN_CANTON_JWT_TOKEN="$(canton-login canton-devnet)"

go run ./scripts/archive_active_canton_contracts \
  --grpc-url canton-devnet.bcy-v.metalhosts.com:443 \
  --party 'ccipOwner::122009b21cd8f52c55fc9854f7b1771837fa661f47529956ed40708bf656905b5721' \
  --template '#ccip-rmn:CCIP.RMNRemote:RMNRemote' \
  --template '#ccip-common:CCIP.GlobalConfig:GlobalConfig' \
  --template '#ccip-tokenadminregistry:CCIP.TokenAdminRegistry:TokenAdminRegistry' \
  --template '#ccip-feequoter:CCIP.FeeQuoter:FeeQuoter' \
  --template '#ccip-offramp:CCIP.OffRamp:OffRamp' \
  --template '#ccip-onramp:CCIP.OnRamp:OnRamp' \
  --template '#ccip-perpartyrouter:CCIP.PerPartyRouter:PerPartyRouterFactory' \
  --template '#ccip-committeeverifier:CCIP.CommitteeVerifier:CommitteeVerifier' \
  --template '#ccip-executor:CCIP.Executor:Executor' \
  --dry-run
```

### Archive All Standard CCIP Canton Contracts

```sh
cd /Users/sish/Desktop/chainlink-canton
export ONCHAIN_CANTON_JWT_TOKEN="$(canton-login canton-devnet)"

go run ./scripts/archive_active_canton_contracts \
  --grpc-url canton-devnet.bcy-v.metalhosts.com:443 \
  --party 'ccipOwner::122009b21cd8f52c55fc9854f7b1771837fa661f47529956ed40708bf656905b5721' \
  --template '#ccip-rmn:CCIP.RMNRemote:RMNRemote' \
  --template '#ccip-common:CCIP.GlobalConfig:GlobalConfig' \
  --template '#ccip-tokenadminregistry:CCIP.TokenAdminRegistry:TokenAdminRegistry' \
  --template '#ccip-feequoter:CCIP.FeeQuoter:FeeQuoter' \
  --template '#ccip-offramp:CCIP.OffRamp:OffRamp' \
  --template '#ccip-onramp:CCIP.OnRamp:OnRamp' \
  --template '#ccip-perpartyrouter:CCIP.PerPartyRouter:PerPartyRouterFactory' \
  --template '#ccip-committeeverifier:CCIP.CommitteeVerifier:CommitteeVerifier' \
  --template '#ccip-executor:CCIP.Executor:Executor'
```

### Archive Only Specific ACS Entries

If you want to limit the helper to exact contracts, add one or more `--contract-id` filters:

```sh
go run ./scripts/archive_active_canton_contracts \
  --grpc-url canton-devnet.bcy-v.metalhosts.com:443 \
  --party 'ccipOwner::...' \
  --template '#ccip-committeeverifier:CCIP.CommitteeVerifier:CommitteeVerifier' \
  --contract-id '0032143d1758259a03b8b5fcc65bd25b1bdb9340c59913649ba54e386925775b10ca1212202670743759abcd69e3b196705b1729d1cc33dbe054a9bd856b034ce27dffdd86'
```

### Notes

- The package must still be vetted when you archive the contracts.
- `participant.dars.remove(...)` does not archive contracts. It only removes DARs after the relevant active contracts are gone.
- Use the exact party that owns the contracts, usually `ccipOwner::...`.
- If needed, you can use `--user-id` instead of `--party` and the helper will resolve the user's primary party.

## Optional: Remove DARs After Archiving

This is a separate step and must be done in the Canton participant admin console, not in this Go helper.

List the remaining CCIP DARs:

```scala
participant.dars.list().filter(_.name.startsWith("ccip-"))
```

Remove a DAR by its current `mainPackageId`:

```scala
participant.dars.remove("<mainPackageId>")
```

Only remove package IDs that still appear in `participant.dars.list()`. If Canton returns `DAR_NOT_FOUND`, that DAR is already gone or you are using an old package ID.
