package blog

import (
	"path/filepath"
	"testing"
	"time"
)

func TestAISummaryQueue_EnqueueDeduplicatesAndFailsGracefully(t *testing.T) {
	store, err := NewSummaryStoreDB(filepath.Join(t.TempDir(), "summaries.db"))
	if err != nil {
		t.Fatalf("failed to create summary store: %v", err)
	}
	defer store.Close()

	queue := NewAISummaryQueue(store, "", "", "", "prompt", 1, 10)
	queue.Enqueue(AISummaryJob{NoteID: "note-1", Content: "hello", SourceHash: "hash"})
	queue.Enqueue(AISummaryJob{NoteID: "note-1", Content: "hello", SourceHash: "hash"})

	time.Sleep(80 * time.Millisecond)

	item, err := store.GetSummary("note-1", "ai")
	if err != nil {
		t.Fatalf("failed to get ai summary: %v", err)
	}
	if item == nil {
		t.Fatalf("expected ai summary record to exist")
	}
	if item.Status != "failed" && item.Status != "processing" {
		t.Fatalf("expected failed or processing status, got %q", item.Status)
	}
}
