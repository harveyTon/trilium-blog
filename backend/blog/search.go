package blog

import (
	"sort"
	"strings"
	"unicode"

	"github.com/harveyTon/trilium-blog/backend/etapi"
)

type searchCandidate struct {
	Post            Post
	PlainText       string
	NormalizedTitle string
	NormalizedBody  string
	Score           int
}

func normalizeSearchText(text string) string {
	var b strings.Builder
	lastSpace := false
	for _, r := range strings.ToLower(text) {
		switch {
		case unicode.IsLetter(r), unicode.IsDigit(r), unicode.Is(unicode.Han, r):
			b.WriteRune(r)
			lastSpace = false
		case unicode.IsSpace(r), unicode.IsPunct(r), unicode.IsSymbol(r):
			if !lastSpace {
				b.WriteRune(' ')
				lastSpace = true
			}
		}
	}
	return strings.TrimSpace(b.String())
}

func containsSubsequence(haystack, needle string) bool {
	if needle == "" {
		return true
	}
	hr := []rune(haystack)
	nr := []rune(needle)
	j := 0
	for _, r := range hr {
		if r == nr[j] {
			j++
			if j == len(nr) {
				return true
			}
		}
	}
	return false
}

func searchScore(titleNorm, bodyNorm, queryNorm string) int {
	if queryNorm == "" {
		return 0
	}

	score := 0
	switch {
	case strings.Contains(titleNorm, queryNorm):
		score += 120
	case containsSubsequence(titleNorm, strings.ReplaceAll(queryNorm, " ", "")):
		score += 70
	}

	switch {
	case strings.Contains(bodyNorm, queryNorm):
		score += 80
	case containsSubsequence(bodyNorm, strings.ReplaceAll(queryNorm, " ", "")):
		score += 40
	}

	for _, token := range strings.Fields(queryNorm) {
		if token == "" {
			continue
		}
		if strings.Contains(titleNorm, token) {
			score += 20
		}
		if strings.Contains(bodyNorm, token) {
			score += 10
		}
	}

	return score
}

func extractSnippet(text, query string) string {
	plain := strings.TrimSpace(htmlEntityDecode(text))
	if plain == "" {
		return ""
	}

	query = strings.TrimSpace(query)
	runes := []rune(plain)
	if len(runes) <= 140 {
		return plain
	}

	lowerRunes := []rune(strings.ToLower(plain))
	lowerQueryRunes := []rune(strings.ToLower(query))
	if idx := runeIndex(lowerRunes, lowerQueryRunes); idx >= 0 {
		start := idx - 40
		if start < 0 {
			start = 0
		}
		end := idx + len(lowerQueryRunes) + 80
		if end > len(runes) {
			end = len(runes)
		}
		snippet := strings.TrimSpace(string(runes[start:end]))
		if start > 0 {
			snippet = "..." + snippet
		}
		if end < len(runes) {
			snippet += "..."
		}
		return snippet
	}

	return strings.TrimSpace(string(runes[:140])) + "..."
}

func runeIndex(haystack, needle []rune) int {
	if len(needle) == 0 {
		return 0
	}
	for i := 0; i <= len(haystack)-len(needle); i++ {
		match := true
		for j := range needle {
			if haystack[i+j] != needle[j] {
				match = false
				break
			}
		}
		if match {
			return i
		}
	}
	return -1
}

func buildSearchMatch(post Post, plainText, query string) SearchMatch {
	titleLower := strings.ToLower(post.Title)
	queryLower := strings.ToLower(query)
	titleMatched := queryLower != "" && strings.Contains(titleLower, queryLower)
	return SearchMatch{
		TitleMatched: titleMatched,
		Snippet:      extractSnippet(plainText, query),
	}
}

func isSearchMatch(note etapi.Note, content, query string) (searchCandidate, bool) {
	sanitized := sanitizeSearchContent(content)
	plainText := strings.TrimSpace(htmlToPlainText(sanitized))
	post := Post{
		NoteID:       note.NoteID,
		Title:        note.Title,
		DateModified: note.DateModified,
		Summary:      extractSearchSummary(plainText),
	}

	titleNorm := normalizeSearchText(note.Title)
	bodyNorm := normalizeSearchText(plainText)
	queryNorm := normalizeSearchText(query)
	score := searchScore(titleNorm, bodyNorm, queryNorm)
	if score <= 0 {
		return searchCandidate{}, false
	}

	return searchCandidate{
		Post:            post,
		PlainText:       plainText,
		NormalizedTitle: titleNorm,
		NormalizedBody:  bodyNorm,
		Score:           score,
	}, true
}

func sortSearchCandidates(items []searchCandidate) {
	sort.SliceStable(items, func(i, j int) bool {
		if items[i].Score == items[j].Score {
			return parseDate(items[i].Post.DateModified).After(parseDate(items[j].Post.DateModified))
		}
		return items[i].Score > items[j].Score
	})
}
