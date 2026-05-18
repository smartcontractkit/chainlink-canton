package main

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"math"
	"sort"
	"strings"
	"time"
)

const flatMaxDepth = 64

func joinPath(prefix, key string) string {
	if prefix == "" {
		return key
	}
	return prefix + "." + key
}

func joinIndex(prefix string, i int) string {
	return fmt.Sprintf("%s[%d]", prefix, i)
}

// splitEnvelopeForHTML returns the DAML create-argument map when the tool printed the JSON envelope;
// otherwise returns out unchanged and no metadata map.
func splitEnvelopeForHTML(out any) (state any, meta map[string]any) {
	m, ok := out.(map[string]any)
	if !ok {
		return out, nil
	}
	payload, has := m["payload"]
	if !has {
		return out, nil
	}
	meta = make(map[string]any)
	for _, k := range []string{"contractId", "createdAt", "synchronizerId", "signatories", "observers", "templateId"} {
		if v, ok := m[k]; ok {
			meta[k] = v
		}
	}
	return payload, meta
}

type flatRow struct {
	Path  string
	Value string
	Kind  string // str, num, bool, null, json, empty
}

func sortedKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func isGenMap(m map[string]any) bool {
	if len(m) != 2 {
		return false
	}
	typ, ok1 := m["_type"].(string)
	_, ok2 := m["entries"]
	return ok1 && ok2 && typ == "genmap"
}

func appendRow(rows *[]flatRow, path, value, kind string) {
	if path == "" {
		path = "."
	}
	*rows = append(*rows, flatRow{Path: path, Value: value, Kind: kind})
}

func formatFloat(f float64) string {
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return fmt.Sprintf("%v", f)
	}
	if f == math.Trunc(f) && f >= math.MinInt64 && f <= math.MaxInt64 {
		return fmt.Sprintf("%.0f", f)
	}
	return fmt.Sprintf("%g", f)
}

func flattenValue(path string, v any, depth int, rows *[]flatRow) {
	if depth > flatMaxDepth {
		appendRow(rows, path, "… (max depth)", "json")
		return
	}

	switch t := v.(type) {
	case nil:
		appendRow(rows, path, "null", "null")
	case bool:
		appendRow(rows, path, fmt.Sprintf("%t", t), "bool")
	case float64:
		appendRow(rows, path, formatFloat(t), "num")
	case json.Number:
		appendRow(rows, path, t.String(), "num")
	case string:
		appendRow(rows, path, t, "str")
	case []any:
		if len(t) == 0 {
			appendRow(rows, path, "[]", "empty")
			return
		}
		for i, item := range t {
			flattenValue(joinIndex(path, i), item, depth+1, rows)
		}
	case map[string]any:
		if isGenMap(t) {
			flattenGenMap(path, t, depth, rows)
			return
		}
		if len(t) == 0 {
			appendRow(rows, path, "{}", "empty")
			return
		}
		for _, k := range sortedKeys(t) {
			flattenValue(joinPath(path, k), t[k], depth+1, rows)
		}
	default:
		b, err := json.MarshalIndent(t, "", "  ")
		if err != nil {
			appendRow(rows, path, fmt.Sprint(t), "json")
			return
		}
		appendRow(rows, path, string(b), "json")
	}
}

func flattenGenMap(path string, m map[string]any, depth int, rows *[]flatRow) {
	entries, ok := m["entries"].([]any)
	if !ok {
		b, _ := json.MarshalIndent(m, "", "  ")
		appendRow(rows, path, string(b), "json")
		return
	}
	if len(entries) == 0 {
		appendRow(rows, path, "{}", "empty")
		return
	}
	for i, item := range entries {
		em, ok := item.(map[string]any)
		if !ok {
			flattenValue(joinIndex(path, i), item, depth+1, rows)
			continue
		}
		kv, hasK := em["key"]
		val, hasV := em["value"]
		if hasK && hasV {
			keyStr := strings.TrimSpace(fmt.Sprint(kv))
			subPath := path + "[" + keyStr + "]"
			flattenValue(subPath, val, depth+1, rows)
			continue
		}
		b, _ := json.MarshalIndent(em, "", "  ")
		appendRow(rows, joinIndex(path, i), string(b), "json")
	}
}

