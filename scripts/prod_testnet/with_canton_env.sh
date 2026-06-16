#!/usr/bin/env bash
# Run prod_testnet Go scripts with a clean Canton auth env (avoids stale JWT / static overrides).
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
cd "$ROOT"

# Shell exports from debugging override scripts/prod_testnet/.env and cause PermissionDenied.
unset PROD_TESTNET_CANTON_AUTH_TYPE PROD_TESTNET_CANTON_JWT TOKEN 2>/dev/null || true

set -a
# shellcheck disable=SC1091
source "$ROOT/scripts/prod_testnet/.env"
set +a

exec go run -mod=mod "$@"
