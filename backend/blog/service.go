package blog

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"unicode"

	"github.com/PuerkitoBio/goquery"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/harveyTon/trilium-blog/backend/etapi"
	"github.com/microcosm-cc/bluemonday"
)

type Store interface {
	Get(key string) (string, error)
	Set(key string, value string, ttlSeconds int) error
}

type Service struct {
	etapiClient       *etapi.Client
	store             Store
	summaryStore      SummaryStore
	aiQueue           *AISummaryQueue
	blogTitle         string
	blogSubtitle      string
	domain            string
	pageSize          int
	imageProxyEnabled bool
	imageProxyBaseUrl string
	aiEnabled         bool
}

type ServiceOption func(*Service)

const (
	notesCacheTTLSeconds      = 300
	noteCacheTTLSeconds       = 300
	noteContentTTLSeconds     = 300
	attachmentCacheTTLSeconds = 3600
)

func WithBlogTitle(title string) ServiceOption {
	return func(s *Service) { s.blogTitle = title }
}

func WithBlogSubtitle(subtitle string) ServiceOption {
	return func(s *Service) { s.blogSubtitle = subtitle }
}

func WithDomain(domain string) ServiceOption {
	return func(s *Service) { s.domain = domain }
}

func WithPageSize(size int) ServiceOption {
	return func(s *Service) { s.pageSize = size }
}

func WithImageProxyEnabled(enabled bool) ServiceOption {
	return func(s *Service) { s.imageProxyEnabled = enabled }
}

func WithImageProxyBaseUrl(baseUrl string) ServiceOption {
	return func(s *Service) { s.imageProxyBaseUrl = baseUrl }
}

func WithSummaryStore(store SummaryStore) ServiceOption {
	return func(s *Service) { s.summaryStore = store }
}

func WithAISummaryQueue(queue *AISummaryQueue) ServiceOption {
	return func(s *Service) { s.aiQueue = queue }
}

func WithAISummaryEnabled(enabled bool) ServiceOption {
	return func(s *Service) { s.aiEnabled = enabled }
}

