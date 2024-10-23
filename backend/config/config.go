package config

import (
	"encoding/json"

	"os"
	"path/filepath"

	"github.com/harveyTon/trilium-blog/backend/pkg/logger"
)

// AppConfig 存储应用程序的所有配置
type AppConfig struct {
	TriliumApiUrl   string `json:"triliumApiUrl"`
	TriliumToken    string `json:"triliumToken"`
	ArticlesPerPage int    `json:"articlesPerPage"`
	BlogName        string `json:"blogName"`
	BlogTitle       string `json:"blogTitle"`
	Domain          string `json:"domain"`
}

// Config 是全局配置变量
var Config AppConfig

// LoadConfig 从 JSON 文件加载配置
func LoadConfig() {
	// 获取当前工作目录
	wd, err := os.Getwd()
	if err != nil {
		logger.Fatal("Unable to get working directory:", err)
	}

	// 构建配置文件路径
	configPath := filepath.Join(wd, "config.json")

	// 读取配置文件
	file, err := os.ReadFile(configPath)
	if err != nil {
		logger.Fatal("Unable to read configuration file:", err)
	}

	// 解析 JSON 到 Config 结构体
	err = json.Unmarshal(file, &Config)
	if err != nil {
		logger.Fatal("Unable to parse configuration file:", err)
	}
}
