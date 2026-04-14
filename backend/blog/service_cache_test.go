package blog

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/harveyTon/trilium-blog/backend/etapi"
)

type memoryStore struct {
	mu   sync.Mutex
	data map[string]string
}

func newMemoryStore() *memoryStore {
	return &memoryStore{data: make(map[string]string)}
}

func (s *memoryStore) Get(key string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	value, ok := s.data[key]
	if !ok {
		return "", ErrCacheMiss
	}
	return value, nil
}

func (s *memoryStore) Set(key string, value string, ttlSeconds int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
	return nil
}

func (s *memoryStore) TTL(key string) (time.Duration, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.data[key]; !ok {
		return 0, ErrCacheMiss
	}
	return 60 * time.Second, nil
}

func (s *memoryStore) Del(keys ...string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, k := range keys {
		delete(s.data, k)
	}
	return nil
}

func (s *memoryStore) Keys(pattern string) ([]string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var result []string
	for k := range s.data {
		if matchPattern(k, pattern) {
			result = append(result, k)
		}
	}
	return result, nil
}

func matchPattern(key, pattern string) bool {
	if pattern == "*" {
		return true
	}
	prefix := strings.TrimSuffix(pattern, "*")
	return strings.HasPrefix(key, prefix)
}

func TestGetPostCachesResolvedPost(t *testing.T) {
	noteID := "cached-post"
	content := "<h1>Hello</h1><p>This article should be served from cache on the second request.</p>"

	var mu sync.Mutex
	noteCalls := 0
	contentCalls := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/etapi/notes/" + noteID:
			mu.Lock()
			noteCalls++
			mu.Unlock()
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"noteId":"%s","title":"Cached Post","dateModified":"2026-04-13T12:00:00Z","type":"text","mime":"text/html","attributes":[{"type":"label","name":"blog","value":"true"}]}`, noteID)
		case "/etapi/notes/" + noteID + "/content":
			mu.Lock()
			contentCalls++
			mu.Unlock()
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, _ = w.Write([]byte(content))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	service := NewService(
		etapi.NewClient(server.URL, "token"),
		newMemoryStore(),
	)

	first, err := service.GetPost(noteID)
	if err != nil {
		t.Fatalf("first GetPost failed: %v", err)
	}
	second, err := service.GetPost(noteID)
	if err != nil {
		t.Fatalf("second GetPost failed: %v", err)
	}

	if first == nil || second == nil {
		t.Fatalf("expected cached posts to be returned")
	}
	if first.NoteID != second.NoteID || first.ContentHTML != second.ContentHTML {
		t.Fatalf("expected cached post content to match original fetch")
	}

	mu.Lock()
	defer mu.Unlock()
	if noteCalls != 1 {
		t.Fatalf("expected note metadata to be fetched once, got %d", noteCalls)
	}
	if contentCalls != 1 {
		t.Fatalf("expected note content to be fetched once, got %d", contentCalls)
	}
}

func TestGetAssetCachesAttachmentContent(t *testing.T) {
	attachmentID := "asset-1"
	content := []byte("asset-bytes")
	contentType := "image/png"

	var mu sync.Mutex
	attachmentMetaCalls := 0
	noteCalls := 0
	assetCalls := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/etapi/attachments/" + attachmentID:
			mu.Lock()
			attachmentMetaCalls++
			mu.Unlock()
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"ownerId":"note-1","mime":"%s"}`, contentType)
		case "/etapi/notes/note-1":
			mu.Lock()
			noteCalls++
			mu.Unlock()
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{"noteId":"note-1","title":"Attachment Owner","dateModified":"2026-04-13T12:00:00Z","type":"text","mime":"text/html","attributes":[{"type":"label","name":"blog","value":"true"}]}`)
		case "/etapi/attachments/" + attachmentID + "/content":
			mu.Lock()
			assetCalls++
			mu.Unlock()
			w.Header().Set("Content-Type", contentType)
			_, _ = w.Write(content)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	service := NewService(
		etapi.NewClient(server.URL, "token"),
		newMemoryStore(),
	)

	firstContent, firstType, err := service.GetAsset(attachmentID)
	if err != nil {
		t.Fatalf("first GetAsset failed: %v", err)
	}
	secondContent, secondType, err := service.GetAsset(attachmentID)
	if err != nil {
		t.Fatalf("second GetAsset failed: %v", err)
	}

	if string(firstContent) != string(content) || string(secondContent) != string(content) {
		t.Fatalf("expected cached asset content to match original response")
	}
	if firstType != contentType || secondType != contentType {
		t.Fatalf("expected cached content type %q, got %q and %q", contentType, firstType, secondType)
	}

	mu.Lock()
	defer mu.Unlock()
	if attachmentMetaCalls != 1 {
		t.Fatalf("expected attachment metadata to be fetched once, got %d", attachmentMetaCalls)
	}
	if noteCalls != 1 {
		t.Fatalf("expected owner note metadata to be fetched once, got %d", noteCalls)
	}
	if assetCalls != 1 {
		t.Fatalf("expected attachment content to be fetched once, got %d", assetCalls)
	}
}

func TestListPostsCachesNotesAndContent(t *testing.T) {
	noteID := "listed-post"
	content := "<p>List posts should reuse cached note metadata and content.</p>"

	var mu sync.Mutex
	listCalls := 0
	contentCalls := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/etapi/notes":
			mu.Lock()
			listCalls++
			mu.Unlock()
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"results":[{"noteId":"%s","title":"Listed Post","dateModified":"2026-04-13T12:00:00Z","type":"text","mime":"text/html","attributes":[{"type":"label","name":"blog","value":"true"}]}]}`, noteID)
		case "/etapi/notes/" + noteID + "/content":
			mu.Lock()
			contentCalls++
			mu.Unlock()
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, _ = w.Write([]byte(content))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	service := NewService(
		etapi.NewClient(server.URL, "token"),
		newMemoryStore(),
	)

	first, err := service.ListPosts(1)
	if err != nil {
		t.Fatalf("first ListPosts failed: %v", err)
	}
	second, err := service.ListPosts(1)
	if err != nil {
		t.Fatalf("second ListPosts failed: %v", err)
	}

	if len(first.Items) != 1 || len(second.Items) != 1 {
		t.Fatalf("expected one post in both responses")
	}
	if first.Items[0].Summary == "" || second.Items[0].Summary == "" {
		t.Fatalf("expected summaries to be populated from cached content")
	}

	mu.Lock()
	defer mu.Unlock()
	if listCalls != 1 {
		t.Fatalf("expected note list to be fetched once, got %d", listCalls)
	}
	if contentCalls != 1 {
		t.Fatalf("expected listed note content to be fetched once, got %d", contentCalls)
	}
}

func TestGetAssetRejectsAttachmentsFromNonBlogNotes(t *testing.T) {
	attachmentID := "asset-private"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/etapi/attachments/" + attachmentID:
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{"ownerId":"note-private","mime":"image/png"}`)
		case "/etapi/notes/note-private":
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{"noteId":"note-private","title":"Private Note","dateModified":"2026-04-13T12:00:00Z","type":"text","mime":"text/html","attributes":[]}`)
		case "/etapi/attachments/" + attachmentID + "/content":
			t.Fatalf("attachment content should not be fetched for non-blog note")
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	service := NewService(
		etapi.NewClient(server.URL, "token"),
		newMemoryStore(),
	)

	if _, _, err := service.GetAsset(attachmentID); err == nil {
		t.Fatalf("expected non-blog attachment fetch to fail")
	}
}
