package main

import (
	"io"
	"net/http"
	"os"
)

func (s *server) handleRecordStart(w http.ResponseWriter, r *http.Request) {
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
	path := "/tmp/recording_" + id + ".wav"
	ac.cm.StartRecording(path)
	writeJSON(w, 200, map[string]any{"status": "recording", "path": path})
}

func (s *server) handleRecordStop(w http.ResponseWriter, r *http.Request) {
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
	path := ac.cm.StopRecording()
	writeJSON(w, 200, map[string]any{"status": "stopped", "path": path})
}

func (s *server) handleRecordDownload(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		writeJSON(w, 400, map[string]string{"error": "path required"})
		return
	}
	f, err := os.Open(path)
	if err != nil {
		writeJSON(w, 404, map[string]string{"error": "file not found"})
		return
	}
	defer f.Close()
	w.Header().Set("Content-Type", "audio/wav")
	w.Header().Set("Content-Disposition", "attachment; filename=recording.wav")
	_, _ = io.Copy(w, f)
}
