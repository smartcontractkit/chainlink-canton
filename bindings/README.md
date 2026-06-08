# Go bindings & contract releases

Generated Go bindings and DAR artifacts share a **dev** vs **frozen release** split for DARs. Bindings are **not** versioned per release in this repo — all in-repo code uses `latest/`.

| Directory | Purpose | Updated by |
|-----------|---------|------------|
| `contracts/dars/current/` | Dev DARs (`*-current.dar`) | `make compile-contracts` / `make contracts` |
| `bindings/generated/latest/` | Dev bindings (sole in-repo import path) | `make generate-bindings` / `make contracts` |
| `contracts/dars/v1_0_0/`, `v2_0_0/`, … | Frozen release DARs | `make freeze-release VERSION=…` |
| `bindings/generated/v1_0_0/` | Legacy frozen bindings (historical only) | Do not add new `v*_*/` binding trees |
| `bindings/generated/mcms/` | Legacy MCMS SDK import path (single file) | `make freeze-release VERSION=…` when `latest/mcms/mcms.go` exists |

**Rules**

- In-repo application code (**deployment**, **integration-tests**, **EDS**, **CCIP devenv**, **party-ceremony** Go) imports **`bindings/generated/latest/...`** and must stay aligned with **`contracts/dars/current/`** (run `make contracts` after DAML changes).
- Frozen **`contracts/dars/v*_*/`** snapshots are for audited releases and components that load pinned DARs (e.g. **party-ceremony** production ceremony via `ReleaseDir`). Downstream repos pin a **chainlink-canton git tag/SHA** and use `bindings/generated/latest/` at that commit — not a versioned bindings folder.
- `latest/` and `current/` are regenerated during development and are not guaranteed stable or audited until tagged.
- **Do not rewrite** existing files under `contracts/dars/v*_*/` or legacy `bindings/generated/v1_0_0/` in normal PRs. CI blocks **modifications** to frozen paths that already exist on the base branch (`make check-frozen-release-artifacts`). After DAML changes, run `make contracts` and commit only `current/` + `latest/`.
- Frozen DAR snapshots change only via **`make freeze-release VERSION=…`** on a dedicated release PR. Add the GitHub label **`release-artifacts`** so CI allows those paths to change.
- Git tags for releases: `contracts-canton-v<x.y.z>` (e.g. `contracts-canton-v2.0.0`).
- Snapshot folder names use underscores: release `2.0.0` → `v2_0_0/`.

There are two notions of “version”:

1. **Release snapshot** — whole audited DAR cut (`v1_0_0`, `v2_0_0`, …). Published for external pin; in-repo Go always uses `latest/` at the tagged commit.
2. **Package semver** — per-package `version:` in each `daml.yaml` (e.g. `mcms-core-1.0.0.dar`, `ccip-core-2.0.0.dar`). Multiple package versions can live in the **same** release folder.

---

## Day-to-day development

From the repo root:

```bash
make contracts
```

Equivalent to:

```bash
make compile-contracts   # → contracts/dars/current/
make generate-bindings   # → bindings/generated/latest/
```

- Bindings are generated from **`current`** DARs (`contracts.GetDar(pkg, "current")`).
- `make contracts` does **not** update `contracts/dars/v*_*/` or `bindings/generated/mcms/`.

Use `contracts.GetDar(pkg, "current")` in dev tooling; use a pinned package version string (e.g. `"2.0.0"`) with `GetDar` for release-aligned behavior once `ReleaseDir` points at that snapshot.

---

## Releasing a new version (e.g. 2.0.0)

Example: cut release **2.0.0** after DAML/bindings work on `main` or your branch.

### 1. Build dev artifacts

```bash
make contracts
```

Commit any intentional DAML / `current` / `latest` changes from active development.

### 2. Freeze the DAR snapshot

```bash
make freeze-release VERSION=2.0.0
```

This runs `contracts/scripts/freeze-release.sh`, which:

- Copies every `daml.yaml` package from `dars/current/{name}-current.dar` → `dars/v2_0_0/{name}-{version}.dar` (version from each `daml.yaml`)
- Optionally copies `latest/mcms/mcms.go` → `bindings/generated/mcms/mcms.go` (external [mcms SDK](https://github.com/smartcontractkit/mcms) Canton import path)

It does **not** create `bindings/generated/v2_0_0/`. Bindings at the release commit are whatever is in `latest/`.

Review the diff, then commit the new `contracts/dars/v2_0_0/` tree (and `bindings/generated/mcms/mcms.go` if updated).

### 3. Point DAR consumers at the new release (manual)

`freeze-release` does **not** switch consumers automatically. Update:

| Location | Change |
|----------|--------|
| `contracts/contracts.go` | `ReleaseDir = "v2_0_0"` |
| `party-ceremony/ceremony/ops/ledger/deps.go` | `releaseDir = "v2_0_0"` (must match `ReleaseDir`) |

**Do not** migrate Go imports to `bindings/generated/v2_0_0/` — keep `bindings/generated/latest/`.

### 4. Verify

```bash
make contracts
go test ./contracts/...
go test ./deployment/...
# other modules as needed
```

`contracts` package tests iterate embedded DARs under `ReleaseDir`; after bumping `ReleaseDir`, embedded `v2_0_0/` must contain all packages listed in `contracts.Versions`.

### 5. Tag and push

```bash
git tag contracts-canton-v2.0.0
git push origin contracts-canton-v2.0.0
```

Open a PR with the freeze commit and `ReleaseDir` bump. Add label **`release-artifacts`**.

**Note:** Older DAR snapshots (`v1_0_0/`, …) stay in the tree so existing deployments can keep using them. Legacy `bindings/generated/v1_0_0/` remains for historical reference only.

---

## Legacy `bindings/generated/mcms/`

MCMS spans `mcms-api` and `mcms-core` DAML packages; generated code may also provide a unified facade at `latest/mcms/mcms.go`.

The flat path `github.com/smartcontractkit/chainlink-canton/bindings/generated/mcms` exists for **`github.com/smartcontractkit/mcms/sdk/canton`**, which predates versioned bindings. Only `mcms.go` lives there; it is refreshed on **`make freeze-release`** when present in `latest/`, not on `make contracts`.

In-repo Canton code should use `bindings/generated/latest/mcms`, not the flat shim.

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

# New release 2.0.0
make contracts
make freeze-release VERSION=2.0.0
# → edit ReleaseDir + party-ceremony releaseDir
git tag contracts-canton-v2.0.0
```
