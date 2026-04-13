package blog

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/harveyTon/trilium-blog/backend/etapi"
)

type failingSummaryStore struct{}

func (f failingSummaryStore) GetSummary(noteID, summaryType string) (*StoredSummary, error) {
	return nil, errors.New("summary store unavailable")
}

func (f failingSummaryStore) UpsertSummary(item StoredSummary) error {
	return errors.New("summary store unavailable")
}

func newBlogTestServer(t *testing.T, noteID, content string) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/etapi/notes/" + noteID:
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"noteId":"%s","title":"Test Post","dateModified":"2026-04-13T12:00:00Z","type":"text","mime":"text/html","attributes":[{"type":"label","name":"blog","value":"true"}]}`, noteID)
		case "/etapi/notes/" + noteID + "/content":
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, _ = w.Write([]byte(content))
		default:
			http.NotFound(w, r)
		}
	}))
}

func TestGetPostIgnoresSummaryStoreFailures(t *testing.T) {
	noteID := "note-1"
	content := `<h1>Hello</h1><p>This is the article body used for the summary.</p><pre><code class="language-javascript">const answer = 42;</code></pre>`
	server := newBlogTestServer(t, noteID, content)
	defer server.Close()

	service := NewService(
		etapi.NewClient(server.URL, "token"),
		&NoopStore{},
		WithSummaryStore(failingSummaryStore{}),
	)

	post, err := service.GetPost(noteID)
	if err != nil {
		t.Fatalf("expected article fetch to succeed, got error: %v", err)
	}
	if post == nil {
		t.Fatalf("expected post to be returned")
	}
	if !strings.Contains(post.ContentHTML, "article body") {
		t.Fatalf("expected article HTML to be preserved, got %q", post.ContentHTML)
	}
	if len(post.CodeBlocks) != 1 {
		t.Fatalf("expected 1 code block, got %d", len(post.CodeBlocks))
	}
	if post.CodeBlocks[0].Index != 0 {
		t.Fatalf("expected first code block index 0, got %d", post.CodeBlocks[0].Index)
	}
	if post.CodeBlocks[0].LanguageID != "javascript" {
		t.Fatalf("expected javascript language id, got %q", post.CodeBlocks[0].LanguageID)
	}
	if strings.TrimSpace(post.Summary) == "" {
		t.Fatalf("expected fallback summary text when summary store fails")
	}
}

func TestGetPost_CodeDetectionFailureDoesNotBreakArticle(t *testing.T) {
	noteID := "note-fallback"
	content := "<h1>Hello</h1><p>Body stays readable.</p><pre><code>\u0000\u0001\u0002</code></pre>"
	server := newBlogTestServer(t, noteID, content)
	defer server.Close()

	service := NewService(
		etapi.NewClient(server.URL, "token"),
		&NoopStore{},
		WithSummaryStore(failingSummaryStore{}),
	)

	post, err := service.GetPost(noteID)
	if err != nil {
		t.Fatalf("expected article fetch to succeed, got error: %v", err)
	}
	if post == nil {
		t.Fatalf("expected post to be returned")
	}
	if !strings.Contains(post.ContentHTML, "Body stays readable") {
		t.Fatalf("expected article body to be preserved, got %q", post.ContentHTML)
	}
	if len(post.CodeBlocks) != 1 {
		t.Fatalf("expected fallback code block metadata, got %d items", len(post.CodeBlocks))
	}
	if post.CodeBlocks[0].LanguageID == "" {
		t.Fatalf("expected fallback language id to be populated")
	}
}

func TestGetPostSummariesReturnsPendingAIWhileGenerationRuns(t *testing.T) {
	noteID := "note-2"
	content := "<p>This content should produce a code summary while the AI summary remains pending.</p>"
	blogServer := newBlogTestServer(t, noteID, content)
	defer blogServer.Close()

	releaseAI := make(chan struct{})
	aiServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-releaseAI
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"choices":[{"message":{"content":"AI summary"}}]}`))
	}))
	defer aiServer.Close()
	defer close(releaseAI)

	store, err := NewSummaryStoreDB(filepath.Join(t.TempDir(), "summaries.db"))
	if err != nil {
		t.Fatalf("failed to create summary store: %v", err)
	}
	defer store.Close()

	queue := NewAISummaryQueue(store, "openai-compatible", aiServer.URL, "token", "model", "prompt", 1, 1, 30000, 2000)
	service := NewService(
		etapi.NewClient(blogServer.URL, "token"),
		&NoopStore{},
		WithSummaryStore(store),
		WithAISummaryQueue(queue),
		WithAISummaryEnabled(true),
	)

	summaries, err := service.GetPostSummaries(noteID)
	if err != nil {
		t.Fatalf("expected summary fetch to succeed, got error: %v", err)
	}
	if summaries == nil || summaries.Code == nil || summaries.Code.Status != "ready" {
		t.Fatalf("expected ready code summary, got %#v", summaries)
	}
	if summaries.NoteID != noteID {
		t.Fatalf("expected note ID %q, got %q", noteID, summaries.NoteID)
	}
	if !summaries.AIEnabled {
		t.Fatalf("expected AI summary to be marked enabled")
	}
	if summaries.AI == nil {
		t.Fatalf("expected AI summary status to be returned")
	}
	if summaries.AI.Status != "pending" {
		t.Fatalf("expected pending AI summary during async generation, got %q", summaries.AI.Status)
	}

	select {
	case <-time.After(20 * time.Millisecond):
	case <-releaseAI:
	}
}
