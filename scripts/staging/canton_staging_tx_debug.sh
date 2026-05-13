#!/usr/bin/env bash

set -euo pipefail

usage() {
  cat <<'EOF'
Usage: canton_staging_tx_debug.sh [options] <message-id> [namespace]

Search staging Canton logs for a message ID on these pods only:
  chainlink-ccv-canton-0, chainlink-ccv-canton-aggregator-1, chainlink-ccv-canton-indexer-1
(Use --all-pods to scan every pod in the namespace.)

Only lines containing the message ID are collected. Identical lines (after stripping
kubectl line-number prefixes) are logged once. A **Message path** section summarizes
detection on canton-0, aggregator evidence (or inferred handoff), then indexer.
A deduplicated flow timeline (by log msg / role) follows.

Options:
- --since <duration>: kubectl log window, for example "24h" or "72h"
- --namespace <name>: Kubernetes namespace, default "chainlink-ccv-canton"
- --context <name>: optional kubectl context override
- --include-previous: also inspect previous container logs
- --all-pods: scan all pods in the namespace instead of canton-0 / aggregator-1 / indexer-1
- --no-flow: print per-pod sections with unique message-ID lines (still deduped); default is compact
- --flow: compact per-pod one-liners (this is the default)
- --verbose: with compact mode, also print every unique message-ID line per pod before the summary
- --no-log-file: do not write under repo logs/ (default is to tee full output to logs/)
- -h, --help: show this help

Environment variables:
- CANTON_TX_DEBUG_SINCE: kubectl log window, default "24h"
- CANTON_TX_DEBUG_CONTEXT: optional kubectl context override
- CANTON_TX_DEBUG_INCLUDE_PREVIOUS: set to "1" to also inspect previous container logs
- CANTON_TX_DEBUG_NO_LOG_FILE: set to "1" to skip writing logs/ (same as --no-log-file)

Full output is copied to chainlink-canton/logs/tx-debug-<id-prefix>-<UTC-time>.txt (gitignored).

Install jq for clearer log .msg extraction in --flow mode (optional; script falls back to regex).
EOF
}

fail() {
  echo "ERROR: $*" >&2
  exit 1
}

warn() {
  echo "WARN: $*" >&2
}

require_cmd() {
  local cmd_name="$1"
  command -v "$cmd_name" >/dev/null 2>&1 || fail "required command not found: $cmd_name"
}

tailscale_status="not checked"

check_tailscale() {
  if ! command -v tailscale >/dev/null 2>&1; then
    tailscale_status="tailscale CLI not found; proceeding with kubectl access checks"
    return 0
  fi

  local status_text
  if ! status_text="$(tailscale status 2>&1)"; then
    tailscale_status="tailscale status failed; proceeding with kubectl access checks"
    return 0
  fi

  if ! printf '%s\n' "$status_text" | rg -qi 'smartcontract\.com|smartcontract'; then
    tailscale_status="tailscale does not appear connected to smartcontract.com; proceeding with kubectl access checks"
    return 0
  fi

  tailscale_status="connected"
}

get_aws_profile_from_kube_context() {
  local config_text
  if ! config_text="$("${kubectl_cmd[@]}" config view --minify --raw 2>/dev/null)"; then
    return 1
  fi

  printf '%s\n' "$config_text" | awk '
    $1 == "-" && $2 == "name:" && $3 == "AWS_PROFILE" { capture = 1; next }
    capture && $1 == "value:" { print $2; exit }
  '
}

maybe_refresh_aws_sso() {
  local error_text="$1"

  if ! printf '%s\n' "$error_text" | rg -qi 'sso session associated with this profile has expired|expired or is otherwise invalid|getting credentials: exec: executable aws failed'; then
    return 1
  fi

  require_cmd aws

  local aws_profile
  aws_profile="$(get_aws_profile_from_kube_context)"
  [[ -n "$aws_profile" ]] || fail "kubectl auth failed due to AWS exec auth, but AWS_PROFILE could not be determined from the active kube context."

  warn "AWS SSO session appears expired for profile ${aws_profile}; attempting aws sso login"
  aws sso login --profile "$aws_profile" \
    || fail "aws sso login failed for profile ${aws_profile}. Refresh staging credentials and retry."

  return 0
}

