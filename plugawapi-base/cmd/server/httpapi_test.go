package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"log/slog"
	"testing"
)

func TestHealth(t *testing.T) {
	dbPath := "test_health.db"
	defer os.Remove(dbPath)
	defer os.Remove(dbPath + "_chat.db")
	srv, err := newServer(context.Background(), dbPath, "", 8, slog.Default())
	if err != nil {
		t.Fatal(err)
	}
	defer srv.sessions.disconnectAll()
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/health", nil)
	srv.routes().ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

func TestSessionCreateAndList(t *testing.T) {
	dbPath := "test_sess.db"
	defer os.Remove(dbPath)
	defer os.Remove(dbPath + "_chat.db")
	srv, err := newServer(context.Background(), dbPath, "", 8, slog.Default())
	if err != nil {
		t.Fatal(err)
	}
	defer srv.sessions.disconnectAll()
	// create
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/sessions", nil)
	srv.routes().ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("create: expected 200, got %d body=%s", rr.Code, rr.Body.String())
	}
	// list
	rr2 := httptest.NewRecorder()
	req2 := httptest.NewRequest("GET", "/api/sessions", nil)
	srv.routes().ServeHTTP(rr2, req2)
	if rr2.Code != http.StatusOK {
		t.Fatalf("list: expected 200, got %d", rr2.Code)
	}
}

func TestSessionNotFound(t *testing.T) {
	dbPath := "test_nf.db"
	defer os.Remove(dbPath)
	defer os.Remove(dbPath + "_chat.db")
	srv, err := newServer(context.Background(), dbPath, "", 8, slogDefault())
	if err != nil {
		t.Fatal(err)
	}
	defer srv.sessions.disconnectAll()
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/sessions/nonexistent/chats", nil)
	srv.routes().ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rr.Code)
	}
}

func slogDefault() *slog.Logger {
	return slog.Default()
}
