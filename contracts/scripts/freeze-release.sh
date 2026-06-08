#!/usr/bin/env bash
#
# Freezes a production release snapshot: DARs only.
#
#   contracts/dars/v<x_y_z>/              from dars/current/
#
# Go bindings are not versioned in this repo — all in-repo code imports
# bindings/generated/latest/. Pin consumers via the chainlink-canton git tag/SHA.
#
# Usage:
#   ./contracts/scripts/freeze-release.sh <release-version>
#
# Example:
#   ./contracts/scripts/freeze-release.sh 2.0.0
#
set -euo pipefail

if [ $# -ne 1 ]; then
  echo "Usage: $0 <release-version>"
  echo "Example: $0 2.0.0"
  exit 1
fi

RELEASE_VERSION="$1"
SNAPSHOT="v${RELEASE_VERSION//./_}"

REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
CONTRACTS_DIR="$REPO_ROOT/contracts"
CURRENT_DARS_DIR="$CONTRACTS_DIR/dars/current"
DARS_SNAPSHOT_DIR="$CONTRACTS_DIR/dars/$SNAPSHOT"
LATEST_BINDINGS_DIR="$REPO_ROOT/bindings/generated/latest"

echo "Freezing release $RELEASE_VERSION -> contracts/dars/$SNAPSHOT/"

# ── DARs ──────────────────────────────────────────────────────────────────────
if [ ! -d "$CURRENT_DARS_DIR" ]; then
  echo "Missing $CURRENT_DARS_DIR — run 'make compile-contracts' first"
  exit 1
fi

mkdir -p "$DARS_SNAPSHOT_DIR"

dar_count=0
while IFS= read -r -d '' damlYaml; do
  pkgVersion="$(grep '^version:' "$damlYaml" | awk '{print $2}')"
  pkgName="$(grep '^name:' "$damlYaml" | awk '{print $2}')"
  src="$CURRENT_DARS_DIR/${pkgName}-current.dar"
  dst="$DARS_SNAPSHOT_DIR/${pkgName}-${pkgVersion}.dar"
  if [ ! -f "$src" ]; then
    echo "Missing current DAR for package $pkgName: $src"
    exit 1
  fi
  cp "$src" "$dst"
  echo "  dars: ${pkgName}-${pkgVersion}.dar"
  dar_count=$((dar_count + 1))
done < <(find "$CONTRACTS_DIR" -name daml.yaml ! -path '*/dependencies/*' -print0)

if [ "$dar_count" -eq 0 ]; then
  echo "No packages found to snapshot"
  exit 1
fi

echo "Snapshotted $dar_count DAR(s) to contracts/dars/$SNAPSHOT/"

# ── Legacy MCMS SDK shim (optional) ───────────────────────────────────────────
LEGACY_MCMS_SRC="$LATEST_BINDINGS_DIR/mcms/mcms.go"
LEGACY_MCMS="$REPO_ROOT/bindings/generated/mcms/mcms.go"
if [ -f "$LEGACY_MCMS_SRC" ]; then
  mkdir -p "$(dirname "$LEGACY_MCMS")"
  cp "$LEGACY_MCMS_SRC" "$LEGACY_MCMS"
  echo "Updated legacy import path bindings/generated/mcms/ for mcms SDK compatibility"
fi

echo "Done. Bindings: use bindings/generated/latest/ at this commit (not snapshotted)."
