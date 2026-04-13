package blog

import (
	"database/sql"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

type SummaryStoreDB struct {
	db *sql.DB
}

type StoredSummary struct {
	NoteID     string
	Type       string
	Status     string
	Content    string
	SourceHash string
	UpdatedAt  string
	Error      string
}

func NewSummaryStoreDB(path string) (*SummaryStoreDB, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(0)
	store := &SummaryStoreDB{db: db}
	if err := store.init(); err != nil {
		_ = db.Close()
		return nil, err
	}
	return store, nil
}

func (s *SummaryStoreDB) init() error {
	for _, pragma := range []string{
		"PRAGMA journal_mode = WAL",
		"PRAGMA synchronous = NORMAL",
		"PRAGMA busy_timeout = 5000",
	} {
		if _, err := s.db.Exec(pragma); err != nil {
			return err
		}
	}

	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS summaries (
			note_id TEXT NOT NULL,
			type TEXT NOT NULL,
			status TEXT NOT NULL,
			content TEXT NOT NULL DEFAULT '',
			source_hash TEXT NOT NULL DEFAULT '',
			updated_at TEXT NOT NULL,
			error TEXT NOT NULL DEFAULT '',
			PRIMARY KEY (note_id, type)
		)
	`)
	return err
}

func (s *SummaryStoreDB) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}

func (s *SummaryStoreDB) GetSummary(noteID, summaryType string) (*StoredSummary, error) {
	row := s.db.QueryRow(`
		SELECT note_id, type, status, content, source_hash, updated_at, error
		FROM summaries
		WHERE note_id = ? AND type = ?
	`, noteID, summaryType)

	var result StoredSummary
	if err := row.Scan(&result.NoteID, &result.Type, &result.Status, &result.Content, &result.SourceHash, &result.UpdatedAt, &result.Error); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}

func (s *SummaryStoreDB) UpsertSummary(item StoredSummary) error {
	if item.UpdatedAt == "" {
		item.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	}
	_, err := s.db.Exec(`
		INSERT INTO summaries (note_id, type, status, content, source_hash, updated_at, error)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(note_id, type) DO UPDATE SET
			status = excluded.status,
			content = excluded.content,
			source_hash = excluded.source_hash,
			updated_at = excluded.updated_at,
			error = excluded.error
	`, item.NoteID, item.Type, item.Status, item.Content, item.SourceHash, item.UpdatedAt, item.Error)
	return err
}
