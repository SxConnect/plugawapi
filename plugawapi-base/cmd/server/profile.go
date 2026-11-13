package main

import (
	"net/http"

	"go.mau.fi/whatsmeow/types"
)

func (s *server) handleProfile(w http.ResponseWriter, r *http.Request) {
	sid := r.PathValue("sid")
	sess, ok := s.sessions.Get(sid)
	if !ok {
		writeJSON(w, 404, map[string]string{"error": "session not found"})
		return
	}
	info, err := sess.client.GetUserInfo(r.Context(), []types.JID{sess.client.Store.ID.ToNonAD()})
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, 200, map[string]any{"info": info})
}

func (s *server) handleBusinessProfile(w http.ResponseWriter, r *http.Request) {
	sid := r.PathValue("sid")
	sess, ok := s.sessions.Get(sid)
	if !ok {
		writeJSON(w, 404, map[string]string{"error": "session not found"})
		return
	}
	jidStr := r.URL.Query().Get("jid")
	jid, _ := types.ParseJID(jidStr)
	bp, err := sess.client.GetBusinessProfile(r.Context(), jid)
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, 200, map[string]any{"profile": bp})
}
