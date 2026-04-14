package blog

import "time"

type Post struct {
	NoteID       string      `json:"noteId"`
	Title        string      `json:"title"`
	DateModified string      `json:"dateModified"`
	Summary      string      `json:"summary,omitempty"`
	Summaries    *Summaries  `json:"summaries,omitempty"`
	TOC          []TOCItem   `json:"toc,omitempty"`
	CodeBlocks   []CodeBlock `json:"codeBlocks,omitempty"`
	ContentHTML  string      `json:"contentHtml,omitempty"`
	PageURL      string      `json:"pageUrl,omitempty"`
}

type CodeBlock struct {
	Index           int    `json:"index"`
	LanguageID      string `json:"languageId"`
	LanguageLabel   string `json:"languageLabel"`
	DetectedBy      string `json:"detectedBy,omitempty"`
	ShowLineNumbers bool   `json:"showLineNumbers"`
}

type SummaryEntry struct {
	Type      string `json:"type,omitempty"`
	Status    string `json:"status,omitempty"`
	Text      string `json:"text,omitempty"`
	UpdatedAt string `json:"updatedAt,omitempty"`
	Error     string `json:"error,omitempty"`
}

type Summaries struct {
	NoteID    string        `json:"noteId,omitempty"`
	AIEnabled bool          `json:"aiEnabled"`
	AI        *SummaryEntry `json:"ai,omitempty"`
	Code      *SummaryEntry `json:"code,omitempty"`
}

type TOCItem struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Level int    `json:"level"`
}

type PostList struct {
	Items      []Post `json:"items"`
	Page       int    `json:"page"`
	PageSize   int    `json:"pageSize"`
	Total      int    `json:"total"`
	TotalPages int    `json:"totalPages"`
}

type SearchResponse struct {
	Query string       `json:"query"`
	Total int          `json:"total"`
	Items []SearchItem `json:"items"`
}

type SearchItem struct {
	Post
	Match SearchMatch `json:"match"`
}

type SearchMatch struct {
	TitleMatched bool   `json:"titleMatched"`
	Snippet      string `json:"snippet"`
}

type Site struct {
	Title      string           `json:"title"`
	Subtitle   string           `json:"subtitle"`
	Domain     string           `json:"domain"`
	Comments   CommentsConfig   `json:"comments"`
	ImageProxy ImageProxyConfig `json:"imageProxy"`
}

type CommentsConfig struct {
	Enabled bool   `json:"enabled"`
	Server  string `json:"server"`
	Site    string `json:"site"`
}

type ImageProxyConfig struct {
	Enabled bool   `json:"enabled"`
	BaseURL string `json:"baseUrl"`
}

type etapiNote struct {
	NoteID       string      `json:"noteId"`
	Title        string      `json:"title"`
	DateModified string      `json:"dateModified"`
	Type         string      `json:"type"`
	Mime         string      `json:"mime"`
	Attributes   []attribute `json:"attributes"`
}

type attribute struct {
	Type  string `json:"type"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (n *etapiNote) hasBlogLabel() bool {
	return hasTrueLabelFromAttributes(n.Attributes, "blog")
}

func (n *etapiNote) hasFeaturedLabel() bool {
	return hasTrueLabelFromAttributes(n.Attributes, "blogtop")
}

func (n *etapiNote) getPageURL() string {
	for _, a := range n.Attributes {
		if a.Type == "label" && a.Name == "pageUrl" {
			return a.Value
		}
	}
	return ""
}

type NoteContent struct {
	Note    *etapiNote `json:"note"`
	Content string     `json:"content"`
}

func parseDate(s string) time.Time {
	t, _ := time.Parse(time.RFC3339, s)
	return t
}
