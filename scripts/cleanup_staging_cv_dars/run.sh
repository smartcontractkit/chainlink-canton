#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$ROOT"

MODE=""
NODE="all"

while [[ $# -gt 0 ]]; do
  case "$1" in
    --dry-run|--apply|--list-dars)
      MODE="$1"
      shift
      ;;
    --node)
      NODE="${2:?--node requires cv0, cv1, cv2, cv3, or all}"
      shift 2
      ;;
    *)
      echo "unknown argument: $1" >&2
      exit 2
      ;;
  esac
done

usage() {
  cat <<EOF
Usage:
  $0 --list-dars [--node cv0|cv1|cv2|cv3|all]
  $0 --dry-run  [--node cv0|cv1|cv2|cv3|all]
  $0 --apply    [--node cv0|cv1|cv2|cv3|all]

One terminal per node — only that node's token is required:

  cv0 terminal (CV0_TOKEN):  $0 --apply --node cv0
  cv1 terminal (CV1_TOKEN):  $0 --apply --node cv1   # archive + cv1 DAR removal
  cv2 terminal (CV2_TOKEN):  $0 --apply --node cv2
  cv3 terminal (CV3_TOKEN):  $0 --apply --node cv3

Archive (shared ledger, per terminal):
  cv1 terminal:  archive ccipOwner + factory contracts
  cv3 terminal:  optional RMNRemote dry-run (--archive-rmn-only; signatory is ccipOwner on staging)
Then remove stale DARs from each participant terminal (--node cv0|cv1|cv3).
EOF
}

[[ -n "$MODE" ]] || { usage; exit 2; }

valid_node() {
  case "$1" in cv0|cv1|cv2|cv3|all) return 0 ;; *) return 1 ;; esac
}
valid_node "$NODE" || { echo "invalid --node: $NODE" >&2; exit 2; }

if [[ "$NODE" == "all" ]]; then
  for cv in cv0 cv1 cv2 cv3; do
    var="CV${cv#cv}_TOKEN"
    if [[ -z "${!var:-}" ]]; then
      echo "${var} is not set (required for --node all)" >&2
      exit 1
    fi
  done
else
  var="CV${NODE#cv}_TOKEN"
  if [[ -z "${!var:-}" ]]; then
    echo "${var} is not set (required for --node $NODE in this terminal)" >&2
    exit 1
  fi
fi

ARGS=(--node "$NODE")
case "$MODE" in
  --list-dars) ARGS=(--list-dars "${ARGS[@]}") ;;
  --dry-run)   ARGS=(--dry-run "${ARGS[@]}") ;;
  --apply)     ;;
esac

go run ./scripts/cleanup_staging_cv_dars "${ARGS[@]}"
