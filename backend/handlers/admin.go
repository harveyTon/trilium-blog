package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *APIHandler) AdminPage(c *gin.Context) {
	if h.adminToken == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "admin not configured"})
		return
	}
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, adminHTML)
}

var adminHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>Cache Admin</title>
<style>
*{margin:0;padding:0;box-sizing:border-box}
body{font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,sans-serif;background:#f5f5f5;color:#333;padding:20px;max-width:720px;margin:0 auto}
h1{font-size:1.4rem;margin-bottom:16px;color:#111}
h2{font-size:1.1rem;margin:20px 0 10px;color:#444}
.card{background:#fff;border:1px solid #ddd;border-radius:6px;padding:16px;margin-bottom:12px}
table{width:100%;border-collapse:collapse;font-size:.875rem}
th,td{text-align:left;padding:6px 10px;border-bottom:1px solid #eee}
th{color:#888;font-weight:500;font-size:.8rem;text-transform:uppercase}
input[type=password]{width:100%;padding:8px 10px;border:1px solid #ccc;border-radius:4px;font-size:.9rem;margin-bottom:10px}
input[type=text]{width:100%;padding:8px 10px;border:1px solid #ccc;border-radius:4px;font-size:.9rem;margin-bottom:10px}
.btn{display:inline-block;padding:7px 14px;border:none;border-radius:4px;cursor:pointer;font-size:.85rem;margin:2px}
.btn-primary{background:#2563eb;color:#fff}.btn-primary:hover{background:#1d4ed8}
.btn-danger{background:#dc2626;color:#fff}.btn-danger:hover{background:#b91c1c}
.btn-secondary{background:#6b7280;color:#fff}.btn-secondary:hover{background:#4b5563}
.btn-sm{padding:4px 10px;font-size:.78rem}
.actions{display:flex;flex-wrap:wrap;gap:6px;margin:10px 0}
.status{display:inline-block;width:8px;height:8px;border-radius:50%;margin-right:6px}
.status-ok{background:#22c55e}.status-err{background:#ef4444}
.msg{padding:8px 12px;border-radius:4px;font-size:.85rem;margin:8px 0;display:none}
.msg-ok{background:#dcfce7;color:#166534;display:block}
.msg-err{background:#fef2f2;color:#991b1b;display:block}
.login-wrap{max-width:360px;margin:80px auto}
.tag{display:inline-block;padding:2px 8px;border-radius:3px;font-size:.75rem;background:#e5e7eb;color:#555;margin-left:4px}
.hidden{display:none!important}
</style>
</head>
<body>

<div id="login" class="login-wrap">
  <h1 style="text-align:center;margin-bottom:20px">Cache Admin</h1>
  <div class="card">
    <input type="password" id="token-input" placeholder="Enter admin token" autofocus>
    <button class="btn btn-primary" style="width:100%" onclick="login()">Login</button>
    <div id="login-msg" class="msg"></div>
  </div>
</div>

<div id="dashboard" class="hidden">
  <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:12px">
    <h1 style="margin:0">Cache Admin</h1>
    <button class="btn btn-secondary btn-sm" onclick="logout()">Logout</button>
  </div>

  <h2>Status</h2>
  <div id="redis-status" class="card" style="font-size:.9rem"></div>

  <h2>Cache Entries</h2>
  <div class="card">
    <table>
      <thead><tr><th>Type</th><th>Keys</th><th>Min TTL</th><th>Max TTL</th><th>Default TTL</th><th>Action</th></tr></thead>
      <tbody id="cache-table"></tbody>
    </table>
  </div>

  <h2>Invalidate</h2>
  <div class="card">
    <div class="actions">
      <button class="btn btn-danger" onclick="invalidateAll()">Clear All Cache</button>
      <button class="btn btn-primary" onclick="triggerPreload()">Trigger Preload</button>
    </div>
    <div style="margin-top:12px">
      <label style="font-size:.8rem;color:#666">Invalidate by note ID</label>
      <div style="display:flex;gap:6px">
        <input type="text" id="note-id-input" placeholder="Note ID" style="flex:1">
        <button class="btn btn-secondary" onclick="invalidateNote()">Clear Note</button>
      </div>
    </div>
    <div style="margin-top:8px">
      <label style="font-size:.8rem;color:#666">Invalidate by attachment ID</label>
      <div style="display:flex;gap:6px">
        <input type="text" id="attach-id-input" placeholder="Attachment ID" style="flex:1">
        <button class="btn btn-secondary" onclick="invalidateAttachment()">Clear Attachment</button>
      </div>
    </div>
    <div id="action-msg" class="msg"></div>
  </div>
</div>

<script>
const LS_KEY = 'cache_admin_token';
let token = localStorage.getItem(LS_KEY) || '';

function api(method, path, body) {
  const opts = {
    method,
    headers: {'Authorization': 'Bearer ' + token},
  };
  if (body) {
    opts.headers['Content-Type'] = 'application/json';
    opts.body = JSON.stringify(body);
  }
  return fetch('/api/admin' + path, opts).then(r => {
    if (r.status === 401 || r.status === 403) {
      localStorage.removeItem(LS_KEY);
      token = '';
      showLogin();
      throw new Error('Unauthorized');
    }
    return r.json().then(j => ({status: r.status, data: j}));
  });
}

function showLogin() {
  document.getElementById('login').classList.remove('hidden');
  document.getElementById('dashboard').classList.add('hidden');
}

function showDash() {
  document.getElementById('login').classList.add('hidden');
  document.getElementById('dashboard').classList.remove('hidden');
  loadStats();
}

function showMsg(id, text, ok) {
  const el = document.getElementById(id);
  el.textContent = text;
  el.className = 'msg ' + (ok ? 'msg-ok' : 'msg-err');
  setTimeout(() => { el.className = 'msg'; }, 3000);
}

function login() {
  token = document.getElementById('token-input').value.trim();
  if (!token) return;
  localStorage.setItem(LS_KEY, token);
  api('GET', '/cache/stats').then(() => showDash()).catch(e => {
    showMsg('login-msg', 'Invalid token', false);
    localStorage.removeItem(LS_KEY);
  });
}

function logout() {
  localStorage.removeItem(LS_KEY);
  token = '';
  showLogin();
}

function loadStats() {
  api('GET', '/cache/stats').then(r => {
    const s = r.data;
    const dot = s.redisConnected ? '<span class="status status-ok"></span>Connected' : '<span class="status status-err"></span>Disconnected';
    document.getElementById('redis-status').innerHTML = dot;
    const tbody = document.getElementById('cache-table');
    tbody.innerHTML = '';
    (s.types || []).forEach(t => {
      const tr = document.createElement('tr');
      tr.innerHTML = '<td><b>' + t.name + '</b></td><td>' + t.keyCount + '</td><td>' + (t.minTTL || '—') + '</td><td>' + (t.maxTTL || '—') + '</td><td>' + t.ttlSeconds + 's</td><td><button class="btn btn-danger btn-sm" onclick="invalidateType(\'' + t.name + '\')">Clear</button></td>';
      tbody.appendChild(tr);
    });
  }).catch(() => {});
}

function invalidateAll() {
  if (!confirm('Clear all cache?')) return;
  api('POST', '/cache/invalidate', {scope: 'all'}).then(r => {
    showMsg('action-msg', 'Cleared ' + r.data.keys_removed + ' keys', true);
    loadStats();
  }).catch(() => {});
}

function invalidateType(name) {
  if (!confirm('Clear cache type: ' + name + '?')) return;
  api('POST', '/cache/invalidate', {scope: 'type', type: name}).then(r => {
    showMsg('action-msg', 'Cleared ' + r.data.keys_removed + ' keys from ' + name, true);
    loadStats();
  }).catch(() => {});
}

function invalidateNote() {
  const id = document.getElementById('note-id-input').value.trim();
  if (!id) { showMsg('action-msg', 'Enter a note ID', false); return; }
  api('POST', '/cache/invalidate', {scope: 'note', id: id}).then(r => {
    showMsg('action-msg', 'Cleared note: ' + id, true);
    document.getElementById('note-id-input').value = '';
    loadStats();
  }).catch(() => {});
}

function invalidateAttachment() {
  const id = document.getElementById('attach-id-input').value.trim();
  if (!id) { showMsg('action-msg', 'Enter an attachment ID', false); return; }
  api('POST', '/cache/invalidate', {scope: 'attachment', id: id}).then(r => {
    showMsg('action-msg', 'Cleared attachment: ' + id, true);
    document.getElementById('attach-id-input').value = '';
    loadStats();
  }).catch(() => {});
}

function triggerPreload() {
  api('POST', '/cache/preload').then(r => {
    showMsg('action-msg', r.data.status || 'Preload triggered', true);
  }).catch(e => {
    if (e.message !== 'Unauthorized') {
      showMsg('action-msg', 'Preload already in progress', false);
    }
  });
}

if (token) {
  showDash();
} else {
  showLogin();
}
</script>
</body>
</html>`