func NewService(client *etapi.Client, store Store, opts ...ServiceOption) *Service {
	s := &Service{
		etapiClient: client,
		store:       store,
		pageSize:    9,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

type cachedAttachment struct {
	Content     []byte `json:"content"`
	ContentType string `json:"contentType"`
}

func (s *Service) cacheString(key string, ttlSeconds int, loader func() (string, error)) (string, error) {
	if s.store != nil {
		if cachedValue, err := s.store.Get(key); err == nil {
			return cachedValue, nil
		}
	}

	value, err := loader()
	if err != nil {
		return "", err
	}

	if s.store != nil {
		_ = s.store.Set(key, value, ttlSeconds)
	}
	return value, nil
}

func cacheJSON[T any](store Store, key string, ttlSeconds int, loader func() (T, error)) (T, error) {
	var zero T

	if store != nil {
		if cachedValue, err := store.Get(key); err == nil {
			var decoded T
			if unmarshalErr := json.Unmarshal([]byte(cachedValue), &decoded); unmarshalErr == nil {
				return decoded, nil
			}
		}
	}

	value, err := loader()
	if err != nil {
		return zero, err
	}

	if store != nil {
		if encoded, marshalErr := json.Marshal(value); marshalErr == nil {
			_ = store.Set(key, string(encoded), ttlSeconds)
		}
	}

	return value, nil
}

func (s *Service) getCachedNotes(search string) ([]etapi.Note, error) {
	return cacheJSON(s.store, fmt.Sprintf("notes:%s", search), notesCacheTTLSeconds, func() ([]etapi.Note, error) {
		return s.etapiClient.GetNotes(search)
	})
}

func (s *Service) getCachedNote(noteID string) (*etapi.Note, error) {
	return cacheJSON(s.store, fmt.Sprintf("note:%s", noteID), noteCacheTTLSeconds, func() (*etapi.Note, error) {
		return s.etapiClient.GetNote(noteID)
	})
}

func (s *Service) getCachedNoteContent(noteID string) (string, error) {
	return s.cacheString(fmt.Sprintf("note-content:%s", noteID), noteContentTTLSeconds, func() (string, error) {
		return s.etapiClient.GetNoteContent(noteID)
	})
}

func (s *Service) getCachedAttachment(attachmentID string) (*cachedAttachment, error) {
	return cacheJSON(s.store, fmt.Sprintf("attachment:%s", attachmentID), attachmentCacheTTLSeconds, func() (*cachedAttachment, error) {
		attachment, err := s.etapiClient.GetAttachment(attachmentID)
		if err != nil {
			return nil, err
		}

		note, err := s.getCachedNote(attachment.OwnerID)
		if err != nil {
			return nil, err
		}
		if !hasBlogLabel(note.Attributes) {
			return nil, ErrNotBlogPost
		}

		content, err := s.etapiClient.GetAttachmentContentBytes(attachmentID)
		if err != nil {
			return nil, err
		}

		return &cachedAttachment{
			Content:     content,
			ContentType: attachment.Mime,
		}, nil
	})
}

func (s *Service) GetSite() Site {
	return Site{
		Title:    s.blogTitle,
		Subtitle: s.blogSubtitle,
		Domain:   s.domain,
		ImageProxy: ImageProxyConfig{
			Enabled: s.imageProxyEnabled,
			BaseURL: s.imageProxyBaseUrl,
		},
	}
}

func (s *Service) ListPosts(page int) (*PostList, error) {
	notes, err := s.getCachedNotes("#blog=true")
	if err != nil {
		return nil, err
	}

	var posts []Post
	for _, n := range notes {
		if n.Type == "text" && hasBlogLabel(n.Attributes) {
			posts = append(posts, Post{
				NoteID:       n.NoteID,
				Title:        n.Title,
				DateModified: n.DateModified,
			})
		}
	}

	total := len(posts)
	pageSize := s.pageSize
	if pageSize <= 0 {
		pageSize = 9
	}
	totalPages := (total + pageSize - 1) / pageSize
	if totalPages == 0 {
		totalPages = 1
	}

	start := (page - 1) * pageSize
	end := start + pageSize
	if end > total {
		end = total
	}
	if start >= total {
		start = total
	}

	pagePosts := posts[start:end]

	var wg sync.WaitGroup
	var mu sync.Mutex
	var fetchErr error

	for i := range pagePosts {
		wg.Add(1)
		// idx is passed as a value to avoid closure issues with the loop variable
		go func(idx int) {
			defer wg.Done()
			content, err := s.getCachedNoteContent(pagePosts[idx].NoteID)
			if err != nil {
				mu.Lock()
				if fetchErr == nil {
					fetchErr = err
				}
				mu.Unlock()
				return
			}
			sanitized := s.sanitizeContent(content)
			summary := s.extractSummary(sanitized)
			mu.Lock()
			pagePosts[idx].Summary = summary
			summaries := s.resolveSummaries(pagePosts[idx].NoteID, pagePosts[idx].Title, content)
			if summaries != nil {
				pagePosts[idx].Summaries = summaries
				pagePosts[idx].Summary = preferredSummaryText(summaries, summary)
			}
			mu.Unlock()
		}(i)
	}
	wg.Wait()

	if fetchErr != nil && len(pagePosts) == 0 {
		return nil, fetchErr
	}

	return &PostList{
		Items:      pagePosts,
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

func (s *Service) SearchPosts(query string, preview bool, limit int) (*SearchResponse, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return &SearchResponse{
			Query: query,
			Total: 0,
			Items: []SearchItem{},
		}, nil
	}

	notes, err := s.getCachedNotes("#blog=true")
	if err != nil {
		return nil, err
	}

	candidates := make([]searchCandidate, 0, len(notes))
	for _, note := range notes {
		if note.Type != "text" || !hasBlogLabel(note.Attributes) {
			continue
		}
		content, err := s.getCachedNoteContent(note.NoteID)
		if err != nil {
			continue
		}
		candidate, ok := isSearchMatch(note, content, query)
		if !ok {
			continue
		}
		summaries, sumErr := s.ensureSummaries(note.NoteID, note.Title, content)
		if sumErr == nil {
			candidate.Post.Summaries = summaries
			candidate.Post.Summary = preferredSummaryText(summaries, candidate.Post.Summary)
		}
		candidates = append(candidates, candidate)
	}

	sortSearchCandidates(candidates)

	total := len(candidates)
	if preview {
		if limit <= 0 {
			limit = 5
		}
		if total > limit {
			candidates = candidates[:limit]
		}
	}

	items := make([]SearchItem, 0, len(candidates))
	for _, candidate := range candidates {
		items = append(items, SearchItem{
			Post:  candidate.Post,
			Match: buildSearchMatch(candidate.Post, candidate.PlainText, query),
		})
	}

	return &SearchResponse{
		Query: query,
		Total: total,
		Items: items,
	}, nil
}

func (s *Service) ListFeaturedPosts() ([]Post, error) {
	notes, err := s.getCachedNotes("#blogtop=true")
	if err != nil {
		return nil, err
	}

	posts := make([]Post, 0, len(notes))
	for _, note := range notes {
		if note.Type != "text" || !hasFeaturedLabel(note.Attributes) {
			continue
		}

		post := Post{
			NoteID:       note.NoteID,
			Title:        note.Title,
			DateModified: note.DateModified,
		}

		content, err := s.getCachedNoteContent(note.NoteID)
		if err == nil {
			post.Summary = s.extractSummary(s.sanitizeContent(content))
			summaries := s.resolveSummaries(note.NoteID, note.Title, content)
			if summaries != nil {
				post.Summaries = summaries
				post.Summary = preferredSummaryText(summaries, post.Summary)
			}
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (s *Service) GenerateSitemap() (string, error) {
	notes, err := s.getCachedNotes("#blog=true")
	if err != nil {
		return "", err
	}

	var posts []Post
	for _, n := range notes {
		if n.Type == "text" && hasBlogLabel(n.Attributes) {
			posts = append(posts, Post{
				NoteID:       n.NoteID,
				Title:        n.Title,
				DateModified: n.DateModified,
			})
		}
	}

	domain := strings.TrimRight(s.domain, "/")

	var sb strings.Builder
	sb.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	sb.WriteString(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">` + "\n")

	for _, p := range posts {
		sb.WriteString("  <url>\n")
		sb.WriteString("    <loc>" + domain + "/post/" + p.NoteID + "</loc>\n")
		sb.WriteString("    <lastmod>" + p.DateModified + "</lastmod>\n")
		sb.WriteString("  </url>\n")
	}

	sb.WriteString("</urlset>")
	return sb.String(), nil
}

func (s *Service) GetPost(noteId string) (*Post, error) {
	note, err := s.getCachedNote(noteId)
	if err != nil {
		return nil, err
	}

	if !hasBlogLabel(note.Attributes) {
		return nil, ErrNotBlogPost
	}

	content, err := s.getCachedNoteContent(noteId)
	if err != nil {
		return nil, err
	}

	sanitized := s.sanitizeContent(content)
	toc, modifiedHtml := s.extractTOC(sanitized)
	processed, codeBlocks := s.processContent(modifiedHtml)
	summaryText := s.extractSummary(sanitized)
	summaries := s.resolveSummaries(note.NoteID, note.Title, content)
	if summaries != nil {
		summaryText = preferredSummaryText(summaries, summaryText)
	}

	return &Post{
		NoteID:       note.NoteID,
		Title:        note.Title,
		DateModified: note.DateModified,
		ContentHTML:  processed,
		CodeBlocks:   codeBlocks,
		TOC:          toc,
		PageURL:      getPageURL(note.Attributes),
		Summary:      summaryText,
		Summaries:    summaries,
	}, nil
}

func (s *Service) GetPostSummaries(noteId string) (*Summaries, error) {
	note, err := s.getCachedNote(noteId)
	if err != nil {
		return nil, err
	}
	if !hasBlogLabel(note.Attributes) {
		return nil, ErrNotBlogPost
	}

	content, err := s.getCachedNoteContent(noteId)
	if err != nil {
		return nil, err
	}

	return s.resolveSummaries(note.NoteID, note.Title, content), nil
}

func (s *Service) GetAsset(attachmentId string) ([]byte, string, error) {
	asset, err := s.getCachedAttachment(attachmentId)
	if err != nil {
		return nil, "", err
	}
	return asset.Content, asset.ContentType, nil
}

var ErrNotBlogPost = &BlogError{Message: "note is not a blog post"}

type BlogError struct {
	Message string
}

func (e *BlogError) Error() string {
	return e.Message
}

func hasBlogLabel(attrs []etapi.Attribute) bool {
	return hasTrueLabel(attrs, "blog")
}

func hasFeaturedLabel(attrs []etapi.Attribute) bool {
	return hasTrueLabel(attrs, "blogtop")
}

func hasTrueLabel(attrs []etapi.Attribute, name string) bool {
	for _, a := range attrs {
		if a.Type == "label" && a.Name == name && a.Value == "true" {
			return true
		}
	}
	return false
}

func hasTrueLabelFromAttributes(attrs []attribute, name string) bool {
	for _, a := range attrs {
		if a.Type == "label" && a.Name == name && a.Value == "true" {
			return true
		}
	}
	return false
}

func getPageURL(attrs []etapi.Attribute) string {
	for _, a := range attrs {
		if a.Type == "label" && a.Name == "pageUrl" {
			return a.Value
		}
	}
	return ""
}

func sanitizeSearchContent(html string) string {
	p := bluemonday.UGCPolicy()
	return string(p.SanitizeBytes([]byte(html)))
}

func htmlToPlainText(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return ""
	}

	doc.Find("pre, code, style, script, svg, canvas, video, audio").Each(func(i int, sel *goquery.Selection) {
		sel.Remove()
	})
	doc.Find("a").Each(func(i int, sel *goquery.Selection) {
		text := strings.TrimSpace(sel.Text())
		sel.ReplaceWithHtml(text)
	})

	return strings.TrimSpace(doc.Text())
}

func extractSearchSummary(text string) string {
	text = strings.TrimSpace(text)
	if text == "" {
		return ""
	}
	runes := []rune(text)
	if len(runes) <= 120 {
		return text
	}
	return strings.TrimSpace(string(runes[:120])) + "..."
}

func (s *Service) sanitizeContent(html string) string {
	p := bluemonday.UGCPolicy()
	p.AllowAttrs("class").OnElements("code")
	return string(p.SanitizeBytes([]byte(html)))
}

func (s *Service) extractTOC(html string) ([]TOCItem, string) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, html
	}

	var toc []TOCItem
	doc.Find("h1, h2, h3").EachWithBreak(func(i int, sel *goquery.Selection) bool {
		text := strings.TrimSpace(sel.Text())
		if text == "" {
			return true
		}
		id := sel.AttrOr("id", "")
		if id == "" {
			id = generateID(text)
			sel.SetAttr("id", id)
		}
		level := 1
		if strings.HasPrefix(sel.Nodes[0].Data, "h2") {
			level = 2
		} else if strings.HasPrefix(sel.Nodes[0].Data, "h3") {
			level = 3
		}
		toc = append(toc, TOCItem{ID: id, Title: text, Level: level})
		return true
	})

	result, _ := doc.Find("body").Html()
	if result == "" {
		return toc, html
	}
	return toc, strings.TrimSpace(result)
}

func (s *Service) extractSummary(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return ""
	}

	// Remove all noise elements
	doc.Find("pre, code, h1, h2, h3, h4, h5, h6, blockquote, img, hr, table, figure, nav, style, script, .toc, .article-toc, iframe, svg, canvas, video, audio").Each(func(i int, sel *goquery.Selection) {
		sel.Remove()
	})

	// Replace links with their text content
	doc.Find("a").Each(func(i int, sel *goquery.Selection) {
		text := strings.TrimSpace(sel.Text())
		sel.ReplaceWithHtml(text)
	})

	// Get plain text
	text := strings.TrimSpace(doc.Text())

	// Decode HTML entities (e.g. &amp; &lt; &gt; &quot; &#39; &nbsp;)
	text = htmlEntityDecode(text)

	// Normalize whitespace: collapse multiple spaces/newlines to single space
	text = collapseWhitespace(text)

	// Remove control characters and other invisible garbage
	text = cleanInvisibleChars(text)

	// Remove consecutive repeated punctuation (e.g. "..." "——" "，，，")
	text = cleanRepeatedPunctuation(text)

	// Try to extract the first meaningful paragraph
	// Split by double newline or single newline that separates blocks
	paragraphs := splitParagraphs(text)
	var bestParagraph string

	for _, p := range paragraphs {
		p = strings.TrimSpace(p)
		if len(p) < 30 {
			continue
		}
		// Skip if looks like a heading or list item (short, ends with : or no terminal punctuation)
		if len(p) < 80 && (strings.HasSuffix(p, ":") || strings.HasSuffix(p, "：")) {
			continue
		}
		// Skip if mostly symbols/numbers
		if isMostlySymbols(p) > 40 {
			continue
		}
		bestParagraph = p
		break
	}

	if bestParagraph == "" && len(text) > 0 {
		// Fallback: use the whole text
		bestParagraph = text
	}

	if bestParagraph == "" {
		return ""
	}

	// Truncate to 90-160 runes, at a natural boundary
	result := truncateAtBoundary(bestParagraph, 90, 160)

	// Final cleanup: strip trailing punctuation and invisible chars
	result = cleanTrailingChars(result)

	// Final check: if result is too short or empty, return empty
	if len([]rune(result)) < 20 {
		return ""
	}

	return result
}

// htmlEntityDecode converts common HTML entities to their unicode characters
func htmlEntityDecode(s string) string {
	replacer := strings.NewReplacer(
		"&amp;", "&",
		"&lt;", "<",
		"&gt;", ">",
		"&quot;", `"`,
		"&#39;", "'",
		"&apos;", "'",
		"&nbsp;", " ",
		"&mdash;", "—",
		"&ndash;", "–",
		"&hellip;", "…",
		"&copy;", "©",
		"&reg;", "®",
		"&trade;", "™",
		"&#x27;", "'",
		"&#x2F;", "/",
		"&lsquo;", "'",
		"&rsquo;", "'",
		"&ldquo;", "\u201c",
		"&rdquo;", "\u201d",
	)
	return replacer.Replace(s)
}

// collapseWhitespace replaces sequences of whitespace with a single space
func collapseWhitespace(s string) string {
	var result strings.Builder
	result.Grow(len(s))
	inSpace := false
	for _, r := range s {
		if unicode.IsSpace(r) {
			if !inSpace && result.Len() > 0 {
				result.WriteRune(' ')
				inSpace = true
			}
		} else {
			result.WriteRune(r)
			inSpace = false
		}
	}
	return result.String()
}

// cleanInvisibleChars removes control characters, zero-width chars, and other invisible garbage
func cleanInvisibleChars(s string) string {
	var result strings.Builder
	for _, r := range s {
		// Remove control chars (C0 and C1 except tab, newline, CR)
		if r < 32 && r != '\t' && r != '\n' && r != '\r' {
			continue
		}
		// Remove zero-width characters
		if r == '\u200b' || r == '\u200c' || r == '\u200d' || r == '\ufeff' || r == '\u00ad' {
			continue
		}
		// Remove replacement character that indicates encoding errors
		if r == '\ufffd' {
			continue
		}
		result.WriteRune(r)
	}
	return result.String()
}

// cleanRepeatedPunctuation collapses 3+ consecutive identical punctuation marks
// into a single instance. Only collapses same-char runs, not mixed punctuation.
func cleanRepeatedPunctuation(s string) string {
	var result strings.Builder
	result.Grow(len(s))
	runes := []rune(s)
	i := 0
	for i < len(runes) {
		r := runes[i]
		if isPunct(r) {
			count := 1
			for i+count < len(runes) && runes[i+count] == r {
				count++
			}
			if count >= 3 {
				result.WriteRune(r)
			} else {
				for j := 0; j < count; j++ {
					result.WriteRune(r)
				}
			}
			i += count
		} else {
			result.WriteRune(r)
			i++
		}
	}
	return collapseWhitespace(result.String())
}

// isPunct reports whether r is a punctuation character that should be
// considered for repetition collapsing.
func isPunct(r rune) bool {
	return r == '.' || r == ',' || r == ';' || r == ':' || r == '!' || r == '?' ||
		r == '。' || r == '，' || r == '、' || r == '；' || r == '：' ||
		r == '！' || r == '？' || r == '…' || r == '—' || r == '–'
}

// isMostlySymbols returns true if more than 40% of characters are symbols/punctuation
func isMostlySymbols(s string) int {
	var symbols, total int
	for _, r := range s {
		if unicode.IsPunct(r) || unicode.IsSymbol(r) {
			symbols++
		}
		total++
	}
	if total == 0 {
		return 0
	}
	return symbols * 100 / total
}

// splitParagraphs splits text into paragraphs using blank lines
func splitParagraphs(s string) []string {
	// Split by double newline or single newline followed by significant content
	parts := strings.Split(s, "\n\n")
	var paragraphs []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if len(p) > 0 {
			paragraphs = append(paragraphs, p)
		}
	}
	// Also try single newline split for inline paragraph breaks
	if len(paragraphs) == 1 && strings.Count(s, "\n") > 2 {
		lines := strings.Split(s, "\n")
		var current strings.Builder
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				if current.Len() > 0 {
					paragraphs = append(paragraphs, current.String())
					current.Reset()
				}
			} else {
				if current.Len() > 0 {
					current.WriteString(" ")
				}
				current.WriteString(line)
			}
		}
		if current.Len() > 0 {
			paragraphs = append(paragraphs, current.String())
		}
	}
	return paragraphs
}

