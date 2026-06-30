#!/usr/bin/env bash
# Staging DevNet upgrade without archiving old contracts (lanes not live).
#
# Flow:
#   1. Upload v2 modular DARs alongside legacy packages (no remove).
#   2. Redeploy new contract stack via MCMS (chainlink-deployments + mcms-tools).
#   3. Disable old lanes on old GlobalConfig; enable on new stack (MCMS proposals).
#
# Old contracts + stale DARs may remain on ledger/package store until a later cleanup.
#
# Usage:
#   export CV0_TOKEN=... CV1_TOKEN=...  # per node, see refresh-token.sh
#   ./upgrade-skip-archive.sh list [--node cv0|cv1|cv2|cv3|all]
#   ./upgrade-skip-archive.sh upload [--node cv0|cv1|cv2|cv3|all]
#   ./upgrade-skip-archive.sh plan   # print MCMS redeploy pointers
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
REPO_ROOT="$(cd "${ROOT}/../../.." && pwd)"
PARTY_CEREMONY="$(cd "${ROOT}/../.." && pwd)"
DAR_DIR="${REPO_ROOT}/contracts/dars/current"

# Production CCIP/MCMS packages for staging (matches cleanup_staging_cv_dars desired set).
UPLOAD_DARS=(
  chainlink-api-current.dar
  mcms-api-current.dar
  link-current.dar
  ccip-core-current.dar
  ccip-extension-api-current.dar
  ccip-runtime-current.dar
  ccip-sender-current.dar
  ccip-receiver-current.dar
  ccip-executor-current.dar
  ccip-committee-verifier-current.dar
  ccip-lock-release-token-pool-current.dar
  ccip-burn-mint-token-pool-current.dar
  ccip-factory-current.dar
  mcms-core-current.dar
)

usage() {
  cat <<EOF
Usage:
  $0 list   [--node cv0|cv1|cv2|cv3|all]
  $0 upload [--node cv0|cv1|cv2|cv3|all]
  $0 plan

  list   — show installed DARs (no archive, no remove)
  upload — upload v2 modular DARs from contracts/dars/current (coexists with legacy)
  plan   — MCMS redeploy / lane-disable next steps

Requires CVN_TOKEN in env; refresh via ./refresh-token.sh cvN
EOF
}

cmd="${1:-}"
shift || true
NODE="all"

while [[ $# -gt 0 ]]; do
  case "$1" in
    --node)
      NODE="${2:?--node requires cv0, cv1, cv2, cv3, or all}"
      shift 2
      ;;
    *)
      echo "unknown argument: $1" >&2
      usage
      exit 2
      ;;
  esac
done

valid_node() {
  case "$1" in cv0|cv1|cv2|cv3|all) return 0 ;; *) return 1 ;; esac
}

nodes_for() {
  case "$1" in
    all) echo cv0 cv1 cv2 cv3 ;;
    cv0|cv1|cv2|cv3) echo "$1" ;;
    *) echo "invalid --node: $1" >&2; exit 2 ;;
  esac
}

require_tokens() {
  local n
  for n in $(nodes_for "$NODE"); do
    local var="CV${n#cv}_TOKEN"
    if [[ -z "${!var:-}" ]]; then
      echo "${var} is not set — run ./refresh-token.sh ${n} after canton-login" >&2
      exit 1
    fi
  done
}

run_list() {
  require_tokens
  cd "${REPO_ROOT}"
  go run ./scripts/cleanup_staging_cv_dars --list-dars --node "$NODE"
}

run_upload() {
  require_tokens
  if [[ ! -d "${DAR_DIR}" ]]; then
    echo "missing ${DAR_DIR} — run: make contracts" >&2
    exit 1
  fi
  for dar in "${UPLOAD_DARS[@]}"; do
    if [[ ! -f "${DAR_DIR}/${dar}" ]]; then
      echo "missing ${DAR_DIR}/${dar} — run: make contracts" >&2
      exit 1
    fi
  done

  local n
  for n in $(nodes_for "$NODE"); do
    local var="CV${n#cv}_TOKEN"
    export "${var?}"
    "${ROOT}/scripts/refresh-token.sh" "$n"
    echo ""
    echo "=== upload v2 DARs → ${n} ==="
    for dar in "${UPLOAD_DARS[@]}"; do
      echo "  uploading ${dar}..."
      (cd "${REPO_ROOT}/party-ceremony" && go run ./cmd/upload-dar \
        -config "${ROOT}/nodes/${n}.participant-config.json" \
        -dar "${DAR_DIR}/${dar}")
    done
  done
  echo ""
  echo "Upload complete. Run: $0 plan"
}

run_plan() {
  cat <<EOF

=== Next: redeploy via MCMS (no archive) ===

Parties (cv1 ledger):
  ccipOwner::1220644bd9e52834e8fba90d4607beed37b65991cc2b5377d5d40d07d3db36d4ed51
  ccipBootstrapOwner::1220a9854ea6590622988af59864d2b1588e004ac9850c140761f1038dd937e8f88d

Auth on cv1:
  actAs  = ccipBootstrapOwner (local)
  readAs = ccipOwner

1. Generate + sign MCMS proposals (chainlink-deployments-staging-testnet-canton-init):
   deploy new factory / core contracts with NEW instanceIds

2. Execute from cv1 (mcms-tools):
   set-root → execute-chain → timelock-execute-chain

3. Lane cutover (lanes not live yet — disable old / enable new is precautionary):
   ApplyDestChainConfigUpdates / ApplySourceChainConfigUpdates with isEnabled=false on OLD GlobalConfig
   then deploy new stack and enable lanes on NEW GlobalConfig

4. Verify:
   $0 list --node all
   (Missing desired: should clear after upload; REMOVE stale rows can stay until later)

5. Optional later: archive + DAR remove when you want package hygiene

See: .for-agents/upgrades-testing.md Phase 2
EOF
}

[[ -n "$cmd" ]] || { usage; exit 2; }

case "$cmd" in
  list)   run_list ;;
  upload) run_upload ;;
  plan)   run_plan ;;
  *)
    echo "unknown command: $cmd" >&2
    usage
    exit 2
    ;;
esac
