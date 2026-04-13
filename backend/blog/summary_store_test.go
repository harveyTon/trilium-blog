package blog

import (
	"path/filepath"
	"strings"
	"sync"
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

func TestSummaryStoreDB_ConcurrentAccessDoesNotReturnBusy(t *testing.T) {
	store, err := NewSummaryStoreDB(filepath.Join(t.TempDir(), "summaries.db"))
	if err != nil {
		t.Fatalf("failed to create summary store: %v", err)
	}
	defer store.Close()

	var wg sync.WaitGroup
	errCh := make(chan error, 200)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(worker int) {
			defer wg.Done()
			for j := 0; j < 20; j++ {
				item := StoredSummary{
					NoteID:     "note-1",
					Type:       "code",
					Status:     "ready",
					Content:    "hello",
					SourceHash: "hash",
				}
				if worker%2 == 0 {
					item.Type = "ai"
					item.Status = "processing"
				}
				if err := store.UpsertSummary(item); err != nil {
					errCh <- err
					return
				}
				if _, err := store.GetSummary("note-1", item.Type); err != nil {
					errCh <- err
					return
				}
			}
		}(i)
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		if strings.Contains(err.Error(), "SQLITE_BUSY") || strings.Contains(strings.ToLower(err.Error()), "database is locked") {
			t.Fatalf("unexpected sqlite busy error: %v", err)
		}
		if err != nil {
			t.Fatalf("unexpected error during concurrent access: %v", err)
		}
	}
}