// splitGroupAndRelative splits fullPath on the first '.' at bracket depth 0
// (e.g. destChainConfigs[1601…].a.b → group "destChainConfigs[1601…]", rel "a.b").
// If there is no such dot, group is the full path and rel is empty.
func splitGroupAndRelative(fullPath string) (group, rel string) {
	depth := 0
	for i := 0; i < len(fullPath); i++ {
		switch fullPath[i] {
		case '[':
			depth++
		case ']':
			if depth > 0 {
				depth--
			}
		case '.':
			if depth == 0 {
				return fullPath[:i], fullPath[i+1:]
			}
		}
	}
	return fullPath, ""
}

func writeGroupedFlatTables(w io.Writer, rows []flatRow) error {
	type rowItem struct {
		fullPath string
		relPath  string
		row      flatRow
	}
	groups := make(map[string][]rowItem)
	for _, r := range rows {
		g, rel := splitGroupAndRelative(r.Path)
		if rel == "" {
			rel = r.Path
		}
		groups[g] = append(groups[g], rowItem{fullPath: r.Path, relPath: rel, row: r})
	}
	groupKeys := make([]string, 0, len(groups))
	for g := range groups {
		groupKeys = append(groupKeys, g)
	}
	sort.Strings(groupKeys)

	for _, g := range groupKeys {
		items := groups[g]
		sort.Slice(items, func(i, j int) bool { return items[i].relPath < items[j].relPath })

		if _, err := fmt.Fprintf(w, `<section class="kv-group"><h3 class="kv-group-title"><code>%s</code></h3>`, html.EscapeString(g)); err != nil {
			return err
		}
		if _, err := io.WriteString(w, `<div class="table-wrap"><table class="flat"><thead><tr><th class="col-path">Field path</th><th class="col-value">Value</th></tr></thead><tbody>`); err != nil {
			return err
		}
		for _, it := range items {
			r := it.row
			kind := r.Kind
			if kind == "" {
				kind = "str"
			}
			if _, err := fmt.Fprintf(w, `<tr class="pickable" data-path="%s"><td class="col-path"><code>%s</code></td><td class="col-value val-%s">%s</td></tr>`,
				html.EscapeString(it.fullPath),
				html.EscapeString(it.relPath),
				html.EscapeString(kind),
				formatValueCell(r.Value, kind),
			); err != nil {
				return err
			}
		}
		if _, err := io.WriteString(w, `</tbody></table></div></section>`); err != nil {
			return err
		}
	}
	return nil
}

func formatValueCell(raw, kind string) string {
	esc := html.EscapeString(raw)
	switch kind {
	case "bool":
		if raw == "true" {
			return `<span class="pill pill-true">true</span>`
		}
		return `<span class="pill pill-false">false</span>`
	case "null":
		return `<span class="pill-null">null</span>`
	default:
		return esc
	}
}

