#!/usr/bin/env bash
# Remove stale CCIP/MCMS DARs on DevNet cv0–cv3.
# One terminal per node — each shell only needs that node's CV*_TOKEN.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$ROOT"

NODE="${1:-all}"
DRY_RUN=false

shift || true
while [[ $# -gt 0 ]]; do
  case "$1" in
    --dry-run) DRY_RUN=true; shift ;;
    cv0|cv1|cv2|cv3|all)
      NODE="$1"; shift ;;
    --skip-archive|--with-archive)
      echo "Note: $1 is set automatically per --node (archive only on cv1). Pass through to go run if needed." >&2
      shift ;;
    *)
      echo "Usage: $0 [cv0|cv1|cv2|cv3|all] [--dry-run]" >&2
      exit 2
      ;;
  esac
done

require_token() {
  local cv="$1"
  local var="CV${cv#cv}_TOKEN"
  if [[ -z "${!var:-}" ]]; then
    echo "Missing $var — login with that node's Okta client in this terminal." >&2
    exit 1
  fi
}

if [[ "$NODE" == "all" ]]; then
  for cv in cv0 cv1 cv2 cv3; do
    require_token "$cv"
  done
else
  require_token "$NODE"
fi

ARGS=(--node "$NODE")
$DRY_RUN && ARGS=(--dry-run "${ARGS[@]}")

echo "==> cleanup_staging_cv_dars ${ARGS[*]}"
go run ./scripts/cleanup_staging_cv_dars "${ARGS[@]}"

if $DRY_RUN; then
  echo ""
  echo "Dry run done. Apply with:"
  echo "  $0 $NODE"
fi
