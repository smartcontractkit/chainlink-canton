#!/usr/bin/env bash
# Upload dev DARs (*-current.dar) from contracts/dars/current to local compose participant1.
# Prereqs: docker compose up in compose/localnet; grpcurl; python3; run `make compile-contracts` first.
# Override: DARS_DIR=.../contracts/dars/v1_0_0 and DAR_GLOB='*-1.0.0.dar' for frozen release DARs.
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
DARS_DIR="${DARS_DIR:-${REPO_ROOT}/contracts/dars/current}"
DAR_GLOB="${DAR_GLOB:-*-current.dar}"
ADMIN="${CANTON_ADMIN_API:-participant1.admin-api.localhost:8080}"
CONFIG="${CANTON_CONFIG:-${REPO_ROOT}/integration-tests/local-docker-compose.toml}"

if [[ -z "${ONCHAIN_CANTON_JWT_TOKEN:-}" ]]; then
  ONCHAIN_CANTON_JWT_TOKEN="$(python3 -c "
import re, pathlib
t = pathlib.Path('${CONFIG}').read_text()
m = re.search(r'^jwt=\"([^\"]+)\"', t, re.M)
print(m.group(1) if m else '')
")"
  export ONCHAIN_CANTON_JWT_TOKEN
fi
if [[ -z "${ONCHAIN_CANTON_JWT_TOKEN}" ]]; then
  echo "Set ONCHAIN_CANTON_JWT_TOKEN or add jwt to ${CONFIG}" >&2
  exit 1
fi

shopt -s nullglob
dars=( "${DARS_DIR}"/${DAR_GLOB} )
if [[ ${#dars[@]} -eq 0 ]]; then
  echo "No DARs matching ${DARS_DIR}/${DAR_GLOB} — run 'make compile-contracts' in chainlink-canton-fcr" >&2
  exit 1
fi
for f in "${dars[@]}"; do
  echo "=== $(basename "$f") ==="
  python3 -c "
import json, base64, sys
b = base64.b64encode(open(sys.argv[1], 'rb').read()).decode()
sys.stdout.write(json.dumps({
    'dars': [{'bytes': b}],
    'vetAllPackages': True,
    'synchronizeVetting': True,
}))
" "$f" | grpcurl -plaintext \
    -H "authorization: Bearer ${ONCHAIN_CANTON_JWT_TOKEN}" \
    -d @ "${ADMIN}" \
    com.digitalasset.canton.admin.participant.v30.PackageService/UploadDar
done

echo "Done: uploaded ${#dars[@]} DAR(s) from ${DARS_DIR} to ${ADMIN}"
