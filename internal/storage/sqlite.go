package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

type SQLiteStore struct {
	db *sql.DB
}

func NewSQLiteStore(dbPath string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite", dbPath+"?_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}
	db.SetMaxOpenConns(1)

	if err := createSchema(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("creating schema: %w", err)
	}

	return &SQLiteStore{db: db}, nil
}

func createSchema(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS entries (
		id               INTEGER PRIMARY KEY AUTOINCREMENT,
		method           TEXT NOT NULL,
		scheme           TEXT NOT NULL,
		host             TEXT NOT NULL,
		path             TEXT NOT NULL,
		query            TEXT DEFAULT '',
		request_headers  TEXT NOT NULL,
		request_body     BLOB,
		status_code      INTEGER DEFAULT 0,
		response_headers TEXT DEFAULT '{}',
		response_body    BLOB,
		created_at       DATETIME DEFAULT CURRENT_TIMESTAMP,
		duration_ms      INTEGER DEFAULT 0
	);
	CREATE INDEX IF NOT EXISTS idx_entries_method ON entries(method);
	CREATE INDEX IF NOT EXISTS idx_entries_host ON entries(host);
	CREATE INDEX IF NOT EXISTS idx_entries_status ON entries(status_code);
	CREATE INDEX IF NOT EXISTS idx_entries_created ON entries(created_at);
	`
	_, err := db.Exec(schema)
	return err
}

func (s *SQLiteStore) Save(entry *Entry) error {
	reqHeaders, err := json.Marshal(entry.RequestHeaders)
	if err != nil {
		return fmt.Errorf("marshaling request headers: %w", err)
	}

	respHeaders, err := json.Marshal(entry.ResponseHeaders)
	if err != nil {
		return fmt.Errorf("marshaling response headers: %w", err)
	}

	ts := entry.Timestamp
	if ts.IsZero() {
		ts = time.Now()
	}

	result, err := s.db.Exec(
		`INSERT INTO entries (method, scheme, host, path, query, request_headers, request_body, status_code, response_headers, response_body, created_at, duration_ms)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		entry.Method,
		entry.Scheme,
		entry.Host,
		entry.Path,
		entry.Query,
		string(reqHeaders),
		entry.RequestBody,
		entry.StatusCode,
		string(respHeaders),
		entry.ResponseBody,
		ts.UTC().Format(time.DateTime),
		entry.DurationMs,
	)
	if err != nil {
		return fmt.Errorf("inserting entry: %w", err)
	}

	entry.ID, _ = result.LastInsertId()
	return nil
}

func (s *SQLiteStore) Get(id int64) (*Entry, error) {
	row := s.db.QueryRow(
		`SELECT id, method, scheme, host, path, query, request_headers, request_body, status_code, response_headers, response_body, created_at, duration_ms
		 FROM entries WHERE id = ?`, id,
	)
	return scanEntry(row)
}

func (s *SQLiteStore) List(limit, offset int) ([]*Entry, error) {
	if limit <= 0 {
		limit = 50
	}
	rows, err := s.db.Query(
		`SELECT id, method, scheme, host, path, query, request_headers, request_body, status_code, response_headers, response_body, created_at, duration_ms
		 FROM entries ORDER BY id DESC LIMIT ? OFFSET ?`, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("querying entries: %w", err)
	}
	defer rows.Close()
	return scanEntries(rows)
}

func (s *SQLiteStore) ListAfter(afterID int64, limit int) ([]*Entry, error) {
	if limit <= 0 {
		limit = 100
	}
	rows, err := s.db.Query(
		`SELECT id, method, scheme, host, path, query, request_headers, request_body, status_code, response_headers, response_body, created_at, duration_ms
		 FROM entries WHERE id > ? ORDER BY id ASC LIMIT ?`, afterID, limit,
	)
	if err != nil {
		return nil, fmt.Errorf("querying entries: %w", err)
	}
	defer rows.Close()
	return scanEntries(rows)
}

func (s *SQLiteStore) Clear() error {
	_, err := s.db.Exec("DELETE FROM entries")
	if err != nil {
		return fmt.Errorf("clearing entries: %w", err)
	}
	return nil
}

func (s *SQLiteStore) Search(params SearchParams) ([]*Entry, error) {
	var conditions []string
	var args []any

	if params.Method != "" {
		conditions = append(conditions, "method = ?")
		args = append(args, strings.ToUpper(params.Method))
	}

	if params.Host != "" {
		if strings.Contains(params.Host, "*") {
			pattern := strings.ReplaceAll(params.Host, "*", "%")
			conditions = append(conditions, "host LIKE ?")
			args = append(args, pattern)
		} else {
			conditions = append(conditions, "host = ?")
			args = append(args, params.Host)
		}
	}

	for _, domain := range params.Domains {
		domain = strings.TrimPrefix(domain, ".")
		conditions = append(conditions, "(host LIKE ? OR host = ?)")
		args = append(args, "%."+domain, domain)
	}

	if params.Path != "" {
		if strings.Contains(params.Path, "*") {
			pattern := strings.ReplaceAll(params.Path, "*", "%")
			conditions = append(conditions, "path LIKE ?")
			args = append(args, pattern)
		} else {
			conditions = append(conditions, "path LIKE ?")
			args = append(args, params.Path+"%")
		}
	}

	if params.Status > 0 {
		conditions = append(conditions, "status_code = ?")
		args = append(args, params.Status)
	}

	query := `SELECT id, method, scheme, host, path, query, request_headers, request_body, status_code, response_headers, response_body, created_at, duration_ms FROM entries`
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	query += " ORDER BY id DESC"

	limit := params.Limit
	if limit <= 0 {
		limit = 100
	}
	query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, params.Offset)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("searching entries: %w", err)
	}
	defer rows.Close()
	return scanEntries(rows)
}

func (s *SQLiteStore) Close() error {
	return s.db.Close()
}

type scanner interface {
	Scan(dest ...any) error
}

func scanEntry(row scanner) (*Entry, error) {
	var e Entry
	var reqHeaders, respHeaders string
	var createdAt string
	var reqBody, respBody []byte

	err := row.Scan(
		&e.ID, &e.Method, &e.Scheme, &e.Host, &e.Path, &e.Query,
		&reqHeaders, &reqBody, &e.StatusCode, &respHeaders, &respBody,
		&createdAt, &e.DurationMs,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("entry not found")
		}
		return nil, fmt.Errorf("scanning entry: %w", err)
	}

	e.RequestBody = reqBody
	e.ResponseBody = respBody

	if err := json.Unmarshal([]byte(reqHeaders), &e.RequestHeaders); err != nil {
		e.RequestHeaders = http.Header{}
	}
	if err := json.Unmarshal([]byte(respHeaders), &e.ResponseHeaders); err != nil {
		e.ResponseHeaders = http.Header{}
	}

	e.Timestamp, _ = time.Parse(time.DateTime, createdAt)

	return &e, nil
}

func scanEntries(rows *sql.Rows) ([]*Entry, error) {
	var entries []*Entry
	for rows.Next() {
		e, err := scanEntry(rows)
		if err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, rows.Err()
}
