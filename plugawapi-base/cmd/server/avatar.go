package main

import (
	"net/http"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
)

func (s *server) handleAvatar(w http.ResponseWriter, r *http.Request) {
	sid := r.PathValue("sid")
	sess, ok := s.sessions.Get(sid)
	if !ok {
		writeJSON(w, 404, map[string]string{"error": "session not found"})
		return
	}
	jidStr := r.URL.Query().Get("jid")
	jid, _ := types.ParseJID(jidStr)
	info, err := sess.client.GetProfilePictureInfo(r.Context(), jid, &whatsmeow.GetProfilePictureParams{})
	if err != nil {
		writeJSON(w, 200, map[string]any{"url": ""})
		return
	}
	if info == nil {
		writeJSON(w, 200, map[string]any{"url": ""})
		return
	}
	writeJSON(w, 200, map[string]any{"url": info.URL})
}
