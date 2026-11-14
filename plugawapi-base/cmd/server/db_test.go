package main

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestOpenDBConcurrencyConfig(t *testing.T) {
	db, err := openDB(filepath.Join(t.TempDir(), "concurrency.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	if got := db.Stats().MaxOpenConnections; got != 4 {
		t.Fatalf("expected pool capped to 4 connections, got %d", got)
	}

	var mode string
	if err := db.QueryRow("PRAGMA journal_mode").Scan(&mode); err != nil {
		t.Fatal(err)
	}
	if !strings.EqualFold(mode, "wal") {
		t.Fatalf("expected WAL journal mode, got %q", mode)
	}
}
