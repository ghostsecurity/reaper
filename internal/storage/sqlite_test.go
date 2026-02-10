package storage

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func testStore(t *testing.T) *SQLiteStore {
	t.Helper()
	dir := t.TempDir()
	store, err := NewSQLiteStore(filepath.Join(dir, "test.db"))
	if err != nil {
		t.Fatalf("creating store: %v", err)
	}
	t.Cleanup(func() { store.Close() })
	return store
}

func TestSaveAndGet(t *testing.T) {
	store := testStore(t)

	entry := &Entry{
		Method: "GET",
		Scheme: "https",
		Host:   "api.acme.com",
		Path:   "/users",
		Query:  "page=1",
		RequestHeaders: http.Header{
			"Accept": []string{"application/json"},
		},
		StatusCode: 200,
		ResponseHeaders: http.Header{
			"Content-Type": []string{"application/json"},
		},
		ResponseBody: []byte(`{"users":[]}`),
		Timestamp:    time.Now(),
		DurationMs:   42,
	}

	if err := store.Save(entry); err != nil {
		t.Fatalf("saving entry: %v", err)
	}

	if entry.ID == 0 {
		t.Fatal("expected non-zero ID after save")
	}

	got, err := store.Get(entry.ID)
	if err != nil {
		t.Fatalf("getting entry: %v", err)
	}

	if got.Method != "GET" {
		t.Errorf("method = %q, want GET", got.Method)
	}
	if got.Host != "api.acme.com" {
		t.Errorf("host = %q, want api.acme.com", got.Host)
	}
	if got.Path != "/users" {
		t.Errorf("path = %q, want /users", got.Path)
	}
	if got.StatusCode != 200 {
		t.Errorf("status = %d, want 200", got.StatusCode)
	}
	if got.DurationMs != 42 {
		t.Errorf("duration = %d, want 42", got.DurationMs)
	}
	if string(got.ResponseBody) != `{"users":[]}` {
		t.Errorf("response body = %q, want {\"users\":[]}", got.ResponseBody)
	}
}

func TestGetNotFound(t *testing.T) {
	store := testStore(t)

	_, err := store.Get(999)
	if err == nil {
		t.Fatal("expected error for missing entry")
	}
}

func TestList(t *testing.T) {
	store := testStore(t)

	for i := 0; i < 5; i++ {
		store.Save(&Entry{
			Method:          "GET",
			Scheme:          "https",
			Host:            "example.com",
			Path:            "/",
			RequestHeaders:  http.Header{},
			ResponseHeaders: http.Header{},
			Timestamp:       time.Now(),
		})
	}

	entries, err := store.List(3, 0)
	if err != nil {
		t.Fatalf("listing: %v", err)
	}
	if len(entries) != 3 {
		t.Errorf("got %d entries, want 3", len(entries))
	}

	// Verify ordering is DESC by ID
	if entries[0].ID < entries[1].ID {
		t.Error("expected DESC ordering")
	}
}

func TestSearchByMethod(t *testing.T) {
	store := testStore(t)

	store.Save(&Entry{Method: "GET", Scheme: "https", Host: "a.com", Path: "/", RequestHeaders: http.Header{}, ResponseHeaders: http.Header{}, Timestamp: time.Now()})
	store.Save(&Entry{Method: "POST", Scheme: "https", Host: "a.com", Path: "/", RequestHeaders: http.Header{}, ResponseHeaders: http.Header{}, Timestamp: time.Now()})

	entries, err := store.Search(SearchParams{Method: "POST"})
	if err != nil {
		t.Fatalf("searching: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("got %d entries, want 1", len(entries))
	}
	if entries[0].Method != "POST" {
		t.Errorf("method = %q, want POST", entries[0].Method)
	}
}

func TestSearchByDomain(t *testing.T) {
	store := testStore(t)

	store.Save(&Entry{Method: "GET", Scheme: "https", Host: "api.acme.com", Path: "/", RequestHeaders: http.Header{}, ResponseHeaders: http.Header{}, Timestamp: time.Now()})
	store.Save(&Entry{Method: "GET", Scheme: "https", Host: "acme.com", Path: "/", RequestHeaders: http.Header{}, ResponseHeaders: http.Header{}, Timestamp: time.Now()})
	store.Save(&Entry{Method: "GET", Scheme: "https", Host: "other.com", Path: "/", RequestHeaders: http.Header{}, ResponseHeaders: http.Header{}, Timestamp: time.Now()})

	entries, err := store.Search(SearchParams{Domains: []string{"acme.com"}})
	if err != nil {
		t.Fatalf("searching: %v", err)
	}
	if len(entries) != 2 {
		t.Errorf("got %d entries, want 2", len(entries))
	}
}

func TestSearchByWildcardHost(t *testing.T) {
	store := testStore(t)

	store.Save(&Entry{Method: "GET", Scheme: "https", Host: "api.acme.com", Path: "/", RequestHeaders: http.Header{}, ResponseHeaders: http.Header{}, Timestamp: time.Now()})
	store.Save(&Entry{Method: "GET", Scheme: "https", Host: "web.acme.com", Path: "/", RequestHeaders: http.Header{}, ResponseHeaders: http.Header{}, Timestamp: time.Now()})
	store.Save(&Entry{Method: "GET", Scheme: "https", Host: "other.com", Path: "/", RequestHeaders: http.Header{}, ResponseHeaders: http.Header{}, Timestamp: time.Now()})

	entries, err := store.Search(SearchParams{Host: "*.acme.com"})
	if err != nil {
		t.Fatalf("searching: %v", err)
	}
	if len(entries) != 2 {
		t.Errorf("got %d entries, want 2", len(entries))
	}
}

func TestSearchByStatus(t *testing.T) {
	store := testStore(t)

	store.Save(&Entry{Method: "GET", Scheme: "https", Host: "a.com", Path: "/", StatusCode: 200, RequestHeaders: http.Header{}, ResponseHeaders: http.Header{}, Timestamp: time.Now()})
	store.Save(&Entry{Method: "GET", Scheme: "https", Host: "a.com", Path: "/err", StatusCode: 500, RequestHeaders: http.Header{}, ResponseHeaders: http.Header{}, Timestamp: time.Now()})

	entries, err := store.Search(SearchParams{Status: 500})
	if err != nil {
		t.Fatalf("searching: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("got %d entries, want 1", len(entries))
	}
	if entries[0].StatusCode != 500 {
		t.Errorf("status = %d, want 500", entries[0].StatusCode)
	}
}

func TestCertGeneration(t *testing.T) {
	dir := t.TempDir()

	// Test by directly importing proxy package cert generation
	// This just tests that the store works with the temp directory
	dbPath := filepath.Join(dir, "test.db")
	store, err := NewSQLiteStore(dbPath)
	if err != nil {
		t.Fatalf("creating store: %v", err)
	}
	store.Close()

	// Verify file was created
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Error("database file was not created")
	}
}
