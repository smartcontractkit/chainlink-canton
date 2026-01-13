#!/bin/bash
# Script to generate Go code from DAML contracts into a single coin package
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
COIN_OUTPUT_DIR="$PROJECT_ROOT/generated/coin"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if godaml is available in PATH
if ! command -v godaml &> /dev/null; then
    # Try local bin directory as fallback
    if [ -f "$PROJECT_ROOT/bin/godaml" ]; then
        GODAML_BIN="$PROJECT_ROOT/bin/godaml"
    else
        log_error "godaml not found in PATH or at $PROJECT_ROOT/bin/godaml"
        log_error "Please install godaml or add it to your PATH"
        exit 1
    fi
else
    GODAML_BIN="godaml"
fi

log_info "Using godaml: $GODAML_BIN"

# Create output directory
mkdir -p "$COIN_OUTPUT_DIR"

log_info "Generating all types into coin package..."

# Generate Splice API Token Holding (into coin package)
log_info "Generating splice-api-token-holding-v1..."
"$GODAML_BIN" \
    --dar "$PROJECT_ROOT/contracts/dependencies/splice-api-token-holding-v1-1.0.0.dar" \
    --output "$COIN_OUTPUT_DIR" \
    --go_package coin 2>&1 | grep -v "command not found" || {
    log_warn "Failed to generate splice-api-token-holding-v1, continuing..."
}

# Generate Splice API Token Metadata (into coin package)
log_info "Generating splice-api-token-metadata-v1..."
"$GODAML_BIN" \
    --dar "$PROJECT_ROOT/contracts/dependencies/splice-api-token-metadata-v1-1.0.0.dar" \
    --output "$COIN_OUTPUT_DIR" \
    --go_package coin 2>&1 | grep -v "command not found" || {
    log_warn "Failed to generate splice-api-token-metadata-v1, continuing..."
}

# Generate Splice API Token Transfer Instruction (into coin package)
log_info "Generating splice-api-token-transfer-instruction-v1..."
"$GODAML_BIN" \
    --dar "$PROJECT_ROOT/contracts/dependencies/splice-api-token-transfer-instruction-v1-1.0.0.dar" \
    --output "$COIN_OUTPUT_DIR" \
    --go_package coin 2>&1 | grep -v "command not found" || {
    log_warn "Failed to generate splice-api-token-transfer-instruction-v1, continuing..."
}

# Generate Splice API Token Burn Mint (into coin package)
# This MUST be generated before coin DAR because coin uses BurnMintOutput
log_info "Generating splice-api-token-burn-mint-v1..."
"$GODAML_BIN" \
    --dar "$PROJECT_ROOT/contracts/dependencies/splice-api-token-burn-mint-v1-1.0.0.dar" \
    --output "$COIN_OUTPUT_DIR" \
    --go_package coin 2>&1 | grep -v "command not found" || {
    log_warn "Failed to generate splice-api-token-burn-mint-v1, continuing..."
}

log_info "Generating Go code from coin DAR..."

# Generate coin package
COIN_DAR="$PROJECT_ROOT/contracts/coin/.daml/dist/coin-0.0.1.dar"
if [ ! -f "$COIN_DAR" ]; then
    log_error "Coin DAR file not found at $COIN_DAR"
    log_info "Please build the coin package first: cd contracts/coin && dpm build"
    exit 1
fi

# Generate coin package - continue even if it fails (we'll fix the errors)
log_warn "Note: godaml may report errors due to unresolved types, but we'll fix them..."
"$GODAML_BIN" \
    --dar "$COIN_DAR" \
    --output "$COIN_OUTPUT_DIR" \
    --go_package coin 2>&1 | grep -v "command not found" | grep -v "Usage:" | grep -v "Flags:" | grep -v "Examples:" || {
    log_warn "godaml reported errors, but continuing to fix generated code..."
}

# Check if any .go files were generated
if ! find "$COIN_OUTPUT_DIR" -name "*.go" -type f | grep -q .; then
    log_error "No Go files were generated. Cannot continue."
    exit 1
fi

log_info "Cleaning up generated files..."

# Remove unnecessary standard library files that godaml generates
# These are empty files from daml-script and other standard library dependencies
find "$COIN_OUTPUT_DIR" -name "daml_script_*.go" -type f -delete 2>/dev/null || true
find "$COIN_OUTPUT_DIR" -name "daml_*.go" -type f ! -name "*.go" -exec sh -c 'if [ $(wc -l < "$1") -lt 20 ]; then rm "$1"; fi' _ {} \; 2>/dev/null || true

log_info "Code generation complete!"
log_info "Generated files in $COIN_OUTPUT_DIR:"
find "$COIN_OUTPUT_DIR" -name "*.go" -type f | while read -r file; do
    log_info "  - $(basename "$file")"
done
