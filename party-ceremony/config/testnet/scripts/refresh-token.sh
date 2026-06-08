#!/usr/bin/env bash
set -euo pipefail

# Inject a canton-login JWT into node participant-config.json.
#
# Usage (one node at a time — each cv has its own Okta client):
#   export CNT_TOKEN=$(canton-login --ci canton-testnet)   # with that cv's CLIENT_ID/SECRET
#   ./refresh-token.sh cv0
#
# Or refresh all configs with the same token (only if intentional):
#   ./refresh-token.sh all

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TARGET="${1:?target required: cv0|cv1|cv2|cv3|all}"

if [[ -z "${CNT_TOKEN:-}" ]]; then
  echo "CNT_TOKEN is not set. Run: export CNT_TOKEN=\$(canton-login canton-testnet)" >&2
  exit 1
fi

refresh_one() {
  local cv="$1"
  local cfg="${ROOT}/nodes/${cv}.participant-config.json"
  if [[ ! -f "$cfg" ]]; then
    echo "missing config: $cfg" >&2
    exit 1
  fi
  local tmp
  tmp="$(mktemp)"
  jq --arg t "$CNT_TOKEN" '.admin_jwt = $t | .ledger_jwt = $t' "$cfg" > "$tmp"
  mv "$tmp" "$cfg"
  echo "updated token: $cfg"
}

if [[ "$TARGET" == "all" ]]; then
  for cv in cv0 cv1 cv2 cv3; do
    refresh_one "$cv"
  done
else
  refresh_one "$TARGET"
fi
