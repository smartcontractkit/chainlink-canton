#!/usr/bin/env bash
set -euo pipefail

# Multiparty archive for all staging ccipOwner legacy contracts (cv1 coordinator + cv3 signer).
#
# Archives every active contract matching the TEMPLATES list below (all package versions
# for those template names). batch-size 1 = one contract per prepare/sign/execute round.
#
# Prerequisites:
#   - VPN + refreshed JWT in nodes/cv1.participant-config.json (CanActAs ccipOwner)
#   - Vault signing access on each hosting node
#   - Shared ceremony state between cv1 and cv3 (copy reports.json both ways before resume)
#   - kms_protocol_key_id in each participant-config.json when vault uses KMS/HSM
#
# Usage:
#   ./archive-ccip-owner.sh check cv1           # same as dry-run — how many left?
#   ./archive-ccip-owner.sh dry-run cv1         # list every contract that will be archived
#   ./archive-ccip-owner.sh init cv1            # start ceremony (exit 2 = need cv3 sign)
#   ./archive-ccip-owner.sh resume cv3 <id>     # sign on cv3, sync reports to cv1
#   ./archive-ccip-owner.sh resume cv1 <id>     # execute batch; exit 0 or exit 2 again
#
# Archive everything (outer loop):
#
#   1. ./archive-ccip-owner.sh check cv1
#      → 0 contracts: done. Proceed to DAR cleanup + v2 upload.
#
#   2. ./archive-ccip-owner.sh init cv1          # note Ceremony ID
#      → exit 2 is normal (cv1 signed, waiting for cv3)
#
#   3. On cv3: resume cv3 <id>  — copy reports.json back to cv1
#      On cv1: resume cv1 <id>
#
#   4. Repeat step 3 with the SAME ceremony id until cv1 resume exits 0.
#      Each round archives one contract (batch-size 1).
#
#   5. ./archive-ccip-owner.sh check cv1 again.
#      → If contracts remain: start a NEW ceremony (init cv1 — new id). Do not reuse
#        a finished ceremony. Go to step 2.
#      → Repeat until check shows nothing left.
#
# CCIP App token (CanActAs ccipOwner): refresh cv1/cv3 config first, or set ledger_jwt manually.
# Export AWS_PROFILE / kms:Sign creds before init/resume when using kms_protocol_key_id.

STAGING_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
PARTY_CEREMONY_ROOT="$(cd "$STAGING_ROOT/../.." && pwd)"
STATE_DIR="${ARCHIVE_STATE_DIR:-$STAGING_ROOT/ceremonies/archive-ccip-owner}"
CCIP_OWNER_PARTY="${CCIP_OWNER_PARTY:-ccipOwner::1220644bd9e52834e8fba90d4607beed37b65991cc2b5377d5d40d07d3db36d4ed51}"
SYNCHRONIZER_ID="${SYNCHRONIZER_ID:-global}"

TEMPLATES=(
  "#ccip-common:CCIP.GlobalConfig:GlobalConfig"
  "#ccip-common:CCIP.RateLimiter:RateLimiter"
  "#ccip-common:CCIP.SendingMessageV1:SendingMessageV1"
  "#ccip-common:CCIP.ExecutingMessageV1:ExecutingMessageV1"
  "#ccip-common:CCIP.Tickets:TokenReceiveTicket"
  "#ccip-common:CCIP.Events:CCIPMessageSent"
  "#ccip-common:CCIP.Events:ExecutionStateChanged"
  "#ccip-common:CCIP.Events:TokenReceiveTicketClaimed"
  "#ccip-tokenadminregistry:CCIP.TokenAdminRegistry:TokenAdminRegistry"
  "#ccip-tokenadminregistry:CCIP.TokenAdminRegistry:TokenConfig"
  "#ccip-feequoter:CCIP.FeeQuoter:FeeQuoter"
  "#ccip-offramp:CCIP.OffRamp:OffRamp"
  "#ccip-onramp:CCIP.OnRamp:OnRamp"
  "#ccip-perpartyrouter:CCIP.PerPartyRouter:PerPartyRouterFactory"
  "#ccip-perpartyrouter:CCIP.PerPartyRouter:PerPartyRouter"
  "#ccip-perpartyrouter:CCIP.PerPartyRouter:ArchivedExecutedMessages"
  "#ccip-committeeverifier:CCIP.CommitteeVerifier:CommitteeVerifier"
  "#ccip-executor:CCIP.Executor:Executor"
  "#ccip-sender:CCIP.CCIPSender:CCIPSender"
  "#ccip-receiver:CCIP.CCIPReceiver:CCIPReceiver"
  "#ccip-receiver:CCIP.CCIPReceiver:CCIPMessageReceived"
  "#ccip-lockreleasetokenpool:CCIP.LockReleaseTokenPool:LockReleaseTokenPool"
  "#ccip-rmn:CCIP.RMNRemote:RMNRemote"
  "#ccip-factory:CCIP.Factory:CCIPFactory"
  "#mcms:MCMS.Main:MCMS"
)

config_for_cv() {
  case "$1" in
    cv0|cv1|cv2|cv3) echo "$STAGING_ROOT/nodes/$1.participant-config.json" ;;
    *) echo "unknown cv: $1" >&2; return 1 ;;
  esac
}

run_init() {
  local cv="$1"
  local dry_run="${2:-false}"
  local config
  config="$(config_for_cv "$cv")"
  cd "$PARTY_CEREMONY_ROOT"
  local -a args=(
    go run .
    init archive-contracts
    --decentralized-party-id "$CCIP_OWNER_PARTY"
    --synchronizer-id "$SYNCHRONIZER_ID"
    --config "$config"
    --state-dir "$STATE_DIR"
    --batch-size 1
  )
  local t
  for t in "${TEMPLATES[@]}"; do
    args+=(--template "$t")
  done
  if [[ "$dry_run" == "true" ]]; then
    args+=(--dry-run)
  fi
  "${args[@]}"
}

cmd="${1:?usage: $0 check|dry-run|init|resume <cv> [ceremony-id]}"
shift

case "$cmd" in
  check|dry-run|init)
    cv="${1:?cv required: cv1|cv3}"
    if [[ "$cmd" == "check" || "$cmd" == "dry-run" ]]; then
      run_init "$cv" true
    else
      run_init "$cv" false
    fi
    ;;
  resume)
    cv="${1:?cv required: cv1|cv3}"
    ceremony_id="${2:?ceremony-id required}"
    config="$(config_for_cv "$cv")"
    cd "$PARTY_CEREMONY_ROOT"
    go run . \
      resume "$ceremony_id" \
      --config "$config" \
      --state-dir "$STATE_DIR"
    ;;
  *)
    echo "unknown command: $cmd" >&2
    exit 2
    ;;
esac
