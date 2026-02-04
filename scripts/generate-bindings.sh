#!/bin/bash
# Script to generate Go code from DAML contracts into generated/{coin,ccip,mcms}
set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

log_info() { echo -e "${GREEN}[INFO]${NC} $1"; }
log_warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }
die() { log_error "$1"; exit 1; }

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

GEN_DIR="$PROJECT_ROOT/bindings"
COIN_OUTPUT_DIR="$GEN_DIR/coin"
CCIP_OUTPUT_DIR="$GEN_DIR/ccip"
MCMS_OUTPUT_DIR="$GEN_DIR/mcms"

# DARs are exported by: go run ./bindings/compile/cmd/export-dars -out artifacts/dars
DAR_DIR="$PROJECT_ROOT/contracts/dars"

# Find godaml
if command -v godaml >/dev/null 2>&1; then
  GODAML_BIN="godaml"
elif [ -f "$PROJECT_ROOT/bin/godaml" ]; then
  GODAML_BIN="$PROJECT_ROOT/bin/godaml"
else
  die "godaml not found in PATH or at $PROJECT_ROOT/bin/godaml"
fi
log_info "Using godaml: $GODAML_BIN"

# Ensure exported DARs exist
[ -d "$DAR_DIR" ] || die "DAR directory not found: $DAR_DIR (run: go run ./bindings/compile/cmd/export-dars -out artifacts/dars)"
log_info "Available DARs in $DAR_DIR:"
ls -1 "$DAR_DIR" | sed 's/^/  - /' || true

# Clean outputs (preserve custom bindings.go)
log_info "Cleaning generated output dirs..."
BINDINGS_GO="$GEN_DIR/bindings.go"
BINDINGS_BACKUP=""
if [ -f "$BINDINGS_GO" ]; then
  BINDINGS_BACKUP="$(mktemp)"
  cp "$BINDINGS_GO" "$BINDINGS_BACKUP"
  log_info "Preserved bindings.go"
fi

rm -rf "$GEN_DIR"
mkdir -p "$COIN_OUTPUT_DIR" "$CCIP_OUTPUT_DIR" "$MCMS_OUTPUT_DIR"

if [ -n "$BINDINGS_BACKUP" ] && [ -f "$BINDINGS_BACKUP" ]; then
  cp "$BINDINGS_BACKUP" "$BINDINGS_GO"
  rm -f "$BINDINGS_BACKUP"
  log_info "Restored bindings.go"
fi

# If required=1 and godaml fails -> exit. If required=0 -> warn and continue.
run_godaml() {
  local dar="$1"
  local out="$2"
  local pkg="$3"
  local required="${4:-1}"

  [ -f "$dar" ] || {
    if [ "$required" -eq 1 ]; then
      die "Missing DAR: $dar"
    else
      log_warn "Missing optional DAR: $dar (skipping)"
      return 0
    fi
  }

  mkdir -p "$out"
  log_info "godaml: $(basename "$dar") -> $out (go_package=$pkg)"

  local tmp
  tmp="$(mktemp)"

  set +e
  "$GODAML_BIN" --dar "$dar" --output "$out" --go_package "$pkg" >"$tmp" 2>&1
  local rc=$?
  set -e

  # Print filtered output; don't let grep affect script success
  grep -v -E "command not found|^Usage:|^Examples:|^Flags:" "$tmp" || true
  rm -f "$tmp"

  if [ "$rc" -ne 0 ]; then
    if [ "$required" -eq 1 ]; then
      die "godaml failed for $(basename "$dar") (rc=$rc)"
    else
      log_warn "godaml failed for $(basename "$dar") (rc=$rc); continuing..."
    fi
  fi

  # Cleanup tiny/empty stdlib junk
  find "$out" -name "daml_script_*.go" -type f -delete 2>/dev/null || true
  find "$out" -name "daml_*.go" -type f -exec sh -c 'if [ "$(wc -l < "$1")" -lt 20 ]; then rm -f "$1"; fi' _ {} \; 2>/dev/null || true
}

# -------------------------
# Coin deps (optional)
# -------------------------
log_info "Generating coin dependency types (optional)..."
run_godaml "$PROJECT_ROOT/contracts/dependencies/splice-api-token-holding-v1-1.0.0.dar" \
  "$COIN_OUTPUT_DIR" "coin" 0
run_godaml "$PROJECT_ROOT/contracts/dependencies/splice-api-token-metadata-v1-1.0.0.dar" \
  "$COIN_OUTPUT_DIR" "coin" 0
run_godaml "$PROJECT_ROOT/contracts/dependencies/splice-api-token-transfer-instruction-v1-1.0.0.dar" \
  "$COIN_OUTPUT_DIR" "coin" 0
run_godaml "$PROJECT_ROOT/contracts/dependencies/splice-api-token-burn-mint-v1-1.0.0.dar" \
  "$COIN_OUTPUT_DIR" "coin" 0

# -------------------------
# Coin (required)
# -------------------------
log_info "Generating coin bindings..."
run_godaml "$DAR_DIR/coin-current.dar" "$COIN_OUTPUT_DIR" "coin" 1
[ -n "$(find "$COIN_OUTPUT_DIR" -name "*.go" -type f -print -quit 2>/dev/null)" ] || die "No Go files generated for coin"

# -------------------------
# CCIP (required)
# -------------------------
log_info "Generating CCIP bindings..."
run_godaml "$DAR_DIR/ccip-common-current.dar"               "$CCIP_OUTPUT_DIR/common"               "common"               1
run_godaml "$DAR_DIR/ccip-feequoter-current.dar"            "$CCIP_OUTPUT_DIR/feequoter"            "feequoter"            1
run_godaml "$DAR_DIR/ccip-tokenadminregistry-current.dar"   "$CCIP_OUTPUT_DIR/tokenadminregistry"   "tokenadminregistry"   1
run_godaml "$DAR_DIR/ccip-receiver-current.dar"             "$CCIP_OUTPUT_DIR/ccipreceiver"         "ccipreceiver"         1
run_godaml "$DAR_DIR/ccip-committeeverifier-current.dar"    "$CCIP_OUTPUT_DIR/ccvs"                 "ccvs"                 1
run_godaml "$DAR_DIR/ccip-onramp-current.dar"               "$CCIP_OUTPUT_DIR/onramp"               "onramp"               1
run_godaml "$DAR_DIR/ccip-offramp-current.dar"              "$CCIP_OUTPUT_DIR/offramp"              "offramp"              1
run_godaml "$DAR_DIR/ccip-perpartyrouter-current.dar"       "$CCIP_OUTPUT_DIR/perpartyrouter"       "perpartyrouter"       1
run_godaml "$DAR_DIR/ccip-tokenpool-interfaces-current.dar" "$CCIP_OUTPUT_DIR/interfaces"           "interfaces"           1
run_godaml "$DAR_DIR/ccip-lockreleasetokenpool-current.dar" "$CCIP_OUTPUT_DIR/lockreleasetokenpool" "lockreleasetokenpool" 1
[ -n "$(find "$CCIP_OUTPUT_DIR" -name "*.go" -type f -print -quit 2>/dev/null)" ] || die "No Go files generated for CCIP"

# -------------------------
# MCMS (required)
# -------------------------
log_info "Generating MCMS bindings..."
run_godaml "$DAR_DIR/mcms-current.dar" "$MCMS_OUTPUT_DIR" "mcms" 1
[ -n "$(find "$MCMS_OUTPUT_DIR" -name "*.go" -type f -print -quit 2>/dev/null)" ] || die "No Go files generated for MCMS"

log_info "All code generation complete!"
