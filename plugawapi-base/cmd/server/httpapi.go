package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"plugawapi/internal/voip/core"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
)

func (s *server) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/sessions", s.handleSessionList)
	mux.HandleFunc("POST /api/sessions", s.handleSessionCreate)
	mux.HandleFunc("DELETE /api/sessions/{sid}", s.handleSessionDelete)
	mux.HandleFunc("POST /api/sessions/{sid}/logout", s.handleSessionLogout)
	mux.HandleFunc("POST /api/sessions/{sid}/pair", s.handleSessionPair)
	mux.HandleFunc("POST /api/sessions/{sid}/pair-code", s.handleSessionPairCode)
	mux.HandleFunc("POST /api/sessions/{sid}/reconnect", s.handleReconnect)
	mux.HandleFunc("POST /api/sessions/{sid}/calls", s.handleStartCall)
	mux.HandleFunc("POST /api/sessions/{sid}/calls/{id}/webrtc", s.handleWebRTC)
	mux.HandleFunc("POST /api/sessions/{sid}/calls/{id}/accept", s.handleAccept)
	mux.HandleFunc("POST /api/sessions/{sid}/calls/{id}/reject", s.handleReject)
	mux.HandleFunc("DELETE /api/sessions/{sid}/calls/{id}", s.handleEndCall)
	mux.HandleFunc("POST /api/sessions/{sid}/calls/{id}/hold", s.handleHold)
	mux.HandleFunc("POST /api/sessions/{sid}/calls/{id}/resume", s.handleResume)
	mux.HandleFunc("POST /api/sessions/{sid}/calls/{id}/mute", s.handleMute)
	mux.HandleFunc("POST /api/sessions/{sid}/calls/{id}/record/start", s.handleRecordStart)
	mux.HandleFunc("POST /api/sessions/{sid}/calls/{id}/record/stop", s.handleRecordStop)
	mux.HandleFunc("POST /api/sessions/{sid}/calls/{id}/transcribe", s.handleTranscribe)
	mux.HandleFunc("GET /api/record/download", s.handleRecordDownload)
	mux.HandleFunc("GET /api/sessions/{sid}/history", s.handleHistory)

	mux.HandleFunc("POST /api/sessions/{sid}/messages", s.handleSendText)
	mux.HandleFunc("POST /api/sessions/{sid}/media", s.handleSendMedia)
	mux.HandleFunc("GET /api/sessions/{sid}/chats", s.handleListChats)
	mux.HandleFunc("GET /api/sessions/{sid}/chats/{jid}/messages", s.handleListMessages)
	mux.HandleFunc("GET /api/sessions/{sid}/contacts", s.handleListContacts)
	mux.HandleFunc("GET /api/sessions/{sid}/contacts/saved", s.handleContactsList)
	mux.HandleFunc("POST /api/sessions/{sid}/contacts/sync", s.handleContactsSync)
	mux.HandleFunc("GET /api/sessions/{sid}/avatar", s.handleAvatar)
	mux.HandleFunc("GET /api/sessions/{sid}/profile", s.handleProfile)
	mux.HandleFunc("GET /api/sessions/{sid}/business-profile", s.handleBusinessProfile)
	mux.HandleFunc("GET /api/sessions/{sid}/groups", s.handleListGroups)
	mux.HandleFunc("POST /api/sessions/{sid}/newsletter/follow", s.handleNewsletterFollow)
	mux.HandleFunc("POST /api/sessions/{sid}/newsletter/unfollow", s.handleNewsletterUnfollow)
	mux.HandleFunc("GET /api/sessions/{sid}/newsletter/info", s.handleNewsletterInfo)
	mux.HandleFunc("GET /api/sessions/{sid}/newsletter/messages", s.handleNewsletterMessages)
	mux.HandleFunc("POST /api/sessions/{sid}/react", s.handleReact)
	mux.HandleFunc("POST /api/sessions/{sid}/presence", s.handlePresenceSubscribe)
	mux.HandleFunc("GET /api/sessions/{sid}/check", s.handleCheckUser)
	mux.HandleFunc("GET /api/webhook", s.handleWebhookConfig)
	mux.HandleFunc("POST /api/webhook", s.handleWebhookConfig)


	mux.HandleFunc("GET /api/events", s.handleEvents)

	if s.staticDir != "" {
		if _, err := os.Stat(s.staticDir); err == nil {
			mux.Handle("/", http.FileServer(http.Dir(s.staticDir)))
		}
	}
	return withCORS(mux)
}

func withCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Client-Id")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func clientID(r *http.Request) string {
	if id := r.Header.Get("X-Client-Id"); id != "" {
		return id
	}
	return r.URL.Query().Get("clientId")
}

func (s *server) sessionByID(w http.ResponseWriter, sid string) *Session {
	sess, ok := s.sessions.Get(sid)
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "no such session"})
		return nil
	}
	return sess
}

func (s *server) handleEvents(w http.ResponseWriter, r *http.Request) {
	s.broker.serveSSE(w, r, clientID(r))
}

func (s *server) handleSessionList(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"sessions": s.sessions.infos()})
}

func (s *server) handleSessionCreate(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name string `json:"name"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)
	name := strings.TrimSpace(body.Name)
	if name == "" {
		name = "Session"
	}
	id, err := s.sessions.Create(name)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"id": id})
}

func (s *server) handleSessionDelete(w http.ResponseWriter, r *http.Request) {
	if err := s.sessions.Delete(r.Context(), r.PathValue("sid")); err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *server) handleSessionLogout(w http.ResponseWriter, r *http.Request) {
	if err := s.sessions.Logout(r.Context(), r.PathValue("sid")); err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *server) handleSessionPair(w http.ResponseWriter, r *http.Request) {
	if err := s.sessions.Pair(r.PathValue("sid")); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *server) handleStartCall(w http.ResponseWriter, r *http.Request) {
	if sess := s.sessionByID(w, r.PathValue("sid")); sess != nil {
		s.doStartCall(sess, w, r)
	}
}

func (s *server) handleWebRTC(w http.ResponseWriter, r *http.Request) {
	if sess := s.sessionByID(w, r.PathValue("sid")); sess != nil {
		s.doWebRTC(sess, w, r)
	}
}

func (s *server) handleAccept(w http.ResponseWriter, r *http.Request) {
	if sess := s.sessionByID(w, r.PathValue("sid")); sess != nil {
		s.doAccept(sess, w, r)
	}
}

func (s *server) handleReject(w http.ResponseWriter, r *http.Request) {
	if sess := s.sessionByID(w, r.PathValue("sid")); sess != nil {
		s.doReject(sess, w, r)
	}
}

func (s *server) handleEndCall(w http.ResponseWriter, r *http.Request) {
	if sess := s.sessionByID(w, r.PathValue("sid")); sess != nil {
		s.doEndCall(sess, w, r)
	}
}

func (s *server) handleHistory(w http.ResponseWriter, r *http.Request) {
	if sess := s.sessionByID(w, r.PathValue("sid")); sess != nil {
		writeJSON(w, http.StatusOK, map[string]any{"rows": s.broker.historyRows(sess.id, 50)})
	}
}

func (s *server) doStartCall(sess *Session, w http.ResponseWriter, r *http.Request) {
	if sess.client.Store.ID == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "not paired"})
		return
	}
	var body struct {
		Phone      string `json:"phone"`
		DurationMs int    `json:"duration_ms"`
		Record     bool   `json:"record"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || strings.TrimSpace(body.Phone) == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "phone required"})
		return
	}
	owner := clientID(r)
	if other := s.broker.ownerActiveCall(owner); other != "" {
		writeJSON(w, http.StatusConflict, map[string]string{"error": "operator already on a call"})
		return
	}
	if max := s.sessions.maxCalls; max > 0 && sess.reg.count() >= max {
		writeJSON(w, http.StatusTooManyRequests, map[string]string{"error": "max concurrent calls"})
		return
	}
	peer := types.NewJID(normalizePhone(body.Phone), types.DefaultUserServer)

	callID, err := sess.startOutgoing(r.Context(), peer, false)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	s.broker.upsertCall(CallRecord{
		SessionID: sess.id, CallID: callID, Owner: &owner, Direction: "outbound", Peer: peer.String(),
		StartedAt: time.Now().UnixMilli(), Status: StatusRinging,
	})
	writeJSON(w, http.StatusOK, map[string]any{"call": map[string]string{"callId": callID}})
}

