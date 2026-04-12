package config

import (
	"encoding/json"

	"os"
	"path/filepath"

	"github.com/harveyTon/trilium-blog/backend/pkg/logger"
)

type ImageProxyConfig struct {
	Enabled bool   `json:"enabled"`
	BaseURL string `json:"baseUrl"`
}

// AppConfig 存储应用程序的所有配置
type AppConfig struct {
	TriliumApiUrl   string           `json:"triliumApiUrl"`
	TriliumToken    string           `json:"triliumToken"`
	ArticlesPerPage int              `json:"articlesPerPage"`
	BlogName        string           `json:"blogName"`
	BlogTitle       string           `json:"blogTitle"`
	Domain          string           `json:"domain"`
	ImageProxy      ImageProxyConfig `json:"imageProxy"`
}

// Config 是全局配置变量
var Config AppConfig

// LoadConfig 从 JSON 文件加载配置
func LoadConfig() {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		wd, err := os.Getwd()
		if err != nil {
			logger.Fatal("Unable to get working directory:", err)
		}
		configPath = filepath.Join(wd, "config.json")
	}

	file, err := os.ReadFile(configPath)
	if err != nil {
		logger.Fatal("Unable to read configuration file:", err)
	}

	err = json.Unmarshal(file, &Config)
	if err != nil {
		logger.Fatal("Unable to parse configuration file:", err)
	}
}