# kubectl auth can-i may print admission warnings on stderr; only the yes/no line is authoritative.
can_i_verdict_from_output() {
  printf '%s\n' "$1" | rg '^yes$|^no$' | tail -n1
}

strip_kubectl_line_prefix() {
  printf '%s\n' "$1" | sed -E 's/^[0-9]+://'
}

# stdin: lines; stdout: first occurrence of each line body after stripping ^digits: prefix
dedupe_identical_log_bodies() {
  awk '{ line=$0; sub(/^[0-9]+:/, "", line); if (line == "") next; if (!seen[line]++) print line }'
}

# ISO-ish timestamp for sorting, or empty.
extract_log_ts() {
  local line="$1"
  local t=""
  if command -v jq >/dev/null 2>&1; then
    t="$(printf '%s\n' "$line" | jq -r 'try fromjson catch empty |
      if .ts == null then empty
      elif (.ts | type) == "string" then .ts
      elif (.ts | type) == "number" then (.ts | tostring)
      else empty end' 2>/dev/null)"
  fi
  if [[ -n "$t" && "$t" != "null" ]]; then
    printf '%s\n' "$t"
    return 0
  fi
  if [[ "$line" =~ ^([0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9:\.]+Z) ]]; then
    printf '%s\n' "${BASH_REMATCH[1]}"
    return 0
  fi
  if [[ "$line" =~ \"ts\":\"([^\"]+)\" ]]; then
    printf '%s\n' "${BASH_REMATCH[1]}"
    return 0
  fi
  return 1
}

# Short human summary for deduplication (JSON .msg when possible).
extract_log_msg() {
  local line="$1"
  local m=""
  local tail_json=""

  if [[ "$line" == *$'\t'* ]]; then
    local msg_col=""
    local last_col=""
    msg_col="$(printf '%s\n' "$line" | awk -F'\t' '{ if (NF >= 2) print $(NF-1) }')"
    last_col="$(printf '%s\n' "$line" | awk -F'\t' '{ print $NF }')"
    if [[ -n "$msg_col" && "$last_col" =~ ^\{.*\}$ ]]; then
      tail_json="$(summarize_json_log_payload "$last_col")"
      if [[ -n "$tail_json" ]]; then
        printf '%s  %s\n' "$msg_col" "$tail_json"
      else
        printf '%s\n' "$msg_col"
      fi
      return 0
    fi
    if [[ -n "$msg_col" && "$msg_col" != "$line" ]]; then
      printf '%s\n' "$msg_col"
      return 0
    fi
  fi

  if command -v jq >/dev/null 2>&1; then
    m="$(printf '%s\n' "$line" | jq -r '
      def kv($k; $v):
        if $v == null or $v == "" or $v == [] then empty
        else "\($k)=\($v|tostring)" end;
      try fromjson catch empty
      | (.msg // empty) as $msg
      | [
          kv("meetsRequirement"; .meetsRequirement?),
          kv("messageBlock"; .messageBlock?),
          kv("latestBlock"; .latestBlock?),
          kv("safeBlock"; .safeBlock?),
          kv("finalizedBlock"; .finalizedBlock?),
          kv("verifierAddress"; .verifierAddress?),
          kv("defaultExecutorAddress"; .defaultExecutorAddress?),
          kv("signer"; .signer?),
          kv("nonce"; .nonce?),
          kv("sourceChain"; .sourceChain?),
          kv("destChain"; .destChain?),
          kv("pending"; .pendingCount?),
          kv("blockNumber"; .blockNumber?),
          kv("seq"; .seqNum?),
          kv("verifierSource"; .verifierSourceAddress?)
        ]
      | map(select(length > 0))
      | if $msg == "" then empty
        elif length == 0 then $msg
        else $msg + "  " + join(" ")
        end' 2>/dev/null)"
  fi
  if [[ -n "$m" && "$m" != "null" ]]; then
    printf '%s\n' "$m"
    return 0
  fi
  m="$(printf '%s\n' "$line" | rg -o '"msg":"[^"]*"' | head -1 | sed 's/^"msg":"//;s/"$//' || true)"
  if [[ -n "$m" ]]; then
    printf '%s\n' "$m"
    return 0
  fi
  printf '%.200s\n' "$line"
}

summarize_json_log_payload() {
  local payload="$1"
  [[ "$payload" =~ ^\{.*\}$ ]] || return 0
  command -v jq >/dev/null 2>&1 || return 0

  printf '%s\n' "$payload" | jq -r '
    def kv($k; $v):
      if $v == null or $v == "" or $v == [] then empty
      else "\($k)=\($v|tostring)" end;
    [
      kv("messageID"; .messageID?),
      kv("verifierSource"; .verifierSourceAddress?),
      kv("sourceSpecifiedCCVs"; ."Source Specified CCVs"?),
      kv("existingVerifications"; ."Exisiting Verifications"?),
      kv("unknownCCVs"; ."Unknown CCVs"?),
      kv("state"; .state?),
      kv("returnData"; .returnData?)
    ]
    | map(select(length > 0))
    | join(" ")' 2>/dev/null
}

pod_to_role() {
  local name="${1#pod/}"
  case "$name" in
    chainlink-ccv-canton-aggregator-1-*)
      printf '%s\n' "aggregator Deployment pod"
      ;;
    chainlink-ccv-canton-indexer-1-*)
      printf '%s\n' "indexer pod"
      ;;
    chainlink-ccv-canton-0-*)
      printf '%s\n' "canton-0 (in-pod default verifier / EVM source reader)"
      ;;
    chainlink-canton-committee-verifier-*)
      printf '%s\n' "standalone committee-verifier"
      ;;
    chainlink-ccv-canton-[0-9]*-*)
      local rest="${name#chainlink-ccv-canton-}"
      local idx="${rest%%-*}"
      printf '%s\n' "canton node ${idx}"
      ;;
    *)
      printf '%s\n' "other (${name%%-*})"
      ;;
  esac
}

