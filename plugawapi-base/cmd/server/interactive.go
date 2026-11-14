package main

import (
	"encoding/json"
	"net/http"

	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

func (s *server) handleSendInteractive(w http.ResponseWriter, r *http.Request) {
	sid := r.PathValue("sid")
	sess, ok := s.sessions.Get(sid)
	if !ok {
		writeJSON(w, 404, map[string]string{"error": "session not found"})
		return
	}
	var body struct {
		To      string `json:"to"`
		Title   string `json:"title"`
		Body    string `json:"body"`
		Footer  string `json:"footer"`
		Buttons []struct {
			ID    string `json:"id"`
			Title string `json:"title"`
		} `json:"buttons"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, 400, map[string]string{"error": "invalid body"})
		return
	}
	if !sess.client.IsConnected() {
		if err := sess.client.Connect(); err != nil {
			writeJSON(w, 500, map[string]string{"error": "connect failed"})
			return
		}
	}
	jid, _ := types.ParseJID(body.To)
	buttons := []*waE2E.InteractiveMessage_NativeFlowMessage_NativeFlowButton{}
	for _, b := range body.Buttons {
		btn := &waE2E.InteractiveMessage_NativeFlowMessage_NativeFlowButton{
			Name:             proto.String("quick_reply"),
			ButtonParamsJSON: proto.String(`{"display_text":"` + b.Title + `","id":"` + b.ID + `"}`),
		}
		buttons = append(buttons, btn)
	}
	msg := &waE2E.Message{
		InteractiveMessage: &waE2E.InteractiveMessage{
			Body: &waE2E.InteractiveMessage_Body{Text: proto.String(body.Body)},
			Footer: &waE2E.InteractiveMessage_Footer{Text: proto.String(body.Footer)},
			Header: &waE2E.InteractiveMessage_Header{Title: proto.String(body.Title)},
			InteractiveMessage: &waE2E.InteractiveMessage_NativeFlowMessage_{
				NativeFlowMessage: &waE2E.InteractiveMessage_NativeFlowMessage{
					Buttons:           buttons,
					MessageParamsJSON: proto.String("{}"),
					MessageVersion:    proto.Int32(3),
				},
			},
		},
	}
	resp, err := sess.client.SendMessage(r.Context(), jid, msg)
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	s.chatStore.save(sid, jid.String(), sess.client.Store.ID.String(), body.Body, "interactive", "", resp.Timestamp, true)
	writeJSON(w, 200, map[string]any{"id": resp.ID, "timestamp": resp.Timestamp})
}
