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

type AppConfig struct {
	TriliumApiUrl   string
	TriliumToken    string
	ArticlesPerPage int
	BlogName        string
	BlogTitle       string
	Domain          string
	ImageProxy      ImageProxyConfig
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
