package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

func (s *server) handleSendMedia(w http.ResponseWriter, r *http.Request) {
	sid := r.PathValue("sid")
	sess, ok := s.sessions.Get(sid)
	if !ok {
		writeJSON(w, 404, map[string]string{"error": "session not found"})
		return
	}
	var body struct {
		To       string `json:"to"`
		URL      string `json:"url"`
		Caption  string `json:"caption"`
		MimeType string `json:"mimetype"`
		Type     string `json:"type"` // image, audio, video, document
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, 400, map[string]string{"error": "invalid body"})
		return
	}
	resp, err := http.Get(body.URL)
	if err != nil {
		writeJSON(w, 400, map[string]string{"error": "download failed"})
		return
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": "read failed"})
		return
	}
	mt := body.MimeType
	if mt == "" {
		mt = resp.Header.Get("Content-Type")
	}
	jid, _ := types.ParseJID(body.To)
	var msg *waE2E.Message
	switch body.Type {
	case "image":
		msg = &waE2E.Message{ImageMessage: &waE2E.ImageMessage{
			Caption: proto.String(body.Caption), Mimetype: proto.String(mt), URL: proto.String(body.URL),
		}}
	case "audio":
		msg = &waE2E.Message{AudioMessage: &waE2E.AudioMessage{Mimetype: proto.String(mt), URL: proto.String(body.URL)}}
	case "video":
		msg = &waE2E.Message{VideoMessage: &waE2E.VideoMessage{
			Caption: proto.String(body.Caption), Mimetype: proto.String(mt), URL: proto.String(body.URL),
		}}
	case "document":
		msg = &waE2E.Message{DocumentMessage: &waE2E.DocumentMessage{
			FileName: proto.String("file"), Mimetype: proto.String(mt), URL: proto.String(body.URL),
		}}
	default:
		writeJSON(w, 400, map[string]string{"error": "invalid type"})
		return
	}
	_ = data
	r2, err := sess.client.SendMessage(r.Context(), jid, msg)
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, 200, map[string]any{"id": r2.ID, "timestamp": r2.Timestamp})
}

func (s *server) handleWebhookConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		writeJSON(w, 200, map[string]string{"webhook": os.Getenv("PLUGAWAPI_WEBHOOK_URL")})
		return
	}
	var body struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, 400, map[string]string{"error": "invalid body"})
		return
	}
	os.Setenv("PLUGAWAPI_WEBHOOK_URL", body.URL)
	writeJSON(w, 200, map[string]string{"webhook": body.URL})
}

var _ = fmt.Sprintf
var _ = strings.TrimSpace
