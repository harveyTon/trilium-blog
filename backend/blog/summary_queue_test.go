package blog

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestAISummaryQueue_EnqueueDeduplicatesAndFailsGracefully(t *testing.T) {
	store, err := NewSummaryStoreDB(filepath.Join(t.TempDir(), "summaries.db"))
	if err != nil {
		t.Fatalf("failed to create summary store: %v", err)
	}
	defer store.Close()

	queue := NewAISummaryQueue(store, "openai-compatible", "", "", "", "prompt", 1, 10, 100, 2000)
	queue.Enqueue(AISummaryJob{NoteID: "note-1", Title: "Title", Content: "hello", SourceHash: "hash"})
	queue.Enqueue(AISummaryJob{NoteID: "note-1", Title: "Title", Content: "hello", SourceHash: "hash"})

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

func TestAISummaryQueue_GenerateIncludesTitleInPrompt(t *testing.T) {
	store, err := NewSummaryStoreDB(filepath.Join(t.TempDir(), "summaries.db"))
	if err != nil {
		t.Fatalf("failed to create summary store: %v", err)
	}
	defer store.Close()

	requestBody := make(chan string, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		var payload struct {
			Messages []struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"messages"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}
		if len(payload.Messages) < 2 {
			t.Fatalf("expected user message to be present")
		}
		requestBody <- payload.Messages[1].Content
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"choices":[{"message":{"content":"summary"}}]}`))
	}))
	defer server.Close()

	queue := NewAISummaryQueue(store, "openai-compatible", server.URL, "token", "model", "prompt", 1, 1, 1000, 2000)
	queue.Enqueue(AISummaryJob{
		NoteID:     "note-1",
		Title:      "My Article",
		Content:    "Body text",
		SourceHash: "hash",
	})

	select {
	case content := <-requestBody:
		if !strings.Contains(content, "Title: My Article") {
			t.Fatalf("expected title in AI summary prompt, got %q", content)
		}
		if !strings.Contains(content, "Content:\nBody text") {
			t.Fatalf("expected content in AI summary prompt, got %q", content)
		}
	case <-time.After(300 * time.Millisecond):
		t.Fatalf("timed out waiting for AI summary request")
	}
}
