package config

import (
	"encoding/json"
	"fmt"
	"os"

	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v2"
)

// FileConfig 配置本地文件日志
type FileConfig struct {
	Enable    bool   `json:"enable" yaml:"enable"`
	FilePath  string `json:"filePath" yaml:"filePath"`
	ErrorPath string `json:"errorPath" yaml:"errorPath"`
}

// LokiConfig 配置 Loki 日志
type LokiConfig struct {
	Enable    bool   `json:"enable" yaml:"enable"`
	URL       string `json:"url" yaml:"url"`
	JobName   string `json:"jobName" yaml:"jobName"`
	BatchWait int    `json:"batchWait" yaml:"batchWait"`
	BatchSize int    `json:"batchSize" yaml:"batchSize"`
}

// LogConfig 统一日志配置
type LogConfig struct {
	Level   string      `json:"level" yaml:"level"`     // 日志级别（字符串格式）
	Console bool        `json:"console" yaml:"console"` // 是否输出到控制台
	File    *FileConfig `json:"file" yaml:"file"`       // 文件日志配置
	Loki    *LokiConfig `json:"loki" yaml:"loki"`       // Loki 配置
}

// ParseLogLevel 解析字符串到 zapcore.Level
func (c *LogConfig) ParseLogLevel() zapcore.Level {
	switch c.Level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

// LoadConfigFromFile 读取 JSON 或 YAML 配置
func LoadConfigFromFile(filename string) (*LogConfig, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("读取日志配置文件失败: %w", err)
	}

	config := &LogConfig{}
	if json.Valid(data) {
		// 解析 JSON
		if err := json.Unmarshal(data, config); err != nil {
			return nil, fmt.Errorf("解析 JSON 失败: %w", err)
		}
	} else {
		// 解析 YAML
		if err := yaml.Unmarshal(data, config); err != nil {
			return nil, fmt.Errorf("解析 YAML 失败: %w", err)
		}
	}

	return config, nil
}
