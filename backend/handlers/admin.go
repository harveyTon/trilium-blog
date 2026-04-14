package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (h *APIHandler) AdminPage(c *gin.Context) {
	if h.adminToken == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "admin not configured"})
		return
	}
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, renderAdminHTML(h.locale))
}

type i18nMap map[string]string

var adminI18n = map[string]i18nMap{
	"title":             {"zh-CN": "缓存管理", "en": "Cache Admin"},
	"enterToken":        {"zh-CN": "输入管理令牌", "en": "Enter admin token"},
	"login":             {"zh-CN": "登录", "en": "Login"},
	"logout":            {"zh-CN": "退出", "en": "Logout"},
	"status":            {"zh-CN": "状态", "en": "Status"},
	"connected":         {"zh-CN": "已连接", "en": "Connected"},
	"disconnected":      {"zh-CN": "未连接", "en": "Disconnected"},
	"cacheEntries":      {"zh-CN": "缓存条目", "en": "Cache Entries"},
	"type":              {"zh-CN": "类型", "en": "Type"},
	"keys":              {"zh-CN": "键数", "en": "Keys"},
	"minTTL":            {"zh-CN": "最小 TTL", "en": "Min TTL"},
	"maxTTL":            {"zh-CN": "最大 TTL", "en": "Max TTL"},
	"defaultTTL":        {"zh-CN": "默认 TTL", "en": "Default TTL"},
	"action":            {"zh-CN": "操作", "en": "Action"},
	"clear":             {"zh-CN": "清除", "en": "Clear"},
	"invalidate":        {"zh-CN": "缓存失效", "en": "Invalidate"},
	"clearAll":          {"zh-CN": "清除全部缓存", "en": "Clear All Cache"},
	"triggerPreload":    {"zh-CN": "触发预加载", "en": "Trigger Preload"},
	"byNoteID":          {"zh-CN": "按笔记 ID 失效", "en": "Invalidate by note ID"},
	"clearNote":         {"zh-CN": "清除笔记", "en": "Clear Note"},
	"byAttachmentID":    {"zh-CN": "按附件 ID 失效", "en": "Invalidate by attachment ID"},
	"clearAttachment":   {"zh-CN": "清除附件", "en": "Clear Attachment"},
	"invalidToken":      {"zh-CN": "无效令牌", "en": "Invalid token"},
	"clearAllConfirm":   {"zh-CN": "确认清除全部缓存？", "en": "Clear all cache?"},
	"clearedKeys":       {"zh-CN": "已清除 %d 个键", "en": "Cleared %d keys"},
	"clearedKeysFrom":   {"zh-CN": "已从 %s 清除 %d 个键", "en": "Cleared %d keys from %s"},
	"clearTypeConfirm":  {"zh-CN": "确认清除缓存类型: %s？", "en": "Clear cache type: %s?"},
	"clearedNote":       {"zh-CN": "已清除笔记: %s", "en": "Cleared note: %s"},
	"enterNoteID":       {"zh-CN": "请输入笔记 ID", "en": "Enter a note ID"},
	"clearedAttachment": {"zh-CN": "已清除附件: %s", "en": "Cleared attachment: %s"},
	"enterAttachID":     {"zh-CN": "请输入附件 ID", "en": "Enter an attachment ID"},
	"preloadTriggered":  {"zh-CN": "预加载已触发", "en": "Preload triggered"},
	"preloadInProgress": {"zh-CN": "预加载正在进行中", "en": "Preload already in progress"},
}

func t(locale, key string) string {
	m, ok := adminI18n[key]
	if !ok {
		return key
	}
	if v, ok := m[locale]; ok {
		return v
	}
	if v, ok := m["en"]; ok {
		return v
	}
	return key
}