// truncateAtBoundary truncates text to be between minRunes and maxRunes,
// preferring to break at sentence/phrase boundaries
func truncateAtBoundary(s string, minRunes, maxRunes int) string {
	runes := []rune(s)
	if len(runes) <= maxRunes {
		return s
	}

	// Candidate truncation points: sentence ends (。！？!?) followed by space or end,
	// or comma/colon/space, or just space
	// We scan from maxRunes backwards to find a good boundary
	bestCut := maxRunes

	// First, look for sentence-ending punctuation within the range [minRunes, maxRunes]
	for i := maxRunes - 1; i >= minRunes; i-- {
		r := runes[i]
		// Sentence-ending punctuation
		if r == '。' || r == '！' || r == '？' || r == '.' || r == '!' || r == '?' {
			// Include the punctuation in the result
			bestCut = i + 1
			goto done
		}
	}

	// Then look for phrase boundary punctuation
	for i := maxRunes - 1; i >= minRunes; i-- {
		r := runes[i]
		if r == '，' || r == ',' || r == '、' || r == '；' || r == ':' || r == '：' || r == '—' || r == '–' || r == '-' {
			bestCut = i
			goto done
		}
	}

	// Fall back to word boundary (space)
	for i := maxRunes - 1; i >= minRunes; i-- {
		if runes[i] == ' ' {
			bestCut = i
			goto done
		}
	}

	// Last resort: hard cut at maxRunes
	bestCut = maxRunes

done:
	result := string(runes[:bestCut])
	return strings.TrimRight(result, " \t\n\r")
}

