package writer

import (
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap/zapcore"
)

// newFileWriter 创建文件 Writer，确保路径存在
func NewFileWriter(logFilePath string) zapcore.WriteSyncer {
	// 获取日志文件所在的目录
	logDir := filepath.Dir(logFilePath)

	// 确保日志目录存在
	if err := os.MkdirAll(logDir, 0755); err != nil {
		panic(fmt.Errorf("无法创建日志目录 %s: %w", logDir, err))
	}

	// 打开或创建日志文件
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(fmt.Errorf("无法打开日志文件 %s: %w", logFilePath, err))
	}

	return zapcore.AddSync(file)
}