func writeHTMLReport(w io.Writer, templateID, instanceIDHint, contractID string, out any, htmlRegenShell string) error {
	state, envMeta := splitEnvelopeForHTML(out)

	pageTitle := "Contract state"
	if strings.TrimSpace(templateID) != "" {
		pageTitle = strings.TrimSpace(templateID)
	}

	var hintParts []string
	if s := strings.TrimSpace(instanceIDHint); s != "" {
		hintParts = append(hintParts, `<span class="chip">instance <code>`+html.EscapeString(s)+`</code></span>`)
	}
	if s := strings.TrimSpace(contractID); s != "" {
		hintParts = append(hintParts, `<span class="chip">lookup <code>`+html.EscapeString(s)+`</code></span>`)
	}
	hints := strings.Join(hintParts, " ")
	if hints == "" {
		hints = `<span class="chip muted-chip">ACS / instance-address match</span>`
	}

	var payloadRows []flatRow
	flattenValue("", state, 0, &payloadRows)

	_, err := fmt.Fprintf(w, `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>%s</title>
<style>
:root {
  --bg: #0b0d12;
  --top: #121622;
  --row: #151a24;
  --row-alt: #1a2130;
  --border: #2a3344;
  --text: #e9eef6;
  --muted: #8b96a8;
  --accent: #5b9cfa;
  --good: #3ecf8e;
  --bad: #f87171;
  --mono: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
}
@media (prefers-color-scheme: light) {
  :root {
    --bg: #eceff5;
    --top: #ffffff;
    --row: #ffffff;
    --row-alt: #f4f6fa;
    --border: #cfd6e4;
    --text: #141a22;
    --muted: #5a6578;
    --accent: #1d4ed8;
  }
}
* { box-sizing: border-box; }
html, body { height: 100%%; margin: 0; }
body {
  font-family: system-ui, -apple-system, "Segoe UI", Roboto, sans-serif;
  background: var(--bg);
  color: var(--text);
  font-size: 14px;
  line-height: 1.45;
}
.app {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
  height: 100vh;
  max-height: 100vh;
}
.topbar {
  flex: 0 0 auto;
  background: var(--top);
  border-bottom: 1px solid var(--border);
  padding: 0.65rem 1rem 0.75rem;
  display: flex;
  flex-wrap: wrap;
  align-items: flex-end;
  gap: 0.5rem 1.25rem;
}
.topbar h1 {
  margin: 0;
  font-size: 1.05rem;
  font-weight: 650;
  flex: 1 1 12rem;
  word-break: break-word;
}
.meta-line {
  font-size: 0.78rem;
  color: var(--muted);
  white-space: nowrap;
}
.chips { display: flex; flex-wrap: wrap; gap: 0.35rem; align-items: center; }
.chip {
  font-size: 0.72rem;
  padding: 0.15rem 0.45rem;
  border-radius: 999px;
  background: var(--row-alt);
  border: 1px solid var(--border);
}
.chip code { font-family: var(--mono); font-size: 0.9em; }
.muted-chip { opacity: 0.85; }
.top-actions { display: flex; flex-wrap: wrap; gap: 0.4rem; }
button.btn {
  font: inherit;
  cursor: pointer;
  border-radius: 6px;
  border: 1px solid var(--border);
  background: var(--row-alt);
  color: var(--text);
  padding: 0.35rem 0.65rem;
  font-size: 0.8rem;
  font-weight: 600;
}
button.btn.primary {
  background: var(--accent);
  border-color: transparent;
  color: #fff;
}
#status { font-size: 0.75rem; color: var(--muted); min-width: 8rem; }
#status.err { color: var(--bad); }
#status.ok { color: var(--good); }

.main {
  flex: 1 1 auto;
  min-height: 0;
  display: flex;
  flex-direction: column;
  padding: 0;
}
.scroll {
  flex: 1 1 auto;
  min-height: 0;
  overflow: auto;
  -webkit-overflow-scrolling: touch;
}

details.env {
  flex: 0 0 auto;
  border-bottom: 1px solid var(--border);
  background: var(--top);
}
details.env > summary {
  cursor: pointer;
  list-style: none;
  padding: 0.5rem 1rem;
  font-size: 0.82rem;
  font-weight: 600;
  color: var(--muted);
  user-select: none;
}
details.env > summary::-webkit-details-marker { display: none; }
details.env[open] > summary { border-bottom: 1px solid var(--border); }

.section-title {
  margin: 0;
  padding: 0.45rem 1rem;
  font-size: 0.75rem;
  font-weight: 700;
  letter-spacing: 0.06em;
  text-transform: uppercase;
  color: var(--muted);
  background: var(--row-alt);
  border-bottom: 1px solid var(--border);
}

.kv-group {
  margin: 0 0 1rem;
  border: 1px solid var(--border);
  border-radius: 8px;
  overflow: hidden;
  background: var(--row);
}
.kv-group-title {
  margin: 0;
  padding: 0.5rem 0.75rem;
  font-size: 0.82rem;
  font-weight: 650;
  background: var(--row-alt);
  border-bottom: 1px solid var(--border);
}
.kv-group-title code {
  font-family: var(--mono);
  color: var(--accent);
  word-break: break-all;
}

.table-wrap { width: 100%%; }
table.flat {
  width: 100%%;
  table-layout: fixed;
  border-collapse: collapse;
  font-size: 0.84rem;
}
table.flat thead th {
  position: sticky;
  top: 0;
  z-index: 1;
  text-align: left;
  padding: 0.45rem 0.65rem;
  background: var(--row-alt);
  border-bottom: 2px solid var(--border);
  font-weight: 700;
  font-size: 0.7rem;
  text-transform: uppercase;
  letter-spacing: 0.04em;
  color: var(--muted);
}
table.flat tbody tr:nth-child(even) { background: rgba(255,255,255,0.02); }
table.flat tbody tr:nth-child(odd) { background: transparent; }
table.flat td {
  padding: 0.45rem 0.65rem;
  border-bottom: 1px solid var(--border);
  vertical-align: top;
}
.col-path {
  width: 34%%;
  font-family: var(--mono);
  font-size: 0.8rem;
  color: var(--muted);
  overflow-wrap: anywhere;
  word-break: break-word;
}
.col-path code { font-family: inherit; color: var(--muted); }
.col-value {
  width: 66%%;
  font-family: var(--mono);
  font-size: 0.82rem;
  white-space: pre-wrap;
  overflow-wrap: anywhere;
  word-break: break-all;
}
.val-num { color: var(--accent); }
.pill {
  display: inline-block;
  padding: 0.1rem 0.4rem;
  border-radius: 4px;
  font-weight: 700;
  font-size: 0.78rem;
}
.pill-true { color: var(--good); border: 1px solid var(--border); }
.pill-false { color: var(--bad); border: 1px solid var(--border); }
.pill-null { color: var(--muted); font-style: italic; }

tr.pickable { cursor: pointer; }
tr.pickable:hover td { background: rgba(91, 156, 250, 0.1); }

.drawer {
  flex: 0 0 auto;
  border-top: 1px solid var(--border);
  background: var(--top);
  padding: 0.6rem 1rem 0.85rem;
  max-height: 38vh;
  overflow: auto;
}
.drawer h2 {
  margin: 0 0 0.4rem;
  font-size: 0.82rem;
  font-weight: 650;
}
.drawer .sub { margin: 0 0 0.5rem; font-size: 0.78rem; color: var(--muted); }
.proposal-grid { display: grid; gap: 0.4rem; }
.proposal-grid label {
  display: flex;
  flex-direction: column;
  gap: 0.2rem;
  font-size: 0.72rem;
  color: var(--muted);
}
.proposal-grid input,
.proposal-grid textarea {
  font: inherit;
  font-size: 0.82rem;
  padding: 0.35rem 0.45rem;
  border-radius: 6px;
  border: 1px solid var(--border);
  background: var(--bg);
  color: var(--text);
}
.proposal-grid textarea { min-height: 3rem; resize: vertical; }
.drawer-actions { display: flex; flex-wrap: wrap; gap: 0.4rem; margin-top: 0.45rem; }
</style>
</head>
<body>
<div class="app">
<header class="topbar">
  <h1>%s</h1>
  <span class="meta-line">Generated %s (UTC)</span>
  <div class="chips">%s</div>
  <div class="top-actions">
    <button type="button" class="btn primary" id="btn-reload">Reload</button>
    <button type="button" class="btn" id="btn-copy-cmd">Copy command</button>
    <span id="status"></span>
  </div>
</header>
<div class="main">
<div class="scroll">
`, html.EscapeString(pageTitle), html.EscapeString(pageTitle), html.EscapeString(time.Now().UTC().Format(time.RFC3339Nano)), hints)
	if err != nil {
		return err
	}

	if len(envMeta) > 0 {
		if _, err := io.WriteString(w, `<details class="env"><summary>Ledger envelope (contract id, parties, template, …)</summary><div class="section-title">Envelope — flat</div>`); err != nil {
			return err
		}
		var envRows []flatRow
		flattenValue("envelope", envMeta, 0, &envRows)
		if err := writeGroupedFlatTables(w, envRows); err != nil {
			return err
		}
		if _, err := io.WriteString(w, `</details>`); err != nil {
			return err
		}
	}

	if _, err := io.WriteString(w, `<div class="section-title">Contract payload — one group per top-level key (Field path | Value)</div>`); err != nil {
		return err
	}
	if err := writeGroupedFlatTables(w, payloadRows); err != nil {
		return err
	}
	if _, err := io.WriteString(w, `</div></div>`); err != nil {
		return err
	}

	regenJSON, jerr := json.Marshal(strings.TrimSpace(htmlRegenShell))
	if jerr != nil {
		regenJSON = []byte(`""`)
	}

	script := fmt.Sprintf(`<div class="drawer">
<h2>Propose a change</h2>
<p class="sub">Click a row above to fill path and current value. Copy proposal puts one JSON line on the clipboard.</p>
<div class="proposal-grid">
  <label>Field path <input id="pp-path" autocomplete="off" placeholder="destChainConfigs[160152…].value.isEnabled" /></label>
  <label>Current value <textarea id="pp-current" placeholder="From table or paste"></textarea></label>
  <label>Desired value <textarea id="pp-proposed" placeholder="What you want instead"></textarea></label>
  <label>Notes <textarea id="pp-note" placeholder="Ticket / PR"></textarea></label>
</div>
<div class="drawer-actions">
  <button type="button" class="btn primary" id="btn-propose">Copy proposal</button>
  <button type="button" class="btn" id="btn-clear-proposal">Clear</button>
</div>
</div>
</div>
<script>
(function(){
  const REGEN_CMD = %s;
  function setStatus(text, cls) {
    var el = document.getElementById('status');
    if (!el) return;
    el.className = cls || '';
    el.textContent = text || '';
  }
  function copyText(t) {
    if (!t) { setStatus('Nothing to copy', 'err'); return; }
    navigator.clipboard.writeText(t).then(function() {
      setStatus('Copied.', 'ok');
    }).catch(function(){ setStatus('Copy failed', 'err'); });
  }
  function saveProposal() {
    var body = {
      field_path: document.getElementById('pp-path').value,
      current: document.getElementById('pp-current').value,
      proposed: document.getElementById('pp-proposed').value,
      note: document.getElementById('pp-note').value
    };
    if (!body.field_path.trim()) { setStatus('Path required', 'err'); return; }
    var line = JSON.stringify({ ts: new Date().toISOString(), field_path: body.field_path, current: body.current, proposed: body.proposed, note: body.note });
    copyText(line + '\n');
  }
  function clearProposal() {
    ['pp-path','pp-current','pp-proposed','pp-note'].forEach(function(id){
      var el = document.getElementById(id);
      if (el) el.value = '';
    });
    setStatus('', '');
  }
  function pickPath(ev) {
    var tr = ev.target.closest('tr.pickable[data-path]');
    if (!tr || !tr.dataset.path) return;
    document.getElementById('pp-path').value = tr.dataset.path;
    var cells = tr.querySelectorAll('td');
    var v = cells.length > 1 ? cells[1].innerText : '';
    document.getElementById('pp-current').value = (v || '').trim();
    setStatus('Loaded row', 'ok');
  }
  document.getElementById('btn-reload').addEventListener('click', function(){ location.reload(); });
  var copyBtn = document.getElementById('btn-copy-cmd');
  if (!REGEN_CMD) { copyBtn.style.display = 'none'; }
  else { copyBtn.addEventListener('click', function(){ copyText(REGEN_CMD); }); }
  document.getElementById('btn-propose').addEventListener('click', saveProposal);
  document.getElementById('btn-clear-proposal').addEventListener('click', clearProposal);
  document.addEventListener('click', pickPath);
})();
</script>
</body></html>`, string(regenJSON))

	if _, err := io.WriteString(w, script); err != nil {
		return err
	}

	return nil
}