# Split one record: ts\x1fpod\x1fclean (clean may contain tabs; do not use TAB as delimiter).
mp_split_ts_pod_clean() {
  local record="$1"
  _mp_ts="${record%%$'\x1f'*}"
  local tail="${record#*$'\x1f'}"
  _mp_pod="${tail%%$'\x1f'*}"
  _mp_clean="${tail#*$'\x1f'}"
}

indent_log_line() {
  printf '%s\n' "$1" | sed 's/^/    /'
}

# Ordered narrative: canton-0 detection → aggregator (pod or inferred) → indexer.
# Prints full kubectl log lines (verbatim aside from optional line-number prefix strip).
emit_message_path_report() {
  local tmp_sorted mp_sep=$'\x1f'
  tmp_sorted="$(mktemp "${TMPDIR:-/tmp}/canton_tx_mp.XXXXXX")"

  while IFS= read -r record || [[ -n "${record:-}" ]]; do
    [[ -z "$record" ]] && continue
    local pod_key log_line clean ts
    pod_key="${record%%"${record_sep}"*}"
    log_line="${record#*"${record_sep}"}"
    clean="$(strip_kubectl_line_prefix "$log_line")"
    ts="$(extract_log_ts "$clean" || true)"
    [[ -z "$ts" ]] && ts="0000-00-00T00:00:00.000Z"
    printf '%s%s%s%s%s\n' "$ts" "$mp_sep" "$pod_key" "$mp_sep" "$clean"
  done <"$all_matching_lines_file" | sort -t"$mp_sep" -k1,1 >"$tmp_sorted"

  local _mp_ts _mp_pod _mp_clean
  local s1_ts="" s1_pod="" s1_line=""
  while IFS= read -r row || [[ -n "${row:-}" ]]; do
    [[ -z "$row" ]] && continue
    mp_split_ts_pod_clean "$row"
    [[ "$_mp_pod" =~ chainlink-ccv-canton-0- ]] || continue
    printf '%s\n' "$_mp_clean" | rg -q 'CCIPMessageSent|Added message to pending queue|OnRamp Event Structure|Event details' || continue
    s1_ts="$_mp_ts"
    s1_pod="$_mp_pod"
    s1_line="$_mp_clean"
    break
  done <"$tmp_sorted"

  local agg_seen=0 agg_ts="" agg_pod="" agg_line=""
  while IFS= read -r row || [[ -n "${row:-}" ]]; do
    [[ -z "$row" ]] && continue
    mp_split_ts_pod_clean "$row"
    [[ "$_mp_pod" =~ chainlink-ccv-canton-aggregator-1- ]] || continue
    agg_seen=1
    agg_ts="$_mp_ts"
    agg_pod="$_mp_pod"
    agg_line="$_mp_clean"
    break
  done <"$tmp_sorted"

  local store_ts="" store_pod="" store_line=""
  while IFS= read -r row || [[ -n "${row:-}" ]]; do
    [[ -z "$row" ]] && continue
    mp_split_ts_pod_clean "$row"
    [[ "$_mp_pod" =~ chainlink-ccv-canton-0- ]] || continue
    printf '%s\n' "$_mp_clean" | rg -q 'Successfully stored CCV data' || continue
    store_ts="$_mp_ts"
    store_pod="$_mp_pod"
    store_line="$_mp_clean"
    break
  done <"$tmp_sorted"

  local idx_ts="" idx_pod="" idx_line=""
  while IFS= read -r row || [[ -n "${row:-}" ]]; do
    [[ -z "$row" ]] && continue
    mp_split_ts_pod_clean "$row"
    [[ "$_mp_pod" =~ chainlink-ccv-canton-indexer-1- ]] || continue
    printf '%s\n' "$_mp_clean" | rg -q 'Found Message' || continue
    idx_ts="$_mp_ts"
    idx_pod="$_mp_pod"
    idx_line="$_mp_clean"
    break
  done <"$tmp_sorted"

  rm -f "$tmp_sorted"

  echo ""
  echo "=== Message path (canton-0 → aggregator → indexer) ==="
  echo "Below: full log lines as returned by kubectl (after stripping optional leading line numbers). Milestones are earliest match in --since=${since_window}."
  echo ""

  echo "(1) chainlink-ccv-canton-0 — detect / queue (EVM side)"
  if [[ -n "$s1_line" ]]; then
    echo "    pod: ${s1_pod}"
    indent_log_line "$(extract_log_msg "$s1_line")"
  else
    echo "    (no matching detection line in window — widen --since or inspect canton-0)"
  fi

  echo ""
  echo "(2) chainlink-ccv-canton-aggregator-1 — aggregator"
  if [[ "$agg_seen" -eq 1 ]]; then
    echo "    pod: ${agg_pod}"
    indent_log_line "$(extract_log_msg "$agg_line")"
  else
    printf "    No stdout/stderr lines with this messageId in the aggregator pod for --since=%s.\n" "${since_window}"
    if [[ -n "$store_line" ]]; then
      echo "    Inferred handoff — canton-0 reported successful CCV write to aggregator service:"
      echo "    pod: ${store_pod}"
      indent_log_line "$(extract_log_msg "$store_line")"
    else
      echo "    No \"Successfully stored CCV data\" line on canton-0 in this window either."
    fi
  fi

  echo ""
  echo "(3) chainlink-ccv-canton-indexer-1 — indexer"
  if [[ -n "$idx_line" ]]; then
    echo "    pod: ${idx_pod}"
    indent_log_line "$(extract_log_msg "$idx_line")"
  else
    echo "    (no \"Found Message\" / indexer hit in window — widen --since or check indexer pods)"
  fi
  echo ""
}

