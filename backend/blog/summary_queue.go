package blog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/harveyTon/trilium-blog/backend/pkg/logger"
)

type AISummaryJob struct {
	NoteID     string
	Title      string
	Content    string
	SourceHash string
}

type AISummaryQueue struct {
	store         SummaryStore
	provider      string
	baseURL       string
	apiKey        string
	model         string
	prompt        string
	concurrency   int
	rateLimit     time.Duration
	maxInputRunes int

	jobs      chan AISummaryJob
	inFlight  map[string]struct{}
	inFlightM sync.Mutex
	client    *http.Client
}

func NewAISummaryQueue(store SummaryStore, provider, baseURL, apiKey, model, prompt string, concurrency, rateLimitMs, timeoutMs, maxInputChars int) *AISummaryQueue {
	if concurrency <= 0 {
		concurrency = 1
	}
	if rateLimitMs <= 0 {
		rateLimitMs = 1200
	}
	if timeoutMs <= 0 {
		timeoutMs = 60000
	}
	if maxInputChars <= 0 {
		maxInputChars = 12000
	}
	q := &AISummaryQueue{
		store:         store,
		provider:      strings.TrimSpace(provider),
		baseURL:       strings.TrimRight(baseURL, "/"),
		apiKey:        apiKey,
		model:         model,
		prompt:        prompt,
		concurrency:   concurrency,
		rateLimit:     time.Duration(rateLimitMs) * time.Millisecond,
		maxInputRunes: maxInputChars,
		jobs:          make(chan AISummaryJob, 128),
		inFlight:      map[string]struct{}{},
		client:        &http.Client{Timeout: time.Duration(timeoutMs) * time.Millisecond},
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
		logger.Logger.Info().Str("note_id", job.NoteID).Msg("Queued AI summary generation")
	default:
		go func() {
			q.jobs <- job
			logger.Logger.Info().Str("note_id", job.NoteID).Msg("Queued AI summary generation after waiting for capacity")
		}()
	}
}

func (q *AISummaryQueue) release(noteID string) {
	q.inFlightM.Lock()
	delete(q.inFlight, noteID)
	q.inFlightM.Unlock()
}

func (q *AISummaryQueue) worker() {
	for job := range q.jobs {
		logger.Logger.Info().Str("note_id", job.NoteID).Msg("Starting AI summary generation")
		_ = q.store.UpsertSummary(StoredSummary{
			NoteID:     job.NoteID,
			Type:       "ai",
			Status:     "processing",
			SourceHash: job.SourceHash,
			Content:    "",
		})
		content, err := q.generate(job.Title, job.Content)
		if err != nil {
			_ = q.store.UpsertSummary(StoredSummary{
				NoteID:     job.NoteID,
				Type:       "ai",
				Status:     "failed",
				SourceHash: job.SourceHash,
				Content:    "",
				Error:      err.Error(),
			})
			logger.Logger.Error().Err(err).Str("note_id", job.NoteID).Msg("AI summary generation failed")
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
		logger.Logger.Info().Str("note_id", job.NoteID).Msg("AI summary generation completed")
		q.release(job.NoteID)
		time.Sleep(q.rateLimit)
	}
}

func (q *AISummaryQueue) generate(title, content string) (string, error) {
	if q.provider != "" && q.provider != "openai-compatible" {
		return "", fmt.Errorf("unsupported ai summary provider: %s", q.provider)
	}
	if q.baseURL == "" || q.apiKey == "" || q.model == "" {
		return "", fmt.Errorf("ai summary provider is not fully configured")
	}
	content = clampSummaryInput(content, q.maxInputRunes)

	requestBody := map[string]any{
		"model": q.model,
		"messages": []map[string]string{
			{"role": "system", "content": q.prompt},
			{"role": "user", "content": buildAISummaryInput(title, content)},
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

func buildAISummaryInput(title, content string) string {
	content = clampSummaryInput(content, 0)
	title = strings.TrimSpace(title)
	if title == "" {
		return content
	}
	return strings.TrimSpace("Title: " + title + "\n\nContent:\n" + content)
}

func clampSummaryInput(content string, maxRunes int) string {
	content = strings.TrimSpace(content)
	if maxRunes <= 0 {
		return content
	}
	runes := []rune(content)
	if len(runes) <= maxRunes {
		return content
	}
	return strings.TrimSpace(string(runes[:maxRunes]))
}
