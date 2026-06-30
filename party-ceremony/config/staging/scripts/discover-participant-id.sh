#!/usr/bin/env bash
set -euo pipefail

# Query the Canton participant UID for a staging DevNet node.
#
# participant_id in nodes/cvN.participant-config.json is the operator label (cv0–cv3).
# Ceremony --participants / --coordinator flags need the full UID printed here.
#
# Usage:
#   export CLIENT_ID='0oatssi17pdaLqfb65d7'
#   export CLIENT_SECRET='<from-vault>'
#   export CV1_TOKEN="$(canton-login --ci canton-devnet)"
#   ./discover-participant-id.sh cv1
#   ./discover-participant-id.sh cv1 --write   # persist participant_uid in config

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
CV="${1:?target required: cv0|cv1|cv2|cv3}"
WRITE=false
if [[ "${2:-}" == "--write" ]]; then
  WRITE=true
fi

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

CFG="${ROOT}/nodes/${CV}.participant-config.json"
if [[ ! -f "$CFG" ]]; then
  echo "missing config: $CFG" >&2
  exit 1
fi

TOKEN_ENV="$(token_env_for_cv "$CV")"
ADMIN_HOST="$(jq -r .admin_host "$CFG")"
ADMIN_PORT="$(jq -r .admin_port "$CFG")"
JWT="$(jq -r .admin_jwt "$CFG")"
if [[ -z "$JWT" || "$JWT" == "null" ]]; then
  JWT="${!TOKEN_ENV:-}"
fi
if [[ -z "$JWT" ]]; then
  echo "admin_jwt is empty and \$$TOKEN_ENV is not set. Example:" >&2
  echo "  export CLIENT_ID='0oatssi17pdaLqfb65d7'" >&2
  echo "  export CLIENT_SECRET='<secret>'" >&2
  echo "  export ${TOKEN_ENV}=\"\$(canton-login --ci canton-devnet)\"" >&2
  exit 1
fi

TARGET="${ADMIN_HOST}:${ADMIN_PORT}"
RESP="$(grpcurl -H "Authorization: Bearer $JWT" \
  "$TARGET" \
  com.digitalasset.canton.topology.admin.v30.IdentityInitializationService/GetId)"
PARTICIPANT_UID="$(echo "$RESP" | jq -r .uniqueIdentifier)"

if [[ -z "$PARTICIPANT_UID" || "$PARTICIPANT_UID" == "null" ]]; then
  echo "GetId returned no uniqueIdentifier:" >&2
  echo "$RESP" >&2
  exit 1
fi

echo "$PARTICIPANT_UID"

if $WRITE; then
  tmp="$(mktemp)"
  jq --arg uid "$PARTICIPANT_UID" '.participant_uid = $uid' "$CFG" > "$tmp"
  mv "$tmp" "$CFG"
  echo "updated participant_uid in $CFG" >&2
fi
