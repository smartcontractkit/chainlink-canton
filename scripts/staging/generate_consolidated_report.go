// Command generate_consolidated_report generates a single HTML report with all CCIP contracts.
// Features a dropdown to toggle between contracts.
//
// Usage:
//   go run ./scripts/staging/generate_consolidated_report.go
//
// Or build and run:
//   go build -o consolidated_report ./scripts/staging/generate_consolidated_report.go
//   ./consolidated_report

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Party constants
const (
	partyCCIPOwner          = "ccipOwner::1220644bd9e52834e8fba90d4607beed37b65991cc2b5377d5d40d07d3db36d4ed51"
	partyCCIPBootstrapOwner = "ccipBootstrapOwner::1220a9854ea6590622988af59864d2b1588e004ac9850c140761f1038dd937e8f88d"
)

// ContractInfo holds the metadata for each contract to fetch
type ContractInfo struct {
	TemplateID      string `json:"template_id"`
	InstanceAddress string `json:"instance_address"`
	InstanceID      string `json:"instance_id"`
	Type            string `json:"type"`
	Party           string `json:"party"`
}

// List of all contracts to fetch
var contractsToFetch = []ContractInfo{
	{
		TemplateID:      "#ccip-common:CCIP.GlobalConfig:GlobalConfig",
		InstanceAddress: `{"address":"0xa95f120fc972c72e75d74c880c26ba982c60b123c74aa9e5b18e138a59e0916a"}`,
		InstanceID:      "globalconfig-szvgb",
		Type:            "GlobalConfig",
		Party:           partyCCIPOwner,
	},
	{
		TemplateID:      "#ccip-feequoter:CCIP.FeeQuoter:FeeQuoter",
		InstanceAddress: `{"address":"0x3891327bf89b1621f67a720f73f8478777f2c106d95e570c5fa388f138bc0728"}`,
		InstanceID:      "feequoter-shywn",
		Type:            "FeeQuoter",
		Party:           partyCCIPOwner,
	},
	{
		TemplateID:      "#ccip-committeeverifier:CCIP.CommitteeVerifier:CommitteeVerifier",
		InstanceAddress: `{"address":"0xf11b7b25ed8ac60beecb78e58fba954dd9b75f13b1b67ff0983b55aab52dfcd1"}`,
		InstanceID:      "committeeverifier-suoid",
		Type:            "CommitteeVerifier",
		Party:           partyCCIPOwner,
	},
	{
		TemplateID:      "#ccip-offramp:CCIP.OffRamp:OffRamp",
		InstanceAddress: `{"address":"0xe9c3534382c638dbd457aa92becdc61cb6c294795e176365baaa06be3dd885fa"}`,
		InstanceID:      "offramp-uaxss",
		Type:            "OffRamp",
		Party:           partyCCIPOwner,
	},
	{
		TemplateID:      "#ccip-onramp:CCIP.OnRamp:OnRamp",
		InstanceAddress: `{"address":"0x92b53bcb058aabfc52cb617230375b5dacf8bc19932de5a9f56df659e4944c7b"}`,
		InstanceID:      "onramp-tlspm",
		Type:            "OnRamp",
		Party:           partyCCIPOwner,
	},
	{
		TemplateID:      "#ccip-executor:CCIP.Executor:Executor",
		InstanceAddress: `{"address":"0xa3fecf9edeb0686bf58e17b4765a5806ff934ff8efb145a42c965a79a32f875c"}`,
		InstanceID:      "executor-zzpfy",
		Type:            "Executor",
		Party:           partyCCIPOwner,
	},
	{
		TemplateID:      "#ccip-lockreleasetokenpool:CCIP.LockReleaseTokenPool:LockReleaseTokenPool",
		InstanceAddress: `{"address":"0x9771c1e34476f3f3468c8bec25b6ac9c67bc1e43a86dc37b97cc3198382a0005"}`,
		InstanceID:      "lockreleasetokenpool-aswyq",
		Type:            "LockReleaseTokenPool",
		Party:           partyCCIPBootstrapOwner,
	},
}

