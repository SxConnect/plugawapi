package main

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

func (s *server) handleMediaProxy(w http.ResponseWriter, r *http.Request) {
	var body struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.URL == "" {
		writeJSON(w, 400, map[string]string{"error": "url required"})
		return
	}
	resp, err := http.Get(body.URL)
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	mt := resp.Header.Get("Content-Type")
	if mt == "" {
		mt = "application/octet-stream"
	}
	dataURL := "data:" + mt + ";base64," + base64.StdEncoding.EncodeToString(data)
	writeJSON(w, 200, map[string]any{"dataUrl": dataURL})
}

var _ = strings.TrimSpace
