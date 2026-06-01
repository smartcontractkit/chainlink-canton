#!/usr/bin/env bash
#
# Fail if a changeset modifies existing frozen release artifacts. Normal PRs
# should only touch dev outputs (contracts/dars/current/, bindings/generated/latest/).
#
# Allowed without the release-artifacts label:
#   - Adding a new release snapshot tree (e.g. first v1_0_0/ layout or v1_1_0/)
#   - Renaming/moving legacy flat DARs into contracts/dars/v*_*/
#   - Syncing bindings/generated/mcms/mcms.go when bootstrapping v1_0_0/mcms/
#
# Blocked:
#   - Modifying or deleting a file that already exists on the base branch under
#     a frozen path (contracts/dars/v*_*/, bindings/generated/v*_*/, mcms/)
#
# Bypass: FROZEN_ARTIFACT_CHANGES_OK=1 or PR label release-artifacts (see CI).
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

# True if path exists in the merge-base tree.
exists_on_base() {
  git cat-file -e "$MERGE_BASE:$1" 2>/dev/null
}

violations=()
bootstrapping_v1_mcms=false
if git diff --name-status "$MERGE_BASE"...HEAD | grep -qE '^A[[:space:]]+bindings/generated/v[0-9]+_[0-9]+_[0-9]+/mcms/mcms\.go'; then
  bootstrapping_v1_mcms=true
fi

while IFS=$'\t' read -r status path1 path2; do
  [ -z "$status" ] && continue

  case "$status" in
    M*)
      if is_frozen_path "$path1" && exists_on_base "$path1"; then
        if [ "$path1" = "bindings/generated/mcms/mcms.go" ] && $bootstrapping_v1_mcms; then
          continue
        fi
        violations+=("$path1 (modified)")
      fi
      ;;
    D*)
      if is_frozen_path "$path1" && exists_on_base "$path1"; then
        violations+=("$path1 (deleted)")
      fi
      ;;
    A*)
      # New files under frozen dirs are OK (new snapshot dir or new package version).
      ;;
    R*|C*)
      # Rename/copy: block only when both sides are frozen (in-place snapshot rewrite).
      if is_frozen_path "$path1" && is_frozen_path "$path2" && exists_on_base "$path1"; then
        violations+=("$path2 (renamed in frozen tree)")
      fi
      ;;
    *)
      if is_frozen_path "$path1" && exists_on_base "$path1"; then
        violations+=("$path1 ($status)")
      fi
      ;;
  esac
done < <(git diff --name-status -M "$MERGE_BASE"...HEAD)

if [ "${#violations[@]}" -eq 0 ]; then
  echo "No existing frozen release artifacts modified (base: $MERGE_BASE)."
  exit 0
fi

echo "ERROR: This change modifies frozen release artifacts that already exist on the base branch." >&2
echo "" >&2
echo "The following changes are not allowed in normal development PRs:" >&2
for f in "${violations[@]}"; do
  echo "  - $f" >&2
done
echo "" >&2
echo "Day-to-day contract work should only update:" >&2
echo "  - contracts/dars/current/     (make compile-contracts)" >&2
echo "  - bindings/generated/latest/  (make generate-bindings)" >&2
echo "" >&2
echo "To refresh an existing frozen snapshot (e.g. rewrite ccip-core-1.0.0.dar in v1_0_0/):" >&2
echo "  1. make contracts && make freeze-release VERSION=<x.y.z>" >&2
echo "  2. Open a PR with label: release-artifacts" >&2
echo "     (or set FROZEN_ARTIFACT_CHANGES_OK=1 locally)" >&2
echo "" >&2
echo "See bindings/README.md for the full release workflow." >&2
exit 1
