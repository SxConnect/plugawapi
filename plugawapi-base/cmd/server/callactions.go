package main

import (
	"encoding/json"
	"net/http"
)

func (s *server) handleHold(w http.ResponseWriter, r *http.Request) {
	sess, ok := s.sessions.Get(r.PathValue("sid"))
	if !ok {
		writeJSON(w, 404, map[string]string{"error": "session not found"})
		return
	}
	id := r.PathValue("id")
	ac, ok := sess.reg.get(id)
	if !ok {
		writeJSON(w, 404, map[string]string{"error": "no such call"})
		return
	}
	moh := true
	var body struct {
		MOH *bool `json:"moh"`
	}
	if json.NewDecoder(r.Body).Decode(&body) == nil && body.MOH != nil {
		moh = *body.MOH
	}
	if err := ac.cm.Hold(moh); err != nil {
		writeJSON(w, 409, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, 200, map[string]any{"status": "ok", "held": true, "moh": moh})
}

func (s *server) handleResume(w http.ResponseWriter, r *http.Request) {
	sess, ok := s.sessions.Get(r.PathValue("sid"))
	if !ok {
		writeJSON(w, 404, map[string]string{"error": "session not found"})
		return
	}
	id := r.PathValue("id")
	ac, ok := sess.reg.get(id)
	if !ok {
		writeJSON(w, 404, map[string]string{"error": "no such call"})
		return
	}
	if err := ac.cm.Resume(); err != nil {
		writeJSON(w, 409, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, 200, map[string]any{"status": "ok", "held": false})
}

func (s *server) handleMute(w http.ResponseWriter, r *http.Request) {
	sess, ok := s.sessions.Get(r.PathValue("sid"))
	if !ok {
		writeJSON(w, 404, map[string]string{"error": "session not found"})
		return
	}
	id := r.PathValue("id")
	ac, ok := sess.reg.get(id)
	if !ok {
		writeJSON(w, 404, map[string]string{"error": "no such call"})
		return
	}
	var body struct {
		Muted bool `json:"muted"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)
	if err := ac.cm.Mute(body.Muted); err != nil {
		writeJSON(w, 409, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, 200, map[string]any{"status": "ok", "muted": body.Muted})
}