// cleanTrailingChars removes trailing invisible garbage and control characters,
// preserving normal sentence-ending punctuation.
func cleanTrailingChars(s string) string {
	s = strings.TrimRight(s, " \t\n\r")
	var result strings.Builder
	for _, r := range s {
		if r == '\u200b' || r == '\u200c' || r == '\u200d' || r == '\ufeff' || r == '\u00ad' || r == '\ufffd' {
			continue
		}
		if r < 32 && r != '\t' && r != '\n' && r != '\r' {
			continue
		}
		result.WriteRune(r)
	}
	return strings.TrimRight(result.String(), " \t")
}

func (s *Service) processContent(html string) (string, []CodeBlock) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return html, nil
	}

	doc.Find("img").Each(func(i int, sel *goquery.Selection) {
		src, _ := sel.Attr("src")
		if attachmentId := extractAttachmentId(src); attachmentId != "" {
			sel.SetAttr("src", "/api/assets/"+attachmentId)
			return
		}
		if !s.imageProxyEnabled {
			return
		}
		if !strings.HasPrefix(src, "http://") && !strings.HasPrefix(src, "https://") {
			return
		}
		if isInternalAssetPath(src, s.domain) {
			return
		}
		if s.imageProxyBaseUrl != "" {
			if strings.HasPrefix(src, s.imageProxyBaseUrl) {
				return
			}
			sel.SetAttr("src", s.imageProxyBaseUrl+"?url="+url.QueryEscape(src))
			return
		}
		sel.SetAttr("src", "/api/imageproxy?url="+url.QueryEscape(src))
	})

	var codeBlocks []CodeBlock

	doc.Find("pre code").Each(func(i int, sel *goquery.Selection) {
		codeText := sel.Text()
		className, _ := sel.Attr("class")
		languageID, detectedBy := detectCodeBlockLanguage(codeText, className)
		setCodeLanguageClass(sel, languageID)

		codeBlocks = append(codeBlocks, CodeBlock{
			Index:           i,
			LanguageID:      languageID,
			LanguageLabel:   friendlyLanguageLabel(languageID),
			DetectedBy:      detectedBy,
			ShowLineNumbers: true,
		})
	})

	result, _ := doc.Find("body").Html()
	if result == "" {
		return html, codeBlocks
	}
	return strings.TrimSpace(result), codeBlocks
}

