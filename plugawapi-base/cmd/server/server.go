package main

import (
	"context"
	"database/sql"
	"log/slog"

	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
	_ "modernc.org/sqlite"
)

type server struct {
	broker    *Broker
	sessions  *SessionManager
	chatStore *chatStore
	sip       *sipGateway
	log       *slog.Logger
	staticDir string
}

func openDB(dbPath string) (*sql.DB, error) {
	dsn := "file:" + dbPath + "?_pragma=foreign_keys(1)&_pragma=busy_timeout(10000)&_pragma=journal_mode(WAL)"
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	return db, nil
}

func newServer(ctx context.Context, dbPath, staticDir string, maxCalls int, log *slog.Logger) (*server, error) {
	db, err := openDB(dbPath)
	if err != nil {
		return nil, err
	}
	chatPath := dbPath
	if len(chatPath) > 3 && chatPath[len(chatPath)-3:] == ".db" {
		chatPath = chatPath[:len(chatPath)-3] + "_chat.db"
	} else {
		chatPath = chatPath + "_chat"
	}
	chatDB, err := openDB(chatPath)
	if err != nil {
		return nil, err
	}
	_, _ = chatDB.Exec(`CREATE TABLE IF NOT EXISTS messages (
		id TEXT PRIMARY KEY,
		session_id TEXT NOT NULL,
		chat_jid TEXT NOT NULL,
		sender_jid TEXT,
		content TEXT,
		media_type TEXT,
		timestamp TEXT,
		is_from_me INTEGER,
		created_at TEXT
	)`)
	_, _ = chatDB.Exec(`CREATE TABLE IF NOT EXISTS contacts (
		session_id TEXT NOT NULL,
		jid TEXT NOT NULL,
		name TEXT,
		PRIMARY KEY (session_id, jid)
	)`)
	cs := newChatStore(chatDB)
	store.SetOSInfo("PlugaWAPI", [3]uint32{0, 1, 0})
	container := sqlstore.NewWithDB(db, "sqlite3", waLog.Noop)
	if err := container.Upgrade(ctx); err != nil {
		return nil, err
	}
	// chat storage

	store, err := newSessionStore(ctx, db)
	if err != nil {
		return nil, err
	}

	waLogger := waLog.Noop
	if log.Enabled(ctx, slog.LevelDebug) {
		waLogger = waLog.Stdout("WA", "INFO", true)
	}

	broker := NewBroker()
	mgr := newSessionManager(ctx, container, broker, store, waLogger, log, maxCalls, cs)
	broker.SnapshotFn = mgr.snapshotEvents

	gw, _ := newSipGateway(log, ":5060")
	go func() { _ = gw.start(ctx) }()
	return &server{broker: broker, sessions: mgr, chatStore: cs, sip: gw, log: log, staticDir: staticDir}, nil
}