// ContractData holds the fetched data for a contract
type ContractData struct {
	Info        ContractInfo           `json:"info"`
	Payload     map[string]interface{} `json:"payload"`
	ContractID  string                 `json:"contract_id"`
	CreatedAt   string                 `json:"created_at"`
	Signatories []string               `json:"signatories"`
	Observers   interface{}            `json:"observers"`
	Template    map[string]interface{} `json:"template"`
	Error       string                 `json:"error,omitempty"`
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get working directory: %w", err)
	}

	outputDir := filepath.Join(wd, "scripts", "staging", "fetch_active_contract_by_instance_address")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	outputFile := filepath.Join(outputDir, "consolidated_report.html")

	fmt.Printf("Fetching %d contracts...\n\n", len(contractsToFetch))

	var contractsData []ContractData

	for i, contract := range contractsToFetch {
		fmt.Printf("[%d/%d] Fetching %s (%s)...\n", i+1, len(contractsToFetch), contract.Type, contract.InstanceID)

		data := fetchContract(contract)
		contractsData = append(contractsData, data)

		if data.Error != "" {
			fmt.Printf("    ERROR: %s\n", data.Error)
		} else {
			fmt.Printf("    OK: %s\n", truncate(data.ContractID, 30))
		}
	}

	fmt.Printf("\nGenerating consolidated report...\n")
	html := generateHTML(contractsData)

	if err := os.WriteFile(outputFile, []byte(html), 0644); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	fmt.Printf("Saved to %s\n", outputFile)
	return nil
}

func fetchContract(contract ContractInfo) ContractData {
	data := ContractData{Info: contract}

	cmd := exec.Command(
		"go", "run",
		"./scripts/staging/fetch_active_contract_by_instance_address",
		"--format", "json",
		"--template", contract.TemplateID,
		"--instance-address", contract.InstanceAddress,
		"--instance-id", contract.InstanceID,
		"--party", contract.Party,
	)

	// Capture only stdout (JSON output), stderr has the logs
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			// Command failed, stderr is in exitErr.Stderr
			stderr := string(exitErr.Stderr)
			data.Error = fmt.Sprintf("fetch failed: %v (stderr: %s)", err, truncate(stderr, 200))
		} else {
			data.Error = fmt.Sprintf("fetch failed: %v", err)
		}
		return data
	}

	var result map[string]interface{}
	if err := json.Unmarshal(output, &result); err != nil {
		preview := truncate(string(output), 200)
		data.Error = fmt.Sprintf("parse JSON failed: %v (output: %s)", err, preview)
		return data
	}

	if envelope, ok := result["envelope"].(map[string]interface{}); ok {
		if cid, ok := envelope["contractId"].(string); ok {
			data.ContractID = cid
		}
		if created, ok := envelope["createdAt"].(string); ok {
			data.CreatedAt = created
		}
		if sigs, ok := envelope["signatories"].([]interface{}); ok {
			for _, s := range sigs {
				if str, ok := s.(string); ok {
					data.Signatories = append(data.Signatories, str)
				}
			}
		}
		data.Observers = envelope["observers"]
		if tmpl, ok := envelope["templateId"].(map[string]interface{}); ok {
			data.Template = tmpl
		}
	}

	// The actual contract payload is nested in result.payload.payload
	// result has: envelope (metadata) and payload (which has contractId, createdAt, payload (actual data))
	if outerPayload, ok := result["payload"].(map[string]interface{}); ok {
		if innerPayload, ok := outerPayload["payload"].(map[string]interface{}); ok {
			data.Payload = innerPayload
		} else {
			// Fallback: use the outer payload if no nested payload
			data.Payload = outerPayload
		}
	} else {
		// Fallback: copy everything except envelope
		data.Payload = make(map[string]interface{})
		for k, v := range result {
			if k != "envelope" {
				data.Payload[k] = v
			}
		}
	}

	return data
}