emit_flow_report() {
  local raw sorted issues_f
  raw="$(mktemp "${TMPDIR:-/tmp}/canton_tx_flow.XXXXXX")"
  sorted="$(mktemp "${TMPDIR:-/tmp}/canton_tx_flow_s.XXXXXX")"
  issues_f="$(mktemp "${TMPDIR:-/tmp}/canton_tx_issues.XXXXXX")"

  while IFS= read -r record || [[ -n "${record:-}" ]]; do
    [[ -z "$record" ]] && continue
    local pod_key log_line clean ts role msg sort_ts
    pod_key="${record%%"${record_sep}"*}"
    log_line="${record#*"${record_sep}"}"
    clean="$(strip_kubectl_line_prefix "$log_line")"
    ts="$(extract_log_ts "$clean" || true)"
    [[ -z "$ts" ]] && ts="0000-00-00T00:00:00.000Z"
    role="$(pod_to_role "$pod_key")"
    msg="$(extract_log_msg "$clean")"
    printf '%s\t%s\t%s\t%s\n' "$ts" "$role" "$msg" "$pod_key"
  done <"$all_matching_lines_file" >"$raw"

  sort -t$'\t' -k1,1 "$raw" >"$sorted"

  echo ""
  echo "=== Flow timeline (deduplicated: same role + same log msg collapsed with counts) ==="
  echo "Ordering: by first-seen timestamp. Expected path: canton-0 (verifier) → CCV stored → indexer discovers message."
  echo "The standalone aggregator Deployment often has no lines with the raw hex messageId; work appears under canton-0's coordinator logger."
  echo ""

  awk -F'\t' 'BEGIN {OFS="\t"}
    {
      key = $2 SUBSEP $3
      if (!(key in first_ts)) first_ts[key] = $1
      count[key]++
      last_ts[key] = $1
      role[key] = $2
      msg[key] = $3
    }
    END {
      for (k in count) print first_ts[k], last_ts[k], count[k], role[k], msg[k]
    }' "$sorted" \
    | sort -t$'\t' -k1,1 \
    | while IFS=$'\t' read -r fts lts c r m || [[ -n "${fts:-}" ]]; do
        if [[ "$c" -gt 1 ]]; then
          printf '%s  %s  %s  (×%s repeats, last %s)\n' "$fts" "$r" "$m" "$c" "$lts"
        else
          printf '%s  %s  %s\n' "$fts" "$r" "$m"
        fi
      done

  echo ""
  echo "=== Message-ID lines that look like errors, warnings, or blocked finality ==="
  local issue_pat='(?i)("level":"(error|warn|dpanic)"|\t(ERROR|WARN)\t|Unauthenticated|authentication failed|\bpanic\b|\bfatal\b|meetsRequirement":false)'
  : >"$issues_f"
  while IFS= read -r record || [[ -n "${record:-}" ]]; do
    [[ -z "$record" ]] && continue
    local pod_key log_line clean ts_line
    pod_key="${record%%"${record_sep}"*}"
    log_line="${record#*"${record_sep}"}"
    clean="$(strip_kubectl_line_prefix "$log_line")"
    if printf '%s\n' "$clean" | rg -qi "$issue_pat"; then
      ts_line="$(extract_log_ts "$clean" || true)"
      printf '%s\t%s\t%s\n' "${ts_line:-?}" "$(pod_to_role "$pod_key")" "$clean" >>"$issues_f"
    fi
  done <"$all_matching_lines_file"
  if [[ ! -s "$issues_f" ]]; then
    echo "(none: no message-ID log line matched error/warn/dpanic/Unauthenticated/panic/fatal/meetsRequirement:false)"
  else
    local raw_lines uniq_total
    raw_lines="$(wc -l <"$issues_f" | tr -d ' ')"
    uniq_total="$(sort -u "$issues_f" | wc -l | tr -d ' ')"
    sort -u "$issues_f" | head -n 80 | while IFS=$'\t' read -r ts_col role_col rest || [[ -n "${ts_col:-}" ]]; do
      printf '%s  %s  %s\n' "$ts_col" "$role_col" "$(extract_log_msg "$rest")"
    done
    if [[ "$uniq_total" -gt 80 ]]; then
      echo "... truncated: showing 80 of ${uniq_total} unique issue-shaped lines (${raw_lines} raw matches before dedup)"
    fi
  fi

  rm -f "$raw" "$sorted" "$issues_f"
}

