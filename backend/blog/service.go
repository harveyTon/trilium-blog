package blog

import (
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/harveyTon/trilium-blog/backend/etapi"
	"github.com/microcosm-cc/bluemonday"
)

type Store interface {
	Get(key string) (string, error)
	Set(key string, value string, ttlSeconds int) error
}

type Service struct {
	etapiClient *etapi.Client
	store       Store
	blogName    string
	blogTitle   string
	domain      string
	pageSize    int
	imageProxy  string
}

type ServiceOption func(*Service)

func WithBlogName(name string) ServiceOption {
	return func(s *Service) { s.blogName = name }
}

func WithBlogTitle(title string) ServiceOption {
	return func(s *Service) { s.blogTitle = title }
}

func WithDomain(domain string) ServiceOption {
	return func(s *Service) { s.domain = domain }
}

func WithPageSize(size int) ServiceOption {
	return func(s *Service) { s.pageSize = size }
}

func WithImageProxy(proxy string) ServiceOption {
	return func(s *Service) { s.imageProxy = proxy }
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

func (s *Service) GetSite() Site {
	return Site{
		Name:   s.blogName,
		Title:  s.blogTitle,
		Domain: s.domain,
		Comments: CommentsConfig{
			Enabled: false,
		},
		ImageProxy: ImageProxyConfig{
			Enabled: s.imageProxy != "",
			BaseURL: s.imageProxy,
		},
	}
}

func (s *Service) ListPosts(page int) (*PostList, error) {
	notes, err := s.etapiClient.GetNotes("#blog=true")
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

	return &PostList{
		Items:      posts[start:end],
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

func (s *Service) GetPost(noteId string) (*Post, error) {
	note, err := s.etapiClient.GetNote(noteId)
	if err != nil {
		return nil, err
	}

	if !hasBlogLabel(note.Attributes) {
		return nil, ErrNotBlogPost
	}

	content, err := s.etapiClient.GetNoteContent(noteId)
	if err != nil {
		return nil, err
	}

	sanitized := s.sanitizeContent(content)
	toc := s.extractTOC(sanitized)
	processed := s.processContent(sanitized)

	return &Post{
		NoteID:       note.NoteID,
		Title:        note.Title,
		DateModified: note.DateModified,
		ContentHTML:  processed,
		TOC:          toc,
		PageURL:      getPageURL(note.Attributes),
	}, nil
}

func (s *Service) GetAsset(attachmentId string) ([]byte, string, error) {
	return s.etapiClient.GetAttachmentContent(attachmentId)
}

var ErrNotBlogPost = &BlogError{Message: "note is not a blog post"}

type BlogError struct {
	Message string
}

func (e *BlogError) Error() string {
	return e.Message
}

func hasBlogLabel(attrs []etapi.Attribute) bool {
	for _, a := range attrs {
		if a.Type == "label" && a.Name == "blog" && a.Value == "true" {
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

func (s *Service) sanitizeContent(html string) string {
	p := bluemonday.UGCPolicy()
	p.AllowAttrs("class").OnElements("code")
	return string(p.SanitizeBytes([]byte(html)))
}

func (s *Service) extractTOC(html string) []TOCItem {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil
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
		return len(toc) < 20
	})

	return toc
}

func (s *Service) processContent(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return html
	}

	doc.Find("img").Each(func(i int, sel *goquery.Selection) {
		src, _ := sel.Attr("src")
		if strings.HasPrefix(src, "/attachments/") {
			attachmentId := strings.TrimPrefix(src, "/attachments/")
			sel.SetAttr("src", "/api/assets/"+attachmentId)
		}
	})

	languageMap := map[string]string{
		"language-application-javascript-env-frontend": "language-javascript",
		"language-application-javascript-env-backend":  "language-javascript",
		"language-text-x-sh":                           "language-bash",
	}

	doc.Find("pre code").Each(func(i int, sel *goquery.Selection) {
		for oldClass, newClass := range languageMap {
			if sel.HasClass(oldClass) {
				sel.RemoveClass(oldClass)
				sel.AddClass(newClass)
				break
			}
		}
		if !sel.HasClass("language-") {
			sel.AddClass("language-text")
		}
	})

	result, _ := doc.Find("body").Html()
	if result == "" {
		return html
	}
	return strings.TrimSpace(result)
}

func generateID(text string) string {
	re := regexp.MustCompile(`[^\p{L}\p{N}]+`)
	id := re.ReplaceAllString(text, "-")
	id = strings.Trim(id, "-")
	return id
}
