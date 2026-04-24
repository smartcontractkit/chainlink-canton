#!/usr/bin/env bash
# scripts/ci/ccv-ref.sh
#
# Manages chainlink-ccv dependency config, go.mod pins, and an optional
# sibling checkout under ../chainlink-ccv.
#
# Usage:
#   ccv-ref.sh pin <commit-sha>    Pin the root module to a specific commit
#   ccv-ref.sh pin-devenv          Pin configured submodules to the root pseudo-version
#   ccv-ref.sh validate            Check yaml, go.mod files, and local checkout are in sync
#   ccv-ref.sh read <key> [<key>]  Emit '<key>=<value>' lines for each key
#                                  (tracking, pinned, module) to $GITHUB_OUTPUT
#                                  (or stdout if unset). Used by CI workflows.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_DIR="$(cd "$SCRIPT_DIR/../.." && pwd)"
CCV_REF_FILE="$SCRIPT_DIR/ccv-ref.yaml"
CCV_DIR="$REPO_DIR/../chainlink-ccv"

if [[ ! -f "$CCV_REF_FILE" ]]; then
  echo "ERROR: ccv-ref.yaml not found" >&2
  exit 1
fi

yaml_scalar() {
  local key="$1"
  awk -F':[[:space:]]*' -v key="$key" '
    $0 ~ "^" key ":" {
      sub("^[^:]+:[[:space:]]*", "", $0)
      print
      exit
    }
  ' "$CCV_REF_FILE"
}

yaml_array() {
  local key="$1"
  awk -v key="$key" '
    $0 ~ "^" key ":" {
      in_array = 1
      next
    }
    in_array && $0 ~ /^[^[:space:]-][^:]*:/ {
      exit
    }
    in_array && $0 ~ /^[[:space:]]*-[[:space:]]*/ {
      sub(/^[[:space:]]*-[[:space:]]*/, "", $0)
      print
    }
  ' "$CCV_REF_FILE"
}

set_yaml_scalar() {
  local key="$1"
  local value="$2"
  local tmp
  tmp=$(mktemp)
  awk -v key="$key" -v value="$value" '
    $0 ~ "^" key ":" {
      print key ": " value
      updated = 1
      next
    }
    { print }
    END {
      if (!updated) {
        print key ": " value
      }
    }
  ' "$CCV_REF_FILE" > "$tmp"
  mv "$tmp" "$CCV_REF_FILE"
}

MODULE=$(yaml_scalar module)
CMD="${1:?Usage: ccv-ref.sh <pin|pin-devenv|validate|read> [args]}"
shift

mapfile -t SUBMODULES < <(yaml_array submodules)

find_gomods() {
  find "$REPO_DIR" -name "go.mod"
}

use_local_checkout() {
  local flag="${CCV_USE_LOCAL_CHECKOUT:-}"
  [[ -n "$flag" && "$flag" != "0" && "$flag" != "false" ]]
}

find_local_ccv_modules() {
  if [[ ! -d "$CCV_DIR" ]]; then
    return 0
  fi

  find "$CCV_DIR" -name "go.mod" | while IFS= read -r gomod; do
    local module_path
    module_path=$(awk '/^module / {print $2; exit}' "$gomod")
    if [[ -n "$module_path" ]]; then
      printf '%s|%s\n' "$module_path" "$(dirname "$gomod")"
    fi
  done
}

apply_local_replaces() {
  if ! use_local_checkout; then
    return 0
  fi
  if [[ ! -d "$CCV_DIR" ]]; then
    echo "ERROR: CCV_USE_LOCAL_CHECKOUT is set but $CCV_DIR does not exist" >&2
    exit 1
  fi

  local gomod dir module_path module_dir
  echo "Applying local chainlink-ccv replaces from $CCV_DIR..."
  while IFS= read -r gomod; do
    dir=$(dirname "$gomod")
    while IFS='|' read -r module_path module_dir; do
      [[ -n "$module_path" && -n "$module_dir" ]] || continue
      (cd "$dir" && go mod edit -replace="${module_path}=${module_dir}")
    done < <(find_local_ccv_modules)
  done < <(find_gomods)
}

has_module() {
  local gomod="$1"
  local module_path="$2"
  grep -q "^[[:space:]]*${module_path} " "$gomod"
}

module_version() {
  local gomod="$1"
  local module_path="$2"
  grep "^[[:space:]]*${module_path} " "$gomod" | awk '{print $2}' | head -n1
}

module_short_sha() {
  local gomod="$1"
  local module_path="$2"
  module_version "$gomod" "$module_path" | grep -oE '[0-9a-f]{12}$' || true
}

