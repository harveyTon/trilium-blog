package blog

import "time"

type Post struct {
	NoteID       string    `json:"noteId"`
	Title        string    `json:"title"`
	DateModified string    `json:"dateModified"`
	TOC          []TOCItem `json:"toc,omitempty"`
	ContentHTML  string    `json:"contentHtml,omitempty"`
	PageURL      string    `json:"pageUrl,omitempty"`
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

type Site struct {
	Name       string           `json:"name"`
	Title      string           `json:"title"`
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
	for _, a := range n.Attributes {
		if a.Type == "label" && a.Name == "blog" && a.Value == "true" {
			return true
		}
	}
	return false
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