func (s *server) doWebRTC(sess *Session, w http.ResponseWriter, r *http.Request) {
	callID := r.PathValue("id")
	ac, ok := sess.reg.get(callID)
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "no such call"})
		return
	}
	var body struct {
		SDPOffer string `json:"sdp_offer"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.SDPOffer == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "sdp_offer required"})
		return
	}
	bridge, answer, err := NewBridge(body.SDPOffer, s.log)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	bridge.OnBrowserPCM = func(pcm []float32) {
		ac.cm.FeedCapturedPCM(pcm)
	}
	bridge.OnTerminalICE = func() {
		go sess.terminateCall(callID, core.EndCallReasonUserEnded)
	}
	sess.setBridge(callID, bridge)
	writeJSON(w, http.StatusOK, map[string]string{"sdp_answer": answer})
}

func (s *server) doAccept(sess *Session, w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	ac, ok := sess.reg.get(id)
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "no such call"})
		return
	}
	owner := clientID(r)
	if other := s.broker.ownerActiveCall(owner); other != "" && other != id {
		writeJSON(w, http.StatusConflict, map[string]string{"error": "operator already on a call"})
		return
	}
	if !s.broker.setOwner(id, owner) {
		writeJSON(w, http.StatusConflict, map[string]string{"error": "claimed by another client"})
		return
	}
	s.broker.emitIncomingClaimed(sess.id, id, owner)
	if err := ac.cm.AcceptCall(r.Context(), id); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"call": map[string]string{"callId": id}})
}

func (s *server) doReject(sess *Session, w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if ac, ok := sess.reg.get(id); ok {
		_ = ac.cm.RejectCall(r.Context(), id, core.EndCallReasonDeclined)
	}
	sess.removeCall(id)
	s.broker.endCall(id, string(core.EndCallReasonDeclined))
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *server) doEndCall(sess *Session, w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if ac, ok := sess.reg.get(id); ok {
		_ = ac.cm.EndCall(r.Context(), core.EndCallReasonUserEnded)
	}
	sess.removeCall(id)
	s.broker.endCall(id, string(core.EndCallReasonUserEnded))
	w.WriteHeader(http.StatusNoContent)
}

func normalizePhone(p string) string {
	p = strings.TrimSpace(p)
	p = strings.TrimPrefix(p, "+")
	var b strings.Builder
	for _, c := range p {
		if c >= '0' && c <= '9' {
			b.WriteRune(c)
		}
	}
	return b.String()
}

func (s *server) handleReconnect(w http.ResponseWriter, r *http.Request) {
	sess, ok := s.sessions.Get(r.PathValue("sid"))
	if !ok {
		writeJSON(w, 404, map[string]string{"error": "session not found"})
		return
	}
	sess.client.Disconnect()
	if err := sess.client.Connect(); err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, 200, map[string]any{"status": "ok"})
}

func (s *server) handleSessionPairCode(w http.ResponseWriter, r *http.Request) {
	sess, ok := s.sessions.Get(r.PathValue("sid"))
	if !ok {
		writeJSON(w, 404, map[string]string{"error": "session not found"})
		return
	}
	var body struct {
		Phone string `json:"phone"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Phone == "" {
		writeJSON(w, 400, map[string]string{"error": "phone required"})
		return
	}
	if !sess.client.IsConnected() {
		if err := sess.client.Connect(); err != nil {
			writeJSON(w, 500, map[string]string{"error": "connect failed: " + err.Error()})
			return
		}
	}
	s.log.Info("pair-code attempt", "phone", body.Phone, "connected", sess.client.IsConnected())
	code, err := sess.client.PairPhone(r.Context(), body.Phone, true, whatsmeow.PairClientChrome, "Chrome (Linux)")
	if err != nil {
		s.log.Error("pair-code failed", "err", err)
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	s.log.Info("pair-code success", "code", code)
	writeJSON(w, 200, map[string]any{"code": code})
}
