# Staging cv0–cv3: archive and remove stale DARs

## Source of truth (desired DARs)

These name+version pairs are **kept**; everything else under `ccip-*`, `mcms*`, `chainlink-api`, or `link` is **removed**:

| Package | Version(s) |
| --- | --- |
| chainlink-api | 2.0.0 |
| mcms-api | 1.0.0 |
| link | 2.0.0 |
| ccip-core | 2.0.0 |
| ccip-extension-api | 2.0.0 |
| ccip-runtime | 2.0.0 |
| ccip-sender | 2.0.0 |
| ccip-receiver | 2.0.0 |
| ccip-executor | 2.0.0 |
| ccip-committee-verifier | 2.0.0 |
| ccip-lock-release-token-pool | 2.0.0 |
| ccip-burn-mint-token-pool | 2.0.0 |
| ccip-factory | 2.0.0 |
| mcms-core | 1.0.0, 2.0.0 |

Examples of **stale** (REMOVE): `mcms@0.0.1`, `ccip-factory@0.0.1`, old `ccip-common@0.0.4`, duplicate factory package IDs.

Splice and platform DARs (utility-*, canton-builtin, etc.) are hidden from list output and never removed.

## List DARs

```bash
go run ./scripts/cleanup_staging_cv_dars --list-dars --node cv1
```

Actions: `DESIRED` | `REMOVE` | `KEEP` (platform/other). Also prints `Missing desired:` for packages not yet installed.

## Cleanup

```bash
go run ./scripts/cleanup_staging_cv_dars --dry-run --node all
go run ./scripts/cleanup_staging_cv_dars --node all
```

## Tokens

See `desiredDARVersions` in `main.go`. Each node needs `CV0_TOKEN` … `CV3_TOKEN` from its Okta client.

## Archive workflow (one terminal per node)

**Staging DevNet:** RMNRemote was deployed with `rmnOwner = ccipOwner` (not a separate `rmnOwner::` party). All CCIP contracts including RMNRemote archive as **decentralized ccipOwner** (cv1 prepare + cv1/cv3 vault sign). Simple JWT `SubmitAndWait` will not work.

1. **cv1 terminal** — dry-run all ccipOwner contracts (includes RMNRemote):
   `go run ./scripts/cleanup_staging_cv_dars --dry-run --skip-remove --node cv1`
2. **Multiparty archive** — party-ceremony InteractiveSubmission (not yet in cleanup script).
3. **Each terminal** — remove stale DARs after ACS is empty:
   `go run ./scripts/cleanup_staging_cv_dars --node cv0` (etc.)

Optional: `--archive-rmn-only --node cv3` only re-runs the RMNRemote template dry-run on cv3 (same signatory).

## Troubleshooting

### `NO_SYNCHRONIZER_ON_WHICH_ALL_SUBMITTERS_CAN_SUBMIT`

**If using `archive_active_canton_contracts` with `ActAs: ccipOwner`:** expected on staging — decentralized party is CONFIRMATION-only on cv1/cv3; don't fix with `CanActAs ccipOwner`.

**If using MCMS (`mcms-tools`):** submit as **local primary party** (`ccipBootstrapOwner` on cv1), `readAs` ccipOwner, plus contract disclosures. See `factoryDeployAndSetOwnerTOMCMs.md` (Nahuel's actAs/readAs split). JWT: `CanActAs` bootstrap, `CanReadAs` ccipOwner — not the reverse.

Parties on staging DevNet:

- `ccipOwner::1220644bd9e52834e8fba90d4607beed37b65991cc2b5377d5d40d07d3db36d4ed51` (decentralized; readAs only for MCMS)
- `ccipBootstrapOwner::1220a9854ea6590622988af59864d2b1588e004ac9850c140761f1038dd937e8f88d` (cv1 local actAs for MCMS submit)

Until `canActAs` is granted, DAR removal will fail with "blocked by active contracts".
