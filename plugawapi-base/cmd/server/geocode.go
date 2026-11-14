package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func (s *server) handleGeocode(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if q == "" {
		writeJSON(w, 400, map[string]string{"error": "q required"})
		return
	}
	u := "https://nominatim.openstreetmap.org/search?format=json&q=" + url.QueryEscape(q)
	req, _ := http.NewRequest("GET", u, nil)
	req.Header.Set("User-Agent", "PlugaWAPI/1.0")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	defer resp.Body.Close()
	var data []map[string]any
	json.NewDecoder(resp.Body).Decode(&data)
	if len(data) == 0 {
		writeJSON(w, 404, map[string]string{"error": "not found"})
		return
	}
	writeJSON(w, 200, map[string]any{"lat": data[0]["lat"], "lon": data[0]["lon"], "display": data[0]["display_name"]})
}

var _ = fmt.Sprintf
