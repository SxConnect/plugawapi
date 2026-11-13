package main

import (
	"encoding/json"
	"net/http"

	"go.mau.fi/whatsmeow/types"
)

func (s *server) handleNewsletterFollow(w http.ResponseWriter, r *http.Request) {
	sid := r.PathValue("sid")
	sess, ok := s.sessions.Get(sid)
	if !ok {
		writeJSON(w, 404, map[string]string{"error": "session not found"})
		return
	}
	var body struct {
		JID string `json:"jid"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, 400, map[string]string{"error": "invalid body"})
		return
	}
	jid, _ := types.ParseJID(body.JID)
	err := sess.client.FollowNewsletter(r.Context(), jid)
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, 200, map[string]any{"status": "ok"})
}

func (s *server) handleNewsletterUnfollow(w http.ResponseWriter, r *http.Request) {
	sid := r.PathValue("sid")
	sess, ok := s.sessions.Get(sid)
	if !ok {
		writeJSON(w, 404, map[string]string{"error": "session not found"})
		return
	}
	var body struct {
		JID string `json:"jid"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, 400, map[string]string{"error": "invalid body"})
		return
	}
	jid, _ := types.ParseJID(body.JID)
	err := sess.client.UnfollowNewsletter(r.Context(), jid)
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, 200, map[string]any{"status": "ok"})
}

func (s *server) handleNewsletterInfo(w http.ResponseWriter, r *http.Request) {
	sid := r.PathValue("sid")
	sess, ok := s.sessions.Get(sid)
	if !ok {
		writeJSON(w, 404, map[string]string{"error": "session not found"})
		return
	}
	jidStr := r.URL.Query().Get("jid")
	jid, _ := types.ParseJID(jidStr)
	info, err := sess.client.GetNewsletterInfo(r.Context(), jid)
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, 200, map[string]any{"info": info})
}

func (s *server) handleNewsletterMessages(w http.ResponseWriter, r *http.Request) {
	sid := r.PathValue("sid")
	sess, ok := s.sessions.Get(sid)
	if !ok {
		writeJSON(w, 404, map[string]string{"error": "session not found"})
		return
	}
	jidStr := r.URL.Query().Get("jid")
	jid, _ := types.ParseJID(jidStr)
	msgs, err := sess.client.GetNewsletterMessages(r.Context(), jid, nil)
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, 200, map[string]any{"messages": msgs})
}
