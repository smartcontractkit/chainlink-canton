# Go bindings & contract releases

Generated Go bindings and DAR artifacts use the same layout: **dev** (`current/` + `latest/`) and **frozen release snapshots** (`v1_0_0/`, `v1_1_0/`, …).

| Directory | Purpose | Updated by |
|-----------|---------|------------|
| `contracts/dars/current/` | Dev DARs (`*-current.dar`) | `make compile-contracts` / `make contracts` |
| `bindings/generated/latest/` | Dev bindings | `make generate-bindings` / `make contracts` |
| `contracts/dars/v1_0_0/` | Frozen release DARs | `make freeze-release VERSION=…` |
| `bindings/generated/v1_0_0/` | Frozen release bindings | `make freeze-release VERSION=…` |
| `bindings/generated/mcms/` | Legacy MCMS SDK import path (single file) | `make freeze-release VERSION=…` only |

**Rules**

- Application code (deployment, integration-tests, EDS, CCIP devenv, etc.) imports **`bindings/generated/v1_0_0/...`** (or whichever release the repo is pinned to). **Do not** import `latest/`.
- `latest/` and `current/` are regenerated during development and are not guaranteed stable or audited.
- **Do not rewrite** existing files under `contracts/dars/v1_0_0/` or `bindings/generated/v1_0_0/` in normal PRs. CI blocks **modifications** to frozen paths that already exist on the base branch (`make check-frozen-release-artifacts`). After DAML changes, run `make contracts` and commit only `current/` + `latest/`.
- Frozen snapshots change only via **`make freeze-release VERSION=…`** on a dedicated release PR. Add the GitHub label **`release-artifacts`** so CI allows those paths to change.
- Git tags for releases: `contracts-canton-v<x.y.z>` (e.g. `contracts-canton-v1.0.0`).
- Snapshot folder names use underscores: release `1.1.0` → `v1_1_0/`.

There are two notions of “version”:

1. **Release snapshot** — whole audited cut (`v1_0_0`, `v1_1_0`). What production pins to.
2. **Package semver** — per-package `version:` in each `daml.yaml` (e.g. `mcms-core-1.0.0.dar`, `globalconfig-2.0.0.dar`). Multiple package versions can live in the **same** release folder.

---

## Day-to-day development

From the repo root:

```bash
make contracts
```

Equivalent to:

```bash
make compile-contracts   # → contracts/dars/current/
make generate-bindings # → bindings/generated/latest/
```

- Bindings are generated from **`current`** DARs (`contracts.GetDar(pkg, "current")`).
- `make contracts` does **not** update `v1_0_0/` or `bindings/generated/mcms/`.

Use `contracts.GetDar(pkg, "current")` in dev tooling; use a pinned package version string (e.g. `"1.0.0"`) with `GetDar` for release-aligned behavior once `ReleaseDir` points at that snapshot.

---

## Releasing a new version (e.g. 1.1.0)

Example: cut release **1.1.0** after DAML/bindings work on `main` or your branch.

### 1. Build dev artifacts

```bash
make contracts
```

Commit any intentional DAML / `current` / `latest` changes from active development.

### 2. Freeze the snapshot

```bash
make freeze-release VERSION=1.1.0
```

This runs `contracts/scripts/freeze-release.sh`, which:

- Copies every `daml.yaml` package from `dars/current/{name}-current.dar` → `dars/v1_1_0/{name}-{version}.dar` (version from each `daml.yaml`)
- Copies `bindings/generated/latest/` → `bindings/generated/v1_1_0/` and rewrites import paths inside those files
- Copies `v1_1_0/mcms/mcms.go` → `bindings/generated/mcms/mcms.go` (external [mcms SDK](https://github.com/smartcontractkit/mcms) Canton import path)

Review the diff, then commit the new `v1_1_0/` trees and updated `bindings/generated/mcms/mcms.go`.

### 3. Point the repo at the new release (manual)

`freeze-release` does **not** switch consumers automatically. Update:

| Location | Change |
|----------|--------|
| `contracts/contracts.go` | `ReleaseDir = "v1_1_0"` |
| `party-ceremony/ceremony/ops/ledger/deps.go` | `releaseDir = "v1_1_0"` (must match `ReleaseDir`) |
| Go imports across the repo | `bindings/generated/v1_0_0/...` → `bindings/generated/v1_1_0/...` |

Search helpers:

```bash
rg 'bindings/generated/v1_0_0' --glob '*.go'
rg 'dars/v1_0_0' 
```

Typical areas: `deployment/`, `integration-tests/`, `eds/`, `ccip/`, `testhelpers/`, `contracts/` (e.g. `instanceAddress.go`, `choiceContext.go`).

Update **deployment operation** imports and any `Version` / `TypeAndVersion` fields in `deployment/operations/` if the deployment framework contract version should match the new Canton release (see each op file).

### 4. Verify

```bash
make contracts          # refresh latest/ only; should not be required for v1_1_0 pin
go test ./contracts/...
go test ./deployment/...
# other modules as needed
```

`contracts` package tests iterate embedded DARs under `ReleaseDir`; after bumping `ReleaseDir`, embedded `v1_1_0/` must contain all packages listed in `contracts.Versions`.

### 5. Tag and push

```bash
git tag contracts-canton-v1.1.0
git push origin contracts-canton-v1.1.0
```

Open a PR with the freeze commit, `ReleaseDir` / import migration, and tag (if tags are applied from the release branch).

**Note:** Older snapshots (`v1_0_0/`) stay in the tree so existing deployments and branches can keep using them until migrated.

---

## Legacy `bindings/generated/mcms/`

MCMS spans `mcms-api` and `mcms-core` DAML packages; generated code also provides a unified facade at `v*/mcms/mcms.go`.

The flat path `github.com/smartcontractkit/chainlink-canton/bindings/generated/mcms` exists for **`github.com/smartcontractkit/mcms/sdk/canton`**, which predates versioned bindings. Only `mcms.go` lives there; it is refreshed on **`make freeze-release`**, not on `make contracts`.

In-repo Canton deployment code should use versioned imports (e.g. `bindings/generated/v1_0_0/mcms`), not the flat shim.

---

## CI

- **Compile Contracts** / **Generate Bindings**: build `current/` and `latest/`.
- **Test Contracts**: embedded `contracts/dars/` via `ReleaseDir`.
- **Party ceremony integration tests**: test-only packages use `current` (e.g. `test-test@current`); production ceremony DARs use the pinned release directory.

---

## Quick reference

```bash
# Dev loop
make contracts

# New release 1.1.0
make contracts
make freeze-release VERSION=1.1.0
# → edit ReleaseDir, party-ceremony releaseDir, imports v1_0_0 → v1_1_0
git tag contracts-canton-v1.1.0
```
