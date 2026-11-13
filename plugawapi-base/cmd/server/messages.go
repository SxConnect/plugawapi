package main

import (
	"encoding/json"
	"net/http"

	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

func (s *server) handleSendText(w http.ResponseWriter, r *http.Request) {
	sid := r.PathValue("sid")
	sess, ok := s.sessions.Get(sid)
	if !ok {
		writeJSON(w, 404, map[string]string{"error": "session not found"})
		return
	}
	var body struct {
		To      string `json:"to"`
		Message string `json:"message"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, 400, map[string]string{"error": "invalid body"})
		return
	}
	jid, _ := types.ParseJID(body.To)
	msg := &waE2E.Message{Conversation: proto.String(body.Message)}
	resp, err := sess.client.SendMessage(r.Context(), jid, msg)
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	s.chatStore.save(sid, jid.String(), sess.client.Store.ID.String(), body.Message, "text", resp.Timestamp, true)
	writeJSON(w, 200, map[string]any{"id": resp.ID, "timestamp": resp.Timestamp})
}

func (s *server) handleListChats(w http.ResponseWriter, r *http.Request) {
	sid := r.PathValue("sid")
	sess, ok := s.sessions.Get(sid)
	if !ok {
		writeJSON(w, 404, map[string]string{"error": "session not found"})
		return
	}
	if s.chatStore == nil { s.log.Error("chatStore nil"); writeJSON(w,500,map[string]string{"error":"nil chatstore"}); return }
	s.log.Info("listChats called", "sid", sid)
	out := s.chatStore.listChats(sid)
	s.log.Info("listChats done", "count", len(out))
	writeJSON(w, 200, map[string]any{"chats": out})
	_ = sess
}

func (s *server) handleListMessages(w http.ResponseWriter, r *http.Request) {
	sid := r.PathValue("sid")
	jidStr := r.PathValue("jid")
	_, ok := s.sessions.Get(sid)
	if !ok {
		writeJSON(w, 404, map[string]string{"error": "session not found"})
		return
	}
	out := s.chatStore.listMessages(sid, jidStr, 200)
	writeJSON(w, 200, map[string]any{"messages": out})
}

func (s *server) handleListContacts(w http.ResponseWriter, r *http.Request) {
	sid := r.PathValue("sid")
	sess, ok := s.sessions.Get(sid)
	if !ok {
		writeJSON(w, 404, map[string]string{"error": "session not found"})
		return
	}
	contacts, err := sess.client.Store.Contacts.GetAllContacts(r.Context())
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	out := make([]map[string]any, 0, len(contacts))
	for jid, info := range contacts {
		name := info.FullName
		if name == "" { name = info.PushName }
		if name == "" { name = info.BusinessName }
		if name == "" { name = jid.User }
		out = append(out, map[string]any{
			"jid":  jid.String(),
			"name": name,
		})
	}
	writeJSON(w, 200, map[string]any{"contacts": out})
}

func (s *server) handleListGroups(w http.ResponseWriter, r *http.Request) {
	sid := r.PathValue("sid")
	sess, ok := s.sessions.Get(sid)
	if !ok {
		writeJSON(w, 404, map[string]string{"error": "session not found"})
		return
	}
	groups, err := sess.client.GetJoinedGroups(r.Context())
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	out := make([]map[string]any, 0, len(groups))
	for _, g := range groups {
		out = append(out, map[string]any{
			"jid":  g.JID.String(),
			"name": g.Name,
		})
	}
	writeJSON(w, 200, map[string]any{"groups": out})
}

