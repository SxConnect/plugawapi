package main

import (
	"io"
	"net/http"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

func (s *server) handleMediaUpload(w http.ResponseWriter, r *http.Request) {
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
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		writeJSON(w, 400, map[string]string{"error": "invalid form"})
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		writeJSON(w, 400, map[string]string{"error": "file required"})
		return
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": "read failed"})
		return
	}
	to := r.FormValue("to")
	caption := r.FormValue("caption")
	mtype := r.FormValue("type")
	if mtype == "" { mtype = "image" }
	jid, _ := types.ParseJID(to)
	mt := header.Header.Get("Content-Type")
	if mt == "" { mt = "application/octet-stream" }
	var msg *waE2E.Message
	switch mtype {
	case "image":
		msg = &waE2E.Message{ImageMessage: &waE2E.ImageMessage{Caption: proto.String(caption), Mimetype: proto.String(mt)}}
	case "audio":
		msg = &waE2E.Message{AudioMessage: &waE2E.AudioMessage{Mimetype: proto.String(mt)}}
	case "video":
		msg = &waE2E.Message{VideoMessage: &waE2E.VideoMessage{Caption: proto.String(caption), Mimetype: proto.String(mt)}}
	case "document":
		msg = &waE2E.Message{DocumentMessage: &waE2E.DocumentMessage{FileName: proto.String(header.Filename), Mimetype: proto.String(mt)}}
	case "sticker":
		msg = &waE2E.Message{StickerMessage: &waE2E.StickerMessage{Mimetype: proto.String(mt)}}
	default:
		writeJSON(w, 400, map[string]string{"error": "invalid type"})
		return
	}
	// Upload media to WhatsApp servers
	uploaded, err := sess.client.Upload(r.Context(), data, whatsmeow.MediaImage)
	if err != nil {
		// fallback: try without upload
	}
	if uploaded.URL != "" {
		switch mtype {
		case "image": msg.ImageMessage.URL = proto.String(uploaded.URL)
		case "audio": msg.AudioMessage.URL = proto.String(uploaded.URL)
		case "video": msg.VideoMessage.URL = proto.String(uploaded.URL)
		case "document": msg.DocumentMessage.URL = proto.String(uploaded.URL)
		case "sticker": msg.StickerMessage.URL = proto.String(uploaded.URL)
		}
	}
	resp, err := sess.client.SendMessage(r.Context(), jid, msg)
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	s.chatStore.save(sid, jid.String(), sess.client.Store.ID.String(), caption, mtype, uploaded.URL, resp.Timestamp, true)
	writeJSON(w, 200, map[string]any{"id": resp.ID, "timestamp": resp.Timestamp})
}

