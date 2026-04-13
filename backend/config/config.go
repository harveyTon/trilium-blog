package config

import (
	"os"
	"strconv"

	"github.com/harveyTon/trilium-blog/backend/pkg/logger"
)

type ImageProxyConfig struct {
	Enabled bool
	BaseURL string
}

type AISummaryConfig struct {
	Enabled       bool
	BaseURL       string
	APIKey        string
	Model         string
	Prompt        string
	Mode          string
	Concurrency   int
	RateLimitMs   int
	DatabasePath  string
}

type AppConfig struct {
	TriliumApiUrl   string
	TriliumToken    string
	ArticlesPerPage int
	BlogName        string
	BlogTitle       string
	Domain          string
	ImageProxy      ImageProxyConfig
	AISummary       AISummaryConfig
}

var Config AppConfig

func LoadConfig() {
	Config = AppConfig{
		TriliumApiUrl:   getEnv("TRILIUM_API_URL", ""),
		TriliumToken:    getEnv("TRILIUM_TOKEN", ""),
		ArticlesPerPage: getEnvInt("ARTICLES_PER_PAGE", 9),
		BlogName:        getEnv("BLOG_NAME", ""),
		BlogTitle:       getEnv("BLOG_TITLE", ""),
		Domain:          getEnv("DOMAIN", ""),
		ImageProxy: ImageProxyConfig{
			Enabled: getEnvBool("IMAGE_PROXY_ENABLED", false),
			BaseURL: getEnv("IMAGE_PROXY_BASE_URL", ""),
		},
		AISummary: AISummaryConfig{
			Enabled:      getEnvBool("AI_SUMMARY_ENABLED", false),
			BaseURL:      getEnv("AI_SUMMARY_BASE_URL", ""),
			APIKey:       getEnv("AI_SUMMARY_API_KEY", ""),
			Model:        getEnv("AI_SUMMARY_MODEL", ""),
			Prompt:       getEnv("AI_SUMMARY_PROMPT", "Summarize the article for a blog reader in concise Chinese."),
			Mode:         getEnv("AI_SUMMARY_MODE", "code"),
			Concurrency:  getEnvInt("AI_SUMMARY_CONCURRENCY", 2),
			RateLimitMs:  getEnvInt("AI_SUMMARY_RATE_LIMIT_MS", 1200),
			DatabasePath: getEnv("AI_SUMMARY_DATABASE_PATH", "./data/summaries.db"),
		},
	}

	if Config.TriliumApiUrl == "" {
		logger.Fatal("TRILIUM_API_URL is required", nil)
	}
	if Config.TriliumToken == "" {
		logger.Fatal("TRILIUM_TOKEN is required", nil)
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return n
}

func getEnvBool(key string, fallback bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return fallback
	}
	return b
}