check_kube_access() {
  local ns="$1"
  local auth_output
  local verdict

  if auth_output="$("${kubectl_cmd[@]}" auth can-i get pods -n "$ns" 2>&1)"; then
    verdict="$(can_i_verdict_from_output "$auth_output")"
    if [[ "$verdict" != "yes" ]]; then
      fail "kubectl cannot get pods in namespace ${ns}. Staging access requires Griddle access, expected via infra-access_team-ccip_engineer."
    fi
  else
    if maybe_refresh_aws_sso "$auth_output"; then
      auth_output="$("${kubectl_cmd[@]}" auth can-i get pods -n "$ns" 2>&1)" \
        || fail "kubectl auth check failed for namespace ${ns} after AWS SSO refresh."
      verdict="$(can_i_verdict_from_output "$auth_output")"
      if [[ "$verdict" != "yes" ]]; then
        fail "kubectl cannot get pods in namespace ${ns} after AWS SSO refresh. Staging access requires Griddle access, expected via infra-access_team-ccip_engineer."
      fi
    else
      fail "kubectl auth check failed for namespace ${ns}. Staging access requires Griddle access, expected via infra-access_team-ccip_engineer."
    fi
  fi

  local get_pods_output
  get_pods_output="$("${kubectl_cmd[@]}" get pods -n "$ns" 2>&1)" \
    || fail "unable to list pods in namespace ${ns}. Fix Kubernetes, AWS auth, or Griddle access before running tx debug."
}

