package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

func (s *server) handleSendLocation(w http.ResponseWriter, r *http.Request) {
	sid := r.PathValue("sid")
	sess, ok := s.sessions.Get(sid)
	if !ok {
		writeJSON(w, 404, map[string]string{"error": "session not found"})
		return
	}
	if !sess.client.IsConnected() {
		if err := sess.client.Connect(); err != nil {
			writeJSON(w, 500, map[string]string{"error": "connect failed: " + err.Error()})
			return
		}
	}
	var body struct {
		To        string  `json:"to"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Name      string  `json:"name"`
		Address   string  `json:"address"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, 400, map[string]string{"error": "invalid body"})
		return
	}
	jid, _ := types.ParseJID(body.To)
	msg := &waE2E.Message{LocationMessage: &waE2E.LocationMessage{
		DegreesLatitude:  proto.Float64(body.Latitude),
		DegreesLongitude: proto.Float64(body.Longitude),
		Name:             proto.String(body.Name),
		Address:          proto.String(body.Address),
	}}
	resp, err := sess.client.SendMessage(r.Context(), jid, msg)
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	s.chatStore.save(sid, jid.String(), sess.client.Store.ID.String(), "[localização]", "location", "", resp.Timestamp, true)
	writeJSON(w, 200, map[string]any{"id": resp.ID, "timestamp": resp.Timestamp})
}

var _ = fmt.Sprintf