func detectCodeBlockLanguage(codeText, className string) (string, string) {
	fallbackLanguageID, fallbackDetectedBy, hasFallback := normalizeCodeLanguageClass(className)
	if hasFallback {
		if resolvedLanguageID := resolveChromaLanguageID(fallbackLanguageID); resolvedLanguageID != "" {
			return resolvedLanguageID, fallbackDetectedBy
		}
	}

	if strings.TrimSpace(stripControlRunes(codeText)) == "" {
		if hasFallback {
			return fallbackLanguageID, fallbackDetectedBy
		}
		return "plaintext", "fallback"
	}

	lexer := lexers.Analyse(codeText)
	if lexer == nil {
		if hasFallback {
			return fallbackLanguageID, fallbackDetectedBy
		}
		return "plaintext", "fallback"
	}

	languageID := normalizeLanguageID(lexer.Config().Name)
	if languageID == "" && len(lexer.Config().Aliases) > 0 {
		languageID = normalizeLanguageID(lexer.Config().Aliases[0])
	}
	if languageID == "" {
		if hasFallback {
			return fallbackLanguageID, fallbackDetectedBy
		}
		return "plaintext", "fallback"
	}

	return languageID, "analyse"
}

func resolveChromaLanguageID(languageID string) string {
	lexer := lexers.Get(languageID)
	if lexer == nil {
		return ""
	}

	resolvedLanguageID := normalizeLanguageID(lexer.Config().Name)
	if resolvedLanguageID == "" && len(lexer.Config().Aliases) > 0 {
		resolvedLanguageID = normalizeLanguageID(lexer.Config().Aliases[0])
	}

	return resolvedLanguageID
}

