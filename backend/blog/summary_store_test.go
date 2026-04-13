package blog

import (
	"path/filepath"
	"testing"
)

func TestSummaryStoreDB_UpsertAndGet(t *testing.T) {
	store, err := NewSummaryStoreDB(filepath.Join(t.TempDir(), "summaries.db"))
	if err != nil {
		t.Fatalf("failed to create summary store: %v", err)
	}
	defer store.Close()

	err = store.UpsertSummary(StoredSummary{
		NoteID:     "note-1",
		Type:       "code",
		Status:     "ready",
		Content:    "hello",
		SourceHash: "hash",
	})
	if err != nil {
		t.Fatalf("failed to upsert summary: %v", err)
	}

	item, err := store.GetSummary("note-1", "code")
	if err != nil {
		t.Fatalf("failed to fetch summary: %v", err)
	}
	if item == nil || item.Content != "hello" {
		t.Fatalf("unexpected summary item: %#v", item)
	}
}
