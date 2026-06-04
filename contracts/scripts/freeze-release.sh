#!/usr/bin/env bash
#
# Freezes a production release snapshot: DARs and Go bindings.
#
#   contracts/dars/v<x_y_z>/              from dars/current/
#   bindings/generated/v<x_y_z>/         from bindings/generated/latest/
#
# Usage:
#   ./contracts/scripts/freeze-release.sh <release-version>
#
# Example:
#   ./contracts/scripts/freeze-release.sh 1.0.0
#
set -euo pipefail

if [ $# -ne 1 ]; then
  echo "Usage: $0 <release-version>"
  echo "Example: $0 1.0.0"
  exit 1
fi

RELEASE_VERSION="$1"
SNAPSHOT="v${RELEASE_VERSION//./_}"

REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
CONTRACTS_DIR="$REPO_ROOT/contracts"
CURRENT_DARS_DIR="$CONTRACTS_DIR/dars/current"
DARS_SNAPSHOT_DIR="$CONTRACTS_DIR/dars/$SNAPSHOT"
LATEST_BINDINGS_DIR="$REPO_ROOT/bindings/generated/latest"
BINDINGS_SNAPSHOT_DIR="$REPO_ROOT/bindings/generated/$SNAPSHOT"

echo "Freezing release $RELEASE_VERSION -> $SNAPSHOT/"

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

# ── Bindings ──────────────────────────────────────────────────────────────────
if [ ! -d "$LATEST_BINDINGS_DIR" ]; then
  echo "Missing $LATEST_BINDINGS_DIR — run 'make generate-bindings' first"
  exit 1
fi

echo "Freezing bindings: latest/ -> $SNAPSHOT/"
rm -rf "$BINDINGS_SNAPSHOT_DIR"
cp -R "$LATEST_BINDINGS_DIR" "$BINDINGS_SNAPSHOT_DIR"

while IFS= read -r -d '' file; do
  sed -i '' "s|bindings/generated/latest|bindings/generated/$SNAPSHOT|g" "$file"
done < <(find "$BINDINGS_SNAPSHOT_DIR" -name '*.go' -print0)

echo "Snapshotted bindings to bindings/generated/$SNAPSHOT/"

LEGACY_MCMS="$REPO_ROOT/bindings/generated/mcms/mcms.go"
mkdir -p "$(dirname "$LEGACY_MCMS")"
cp "$BINDINGS_SNAPSHOT_DIR/mcms/mcms.go" "$LEGACY_MCMS"
echo "Updated legacy import path bindings/generated/mcms/ for mcms SDK compatibility"

echo "Done."