func normalizeCodeLanguageClass(className string) (string, string, bool) {
	if className == "" {
		return "", "", false
	}

	aliasMap := map[string]string{
		"language-application-javascript-env-frontend": "javascript",
		"language-application-javascript-env-backend":  "javascript",
		"language-text-x-sh":                           "bash",
	}

	for _, classToken := range strings.Fields(className) {
		if languageID, ok := aliasMap[classToken]; ok {
			return languageID, "alias", true
		}
		if strings.HasPrefix(classToken, "language-") {
			languageID := normalizeLanguageID(strings.TrimPrefix(classToken, "language-"))
			if languageID != "" {
				return languageID, "class", true
			}
		}
	}

	return "", "", false
}

func normalizeLanguageID(value string) string {
	normalized := strings.ToLower(strings.TrimSpace(value))
	normalized = strings.ReplaceAll(normalized, "_", "-")
	normalized = strings.ReplaceAll(normalized, " ", "-")
	for _, prefix := range []string{
		"text-x-",
		"application-x-",
		"application-",
		"text-",
		"source-",
		"source.",
	} {
		if strings.HasPrefix(normalized, prefix) {
			normalized = strings.TrimPrefix(normalized, prefix)
			break
		}
	}

	switch normalized {
	case "", "text", "plain-text", "plain", "fallback":
		return "plaintext"
	case "shell", "sh", "zsh", "fish", "console":
		return "bash"
	case "js":
		return "javascript"
	case "ts":
		return "typescript"
	case "golang":
		return "go"
	}

	return normalized
}

