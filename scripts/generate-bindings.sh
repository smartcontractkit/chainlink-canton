#!/bin/bash
# Script to generate Go code from DAML contracts into coin and ccip packages
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
COIN_OUTPUT_DIR="$PROJECT_ROOT/generated/coin"
CCIP_OUTPUT_DIR="$PROJECT_ROOT/generated/ccip"

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

log_info "Coin code generation complete!"
log_info "Generated files in $COIN_OUTPUT_DIR:"
find "$COIN_OUTPUT_DIR" -name "*.go" -type f | while read -r file; do
    log_info "  - $(basename "$file")"
done

# ============================================
# CCIP Contract Generation
# ============================================

log_info ""
log_info "============================================"
log_info "Generating CCIP contracts..."
log_info "============================================"

# Create CCIP output directory
mkdir -p "$CCIP_OUTPUT_DIR"

# Function to build a CCIP package if needed
build_ccip_package() {
    local package_dir="$1"
    local package_name="$2"
    local version="$3"
    local dar_path="$PROJECT_ROOT/contracts/ccip/$package_dir/.daml/dist/$package_name-$version.dar"
    
    if [ ! -f "$dar_path" ]; then
        log_info "Building $package_name..."
        cd "$PROJECT_ROOT/contracts/ccip/$package_dir" || exit 1
        if ! dpm build 2>&1; then
            log_error "Failed to build $package_name"
            return 1
        fi
        cd "$PROJECT_ROOT" || exit 1
    else
        log_info "$package_name DAR already exists, skipping build"
    fi
}

# Build CCIP packages in dependency order
log_info "Building CCIP packages in dependency order..."

# 1. Build common first (no CCIP dependencies)
build_ccip_package "common" "ccip-common" "1.0.0" || log_warn "Failed to build ccip-common, but continuing..."

# 2. Build packages that only depend on common
build_ccip_package "feequoter" "ccip-feequoter" "1.0.0" || log_warn "Failed to build ccip-feequoter, but continuing..."
build_ccip_package "tokenAdminRegistry" "ccip-tokenadminregistry" "1.0.0" || log_warn "Failed to build ccip-tokenadminregistry, but continuing..."
build_ccip_package "ccipreceiver" "ccip-receiver" "1.0.0" || log_warn "Failed to build ccip-receiver, but continuing..."
build_ccip_package "ccvs" "ccip-committeeverifier" "1.0.0" || log_warn "Failed to build ccip-committeeverifier, but continuing..."

# 3. Build packages that depend on common + feequoter/tokenAdminRegistry
build_ccip_package "onramp" "ccip-onramp" "1.0.0" || log_warn "Failed to build ccip-onramp, but continuing..."
build_ccip_package "offramp" "ccip-offramp" "1.0.0" || log_warn "Failed to build ccip-offramp, but continuing..."

# 4. Build perpartyrouter (depends on everything)
build_ccip_package "perpartyrouter" "ccip-perpartyrouter" "1.0.0" || log_warn "Failed to build ccip-perpartyrouter, but continuing..."

# 5. Build pools packages
build_ccip_package "pools/interfaces" "ccip-tokenpool-interfaces" "1.0.0" || log_warn "Failed to build ccip-tokenpool-interfaces, but continuing..."
build_ccip_package "pools/lockReleaseTokenPool" "ccip-lockreleasetokenpool" "1.0.0" || log_warn "Failed to build ccip-lockreleasetokenpool, but continuing..."

# Generate Go code from CCIP DAR files
log_info "Generating Go code from CCIP DAR files..."

# List of CCIP packages to generate (in order)
# Format: package_dir:package_name:version:go_package_name
declare -a CCIP_PACKAGES=(
    "common:ccip-common:1.0.0:common"
    "feequoter:ccip-feequoter:1.0.0:feequoter"
    "tokenAdminRegistry:ccip-tokenadminregistry:1.0.0:tokenadminregistry"
    "ccipreceiver:ccip-receiver:1.0.0:ccipreceiver"
    "ccvs:ccip-committeeverifier:1.0.0:ccvs"
    "onramp:ccip-onramp:1.0.0:onramp"
    "offramp:ccip-offramp:1.0.0:offramp"
    "perpartyrouter:ccip-perpartyrouter:1.0.0:perpartyrouter"
    "pools/interfaces:ccip-tokenpool-interfaces:1.0.0:interfaces"
    "pools/lockReleaseTokenPool:ccip-lockreleasetokenpool:1.0.0:lockreleasetokenpool"
)

for package_info in "${CCIP_PACKAGES[@]}"; do
    IFS=':' read -r package_dir package_name version go_package_name <<< "$package_info"
    dar_path="$PROJECT_ROOT/contracts/ccip/$package_dir/.daml/dist/$package_name-$version.dar"
    
    if [ ! -f "$dar_path" ]; then
        log_error "DAR file not found: $dar_path"
        log_error "Please build the package first: cd contracts/ccip/$package_dir && dpm build"
        continue
    fi
    
    # Create subdirectory for this package
    package_output_dir="$CCIP_OUTPUT_DIR/$go_package_name"
    mkdir -p "$package_output_dir"
    
    log_info "Generating $package_name into $go_package_name/..."
    "$GODAML_BIN" \
        --dar "$dar_path" \
        --output "$package_output_dir" \
        --go_package "$go_package_name" 2>&1 | grep -v "command not found" | grep -v "Usage:" | grep -v "Flags:" | grep -v "Examples:" || {
        log_warn "godaml reported errors for $package_name, but continuing..."
    }
    
    # Clean up unnecessary standard library files for this package
    find "$package_output_dir" -name "daml_script_*.go" -type f -delete 2>/dev/null || true
    find "$package_output_dir" -name "daml_*.go" -type f ! -name "*.go" -exec sh -c 'if [ $(wc -l < "$1") -lt 20 ]; then rm "$1"; fi' _ {} \; 2>/dev/null || true
done

# Check if any .go files were generated
if ! find "$CCIP_OUTPUT_DIR" -name "*.go" -type f | grep -q .; then
    log_error "No Go files were generated for CCIP. Cannot continue."
    exit 1
fi

log_info "CCIP code generation complete!"
log_info "Generated packages in $CCIP_OUTPUT_DIR:"
for package_info in "${CCIP_PACKAGES[@]}"; do
    IFS=':' read -r package_dir package_name version go_package_name <<< "$package_info"
    package_output_dir="$CCIP_OUTPUT_DIR/$go_package_name"
    if [ -d "$package_output_dir" ] && find "$package_output_dir" -name "*.go" -type f | grep -q .; then
        log_info "  - $go_package_name/:"
        find "$package_output_dir" -name "*.go" -type f | while read -r file; do
            log_info "      $(basename "$file")"
        done
    fi
done

log_info ""
log_info "============================================"
log_info "All code generation complete!"
log_info "============================================"
