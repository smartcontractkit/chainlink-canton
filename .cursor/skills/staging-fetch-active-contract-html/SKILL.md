---
name: staging-fetch-active-contract-html
description: >-
  Refreshes the HTML report produced by scripts/staging/fetch_active_contract_by_instance_address
  (ledger ACS / contract-id fetch). Use when the user asks to update staging contract HTML,
  refresh GlobalConfig HTML, rerun the fetch script to HTML, or regenerate report.html for that tool.
disable-model-invocation: true
---

# Staging fetch → HTML report

## Goal

Run the Go fetch from the repo root and overwrite the HTML snapshot at `scripts/staging/fetch_active_contract_by_instance_address/report.html`.

## Steps

1. Ensure `ONCHAIN_CANTON_JWT_TOKEN` is set in the environment (same as the JSON workflow for this script).

2. From the repository root, run (adjust template, instance-address, or instance-id only when the user’s target contract differs):

```bash
go run ./scripts/staging/fetch_active_contract_by_instance_address \
  --format html \
  --html-out ./scripts/staging/fetch_active_contract_by_instance_address/report.html \
  --template '#ccip-common:CCIP.GlobalConfig:GlobalConfig' \
  --instance-address '{"address":"0xa95f120fc972c72e75d74c880c26ba982c60b123c74aa9e5b18e138a59e0916a"}' \
  --instance-id globalconfig-szvgb
```

3. Require exit code 0. Open or mention `scripts/staging/fetch_active_contract_by_instance_address/report.html`; the page header includes a UTC `generated` timestamp.

4. On failure (missing JWT, gRPC errors, no matching contract), report stderr to the user. Do not hand-edit the HTML file to “fix” fetch errors.

## Notes

- Default script output remains JSON (`--format json` or omit `--format`).
- For other templates or instances, keep `--format html` and `--html-out` as above (or another explicit path) and adjust `--template`, `--instance-address`, `--instance-id`, or `--contract-id` per the script’s help text.
