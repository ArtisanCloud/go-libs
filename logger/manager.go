package logger

import (
	"github.com/ArtisanCloud/PowerLibs/v3/logger/writer"
	"os"

	"github.com/ArtisanCloud/PowerLibs/v3/logger/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerManager struct {
	logger *zap.Logger
	config config.LogConfig
}

func NewLoggerManager(config *config.LogConfig) *LoggerManager {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	var cores []zapcore.Core

	// Console 日志
	if config.Console {
		consoleCore := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), config.ParseLogLevel())
		cores = append(cores, consoleCore)
	}

	// 文件日志
	if config.File.Enable {
		// info log path
		infoFileCore := zapcore.NewCore(encoder, writer.NewFileWriter(config.File.FilePath), config.ParseLogLevel())
		cores = append(cores, infoFileCore)

		// error log path
		errorFileCore := zapcore.NewCore(encoder, writer.NewFileWriter(config.File.ErrorPath), config.ParseLogLevel())
		cores = append(cores, errorFileCore)
	}

	// Loki 日志
	if config.Loki.Enable {
		lokiCore := zapcore.NewCore(encoder, writer.NewLokiWriter(config.Loki), config.ParseLogLevel())
		cores = append(cores, lokiCore)
	}

	// 组合多个日志 Core
	core := zapcore.NewTee(cores...)

	// 创建 Logger
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return &LoggerManager{
		logger: logger,
		config: *config,
	}
}
