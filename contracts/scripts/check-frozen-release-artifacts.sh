#!/usr/bin/env bash
#
# Fail if a changeset modifies frozen release artifacts. Normal PRs should only
# touch dev outputs (contracts/dars/current/, bindings/generated/latest/).
#
# Frozen paths are updated deliberately via make freeze-release (use PR label
# release-artifacts to allow those changes in CI).
#
# Usage:
#   ./contracts/scripts/check-frozen-release-artifacts.sh [base-ref]
#
# Examples:
#   ./contracts/scripts/check-frozen-release-artifacts.sh origin/main
#   FROZEN_ARTIFACT_CHANGES_OK=1 ./contracts/scripts/check-frozen-release-artifacts.sh
#
set -euo pipefail

if [ "${FROZEN_ARTIFACT_CHANGES_OK:-}" = "1" ]; then
  echo "Skipping frozen-artifact check (FROZEN_ARTIFACT_CHANGES_OK=1)."
  exit 0
fi

REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
cd "$REPO_ROOT"

BASE_REF="${1:-}"
if [ -z "$BASE_REF" ]; then
  if [ -n "${GITHUB_BASE_REF:-}" ]; then
    git fetch origin "${GITHUB_BASE_REF}" --depth=1 2>/dev/null || true
    BASE_REF="origin/${GITHUB_BASE_REF}"
  else
    BASE_REF="origin/main"
  fi
fi

if ! git rev-parse --verify "$BASE_REF" >/dev/null 2>&1; then
  echo "Base ref $BASE_REF not found; skipping frozen-artifact check."
  exit 0
fi

MERGE_BASE="$(git merge-base HEAD "$BASE_REF" 2>/dev/null || echo "$BASE_REF")"

is_frozen_path() {
  local path="$1"
  case "$path" in
    contracts/dars/v[0-9]_*/*) return 0 ;;
    bindings/generated/v[0-9]_*/*) return 0 ;;
    bindings/generated/mcms/*) return 0 ;;
    *) return 1 ;;
  esac
}

violations=()
while IFS= read -r -d '' file; do
  [ -z "$file" ] && continue
  if is_frozen_path "$file"; then
    violations+=("$file")
  fi
done < <(git diff --name-only -z "$MERGE_BASE"...HEAD)

if [ "${#violations[@]}" -eq 0 ]; then
  echo "No frozen release artifacts modified (base: $MERGE_BASE)."
  exit 0
fi

echo "ERROR: This change modifies frozen release artifacts." >&2
echo "" >&2
echo "The following paths must not change in normal development PRs:" >&2
for f in "${violations[@]}"; do
  echo "  - $f" >&2
done
echo "" >&2
echo "Day-to-day contract work should only update:" >&2
echo "  - contracts/dars/current/     (make compile-contracts)" >&2
echo "  - bindings/generated/latest/  (make generate-bindings)" >&2
echo "" >&2
echo "To cut or refresh a release snapshot intentionally:" >&2
echo "  1. make contracts && make freeze-release VERSION=<x.y.z>" >&2
echo "  2. Open a PR with label: release-artifacts" >&2
echo "     (or set FROZEN_ARTIFACT_CHANGES_OK=1 locally)" >&2
echo "" >&2
echo "See bindings/README.md for the full release workflow." >&2
exit 1
