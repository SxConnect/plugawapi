package main

import (
	"encoding/json"
	"net/http"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waCommon"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

func (s *server) handleReact(w http.ResponseWriter, r *http.Request) {
	sid := r.PathValue("sid")
	sess, ok := s.sessions.Get(sid)
	if !ok {
		writeJSON(w, 404, map[string]string{"error": "session not found"})
		return
	}
	var body struct {
		Chat    string `json:"chat"`
		MsgID   string `json:"message_id"`
		Emoji   string `json:"emoji"`
		Sender  string `json:"sender"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, 400, map[string]string{"error": "invalid body"})
		return
	}
	chat, _ := types.ParseJID(body.Chat)
	msg := &waE2E.Message{ReactionMessage: &waE2E.ReactionMessage{
		Key: &waCommon.MessageKey{RemoteJID: proto.String(chat.String()), FromMe: proto.Bool(false), ID: proto.String(body.MsgID)},
		Text: proto.String(body.Emoji),
	}}
	_, err := sess.client.SendMessage(r.Context(), chat, msg)
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, 200, map[string]any{"status": "ok"})
}

func (s *server) handlePresenceSubscribe(w http.ResponseWriter, r *http.Request) {
	sid := r.PathValue("sid")
	sess, ok := s.sessions.Get(sid)
	if !ok {
		writeJSON(w, 404, map[string]string{"error": "session not found"})
		return
	}
	var body struct {
		JID string `json:"jid"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, 400, map[string]string{"error": "invalid body"})
		return
	}
	jid, _ := types.ParseJID(body.JID)
	err := sess.client.SubscribePresence(r.Context(), jid)
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, 200, map[string]any{"status": "ok"})
}

func (s *server) handleCheckUser(w http.ResponseWriter, r *http.Request) {
	sid := r.PathValue("sid")
	sess, ok := s.sessions.Get(sid)
	if !ok {
		writeJSON(w, 404, map[string]string{"error": "session not found"})
		return
	}
	jidStr := r.URL.Query().Get("jid")
	jid, _ := types.ParseJID(jidStr)
	res, err := sess.client.IsOnWhatsApp(r.Context(), []string{jid.User})
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, 200, map[string]any{"results": res})
}

var _ = whatsmeow.NewClient
