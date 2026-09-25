package handlers

import "net/http"

const dashboardHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<title>Compliance Monitor</title>
<style>
  :root { color-scheme: light dark; }
  body { font-family: -apple-system, Segoe UI, Roboto, sans-serif; max-width: 960px; margin: 40px auto; padding: 0 20px; }
  h1 { font-size: 1.4rem; margin-bottom: 4px; }
  .sub { color: #888; margin-bottom: 24px; font-size: 0.9rem; }
  .toolbar { display: flex; gap: 10px; margin-bottom: 20px; flex-wrap: wrap; align-items: center; }
  button { padding: 8px 14px; border: 1px solid #999; border-radius: 6px; background: transparent; cursor: pointer; font-size: 0.9rem; }
  button:hover { background: rgba(128,128,128,0.15); }
  table { width: 100%; border-collapse: collapse; font-size: 0.88rem; }
  th, td { text-align: left; padding: 8px 10px; border-bottom: 1px solid rgba(128,128,128,0.25); vertical-align: top; }
  th { color: #888; font-weight: 600; font-size: 0.78rem; text-transform: uppercase; letter-spacing: 0.03em; }
  .sev-high { color: #d64545; font-weight: 600; }
  .sev-medium { color: #c78a1f; font-weight: 600; }
  .badge { display: inline-block; padding: 2px 8px; border-radius: 10px; font-size: 0.75rem; background: rgba(128,128,128,0.2); }
  .resolved { opacity: 0.45; }
  .empty { color: #888; padding: 30px 0; text-align: center; }
  .status { font-size: 0.8rem; color: #888; margin-bottom: 10px; }
  .resolve-btn { padding: 4px 10px; font-size: 0.78rem; }
</style>
</head>
<body>
  <h1>Compliance Monitor</h1>
  <div class="sub">Live alerts from the transaction monitoring engine.</div>

  <div class="toolbar">
    <button onclick="seed()">Seed demo data</button>
    <button onclick="load()">Refresh</button>
    <label style="font-size:0.85rem; display:flex; align-items:center; gap:6px; margin-left:auto;">
      <input type="checkbox" id="hideResolved" onchange="load()"> Hide resolved
    </label>
  </div>

  <div class="status" id="status">Loading...</div>
  <div id="content"></div>

<script>
async function seed() {
  document.getElementById('status').textContent = 'Seeding...';
  await fetch('/demo/seed', { method: 'POST' });
  setTimeout(load, 600);
}

async function resolveAlert(id) {
  await fetch('/alerts/' + id + '/resolve', { method: 'POST' });
  load();
}

async function load() {
  const statusEl = document.getElementById('status');
  const hideResolved = document.getElementById('hideResolved').checked;
  try {
    const url = hideResolved ? '/alerts?resolved=false' : '/alerts';
    const res = await fetch(url);
    const alerts = await res.json();
    render(alerts);
    statusEl.textContent = alerts.length + ' alert(s) - last updated ' + new Date().toLocaleTimeString();
  } catch (e) {
    statusEl.textContent = 'Could not reach the server.';
  }
}

function render(alerts) {
  const content = document.getElementById('content');
  if (!alerts || alerts.length === 0) {
    content.innerHTML = '<div class="empty">No alerts yet. Click "Seed demo data" to generate some.</div>';
    return;
  }
  alerts.sort((a, b) => new Date(b.created_at) - new Date(a.created_at));
  let rows = alerts.map(a => {
    const sevClass = a.severity === 'high' ? 'sev-high' : 'sev-medium';
    const rowClass = a.resolved ? 'resolved' : '';
    const time = new Date(a.created_at).toLocaleString();
    const action = a.resolved
      ? '<span class="badge">resolved</span>'
      : '<button class="resolve-btn" onclick="resolveAlert(\'' + a.id + '\')">Resolve</button>';
    return '<tr class="' + rowClass + '">' +
      '<td>' + time + '</td>' +
      '<td>' + escapeHTML(a.user_id) + '</td>' +
      '<td><span class="badge">' + escapeHTML(a.rule_name) + '</span></td>' +
      '<td class="' + sevClass + '">' + escapeHTML(a.severity) + '</td>' +
      '<td>' + escapeHTML(a.reason) + '</td>' +
      '<td>' + action + '</td>' +
      '</tr>';
  }).join('');

  content.innerHTML = '<table><thead><tr>' +
    '<th>Time</th><th>User</th><th>Rule</th><th>Severity</th><th>Reason</th><th></th>' +
    '</tr></thead><tbody>' + rows + '</tbody></table>';
}

function escapeHTML(s) {
  const d = document.createElement('div');
  d.innerText = s == null ? '' : s;
  return d.innerHTML;
}

load();
setInterval(load, 5000);
</script>
</body>
</html>`

func (a *API) Dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(dashboardHTML))
}