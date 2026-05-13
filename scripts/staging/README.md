# Staging Scripts

Run everything from:

```sh
cd /Desktop/chainlink-canton
```

## Repo compatibility

Staging CLIs use **`chainlink-deployments-framework/chain/canton/provider`** for ledger gRPC, **`testhelpers/eds`** for hosted EDS, and **`deployment/utils/operations/contract`** for ledger lookups. They do **not** import **`party-ceremony/internal/client`** or **`deployment/operations/ccip/factory`** (e.g. factory MCMS changesets). Commits that only touch those areas—such as **`44de5367fc889c250b61aeded652f2645ba11dc1`** (factory changesets + TLS on port 443 for party-ceremony gRPC)—do **not** require staging script changes.

## Prerequisite

Set up `scripts/staging/.env` first.

- copy your working staging env into `scripts/staging/.env`
- make sure it includes the auth, RPC, selector, router, on-ramp/off-ramp, token, and token-pool values needed for the lane you are running
- for EVM -> Canton manual execute, set `STAGING_EVM_TO_CANTON_INDEXER_URL` and `STAGING_EVM_TO_CANTON_CCV` in `scripts/staging/.env`
- for the current DRIP_V2 EVM -> Canton token lane, `STAGING_EVM_TO_CANTON_DEST_TOKEN` should match the live pool remote token

## 1. Hosted EDS

These staging scripts use the staging EDS URL from `scripts/staging/.env` by default. You do not need to run a local EDS for the normal staging send/manual-execute flows. If you do want to override it for local debugging, set `STAGING_CANTON_EDS_URL` or pass `-eds-url`.

## 2. Send with `ccipSend`

### EVM -> Canton message

```sh
go run -mod=mod ./scripts/staging/send_staging_evm_to_canton \
  -data 'hello from evm to canton'
```

### EVM -> Canton token transfer

`requested-finality` defaults to `1`.

```sh
go run -mod=mod ./scripts/staging/send_staging_evm_to_canton \
  -data 'hello from evm to canton with token' \
  -token 0x7fa86664F404f3D63813cD9a6d07d66dd9085691 \
  -token-amount 1000000000000
```

This command uses the current DRIP_V2 token lane values from `scripts/staging/.env`. The bare command without `-token-amount` is message-only unless you also set `STAGING_EVM_TO_CANTON_TOKEN_AMOUNT` in `scripts/staging/.env`.

### Canton -> EVM message

Start on `smartcontracts.com`. The script fetches hosted EDS send disclosures first, then pauses so you can switch to `chainlink_legacy` before Canton and scan-proxy work begins.

If the party has never had a **PerPartyRouter**, the send script will create one: it finds **`PerPartyRouterFactory`** on the Canton ledger (ACS for your participant party), builds the same single-factory disclosure that hosted EDS would return, then exercises **`CreateRouter`**. If your party cannot see the factory template in ACS, set **`STAGING_CANTON_PER_PARTY_ROUTER_FACTORY_CID`** or **`-per-party-router-factory-cid`** to the factory contract id from ops / datastore. The router DAML `instanceId` defaults to **`STAGING_CANTON_TO_EVM_SENDER_INSTANCE_ID`** unless you set **`STAGING_CANTON_TO_EVM_ROUTER_INSTANCE_ID`** / **`-router-instance-id`**. You still need an on-ledger **CCIPSender** (and executor / committee verifier) from deploy—only the router is bootstrapped here. **CCIP send** still uses hosted EDS (`STAGING_CANTON_EDS_URL`) for send disclosures after the router exists.

```sh
go run -mod=mod ./scripts/staging/send_staging_canton_to_evm \
  -data 'hello from canton to evm'
```

### Canton -> EVM token transfer

Start on `smartcontracts.com`. The script fetches hosted EDS send disclosures first, then pauses so you can switch to `chainlink_legacy` before Canton and scan-proxy work begins.

```sh
go run -mod=mod ./scripts/staging/send_staging_canton_to_evm_token
```

Each send command prints a `messageID`. Keep it for manual execution.

## 3. Manual Execute

### Manual execute EVM -> Canton

```sh
go run -mod=mod ./scripts/staging/manual_execute_staging_evm_to_canton \
  -message-id 0x<message-id>
```

This uses `STAGING_EVM_TO_CANTON_INDEXER_URL` from `scripts/staging/.env` by default.

For EVM -> Canton token executes, `STAGING_EVM_TO_CANTON_CCV` should be set to the Canton CommitteeVerifier instance address used by EDS execute disclosures.

This script intentionally pauses after auth, indexer lookup, and hosted EDS execute-disclosure lookup so you can switch VPNs to `chainlink_legacy` from `smartcontracts.com` before the Canton ledger transaction begins.

Why:

- indexer/hosted EDS access and Canton ledger write access may require different network paths in staging
- the script waits at that handoff point so you can move to the VPN that allows the Canton transaction to succeed

For token transfers, this now also logs:

- `receiverHoldingsBefore`
- `receiverHoldingsAfter`
- `receiverHoldingsDelta`

Tiny transfers such as `0.0000010000` may still show `receiverHoldingsDelta: 0` due to log precision. Use `-query-update-id` and check for `TokenReceiveTicketClaimed` plus `Amulet` events to confirm receipt.

### Manual execute Canton -> EVM

Always use high gas. This path is EVM-side execution and does not need the Canton-side VPN handoff used by the EVM -> Canton manual execute flow:

```sh
go run -mod=mod ./scripts/staging/manual_execute_staging_canton_to_evm \
  -message-id 0x<message-id> \
  -tx-gas-limit 5000000 \
  -gas-limit-override 5000000
```

## Notes

- `scripts/staging/.env` is used for staging defaults.
- Commands that need EDS now default to the hosted staging EDS URL. Use `STAGING_CANTON_EDS_URL` or `-eds-url` only if you want to override it.
- `send_staging_evm_to_canton` becomes a token transfer only when `-token-amount` is set.