since_window="${CANTON_TX_DEBUG_SINCE:-24h}"
namespace="chainlink-ccv-canton"
context_override="${CANTON_TX_DEBUG_CONTEXT:-}"
include_previous="${CANTON_TX_DEBUG_INCLUDE_PREVIOUS:-0}"
flow_mode=1
verbose_mode=0
no_log_file=0
scan_all_pods=0
message_id_raw=""
positionals=()

while [[ $# -gt 0 ]]; do
  case "$1" in
    --since)
      [[ $# -ge 2 ]] || fail "--since requires a value"
      since_window="$2"
      shift 2
      ;;
    --since=*)
      since_window="${1#*=}"
      shift
      ;;
    --namespace)
      [[ $# -ge 2 ]] || fail "--namespace requires a value"
      namespace="$2"
      shift 2
      ;;
    --namespace=*)
      namespace="${1#*=}"
      shift
      ;;
    --context)
      [[ $# -ge 2 ]] || fail "--context requires a value"
      context_override="$2"
      shift 2
      ;;
    --context=*)
      context_override="${1#*=}"
      shift
      ;;
    --include-previous)
      include_previous="1"
      shift
      ;;
    --flow)
      flow_mode=1
      shift
      ;;
    --no-flow)
      flow_mode=0
      shift
      ;;
    --all-pods)
      scan_all_pods=1
      shift
      ;;
    --verbose)
      verbose_mode=1
      shift
      ;;
    --no-log-file)
      no_log_file=1
      shift
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    --)
      shift
      while [[ $# -gt 0 ]]; do
        positionals+=("$1")
        shift
      done
      ;;
    -*)
      fail "unknown option: $1"
      ;;
    *)
      positionals+=("$1")
      shift
      ;;
  esac
done

