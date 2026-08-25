# Go bindings & contract releases

While the DAR artifacts are split in **dev** and **released** artifacts, bindings are always generated for the **dev**
version of contracts.

External consumers should always use versioned releases for `contracts/v<x.y.z>`.
See [Releases](https://github.com/smartcontractkit/chainlink-canton/releases?q=contracts%2F&expanded=true) for the
latest available version.

| Directory                       | Purpose                             | Updated by                                               |
|---------------------------------|-------------------------------------|----------------------------------------------------------|
| `contracts/dars/dev/`           | Dev DARs (`*-dev.dar`)              | `make compile-contracts` / `make contracts`              |
| `contracts/dars/released/`      | All released DARs, grown additively | `make contracts` (when version is in `ReleasedVersions`) |
| `contracts/bindings/generated/` | Go bindings                         | `make generate-bindings` / `make contracts`              |

**Rules**

- In-repo application code (**deployment**, **integration-tests**, **EDS**, **CCIP devenv**, **party-ceremony** Go)
  imports **`contracts/bindings/generated/...`** and uses a `go.mod` replace directive to always point to the latest
  available version.
- `contracts/dars/released/` is **append-only**: new versioned DAR files are added alongside existing ones; existing
  files are never modified. CI blocks modifications to already-committed artifacts there.
- `dev/` DARs are regenerated during development and are not guaranteed to be stable or audited until tagged.
- Consumers should always use a released version of both the contract artifacts + bindings.
  See [Releases](https://github.com/smartcontractkit/chainlink-canton/releases?q=contracts%2F&expanded=true) for the
  latest available version.
- Git tags for releases follow `contracts/v<x.y.z>` (e.g. `contracts/v2.0.0`).
- **Do not rewrite** existing files under `contracts/dars/released/` in normal PRs. After DAML changes, run
  `make contracts` and commit only `dev/` (and `released/` only when intentionally publishing a new version).

There are two notions of "version":

1. **Package semver** — per-package `version:` in each `daml.yaml` (e.g. `mcms-core-2.0.0.dar`). Multiple package
   versions can coexist in `contracts/dars/released/`.
2. **Repo release tag** — a `contracts/v<x.y.z>` git tag marking a stable, audited snapshot of the whole repo at a
   point in time.

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
- `make contracts` does **not** touch `contracts/dars/released/` unless a package's current version is listed in
  `ReleasedVersions` in `contracts/contracts.go`.

---

## Releasing a package version

### Staging changes without releasing

You can bump a package's version in its `daml.yaml` **without** immediately committing a release artifact. This is
useful to accumulate multiple changes under one version before publishing.

1. Bump `version:` in the package's `daml.yaml`.
2. Update `upgrades:` in the package's `daml.yaml` to point at the previous released version in `dars/released/`
3. Run `make contracts`.

Result: only `dars/dev/{name}-dev.dar` is updated. Nothing is written to `dars/released/`. You can keep iterating.

### Publishing a release

Once changes are ready to ship:

1. **Bump the version** in the package's `daml.yaml` (if not already done).
2. **Update `upgrades:`** in the package's `daml.yaml` to point at the previous released version in `dars/released/`.
   (if not already done)
3. **Add the new version** to `ReleasedVersions` in `contracts/contracts.go`:
   ```go
   CCIPCoreV2: []string{"2.0.0", "3.0.0"},
   ```
4. **Run `make contracts`** — the new `{name}-{version}.dar` is written to `dars/released/` alongside existing versions.
5. **Commit** both the updated `contracts/contracts.go` and the new file (s) under `dars/released/`.
6. **Push** the commit and merge into `main`.
7. **Release-Please** will automatically create a release-PR that when merged will create a new GitHub release.

CI will fail if any already-committed file under `dars/released/` is modified — only additions are allowed.

---

## CI

- **Compile Contracts** / **Generate Bindings**: build `dev/` DARs and `latest/` bindings.
- **Test Contracts**: run against embedded `contracts/dars/` artifacts.
- **Frozen artifact check**: fails if any existing file under `contracts/dars/released/` is modified or deleted
  (additions are allowed).

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
git commit -m "feat: bump ccip-core-v2 to 3.0.0"
# This will release the contracts under a minor version bump.
# To bump to a specific version, include `Release-As: x.y.z` in the commit message.
```
