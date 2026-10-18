package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Endpoint string `json:"endpoint"`
	APIKey   string `json:"api_key"`
	Timeout  int    `json:"timeout"`
}

func LoadForEnv(env string) (*Config, error) {
	path := ""
	switch env {
	case "dev": path = "config.dev.json"
	case "staging": path = "config.staging.json"
	case "prod": path = "config.prod.json"
	}
	if path != "" {
		if _, err := os.Stat(path); err == nil {
			return Load(path)
		}
	}
	return Load("")
}

func Load(path string) (*Config, error) {
	cfg := &Config{Timeout: 30}
	// 1. try file
	if path != "" {
		if b, err := os.ReadFile(path); err == nil {
			if err := json.Unmarshal(b, cfg); err != nil {
				return nil, fmt.Errorf("配置文件解析失败: %w", err)
			}
		}
	} else if b, err := os.ReadFile("config.json"); err == nil {
		if err := json.Unmarshal(b, cfg); err != nil {
			return nil, fmt.Errorf("配置文件解析失败: %w", err)
		}
	}
	// 2. env override
	if ep := os.Getenv("TARS_MCP_ENDPOINT"); ep != "" { cfg.Endpoint = ep }
	if ep := os.Getenv("DIFY_ENDPOINT"); ep != "" && cfg.Endpoint == "" { cfg.Endpoint = ep }
	if ak := os.Getenv("TARS_MCP_API_KEY"); ak != "" { cfg.APIKey = ak }
	if ak := os.Getenv("DIFY_API_KEY"); ak != "" && cfg.APIKey == "" { cfg.APIKey = ak }
	if v := os.Getenv("TARS_MCP_TIMEOUT"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 { cfg.Timeout = n }
	}
	if cfg.Endpoint == "" {
		return nil, fmt.Errorf("环境变量 TARS_MCP_ENDPOINT 或 DIFY_ENDPOINT 未设置，且未提供配置文件")
	}
	cfg.Endpoint = strings.TrimRight(cfg.Endpoint, "/")
	return cfg, nil
}
