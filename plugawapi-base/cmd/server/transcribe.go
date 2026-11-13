package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

func (s *server) handleTranscribe(w http.ResponseWriter, r *http.Request) {
	sid := r.PathValue("sid")
	id := r.PathValue("id")
	sess, ok := s.sessions.Get(sid)
	if !ok {
		writeJSON(w, 404, map[string]string{"error": "session not found"})
		return
	}
	ac, ok := sess.reg.get(id)
	if !ok {
		writeJSON(w, 404, map[string]string{"error": "no such call"})
		return
	}
	path := "/tmp/recording_" + id + ".wav"
	f, err := os.Open(path)
	if err != nil {
		writeJSON(w, 404, map[string]string{"error": "recording not found"})
		return
	}
	defer f.Close()
	data, _ := io.ReadAll(f)

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		writeJSON(w, 400, map[string]string{"error": "OPENAI_API_KEY not set"})
		return
	}
	// Whisper
	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	part, _ := mw.CreateFormFile("file", "recording.wav")
	part.Write(data)
	mw.WriteField("model", "whisper-1")
	mw.Close()
	req, _ := http.NewRequest("POST", "https://api.openai.com/v1/audio/transcriptions", &buf)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+apiKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		writeJSON(w, 500, map[string]string{"error": err.Error()})
		return
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	var wres struct {
		Text string `json:"text"`
	}
	json.Unmarshal(b, &wres)
	transcript := wres.Text

	// Analysis via LLM
	baseURL := os.Getenv("AI_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}
	model := os.Getenv("AI_MODEL")
	if model == "" {
		model = "gpt-4o-mini"
	}
	prompt := os.Getenv("AI_PROMPT")
	if prompt == "" {
		prompt = "Analyze this call transcript and provide insights."
	}
	analysis := ""
	if transcript != "" {
		body := map[string]any{
			"model": model,
			"messages": []map[string]string{
				{"role": "system", "content": prompt + "\n\nTranscript:\n" + transcript},
				{"role": "user", "content": "Provide insights."},
			},
		}
		bb, _ := json.Marshal(body)
		r2, err := http.Post(strings.TrimRight(baseURL, "/")+"/chat/completions", "application/json", bytes.NewReader(bb))
		if err == nil {
			defer r2.Body.Close()
			bb2, _ := io.ReadAll(r2.Body)
			var cres struct {
				Choices []struct {
					Message struct {
						Content string `json:"content"`
					} `json:"message"`
				} `json:"choices"`
			}
			json.Unmarshal(bb2, &cres)
			if len(cres.Choices) > 0 {
				analysis = cres.Choices[0].Message.Content
			}
		}
	}
	writeJSON(w, 200, map[string]any{"transcript": transcript, "analysis": analysis})
	_ = ac
}

var _ = fmt.Sprintf
