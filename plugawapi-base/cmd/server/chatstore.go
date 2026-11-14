package main

import (
	"database/sql"
	"time"
)

type chatStore struct {
	db *sql.DB
}

func newChatStore(db *sql.DB) *chatStore {
	return &chatStore{db: db}
}

func (c *chatStore) save(sessionID, chatJID, senderJID, content, mediaType, mediaURL string, ts time.Time, fromMe bool) {
	from := 0
	if fromMe {
		from = 1
	}
	_, _ = c.db.Exec(`INSERT OR REPLACE INTO messages (id, session_id, chat_jid, sender_jid, content, media_type, media_url, timestamp, is_from_me, created_at) VALUES (?,?,?,?,?,?,?,?,?,?)`,
		"m_"+ts.UTC().Format("20060102150405.000000000")+senderJID, sessionID, chatJID, senderJID, content, mediaType, mediaURL, ts.UTC().Format(time.RFC3339), from, time.Now().UTC().Format(time.RFC3339))
}

func (c *chatStore) listChats(sessionID string) []map[string]any {
	res, err := c.db.Query(`SELECT chat_jid, MAX(timestamp) FROM messages WHERE session_id=? GROUP BY chat_jid ORDER BY MAX(timestamp) DESC`, sessionID)
	if err != nil {
		return nil
	}
	defer res.Close()
	out := []map[string]any{}
	for res.Next() {
		var jid, ts string
		if err := res.Scan(&jid, &ts); err != nil {
			continue
		}
		var lastMsg sql.NullString
		_ = c.db.QueryRow(`SELECT content FROM messages WHERE session_id=? AND chat_jid=? ORDER BY timestamp DESC LIMIT 1`, sessionID, jid).Scan(&lastMsg)
		out = append(out, map[string]any{"jid": jid, "lastTime": ts, "lastMessage": lastMsg.String})
	}
	return out
}

func (c *chatStore) listMessages(sessionID, chatJID string, limit int) []map[string]any {
	rows, err := c.db.Query(`SELECT id, sender_jid, content, media_type, media_url, timestamp, is_from_me FROM messages WHERE session_id=? AND chat_jid=? ORDER BY timestamp ASC LIMIT ?`, sessionID, chatJID, limit)
	if err != nil {
		return nil
	}
	defer rows.Close()
	out := []map[string]any{}
	for rows.Next() {
		var id, sender, content, mediaType, mediaURL, ts string
		var fromMe int
		if err := rows.Scan(&id, &sender, &content, &mediaType, &mediaURL, &ts, &fromMe); err != nil {
			continue
		}
		out = append(out, map[string]any{
			"id":        id,
			"sender":    sender,
			"content":   content,
			"mediaType": mediaType,
			"mediaUrl":  mediaURL,
			"timestamp": ts,
			"fromMe":    fromMe == 1,
		})
	}
	return out
}

func (c *chatStore) saveContacts(sessionID string, contacts []map[string]string) {
	for _, ct := range contacts {
		_, _ = c.db.Exec(`INSERT OR REPLACE INTO contacts (session_id, jid, name) VALUES (?,?,?)`, sessionID, ct["jid"], ct["name"])
	}
}

func (c *chatStore) listContacts(sessionID string) []map[string]string {
	rows, _ := c.db.Query(`SELECT jid, name FROM contacts WHERE session_id=? ORDER BY name`, sessionID)
	defer rows.Close()
	out := []map[string]string{}
	for rows.Next() {
		var jid, name string
		rows.Scan(&jid, &name)
		out = append(out, map[string]string{"jid": jid, "name": name})
	}
	return out
}
