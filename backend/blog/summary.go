package blog

import (
	"crypto/sha1"
	"encoding/hex"
	"strings"
	"time"
)

type SummaryStore interface {
	GetSummary(noteID, summaryType string) (*StoredSummary, error)
	UpsertSummary(item StoredSummary) error
}

func contentHash(content string) string {
	sum := sha1.Sum([]byte(content))
	return hex.EncodeToString(sum[:])
}

func (s *Service) ensureSummaries(noteID, content string) (*Summaries, error) {
	if s.summaryStore == nil {
		codeSummary := s.extractSummary(s.sanitizeContent(content))
		return &Summaries{
			Code: &SummaryEntry{
				Type:      "code",
				Status:    "ready",
				Text:      codeSummary,
				UpdatedAt: time.Now().UTC().Format(time.RFC3339),
			},
			Fallback: codeSummary,
		}, nil
	}

	hash := contentHash(content)
	codeSummary := s.extractSummary(s.sanitizeContent(content))
	codeStored, err := s.summaryStore.GetSummary(noteID, "code")
	if err != nil {
		return nil, err
	}
	if codeStored == nil || codeStored.SourceHash != hash || codeStored.Content == "" {
		codeStored = &StoredSummary{
			NoteID:     noteID,
			Type:       "code",
			Status:     "ready",
			Content:    codeSummary,
			SourceHash: hash,
		}
		if err := s.summaryStore.UpsertSummary(*codeStored); err != nil {
			return nil, err
		}
	}

	aiStored, err := s.summaryStore.GetSummary(noteID, "ai")
	if err != nil {
		return nil, err
	}

	if s.aiQueue != nil && s.aiEnabled {
		if aiStored == nil || aiStored.SourceHash != hash || aiStored.Status == "" {
			_ = s.summaryStore.UpsertSummary(StoredSummary{
				NoteID:     noteID,
				Type:       "ai",
				Status:     "pending",
				Content:    "",
				SourceHash: hash,
				Error:      "",
			})
			s.aiQueue.Enqueue(AISummaryJob{
				NoteID:     noteID,
				Content:    content,
				SourceHash: hash,
			})
			aiStored, _ = s.summaryStore.GetSummary(noteID, "ai")
		}
	}

	result := &Summaries{
		Fallback: codeStored.Content,
		Code: &SummaryEntry{
			Type:      "code",
			Status:    codeStored.Status,
			Text:      codeStored.Content,
			UpdatedAt: codeStored.UpdatedAt,
			Error:     codeStored.Error,
		},
	}

	if aiStored != nil {
		result.AI = &SummaryEntry{
			Type:      "ai",
			Status:    aiStored.Status,
			Text:      aiStored.Content,
			UpdatedAt: aiStored.UpdatedAt,
			Error:     aiStored.Error,
		}
	}

	return result, nil
}

func preferredSummaryText(summaries *Summaries, fallback string) string {
	if summaries != nil {
		if summaries.AI != nil && summaries.AI.Status == "ready" && strings.TrimSpace(summaries.AI.Text) != "" {
			return summaries.AI.Text
		}
		if summaries.Code != nil && summaries.Code.Status == "ready" && strings.TrimSpace(summaries.Code.Text) != "" {
			return summaries.Code.Text
		}
		if strings.TrimSpace(summaries.Fallback) != "" {
			return summaries.Fallback
		}
	}
	return fallback
}
