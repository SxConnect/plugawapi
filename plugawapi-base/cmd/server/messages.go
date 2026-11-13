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
	writeJSON(w, 200, map[string]any{"id": resp.ID, "timestamp": resp.Timestamp})
}

func (s *server) handleListChats(w http.ResponseWriter, r *http.Request) {
	sid := r.PathValue("sid")
	_, ok := s.sessions.Get(sid)
	if !ok {
		writeJSON(w, 404, map[string]string{"error": "session not found"})
		return
	}
	writeJSON(w, 501, map[string]string{"error": "chat history not implemented; use GOWA for history"})
}

func (s *server) handleListMessages(w http.ResponseWriter, r *http.Request) {
	sid := r.PathValue("sid")
	jidStr := r.PathValue("jid")
	_, ok := s.sessions.Get(sid)
	if !ok {
		writeJSON(w, 404, map[string]string{"error": "session not found"})
		return
	}
	_ = jidStr
	writeJSON(w, 501, map[string]string{"error": "message history not implemented; use GOWA for history"})
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
		out = append(out, map[string]any{
			"jid":  jid.String(),
			"name": info.FullName,
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