func generateHTML(contracts []ContractData) string {
	contractsJSONBytes, _ := json.Marshal(contracts)
	contractsJSON := string(contractsJSONBytes)
	generatedTime := time.Now().Format("2006-01-02 15:04:05")

	var sb strings.Builder

	// HTML head and styles
	sb.WriteString(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>CCIP Contracts - Consolidated Report</title>
<style>
:root {
  --bg: #0b0d12; --top: #121622; --row: #151a24; --row-alt: #1a2130;
  --border: #2a3344; --text: #e9eef6; --muted: #8b96a8; --accent: #5b9cfa;
  --good: #3ecf8e; --bad: #f87171; --mono: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
}
@media (prefers-color-scheme: light) {
  :root { --bg: #f5f7fa; --top: #fff; --row: #fff; --row-alt: #f4f6fa; --border: #cfd6e4; --text: #141a22; --muted: #5a6578; --accent: #1d4ed8; --good: #059669; --bad: #dc2626; }
}
* { box-sizing: border-box; margin: 0; padding: 0; }
html, body { height: 100%; font-family: system-ui, -apple-system, "Segoe UI", Roboto, sans-serif; background: var(--bg); color: var(--text); font-size: 14px; line-height: 1.5; }
.header { background: var(--top); border-bottom: 1px solid var(--border); padding: 1rem 1.5rem; position: sticky; top: 0; z-index: 100; }
.header-content { max-width: 1400px; margin: 0 auto; display: flex; align-items: center; gap: 1rem; flex-wrap: wrap; }
.header h1 { font-size: 1.25rem; font-weight: 600; }
.header-meta { color: var(--muted); font-size: 0.85rem; margin-left: auto; }
.dropdown-container { display: flex; align-items: center; gap: 0.5rem; }
.dropdown-label { font-weight: 500; color: var(--muted); }
select#contractSelect { background: var(--row); color: var(--text); border: 1px solid var(--border); border-radius: 6px; padding: 0.5rem 2rem 0.5rem 0.75rem; font-size: 0.95rem; cursor: pointer; min-width: 280px; }
select#contractSelect:focus { outline: none; border-color: var(--accent); }
.main { max-width: 1400px; margin: 0 auto; padding: 1.5rem; }
.contract-card { background: var(--row); border: 1px solid var(--border); border-radius: 12px; overflow: hidden; display: none; }
.contract-card.active { display: block; }
.contract-card.error { border-color: var(--bad); }
.card-header { background: var(--row-alt); padding: 1rem 1.25rem; border-bottom: 1px solid var(--border); display: flex; align-items: center; gap: 0.75rem; flex-wrap: wrap; }
.contract-type { background: var(--accent); color: white; padding: 0.35rem 0.75rem; border-radius: 6px; font-size: 0.8rem; font-weight: 600; text-transform: uppercase; }
.contract-id { font-family: var(--mono); font-size: 0.85rem; color: var(--muted); }
.error-badge { background: var(--bad); color: white; padding: 0.35rem 0.75rem; border-radius: 6px; font-size: 0.8rem; font-weight: 600; }
.card-body { padding: 1.25rem; }
.meta-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(300px, 1fr)); gap: 1rem; margin-bottom: 1.5rem; }
.meta-item { background: var(--bg); border: 1px solid var(--border); border-radius: 8px; padding: 0.875rem 1rem; }
.meta-label { font-size: 0.75rem; color: var(--muted); text-transform: uppercase; letter-spacing: 0.5px; margin-bottom: 0.35rem; }
.meta-value { font-family: var(--mono); font-size: 0.85rem; word-break: break-all; }
.section { margin-top: 1.5rem; }
.section-title { font-size: 0.9rem; font-weight: 600; color: var(--accent); text-transform: uppercase; letter-spacing: 0.5px; margin-bottom: 0.75rem; padding-bottom: 0.5rem; border-bottom: 1px solid var(--border); }
.data-table { width: 100%; border-collapse: collapse; font-size: 13px; }
.data-table th, .data-table td { text-align: left; padding: 10px 12px; border-bottom: 1px solid var(--border); }
.data-table th { background: var(--row-alt); font-weight: 600; color: var(--muted); font-size: 0.75rem; text-transform: uppercase; }
.data-table tr:hover { background: var(--bg); }
.data-table td { font-family: var(--mono); font-size: 0.85rem; vertical-align: top; }
.val-str { color: var(--text); font-weight: 600; } .val-num { color: #f6ad55; } .val-bool { color: #63b3ed; } .val-null { color: var(--muted); font-style: italic; }
.val-empty { color: var(--muted); } .val-json { color: #b794f4; }
.nested-value { margin: 4px 0; }
.nested-value .val-json { opacity: 0.7; }
.error-message { background: rgba(248, 113, 113, 0.1); border: 1px solid var(--bad); color: var(--bad); padding: 1rem; border-radius: 8px; margin: 1rem 0; }
.summary-bar { display: flex; gap: 1rem; margin-bottom: 1rem; flex-wrap: wrap; }
.summary-stat { background: var(--row); border: 1px solid var(--border); border-radius: 8px; padding: 0.75rem 1rem; display: flex; align-items: center; gap: 0.5rem; }
.summary-stat .number { font-size: 1.25rem; font-weight: 600; color: var(--accent); }
.summary-stat .label { font-size: 0.8rem; color: var(--muted); }
@media (max-width: 768px) { .header-content { flex-direction: column; align-items: flex-start; } .meta-grid { grid-template-columns: 1fr; } .main { padding: 1rem; } }
</style>
</head>
<body>
<div class="header">
  <div class="header-content">
    <h1>CCIP Contracts Report</h1>
    <div class="dropdown-container">
      <span class="dropdown-label">Contract:</span>
      <select id="contractSelect"><option value="">Select a contract...</option></select>
    </div>
    <div class="header-meta">Generated `)
	sb.WriteString(generatedTime)
	sb.WriteString(`</div>
  </div>
</div>
<div class="main">
  <div class="summary-bar" id="summaryBar"></div>
  <div id="contractsContainer"></div>
</div>
<script>
const contractsData = `)
	sb.WriteString(contractsJSON)
	sb.WriteString(`;

function init() {
  populateDropdown();
  renderContracts();
  updateSummary();
  const firstSuccess = contractsData.find(function(c) { return !c.error; });
  if (firstSuccess) selectContract(firstSuccess.info.instance_id);
}

function populateDropdown() {
  const select = document.getElementById('contractSelect');
  contractsData.forEach(function(contract) {
    const option = document.createElement('option');
    option.value = contract.info.instance_id;
    option.textContent = contract.info.type + ' (' + contract.info.instance_id + ')' + (contract.error ? ' [ERROR]' : '');
    select.appendChild(option);
  });
  select.addEventListener('change', function() { selectContract(this.value); });
}

function selectContract(instanceId) {
  if (!instanceId) return;
  document.getElementById('contractSelect').value = instanceId;
  document.querySelectorAll('.contract-card').forEach(function(card) { card.classList.remove('active'); });
  const selected = document.getElementById('card-' + instanceId);
  if (selected) selected.classList.add('active');
}

function updateSummary() {
  const success = contractsData.filter(function(c) { return !c.error; }).length;
  const fail = contractsData.filter(function(c) { return c.error; }).length;
  let html = '<div class="summary-stat"><span class="number">' + contractsData.length + '</span><span class="label">Total</span></div>';
  html += '<div class="summary-stat"><span class="number" style="color:var(--good)">' + success + '</span><span class="label">Success</span></div>';
  if (fail > 0) html += '<div class="summary-stat"><span class="number" style="color:var(--bad)">' + fail + '</span><span class="label">Failed</span></div>';
  document.getElementById('summaryBar').innerHTML = html;
}

function renderContracts() {
  const container = document.getElementById('contractsContainer');
  contractsData.forEach(function(contract) {
    container.appendChild(createCard(contract));
  });
}

function createCard(contract) {
  const div = document.createElement('div');
  div.id = 'card-' + contract.info.instance_id;
  div.className = 'contract-card' + (contract.error ? ' error' : '');
  if (contract.error) {
    div.innerHTML = '<div class="card-header"><span class="contract-type">' + contract.info.type + '</span><span class="contract-id">' + contract.info.instance_id + '</span><span class="error-badge">Error</span></div><div class="card-body"><div class="error-message"><strong>Failed:</strong> ' + escapeHtml(contract.error) + '</div></div>';
  } else {
    div.innerHTML = renderSuccessCard(contract);
  }
  return div;
}

function renderSuccessCard(contract) {
  let html = '<div class="card-header"><span class="contract-type">' + contract.info.type + '</span><span class="contract-id">' + truncate(contract.contract_id, 40) + '</span></div>';
  html += '<div class="card-body"><div class="meta-grid">';
  html += metaItem('Instance ID', contract.info.instance_id);
  html += metaItem('Instance Address', contract.info.instance_address);
  html += metaItem('Contract ID', contract.contract_id);
  html += metaItem('Created At', formatDate(contract.created_at));
  html += metaItem('Party', contract.signatories && contract.signatories[0] ? contract.signatories[0] : 'N/A');
  html += metaItem('Template', contract.info.template_id);
  html += '</div>';
  html += '<div class="section"><div class="section-title">Contract Payload</div>';
  html += renderPayloadTable(contract.payload, contract.info.type);
  html += '</div></div>';
  return html;
}

function metaItem(label, value) {
  return '<div class="meta-item"><div class="meta-label">' + label + '</div><div class="meta-value">' + escapeHtml(value) + '</div></div>';
}

function renderPayloadTable(payload, contractType) {
  if (!payload || Object.keys(payload).length === 0) return '<p class="val-null">No payload data</p>';
  let html = '<table class="data-table"><thead><tr><th style="width: 30%;">Field</th><th>Value</th></tr></thead><tbody>';
  
  // Show all existing fields
  Object.keys(payload).forEach(function(key) {
    const val = payload[key];
    html += '<tr><td style="vertical-align: top;"><strong>' + escapeHtml(key) + '</strong></td><td style="vertical-align: top;">' + renderValue(val, 0) + '</td></tr>';
  });
  
  // For token pools, show tokenTransferFeeConfigs fields even if empty
  if (contractType && (contractType.includes('TokenPool') || contractType.includes('LockRelease') || contractType.includes('BurnMint'))) {
    const feeConfigs = payload.tokenTransferFeeConfigs;
    const hasEntries = feeConfigs && feeConfigs.entries && feeConfigs.entries.length > 0;
    
    if (!hasEntries) {
      html += '<tr><td style="vertical-align: top;"><strong>tokenTransferFeeConfigs (empty)</strong></td><td style="vertical-align: top;">';
      html += '<div style="margin-left: 10px; color: var(--muted);">';
      html += '<div>isEnabled: <span class="val-bool">false</span></div>';
      html += '<div>destGasOverhead: <span class="val-num">0</span></div>';
      html += '<div>destBytesOverhead: <span class="val-num">0</span></div>';
      html += '<div>feeUSDCents: <span class="val-num">0</span></div>';
      html += '<div>feeBps: <span class="val-num">0</span></div>';
      html += '</div></td></tr>';
    }
  }
  
  html += '</tbody></table>';
  return html;
}

function renderValue(value, indent) {
  if (indent === undefined) indent = 0;
  if (value === null) return '<span class="val-null">null</span>';
  if (typeof value === 'undefined') return '<span class="val-null">undefined</span>';
  const type = typeof value;
  if (type === 'string') {
    if (value === '') return '<span class="val-empty">(empty)</span>';
    return '<span class="val-str">"' + escapeHtml(value) + '"</span>';
  }
  if (type === 'number') return '<span class="val-num">' + value + '</span>';
  if (type === 'boolean') return '<span class="val-bool">' + value + '</span>';
  if (Array.isArray(value)) {
    if (value.length === 0) return '<span class="val-empty">[]</span>';
    // Show array contents
    let html = '<div style="margin-left: ' + (indent * 10) + 'px;">';
    html += '<span class="val-json">[</span>';
    value.forEach(function(item, idx) {
      html += '<div style="margin-left: 20px;">' + renderValue(item, indent + 1) + (idx < value.length - 1 ? '<span class="val-json">,</span>' : '') + '</div>';
    });
    html += '<span class="val-json">]</span></div>';
    return html;
  }
  if (type === 'object') {
    const keys = Object.keys(value);
    if (keys.length === 0) return '<span class="val-empty">{}</span>';
    // Special handling for genmap structure
    if (value._type === 'genmap' && value.entries) {
      return renderGenmap(value.entries, indent);
    }
    // Show object contents
    let html = '<div style="margin-left: ' + (indent * 10) + 'px;">';
    html += '<span class="val-json">{</span>';
    keys.forEach(function(key, idx) {
      html += '<div style="margin-left: 20px;"><span class="val-str">"' + escapeHtml(key) + '":</span> ' + renderValue(value[key], indent + 1) + (idx < keys.length - 1 ? '<span class="val-json">,</span>' : '') + '</div>';
    });
    html += '<span class="val-json">}</span></div>';
    return html;
  }
  return '<span>' + escapeHtml(String(value)) + '</span>';
}

function renderGenmap(entries, indent) {
  if (!entries || entries.length === 0) return '<span class="val-empty">{}</span>';
  let html = '<div style="margin-left: ' + (indent * 10) + 'px;">';
  html += '<span class="val-json">{</span><br>';
  entries.forEach(function(entry, idx) {
    html += '<div style="margin-left: 20px; border-left: 2px solid var(--border); padding-left: 10px; margin-bottom: 5px;">';
    html += '<span class="val-str">Key:</span> ' + renderValue(entry.key, indent + 1) + '<br>';
    html += '<span class="val-str">Value:</span> ' + renderValue(entry.value, indent + 1);
    html += '</div>';
    if (idx < entries.length - 1) html += '<hr style="border: none; border-top: 1px solid var(--border); margin: 5px 0;">';
  });
  html += '<span class="val-json">}</span></div>';
  return html;
}

function escapeHtml(text) {
  if (typeof text !== 'string') return String(text);
  return text.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;');
}

function truncate(str, len) {
  if (!str || str.length <= len) return str;
  return str.substring(0, len) + '...';
}

function formatDate(dateStr) {
  if (!dateStr) return 'N/A';
  try { return new Date(dateStr).toLocaleString(); } catch { return dateStr; }
}

init();
</script>
</body>
</html>`)

	return sb.String()
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