if [[ ${#positionals[@]} -lt 1 || ${#positionals[@]} -gt 2 ]]; then
  usage >&2
  exit 1
fi

message_id_raw="${positionals[0]}"
if [[ ${#positionals[@]} -eq 2 ]]; then
  namespace="${positionals[1]}"
fi

message_id="${message_id_raw#0x}"
message_id="${message_id,,}"

if [[ -z "${message_id}" ]]; then
  echo "message id is empty after normalization" >&2
  exit 1
fi

if [[ -n "${context_override}" ]]; then
  kubectl_cmd=(kubectl --context "$context_override")
else
  kubectl_cmd=(kubectl)
fi

require_cmd kubectl
require_cmd rg
check_tailscale
check_kube_access "$namespace"

if [[ "${CANTON_TX_DEBUG_NO_LOG_FILE:-0}" != "1" && "$no_log_file" -eq 0 ]]; then
  require_cmd tee
  _script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
  _repo_root="$(cd "$_script_dir/.." && pwd)"
  _log_dir="$_repo_root/logs"
  mkdir -p "$_log_dir" || fail "cannot create logs directory: $_log_dir"
  _log_file="$_log_dir/tx-debug-${message_id:0:12}-$(date -u +%Y%m%dT%H%M%SZ).txt"
  exec > >(tee "$_log_file") 2>&1
  echo "Log file: ${_log_file}"
fi

message_pattern="(${message_id}|0x${message_id})"
record_sep=$'\x1f'

# Lines that mention the message ID and look like forward progress / success (heuristic).
success_hint_pattern='(?i)(successfully stored|message verification completed successfully|message signed successfully|message validation passed|meetsRequirement":true|successfully published and tracked tasks|successfully stored ccv|ccv data batch write completed|batch verification completed|found message)'

echo "Namespace: ${namespace}"
echo "Message ID: ${message_id}"
echo "Since: ${since_window}"
echo "Tailscale: ${tailscale_status}"
echo "Access preflight: passed"
if [[ "$scan_all_pods" -eq 1 ]]; then
  echo "Scope: all pods in namespace"
else
  echo "Scope: chainlink-ccv-canton-0, chainlink-ccv-canton-aggregator-1, chainlink-ccv-canton-indexer-1 only"
fi
if [[ "$flow_mode" -eq 1 ]]; then
  echo "Output mode: compact (use --no-flow for per-pod sections; --verbose adds unique message-ID lines)"
else
  echo "Output mode: per-pod unique message-ID lines (--flow for compact)"
fi
echo

if [[ "$scan_all_pods" -eq 1 ]]; then
  mapfile -t pods < <(
    "${kubectl_cmd[@]}" get pods -n "$namespace" -o name | sort
  )
else
  mapfile -t pods < <(
    "${kubectl_cmd[@]}" get pods -n "$namespace" -o name | sort \
      | rg -N 'chainlink-ccv-canton-0-|chainlink-ccv-canton-aggregator-1-|chainlink-ccv-canton-indexer-1-'
  )
fi

if [[ ${#pods[@]} -eq 0 ]]; then
  if [[ "$scan_all_pods" -eq 1 ]]; then
    echo "No pods found in namespace ${namespace}" >&2
  else
    echo "No pods matching canton-0 / aggregator-1 / indexer-1 in namespace ${namespace}" >&2
  fi
  exit 2
fi

found_any=0
declare -A hit_lines_per_pod=()
all_matching_lines_file=""
success_examples_file=""
all_matching_lines_file="$(mktemp "${TMPDIR:-/tmp}/canton_tx_debug.XXXXXX")"
success_examples_file="$(mktemp "${TMPDIR:-/tmp}/canton_tx_debug_succ.XXXXXX")"
trap 'rm -f "${all_matching_lines_file}" "${success_examples_file}"' EXIT

search_logs() {
  local pod_name="$1"
  local previous_flag="$2"
  local title="$3"

  local log_stream
  if ! log_stream="$("${kubectl_cmd[@]}" logs -n "$namespace" "$pod_name" --all-containers=true --since="$since_window" $previous_flag 2>&1)"; then
    if [[ "$flow_mode" -eq 1 && "$verbose_mode" -eq 0 ]]; then
      printf '%s  kubectl logs failed\n' "$pod_name"
    else
      echo "=== ${title}: ${pod_name} ==="
      echo "(kubectl logs failed for ${pod_name})"
      echo
    fi
    return 0
  fi

  if printf '%s\n' "$log_stream" | rg -q '^Error from server|^Unable to retrieve container logs'; then
    if [[ "$flow_mode" -eq 1 && "$verbose_mode" -eq 0 ]]; then
      printf '%s  logs not available\n' "$pod_name"
    else
      echo "=== ${title}: ${pod_name} ==="
      echo "(kubectl logs not available for ${pod_name})"
      printf '%s\n' "$log_stream" | head -n 3 | sed 's/^/  /'
      echo
    fi
    return 0
  fi

  local match_only
  match_only="$(
    ( printf '%s\n' "$log_stream" | rg -ni -- "$message_pattern" || true ) | dedupe_identical_log_bodies
  )"

  if [[ -n "$match_only" ]]; then
    found_any=1
    local n
    n="$(printf '%s\n' "$match_only" | rg -c . || true)"
    [[ -n "$n" ]] || n=0
    hit_lines_per_pod["$pod_name"]=$(( ${hit_lines_per_pod["$pod_name"]:-0} + n ))
    while IFS= read -r line || [[ -n "$line" ]]; do
      [[ -z "$line" ]] && continue
      printf '%s%s%s\n' "$pod_name" "$record_sep" "$line"
    done <<<"$match_only" >>"$all_matching_lines_file"
  fi

  if [[ "$flow_mode" -eq 1 && "$verbose_mode" -eq 0 ]]; then
    if [[ -n "$match_only" ]]; then
      printf '%s  %s unique line(s) with message ID\n' "$pod_name" "$n"
    else
      printf '%s  (no matches)\n' "$pod_name"
    fi
    return 0
  fi

  if [[ "$flow_mode" -eq 1 && "$verbose_mode" -eq 1 ]]; then
    echo "=== ${title}: ${pod_name} (unique message-ID lines) ==="
    if [[ -n "$match_only" ]]; then
      printf '%s\n' "$match_only"
    else
      echo "(no matches)"
    fi
    echo
    return 0
  fi

  echo "=== ${title}: ${pod_name} (unique message-ID lines) ==="
  if [[ -n "$match_only" ]]; then
    printf '%s\n' "$match_only"
  else
    echo "(no matches)"
  fi
  echo
}

if [[ "$flow_mode" -eq 1 && "$verbose_mode" -eq 0 ]]; then
  echo "=== Per-pod scan (compact) ==="
fi
for pod in "${pods[@]}"; do
  search_logs "$pod" "" "current logs"
  if [[ "$include_previous" == "1" ]]; then
    search_logs "$pod" "--previous" "previous logs"
  fi
done
if [[ "$flow_mode" -eq 1 && "$verbose_mode" -eq 0 ]]; then
  echo
fi

echo "=== Cross-pod summary ==="
if [[ "$found_any" == "0" ]]; then
  echo "No log lines contained message ID ${message_id} in namespace ${namespace} for --since=${since_window}."
  echo "Try widening the time window with --since 72h."
  exit 0
fi

echo "Pods with at least one matching log line (unique bodies after strip/dedup):"
for pod in "${!hit_lines_per_pod[@]}"; do
  printf '  %s: %s unique line(s)\n' "$pod" "${hit_lines_per_pod[$pod]}"
done | sort

success_count=0
: >"$success_examples_file"
if [[ -s "$all_matching_lines_file" ]]; then
  while IFS= read -r record || [[ -n "${record:-}" ]]; do
    [[ -z "$record" ]] && continue
    pod_key="${record%%"${record_sep}"*}"
    log_line="${record#*"${record_sep}"}"
    [[ -z "$log_line" ]] && continue
    if printf '%s\n' "$log_line" | rg -qi "$success_hint_pattern"; then
      ((success_count++)) || true
      if [[ "$success_count" -le 5 ]]; then
        printf '%s\n' "${pod_key}: $(extract_log_msg "$log_line")" >>"$success_examples_file"
      fi
    fi
  done <"$all_matching_lines_file"
fi

if [[ "$success_count" -gt 0 ]]; then
  echo "Success-like lines (heuristic: same log line contains message ID + progress/success substring): YES (${success_count} line(s))"
  echo "Examples (up to 5, already deduped at ingest):"
  sed 's/^/  /' "$success_examples_file"
else
  echo "Success-like lines (heuristic): NO — no single line contained both the message ID and a success/progress hint."
  echo "Inspect flow timeline below or use --verbose with --no-flow for per-pod lines."
fi

echo
echo "Total pods scanned: ${#pods[@]}"

emit_message_path_report
emit_flow_report