func friendlyLanguageLabel(languageID string) string {
	switch languageID {
	case "javascript":
		return "JavaScript"
	case "typescript":
		return "TypeScript"
	case "bash":
		return "Shell"
	case "go":
		return "Go"
	case "json":
		return "JSON"
	case "html":
		return "HTML"
	case "css":
		return "CSS"
	case "yaml":
		return "YAML"
	case "markdown":
		return "Markdown"
	case "plaintext", "":
		return "Code"
	default:
		return humanizeLanguageLabel(languageID)
	}
}

func humanizeLanguageLabel(languageID string) string {
	parts := strings.FieldsFunc(languageID, func(r rune) bool {
		return r == '-' || r == '_' || r == '/' || r == '.'
	})
	if len(parts) == 0 {
		return "Code"
	}

	for i, part := range parts {
		switch part {
		case "csharp":
			parts[i] = "C#"
		case "cpp":
			parts[i] = "C++"
		case "jsx":
			parts[i] = "JSX"
		case "tsx":
			parts[i] = "TSX"
		case "sql":
			parts[i] = "SQL"
		default:
			runes := []rune(part)
			if len(runes) == 0 {
				continue
			}
			parts[i] = strings.ToUpper(string(runes[0])) + strings.ToLower(string(runes[1:]))
		}
	}

	return strings.Join(parts, " ")
}

