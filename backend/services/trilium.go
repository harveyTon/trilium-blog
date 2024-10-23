package services

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/harveyTon/trilium-blog/backend/config"
	"github.com/harveyTon/trilium-blog/backend/models"
	"github.com/harveyTon/trilium-blog/backend/pkg/logger"
	"github.com/microcosm-cc/bluemonday"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
)

var redisClient = redis.NewClient(&redis.Options{
	Addr:     "redis:6379",
	Password: "",
	DB:       0,
})

var ctx = context.Background()
var cacheMutex = &sync.Mutex{}
var client = &http.Client{
	Timeout: time.Second * 30,
}

func getCachedData(cacheKey string, fetcher func() ([]byte, error)) ([]byte, error) {
	cachedResult, err := redisClient.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		cacheMutex.Lock()
		defer cacheMutex.Unlock()

		cachedResult, err = redisClient.Get(ctx, cacheKey).Result()
		if err == redis.Nil {
			logger.Infof("Cache miss, fetching from source: %s", cacheKey)
			fetchedData, fetchErr := fetcher()
			if fetchErr != nil {
				return nil, fetchErr
			}
			compressedData, _ := compressData(fetchedData)
			err = redisClient.Set(ctx, cacheKey, compressedData, time.Hour*1).Err()
			if err != nil {
				logger.Infof("Failed to cache data: %v", err)
			}
			return fetchedData, nil
		} else if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	decompressedData, _ := decompressData([]byte(cachedResult))
	return decompressedData, nil
}

func GetArticles(page, pageSize int) ([]models.Article, int, error) {
	cacheKey := fmt.Sprintf("trilium_blog_articles:%d:%d", page, pageSize)

	data, err := getCachedData(cacheKey, func() ([]byte, error) {
		url := fmt.Sprintf("%s/etapi/notes?search=%%23blog%%3Dtrue&orderBy=utcDateModified", config.Config.TriliumApiUrl)
		logger.Infof("Requesting articles from Trilium API: %s", url)

		var response models.APIResponse
		if err := apiRequest(url, &response); err != nil {
			logger.Infof("Error from Trilium API: %v", err)
			return nil, err
		}

		var filteredArticles []models.Article
		for _, article := range response.Results {
			if article.Type == "text" && hasBlogAttribute(article.Attributes) {
				filteredArticles = append(filteredArticles, article)
			}
		}

		totalArticles := len(filteredArticles)
		startIndex := (page - 1) * pageSize
		endIndex := startIndex + pageSize
		if endIndex > totalArticles {
			endIndex = totalArticles
		}
		if startIndex >= totalArticles {
			return json.Marshal(struct {
				Articles []models.Article
				Total    int
			}{[]models.Article{}, totalArticles})
		}

		paginatedArticles := filteredArticles[startIndex:endIndex]
		return json.Marshal(struct {
			Articles []models.Article
			Total    int
		}{paginatedArticles, totalArticles})
	})

	if err != nil {
		return nil, 0, err
	}

	var cachedResponse struct {
		Articles []models.Article
		Total    int
	}
	json.Unmarshal(data, &cachedResponse)
	return cachedResponse.Articles, cachedResponse.Total, nil
}

func GetArticle(noteId string) (*models.Article, string, error) {
	cacheKey := fmt.Sprintf("trilium_blog_article:%s", noteId)

	data, err := getCachedData(cacheKey, func() ([]byte, error) {
		var article models.Article
		var content string
		var wg sync.WaitGroup
		wg.Add(2)

		var metaErr, contentErr error
		go func() {
			defer wg.Done()
			metaUrl := fmt.Sprintf("%s/etapi/notes/%s", config.Config.TriliumApiUrl, noteId)
			metaErr = apiRequest(metaUrl, &article)
		}()

		go func() {
			defer wg.Done()
			contentUrl := fmt.Sprintf("%s/etapi/notes/%s/content", config.Config.TriliumApiUrl, noteId)
			req, _ := http.NewRequest("GET", contentUrl, nil)
			req.Header.Add("Authorization", config.Config.TriliumToken)
			resp, err := client.Do(req)
			if err == nil {
				defer resp.Body.Close()
				body, _ := io.ReadAll(resp.Body)
				html := bluemonday.UGCPolicy().SanitizeBytes(body)
				content = string(html)
				content = regexp.MustCompile(`api/attachments/([^/]+)/image/[^"]+`).ReplaceAllString(content, "/attachments/$1")
				content = replaceImgSrc(content)
				content = removeHtmlHeadBody(content)
			} else {
				contentErr = err
			}
		}()
		wg.Wait()

		if metaErr != nil || contentErr != nil {
			return nil, fmt.Errorf("Failed to retrieve article: %v, %v", metaErr, contentErr)
		}

		if !hasBlogAttribute(article.Attributes) {
			return nil, fmt.Errorf("Article is not a blog post")
		}

		return json.Marshal(struct {
			Article models.Article
			Content string
		}{article, content})
	})

	if err != nil {
		return nil, "", err
	}

	var cachedArticle struct {
		Article models.Article
		Content string
	}
	json.Unmarshal(data, &cachedArticle)
	return &cachedArticle.Article, cachedArticle.Content, nil
}