pin() {
  local sha="${1:?Usage: ccv-ref.sh pin <commit-sha>}"
  echo "Pinning ccv to ${sha:0:12}"

  echo "  Updating ccv-ref.yaml..."
  set_yaml_scalar pinned "$sha"

  echo "  Updating go.mod files for ${MODULE}..."
  while IFS= read -r gomod; do
    if has_module "$gomod" "$MODULE"; then
      local dir
      dir=$(dirname "$gomod")
      echo "    $gomod"
      (cd "$dir" && go get "${MODULE}@${sha}" && go mod tidy)
    fi
  done < <(find_gomods)

  if [[ -d "$CCV_DIR" ]]; then
    echo "  Checking out ../chainlink-ccv..."
    (cd "$CCV_DIR" && git fetch origin && git checkout "$sha" 2>/dev/null) \
      || echo "  WARNING: could not checkout ../chainlink-ccv at $sha" >&2
  fi

  apply_local_replaces

  if [[ "${#SUBMODULES[@]}" -gt 0 ]]; then
    pin_devenv
  fi

  echo "  Validating..."
  validate
}

pin_devenv() {
  if [[ "${#SUBMODULES[@]}" -eq 0 ]]; then
    echo "No ccv submodules configured; nothing to pin."
    return 0
  fi

  local pseudo_ver
  pseudo_ver=$(module_version "$REPO_DIR/go.mod" "$MODULE")
  if [[ -z "$pseudo_ver" ]]; then
    echo "ERROR: could not find ${MODULE} in $REPO_DIR/go.mod — run pin first" >&2
    exit 1
  fi

  local submodule module_path gomod dir
  for submodule in "${SUBMODULES[@]}"; do
    module_path="${MODULE}/${submodule}"
    echo "Pinning ${module_path} to ${pseudo_ver}"
    while IFS= read -r gomod; do
      if has_module "$gomod" "$module_path"; then
        dir=$(dirname "$gomod")
        echo "  $gomod"
        (cd "$dir" && go get "${module_path}@${pseudo_ver}" && go mod tidy)
      fi
    done < <(find_gomods)
  done
}

validate() {
  local pinned
  pinned=$(yaml_scalar pinned)
  if [[ -z "$pinned" || "$pinned" == "null" ]]; then
    echo "ERROR: pinned ref not set in ccv-ref.yaml" >&2
    exit 1
  fi

  local short="${pinned:0:12}"
  local errors=0
  local warnings=0
  local gomod sha submodule module_path

  echo "Checking go.mod files against ccv-ref.yaml pinned ($short)..."
  while IFS= read -r gomod; do
    if has_module "$gomod" "$MODULE"; then
      sha=$(module_short_sha "$gomod" "$MODULE")
      if [[ -n "$sha" && "$sha" == "$short" ]]; then
        echo "  OK: $gomod (${MODULE})"
      elif [[ -n "$sha" ]]; then
        echo "  MISMATCH: $gomod (${MODULE} $sha != $short)" >&2
        errors=1
      fi
    fi

    for submodule in "${SUBMODULES[@]}"; do
      module_path="${MODULE}/${submodule}"
      if has_module "$gomod" "$module_path"; then
        sha=$(module_short_sha "$gomod" "$module_path")
        if [[ -n "$sha" && "$sha" == "$short" ]]; then
          echo "  OK: $gomod (${module_path})"
        elif [[ -n "$sha" ]]; then
          echo "  MISMATCH: $gomod (${module_path} $sha != $short)" >&2
          errors=1
        fi
      fi
    done
  done < <(find_gomods)

  if [[ -d "$CCV_DIR" ]]; then
    echo "Checking ../chainlink-ccv checkout..."
    local ccv_sha
    ccv_sha=$(cd "$CCV_DIR" && git rev-parse HEAD 2>/dev/null || true)
    if [[ -n "$ccv_sha" && "$ccv_sha" == "$pinned" ]]; then
      echo "  OK: ../chainlink-ccv at $short"
    elif [[ -n "$ccv_sha" ]]; then
      echo "  MISMATCH: ../chainlink-ccv at ${ccv_sha:0:12}, expected $short" >&2
      warnings=1
    fi
  fi

  echo ""
  if [[ "$errors" -eq 1 ]]; then
    echo "FAILED: go.mod files out of sync. Run: make pin-ccv-all REF=$pinned" >&2
    exit 1
  elif [[ "$warnings" -eq 1 ]]; then
    echo "PASSED with warnings. Run: make pin-ccv-all REF=$pinned" >&2
  else
    echo "PASSED"
  fi
}

read_keys() {
  local out="${GITHUB_OUTPUT:-}"
  if [[ "$#" -eq 0 ]]; then
    echo "ERROR: read requires at least one key (tracking|pinned|module)" >&2
    exit 1
  fi

  local key value
  for key in "$@"; do
    case "$key" in
      tracking|pinned|module) ;;
      *) echo "ERROR: unknown key '$key' (expected tracking|pinned|module)" >&2; exit 1 ;;
    esac
    value=$(yaml_scalar "$key")
    if [[ -z "$value" || "$value" == "null" ]]; then
      echo "ERROR: .${key} not set in $CCV_REF_FILE" >&2
      exit 1
    fi
    if [[ -n "$out" ]]; then
      echo "${key}=${value}" >> "$out"
    else
      echo "${key}=${value}"
    fi
  done
}

case "$CMD" in
  pin)         pin "$@" ;;
  pin-devenv)  pin_devenv ;;
  validate)    validate ;;
  read)        read_keys "$@" ;;
  *)           echo "Usage: ccv-ref.sh <pin|pin-devenv|validate|read>" >&2; exit 1 ;;
esac
