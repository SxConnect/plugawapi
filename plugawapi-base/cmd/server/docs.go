package main

import (
	"net/http"
)

func (s *server) handleDocs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(`<!doctype html><html><head><title>PlugaWAPI Docs</title><style>body{font-family:sans-serif;background:#07152b;color:#e6f0ff;padding:2rem}h1{color:#00e676}h2{color:#00d2ff;margin-top:1.5rem}pre{background:#0d1f3d;padding:1rem;border-radius:8px;overflow:auto}code{color:#00e676}</style></head><body><h1>PlugaWAPI — Documentação</h1><p>API REST de WhatsApp. Base URL: <code>http://localhost:8080</code></p>
<h2>Sessões</h2><pre>GET  /api/sessions
POST /api/sessions  {name}
DELETE /api/sessions/{sid}
POST /api/sessions/{sid}/pair
POST /api/sessions/{sid}/pair-code  {phone}
POST /api/sessions/{sid}/reconnect
POST /api/sessions/{sid}/logout</pre>
<h2>Mensagens</h2><pre>POST /api/sessions/{sid}/messages  {to, message}
POST /api/sessions/{sid}/media  {to, url, type, caption}
POST /api/sessions/{sid}/media-upload  (multipart: file, to, type)
POST /api/sessions/{sid}/location  {to, latitude, longitude, name}
GET  /api/sessions/{sid}/chats
GET  /api/sessions/{sid}/chats/{jid}/messages
GET  /api/sessions/{sid}/contacts
POST /api/sessions/{sid}/contacts/sync
GET  /api/sessions/{sid}/contacts/saved
POST /api/sessions/{sid}/react  {chat, message_id, emoji}
POST /api/sessions/{sid}/presence  {jid}</pre>
<h2>Chamadas</h2><pre>POST /api/sessions/{sid}/calls  {phone}
POST /api/sessions/{sid}/calls/{id}/webrtc  {sdp_offer}
POST /api/sessions/{sid}/calls/{id}/accept
POST /api/sessions/{sid}/calls/{id}/reject
DELETE /api/sessions/{sid}/calls/{id}
POST /api/sessions/{sid}/calls/{id}/hold  {moh}
POST /api/sessions/{sid}/calls/{id}/resume
POST /api/sessions/{sid}/calls/{id}/mute  {muted}
POST /api/sessions/{sid}/calls/{id}/record/start
POST /api/sessions/{sid}/calls/{id}/record/stop
GET  /api/record/download?path=
POST /api/sessions/{sid}/calls/{id}/transcribe
POST /api/sessions/{sid}/interactive  {to, body, footer, buttons}</pre>
<h2>Perfil / Newsletter</h2><pre>GET /api/sessions/{sid}/profile
GET /api/sessions/{sid}/business-profile?jid=
POST /api/sessions/{sid}/newsletter/follow  {jid}
POST /api/sessions/{sid}/newsletter/unfollow  {jid}
GET  /api/sessions/{sid}/newsletter/info?jid=
GET  /api/sessions/{sid}/newsletter/messages?jid=</pre>
<h2>Utilitários</h2><pre>GET  /api/sessions/{sid}/check?jid=
GET  /api/sessions/{sid}/groups
GET  /api/geocode?q=
POST /api/media/proxy  {url}
GET  /api/events (SSE)
GET/POST /api/webhook</pre>
<p style="margin-top:2rem;color:#8aa0c4">PlugaWAPI — base WaCalls (MIT) + extensões. Licença MIT.</p>
</body></html>`))
}