func GetAttachment(attachmentId string) ([]byte, string, error) {
	cacheKey := fmt.Sprintf("trilium_blog_attachment:%s", attachmentId)

	data, err := getCachedData(cacheKey, func() ([]byte, error) {
		detailUrl := fmt.Sprintf("%s/etapi/attachments/%s", config.Config.TriliumApiUrl, attachmentId)
		var attachmentDetail struct {
			OwnerId string `json:"ownerId"`
			Mime    string `json:"mime"`
		}
		if err := apiRequest(detailUrl, &attachmentDetail); err != nil {
			return nil, fmt.Errorf("Failed to retrieve attachment details: %v", err)
		}

		noteUrl := fmt.Sprintf("%s/etapi/notes/%s", config.Config.TriliumApiUrl, attachmentDetail.OwnerId)
		var note struct {
			Attributes []models.Attribute `json:"attributes"`
		}
		if err := apiRequest(noteUrl, &note); err != nil {
			return nil, fmt.Errorf("Failed to retrieve note information: %v", err)
		}

		if !hasBlogAttribute(note.Attributes) {
			return nil, fmt.Errorf("Attachment belongs to a non-blog post")
		}

		contentUrl := fmt.Sprintf("%s/etapi/attachments/%s/content", config.Config.TriliumApiUrl, attachmentId)
		req, _ := http.NewRequest("GET", contentUrl, nil)
		req.Header.Add("Authorization", config.Config.TriliumToken)
		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("Failed to retrieve attachment content: %v", err)
		}
		defer resp.Body.Close()

		content, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("Failed to read attachment content: %v", err)
		}

		return json.Marshal(struct {
			Content []byte
			Mime    string
		}{content, attachmentDetail.Mime})
	})

	if err != nil {
		return nil, "", err
	}

	var cachedAttachment struct {
		Content []byte
		Mime    string
	}
	json.Unmarshal(data, &cachedAttachment)
	return cachedAttachment.Content, cachedAttachment.Mime, nil
}

func GenerateSitemap() (string, error) {
	cacheKey := "trilium_blog_sitemap"

	data, err := getCachedData(cacheKey, func() ([]byte, error) {
		url := fmt.Sprintf("%s/etapi/notes?search=%%23blog%%3Dtrue&orderBy=utcDateModified", config.Config.TriliumApiUrl)
		var response models.APIResponse
		if err := apiRequest(url, &response); err != nil {
			return nil, fmt.Errorf("Failed to retrieve article list: %v", err)
		}

		sitemap := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
`)
		for _, article := range response.Results {
			if article.Type == "text" && hasBlogAttribute(article.Attributes) {
				sitemap += fmt.Sprintf(`
	<url>
		<loc>%s/articles/%s</loc>
		<lastmod>%s</lastmod>
	</url>`, config.Config.TriliumApiUrl, article.NoteId, article.DateModified)
			}
		}
		sitemap += "</urlset>"

		return []byte(sitemap), nil
	})

	if err != nil {
		return "", err
	}

	return string(data), nil
}

func apiRequest(url string, target interface{}) error {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", config.Config.TriliumToken)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	return json.NewDecoder(resp.Body).Decode(target)
}

func hasBlogAttribute(attributes []models.Attribute) bool {
	for _, attr := range attributes {
		if attr.Type == "label" && attr.Name == "blog" && attr.Value == "true" {
			return true
		}
	}
	return false
}

func replaceImgSrc(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		logger.Infof("Failed to parse HTML: %v", err)
		return html
	}

	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		src, _ := s.Attr("src")
		if strings.HasPrefix(src, "http://") || strings.HasPrefix(src, "https://") {
			s.SetAttr("src", "https://88900.net/api/imageproxy/"+src)
		}
		if strings.HasPrefix(src, "/attachments/") {
			s.SetAttr("src", "https://88900.net/api/imageproxy/"+config.Config.Domain+src)
		}
	})

	result, _ := doc.Html()
	return result
}

func removeHtmlHeadBody(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		logger.Infof("Failed to parse HTML: %v", err)
		return html
	}

	bodyContent, err := doc.Find("body").Html()
	if err != nil {
		logger.Infof("Failed to get body content: %v", err)
		return html
	}

	return strings.TrimSpace(bodyContent)
}

func compressData(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	if _, err := gz.Write(data); err != nil {
		return nil, err
	}
	if err := gz.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func decompressData(data []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	return io.ReadAll(reader)
}
