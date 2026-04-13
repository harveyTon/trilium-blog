package blog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

type AISummaryJob struct {
	NoteID     string
	Content    string
	SourceHash string
}

type AISummaryQueue struct {
	store       SummaryStore
	baseURL     string
	apiKey      string
	model       string
	prompt      string
	concurrency int
	rateLimit   time.Duration

	jobs      chan AISummaryJob
	inFlight  map[string]struct{}
	inFlightM sync.Mutex
	client    *http.Client
}

func NewAISummaryQueue(store SummaryStore, baseURL, apiKey, model, prompt string, concurrency, rateLimitMs int) *AISummaryQueue {
	if concurrency <= 0 {
		concurrency = 1
	}
	if rateLimitMs <= 0 {
		rateLimitMs = 1200
	}
	q := &AISummaryQueue{
		store:       store,
		baseURL:     strings.TrimRight(baseURL, "/"),
		apiKey:      apiKey,
		model:       model,
		prompt:      prompt,
		concurrency: concurrency,
		rateLimit:   time.Duration(rateLimitMs) * time.Millisecond,
		jobs:        make(chan AISummaryJob, 128),
		inFlight:    map[string]struct{}{},
		client:      &http.Client{Timeout: 60 * time.Second},
	}
	for i := 0; i < concurrency; i++ {
		go q.worker()
	}
	return q
}

func (q *AISummaryQueue) Enqueue(job AISummaryJob) {
	q.inFlightM.Lock()
	if _, exists := q.inFlight[job.NoteID]; exists {
		q.inFlightM.Unlock()
		return
	}
	q.inFlight[job.NoteID] = struct{}{}
	q.inFlightM.Unlock()

	select {
	case q.jobs <- job:
	default:
		q.release(job.NoteID)
	}
}

func (q *AISummaryQueue) release(noteID string) {
	q.inFlightM.Lock()
	delete(q.inFlight, noteID)
	q.inFlightM.Unlock()
}

func (q *AISummaryQueue) worker() {
	for job := range q.jobs {
		_ = q.store.UpsertSummary(StoredSummary{
			NoteID:     job.NoteID,
			Type:       "ai",
			Status:     "processing",
			SourceHash: job.SourceHash,
			Content:    "",
		})
		content, err := q.generate(job.Content)
		if err != nil {
			_ = q.store.UpsertSummary(StoredSummary{
				NoteID:     job.NoteID,
				Type:       "ai",
				Status:     "failed",
				SourceHash: job.SourceHash,
				Content:    "",
				Error:      err.Error(),
			})
			q.release(job.NoteID)
			time.Sleep(q.rateLimit)
			continue
		}
		_ = q.store.UpsertSummary(StoredSummary{
			NoteID:     job.NoteID,
			Type:       "ai",
			Status:     "ready",
			SourceHash: job.SourceHash,
			Content:    content,
			Error:      "",
		})
		q.release(job.NoteID)
		time.Sleep(q.rateLimit)
	}
}

func (q *AISummaryQueue) generate(content string) (string, error) {
	if q.baseURL == "" || q.apiKey == "" || q.model == "" {
		return "", fmt.Errorf("ai summary provider is not fully configured")
	}

	requestBody := map[string]any{
		"model": q.model,
		"messages": []map[string]string{
			{"role": "system", "content": q.prompt},
			{"role": "user", "content": content},
		},
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", q.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+q.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := q.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("ai provider returned status %d", resp.StatusCode)
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	if len(result.Choices) == 0 {
		return "", fmt.Errorf("ai provider returned no choices")
	}
	return strings.TrimSpace(result.Choices[0].Message.Content), nil
}