func setCodeLanguageClass(sel *goquery.Selection, languageID string) {
	className, _ := sel.Attr("class")
	var kept []string
	for _, classToken := range strings.Fields(className) {
		if !strings.HasPrefix(classToken, "language-") {
			kept = append(kept, classToken)
		}
	}

	languageClass := "language-" + languageID
	kept = append(kept, languageClass)
	sel.SetAttr("class", strings.Join(kept, " "))
}

func stripControlRunes(value string) string {
	return strings.Map(func(r rune) rune {
		if r == '\n' || r == '\r' || r == '\t' {
			return r
		}
		if unicode.IsControl(r) {
			return -1
		}
		return r
	}, value)
}

func isInternalAssetPath(src, domain string) bool {
	if strings.HasPrefix(src, "/api/assets/") || strings.HasPrefix(src, "/assets/") || strings.HasPrefix(src, "/api/imageproxy") {
		return true
	}
	if strings.HasPrefix(src, "http://") || strings.HasPrefix(src, "https://") {
		parsed, err := url.Parse(src)
		if err != nil {
			return false
		}
		domainParsed, err := url.Parse(domain)
		if err == nil && parsed.Host == domainParsed.Host {
			path := parsed.Path
			if strings.HasPrefix(path, "/api/assets/") || strings.HasPrefix(path, "/assets/") || strings.HasPrefix(path, "/api/imageproxy") {
				return true
			}
		}
	}
	return false
}

func generateID(text string) string {
	re := regexp.MustCompile(`[^\p{L}\p{N}]+`)
	id := re.ReplaceAllString(text, "-")
	id = strings.Trim(id, "-")
	return id
}

var attachmentPathRe = regexp.MustCompile(`^(?:/?(?:api/)?)?attachments/([^/]+)`)

func extractAttachmentId(src string) string {
	m := attachmentPathRe.FindStringSubmatch(src)
	if m == nil {
		return ""
	}
	return m[1]
}
