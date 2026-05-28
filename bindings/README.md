# Go bindings

After a version is released, the corresponding bindings will be moved to a versioned directory (e.g., `v1_0_0`).
Contracts in the versioned directory are guaranteed to be stable and audited.

Any contract bindings in the `latest` directory are not guaranteed to be stable or audited, and are only intended for testing and development, not for mainnet use.

Versions will be tagged in git using the format `contracts-canton-v<x.y.z>`, e.g., `contracts-canton-v1.0.0`.
Versioned directories will be named using the format `v<x_y_z>`, e.g., `v1_0_0`.

## Workflow

```bash
# Regenerate dev bindings from current DARs
make generate-bindings

# Freeze DARs + bindings into a versioned snapshot after release
make freeze-release VERSION=1.0.0
```

Production deployment operations should import from a versioned directory (e.g. `bindings/generated/v1_0_0/...`).
Development code, integration tests, and changesets should import from `bindings/generated/latest/...`.

DAR artifacts follow the same layout under `contracts/dars/`:

```
contracts/dars/
├── current/     ← dev DARs (*-current.dar), rebuilt by make compile-contracts
└── v1_0_0/      ← frozen release DARs, created by make freeze-release
```

Use `contracts.GetDar(pkg, "current")` for dev and `contracts.GetDar(pkg, "1.0.0")` for pinned releases. Packages at other semvers (e.g. `globalconfig` at `2.0.0`) also live in the same release directory (`v1_0_0/`).
