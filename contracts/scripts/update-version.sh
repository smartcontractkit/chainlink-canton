#!/usr/bin/env bash
#
# Updates the contract version across all daml.yaml files, contracts.go, and DAR artifacts.
#
# Usage:
#   ./contracts/scripts/update-version.sh <old-version> <new-version>
#
# Example:
#   ./contracts/scripts/update-version.sh 0.0.1 1.0.0
#
set -euo pipefail

if [ $# -ne 2 ]; then
  echo "Usage: $0 <old-version> <new-version>"
  echo "Example: $0 0.0.1 1.0.0"
  exit 1
fi

OLD_VERSION="$1"
NEW_VERSION="$2"

if [ "$OLD_VERSION" = "$NEW_VERSION" ]; then
  echo "Old and new versions are the same ($OLD_VERSION). Nothing to do."
  exit 0
fi

# Resolve the repo root (script lives in contracts/scripts/)
REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
CONTRACTS_DIR="$REPO_ROOT/contracts"

echo "Updating contract version: $OLD_VERSION -> $NEW_VERSION"
echo "  Repo root: $REPO_ROOT"
echo ""

# ──────────────────────────────────────────────
# 1. Update 'version:' field in all daml.yaml files
# ──────────────────────────────────────────────
echo "Step 1: Updating 'version:' in daml.yaml files..."
count=0
while IFS= read -r -d '' file; do
  if grep -q "^version: ${OLD_VERSION}$" "$file"; then
    sed -i '' "s/^version: ${OLD_VERSION}$/version: ${NEW_VERSION}/" "$file"
    echo "  Updated: ${file#"$REPO_ROOT/"}"
    count=$((count + 1))
  fi
done < <(find "$CONTRACTS_DIR" -name "daml.yaml" -print0)
echo "  -> Updated $count daml.yaml file(s)"
echo ""

# ──────────────────────────────────────────────
# 2. Update .daml/dist/ cross-references in daml.yaml files
#    (Only internal packages, NOT dependencies/ folder)
# ──────────────────────────────────────────────
echo "Step 2: Updating .daml/dist/ cross-references in daml.yaml files..."
count=0
while IFS= read -r -d '' file; do
  # Match lines like:  - ../foo/.daml/dist/bar-0.0.1.dar
  # but NOT lines like: - ../../dependencies/splice-api-foo-1.0.0.dar
  if grep -q "\.daml/dist/.*-${OLD_VERSION}\.dar" "$file"; then
    sed -i '' "s/\(\.daml\/dist\/.*\)-${OLD_VERSION}\.dar/\1-${NEW_VERSION}.dar/g" "$file"
    echo "  Updated: ${file#"$REPO_ROOT/"}"
    count=$((count + 1))
  fi
done < <(find "$CONTRACTS_DIR" -name "daml.yaml" -print0)
echo "  -> Updated cross-references in $count daml.yaml file(s)"
echo ""

# ──────────────────────────────────────────────
# 3. Update the Versions map in contracts/contracts.go
#    (Only internal packages — skip SpliceApi* lines)
# ──────────────────────────────────────────────
CONTRACTS_GO="$CONTRACTS_DIR/contracts.go"
echo "Step 3: Updating Versions map in contracts.go..."
if [ -f "$CONTRACTS_GO" ]; then
  # Replace "OLD_VERSION" with "NEW_VERSION" but only on lines that also contain 'CurrentVersion'
  # (this targets internal packages which have {version, CurrentVersion}, not splice externals which have just {version})
  sed -i '' "/CurrentVersion/s/\"${OLD_VERSION}\"/\"${NEW_VERSION}\"/g" "$CONTRACTS_GO"
  echo "  Updated: contracts/contracts.go"
else
  echo "  WARNING: contracts/contracts.go not found, skipping"
fi
echo ""

# ──────────────────────────────────────────────
# 4. Rename DAR files in contracts/dars/<version-dir>/
# ──────────────────────────────────────────────
OLD_DIR="v${OLD_VERSION//./_}"
NEW_DIR="v${NEW_VERSION//./_}"
DARS_DIR="$CONTRACTS_DIR/dars"
echo "Step 4: Renaming DAR files in contracts/dars/${OLD_DIR}/..."
count=0
if [ -d "$DARS_DIR/$OLD_DIR" ]; then
  for dar in "$DARS_DIR/$OLD_DIR"/*-"${OLD_VERSION}".dar; do
    [ -e "$dar" ] || continue
    new_dar="${dar%-"${OLD_VERSION}".dar}-${NEW_VERSION}.dar"
    mv "$dar" "$new_dar"
    echo "  Renamed: $(basename "$dar") -> $(basename "$new_dar")"
    count=$((count + 1))
  done
  if [ "$OLD_DIR" != "$NEW_DIR" ] && [ -d "$DARS_DIR/$OLD_DIR" ]; then
    mkdir -p "$DARS_DIR/$NEW_DIR"
    mv "$DARS_DIR/$OLD_DIR"/* "$DARS_DIR/$NEW_DIR/" 2>/dev/null || true
    rmdir "$DARS_DIR/$OLD_DIR" 2>/dev/null || true
  fi
fi
echo "  -> Renamed $count DAR file(s)"
echo ""

echo "Done! Version updated from $OLD_VERSION to $NEW_VERSION."
echo ""
echo "Next steps:"
echo "  1. Rebuild contracts:  make compile-contracts"
echo "  2. Regenerate bindings: make generate-bindings"
echo "  3. Freeze DAR snapshot: make freeze-release VERSION=${NEW_VERSION}"
echo "  Or just run:           make contracts"
