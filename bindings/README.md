# Go bindings & contract releases

Generated Go bindings and DAR artifacts share a **dev** vs **released** split. Bindings are **not** versioned per release in this repo — all in-repo code uses `latest/`.

| Directory | Purpose | Updated by |
|-----------|---------|------------|
| `contracts/dars/dev/` | Dev DARs (`*-dev.dar`) | `make compile-contracts` / `make contracts` |
| `contracts/dars/released/` | All released DARs, grown additively | `make contracts` (when version is in `ReleasedVersions`) |
| `bindings/generated/latest/` | Dev bindings (sole in-repo import path) | `make generate-bindings` / `make contracts` |
| `bindings/generated/v1_0_0/` | **Deprecated** — legacy frozen bindings (historical only) | Do not use; see note below |

**Rules**

- In-repo application code (**deployment**, **integration-tests**, **EDS**, **CCIP devenv**, **party-ceremony** Go) imports **`bindings/generated/latest/...`** and must stay aligned with **`contracts/dars/dev/`** (run `make contracts` after DAML changes).
- `contracts/dars/released/` is **append-only**: new versioned DAR files are added alongside existing ones; existing files are never modified. CI blocks modifications to already-committed artifacts there.
- `latest/` and `dev/` DARs are regenerated during development and are not guaranteed stable or audited until tagged.
- **Do not rewrite** existing files under `contracts/dars/released/` in normal PRs. After DAML changes, run `make contracts` and commit only `dev/` + `latest/` (and `released/` only when intentionally publishing a new version).
- Git tags for releases follow `contracts-canton-v<x.y.z>` (e.g. `contracts-canton-v2.0.0`).

There are two notions of "version":

1. **Package semver** — per-package `version:` in each `daml.yaml` (e.g. `mcms-core-2.0.0.dar`). Multiple package versions can coexist in `contracts/dars/released/`.
2. **Repo release tag** — a `contracts-canton-vX.Y.Z` git tag marking a stable, audited snapshot of the whole repo at a point in time.

---

## Day-to-day development

From the repo root:

```bash
make contracts
```

Equivalent to:

```bash
make compile-contracts   # → contracts/dars/dev/
make generate-bindings   # → bindings/generated/latest/
```

- This updates the dev DARs and regenerates bindings from them.
- `make contracts` does **not** touch `contracts/dars/released/` unless a package's current version is listed in `ReleasedVersions` in `contracts/contracts.go`.

---

## Releasing a package version

### Staging changes without releasing

You can bump a package's version in its `daml.yaml` **without** immediately committing a release artifact. This is useful to accumulate multiple changes under one version before publishing.

1. Bump `version:` in the package's `daml.yaml`.
2. Update `upgrades:` in the package's `daml.yaml` to point at the previous released version in `dars/released/`
3. Run `make contracts`.

Result: only `dars/dev/{name}-dev.dar` is updated. Nothing is written to `dars/released/`. You can keep iterating.

### Publishing a release

Once changes are ready to ship:

1. **Bump the version** in the package's `daml.yaml` (if not already done).
2. **Update `upgrades:`** in the package's `daml.yaml` to point at the previous released version in `dars/released/`. (if not already done)
3. **Add the new version** to `ReleasedVersions` in `contracts/contracts.go`:
   ```go
   CCIPCoreV2: []string{"2.0.0", "3.0.0"},
   ```
4. **Run `make contracts`** — the new `{name}-{version}.dar` is written to `dars/released/` alongside existing versions.
5. **Commit** both the updated `contracts/contracts.go` and the new file(s) under `dars/released/`.
6. **Tag and push**:
   ```bash
   git tag contracts-canton-v3.0.0
   git push origin contracts-canton-v3.0.0
   ```

CI will fail if any already-committed file under `dars/released/` is modified — only additions are allowed.

---

## CI

- **Compile Contracts** / **Generate Bindings**: build `dev/` DARs and `latest/` bindings.
- **Test Contracts**: run against embedded `contracts/dars/` artifacts.
- **Frozen artifact check**: fails if any existing file under `contracts/dars/released/` is modified or deleted (additions are allowed).

---

## Legacy: `bindings/generated/v1_0_0/`

> ⚠️ **Deprecated.** The `bindings/generated/v1_0_0/` directory is a historical snapshot and is no longer maintained. Do not rely on it. Migrate any remaining consumers to `bindings/generated/latest/`.

---


## Quick reference

```bash
# Dev loop — update dev DARs and latest bindings
make contracts

# Stage a version bump without releasing
# (edit daml.yaml version only, do NOT add to ReleasedVersions yet)
make contracts
# → only dars/dev/ is updated

# Publish a release
# 1. Edit daml.yaml: bump version, update upgrades: field
# 2. Add version to ReleasedVersions in contracts/contracts.go
make contracts
# → dars/released/{name}-{version}.dar added
git add contracts/contracts.go contracts/dars/released/
git commit -m "release: bump ccip-core-v2 to 3.0.0"
git tag contracts-canton-v3.0.0
git push origin contracts-canton-v3.0.0
```