func renderAdminHTML(locale string) string {
	if locale == "" {
		locale = "en"
	}
	lang := "en"
	if strings.HasPrefix(locale, "zh") {
		lang = "zh-CN"
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="%s">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>%s</title>
<style>
*{margin:0;padding:0;box-sizing:border-box}
body{font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,sans-serif;background:#f5f5f5;color:#333;padding:20px;max-width:720px;margin:0 auto}
h1{font-size:1.4rem;margin-bottom:16px;color:#111}
h2{font-size:1.1rem;margin:20px 0 10px;color:#444}
.card{background:#fff;border:1px solid #ddd;border-radius:6px;padding:16px;margin-bottom:12px}
table{width:100%%;border-collapse:collapse;font-size:.875rem}
th,td{text-align:left;padding:6px 10px;border-bottom:1px solid #eee}
th{color:#888;font-weight:500;font-size:.8rem;text-transform:uppercase}
input[type=password]{width:100%%;padding:8px 10px;border:1px solid #ccc;border-radius:4px;font-size:.9rem;margin-bottom:10px}
input[type=text]{width:100%%;padding:8px 10px;border:1px solid #ccc;border-radius:4px;font-size:.9rem;margin-bottom:10px}
.btn{display:inline-block;padding:7px 14px;border:none;border-radius:4px;cursor:pointer;font-size:.85rem;margin:2px}
.btn-primary{background:#2563eb;color:#fff}.btn-primary:hover{background:#1d4ed8}
.btn-danger{background:#dc2626;color:#fff}.btn-danger:hover{background:#b91c1c}
.btn-secondary{background:#6b7280;color:#fff}.btn-secondary:hover{background:#4b5563}
.btn-sm{padding:4px 10px;font-size:.78rem}
.actions{display:flex;flex-wrap:wrap;gap:6px;margin:10px 0}
.status{display:inline-block;width:8px;height:8px;border-radius:50%%;margin-right:6px}
.status-ok{background:#22c55e}.status-err{background:#ef4444}
.msg{padding:8px 12px;border-radius:4px;font-size:.85rem;margin:8px 0;display:none}
.msg-ok{background:#dcfce7;color:#166534;display:block}
.msg-err{background:#fef2f2;color:#991b1b;display:block}
.login-wrap{max-width:360px;margin:80px auto}
.hidden{display:none!important}
</style>
</head>
<body>

<div id="login" class="login-wrap">
  <h1 style="text-align:center;margin-bottom:20px">%s</h1>
  <div class="card">
    <input type="password" id="token-input" placeholder="%s" autofocus>
    <button class="btn btn-primary" style="width:100%%" onclick="login()">%s</button>
    <div id="login-msg" class="msg"></div>
  </div>
</div>

<div id="dashboard" class="hidden">
  <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:12px">
    <h1 style="margin:0">%s</h1>
    <button class="btn btn-secondary btn-sm" onclick="logout()">%s</button>
  </div>

  <h2>%s</h2>
  <div id="redis-status" class="card" style="font-size:.9rem"></div>

  <h2>%s</h2>
  <div class="card">
    <table>
      <thead><tr><th>%s</th><th>%s</th><th>%s</th><th>%s</th><th>%s</th><th>%s</th></tr></thead>
      <tbody id="cache-table"></tbody>
    </table>
  </div>

  <h2>%s</h2>
  <div class="card">
    <div class="actions">
      <button class="btn btn-danger" onclick="invalidateAll()">%s</button>
      <button class="btn btn-primary" onclick="triggerPreload()">%s</button>
    </div>
    <div style="margin-top:12px">
      <label style="font-size:.8rem;color:#666">%s</label>
      <div style="display:flex;gap:6px">
        <input type="text" id="note-id-input" placeholder="Note ID" style="flex:1">
        <button class="btn btn-secondary" onclick="invalidateNote()">%s</button>
      </div>
    </div>
    <div style="margin-top:8px">
      <label style="font-size:.8rem;color:#666">%s</label>
      <div style="display:flex;gap:6px">
        <input type="text" id="attach-id-input" placeholder="Attachment ID" style="flex:1">
        <button class="btn btn-secondary" onclick="invalidateAttachment()">%s</button>
      </div>
    </div>
    <div id="action-msg" class="msg"></div>
  </div>
</div>

<script>
const LS_KEY = 'cache_admin_token';
const i18n = {
  connected: %q,
  disconnected: %q,
  invalidToken: %q,
  clearAllConfirm: %q,
  clearedKeys: %q,
  clearedKeysFrom: %q,
  clearTypeConfirm: %q,
  clearedNote: %q,
  enterNoteID: %q,
  clearedAttachment: %q,
  enterAttachID: %q,
  preloadTriggered: %q,
  preloadInProgress: %q,
  clear: %q,
};
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
    showMsg('login-msg', i18n.invalidToken, false);
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
    const dot = s.redisConnected
      ? '<span class="status status-ok"></span>' + i18n.connected
      : '<span class="status status-err"></span>' + i18n.disconnected;
    document.getElementById('redis-status').innerHTML = dot;
    const tbody = document.getElementById('cache-table');
    tbody.innerHTML = '';
    (s.types || []).forEach(t => {
      const tr = document.createElement('tr');
      tr.innerHTML = '<td><b>' + t.name + '</b></td><td>' + t.keyCount + '</td><td>' + (t.minTTL || '\u2014') + '</td><td>' + (t.maxTTL || '\u2014') + '</td><td>' + t.ttlSeconds + 's</td><td><button class="btn btn-danger btn-sm" onclick="invalidateType(\'' + t.name + '\')">' + i18n.clear + '</button></td>';
      tbody.appendChild(tr);
    });
  }).catch(() => {});
}

function invalidateAll() {
  if (!confirm(i18n.clearAllConfirm)) return;
  api('POST', '/cache/invalidate', {scope: 'all'}).then(r => {
    showMsg('action-msg', i18n.clearedKeys.replace('%%d', r.data.keys_removed), true);
    loadStats();
  }).catch(() => {});
}

function invalidateType(name) {
  if (!confirm(i18n.clearTypeConfirm.replace('%%s', name))) return;
  api('POST', '/cache/invalidate', {scope: 'type', type: name}).then(r => {
    showMsg('action-msg', i18n.clearedKeysFrom.replace('%%d', r.data.keys_removed).replace('%%s', name), true);
    loadStats();
  }).catch(() => {});
}

function invalidateNote() {
  const id = document.getElementById('note-id-input').value.trim();
  if (!id) { showMsg('action-msg', i18n.enterNoteID, false); return; }
  api('POST', '/cache/invalidate', {scope: 'note', id: id}).then(r => {
    showMsg('action-msg', i18n.clearedNote.replace('%%s', id), true);
    document.getElementById('note-id-input').value = '';
    loadStats();
  }).catch(() => {});
}

function invalidateAttachment() {
  const id = document.getElementById('attach-id-input').value.trim();
  if (!id) { showMsg('action-msg', i18n.enterAttachID, false); return; }
  api('POST', '/cache/invalidate', {scope: 'attachment', id: id}).then(r => {
    showMsg('action-msg', i18n.clearedAttachment.replace('%%s', id), true);
    document.getElementById('attach-id-input').value = '';
    loadStats();
  }).catch(() => {});
}

function triggerPreload() {
  api('POST', '/cache/preload').then(r => {
    showMsg('action-msg', r.data.status || i18n.preloadTriggered, true);
  }).catch(e => {
    if (e.message !== 'Unauthorized') {
      showMsg('action-msg', i18n.preloadInProgress, false);
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
</html>`,
		lang,
		t(lang, "title"),
		t(lang, "title"),
		t(lang, "enterToken"),
		t(lang, "login"),
		t(lang, "title"),
		t(lang, "logout"),
		t(lang, "status"),
		t(lang, "cacheEntries"),
		t(lang, "type"),
		t(lang, "keys"),
		t(lang, "minTTL"),
		t(lang, "maxTTL"),
		t(lang, "defaultTTL"),
		t(lang, "action"),
		t(lang, "invalidate"),
		t(lang, "clearAll"),
		t(lang, "triggerPreload"),
		t(lang, "byNoteID"),
		t(lang, "clearNote"),
		t(lang, "byAttachmentID"),
		t(lang, "clearAttachment"),
		t(lang, "connected"),
		t(lang, "disconnected"),
		t(lang, "invalidToken"),
		t(lang, "clearAllConfirm"),
		t(lang, "clearedKeys"),
		t(lang, "clearedKeysFrom"),
		t(lang, "clearTypeConfirm"),
		t(lang, "clearedNote"),
		t(lang, "enterNoteID"),
		t(lang, "clearedAttachment"),
		t(lang, "enterAttachID"),
		t(lang, "preloadTriggered"),
		t(lang, "preloadInProgress"),
		t(lang, "clear"),
	)
}
