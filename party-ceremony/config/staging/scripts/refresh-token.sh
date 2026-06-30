#!/usr/bin/env bash
set -euo pipefail

# Inject a canton-login JWT into staging DevNet node participant-config.json.
#
# Usage (one node at a time — each cv has its own Okta client):
#   export CLIENT_ID='0oatssi17pdaLqfb65d7'
#   export CLIENT_SECRET='<from-vault>'
#   export CV1_TOKEN="$(canton-login --ci canton-devnet)"
#   ./refresh-token.sh cv1
#
# Refresh all (each CVN_TOKEN must be set separately):
#   ./refresh-token.sh all

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TARGET="${1:?target required: cv0|cv1|cv2|cv3|all}"

token_env_for_cv() {
  local cv="$1"
  case "$cv" in
    cv0) echo CV0_TOKEN ;;
    cv1) echo CV1_TOKEN ;;
    cv2) echo CV2_TOKEN ;;
    cv3) echo CV3_TOKEN ;;
    *) echo "unknown cv: $cv" >&2; return 1 ;;
  esac
}

refresh_one() {
  local cv="$1"
  local token_env
  token_env="$(token_env_for_cv "$cv")"
  local token="${!token_env:-}"
  if [[ -z "$token" ]]; then
    echo "\$$token_env is not set. Example:" >&2
    echo "  export CLIENT_ID='<okta-client-id>'" >&2
    echo "  export CLIENT_SECRET='<secret>'" >&2
    echo "  export ${token_env}=\"\$(canton-login --ci canton-devnet)\"" >&2
    exit 1
  fi

  local cfg="${ROOT}/nodes/${cv}.participant-config.json"
  if [[ ! -f "$cfg" ]]; then
    echo "missing config: $cfg" >&2
    exit 1
  fi
  local tmp
  tmp="$(mktemp)"
  jq --arg t "$token" '.admin_jwt = $t | .ledger_jwt = $t' "$cfg" > "$tmp"
  mv "$tmp" "$cfg"
  echo "updated token: $cfg (from \$$token_env)"
}

if [[ "$TARGET" == "all" ]]; then
  for cv in cv0 cv1 cv2 cv3; do
    refresh_one "$cv"
  done
else
  refresh_one "$TARGET"
fi
