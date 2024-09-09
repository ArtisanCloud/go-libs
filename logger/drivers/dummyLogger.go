package drivers

import (
	"context"
	"github.com/ArtisanCloud/PowerLibs/v3/logger/contract"
)

type DummyLogger struct {
	ctx context.Context
}

func (l *DummyLogger) WithContext(ctx context.Context) contract.LoggerInterface {
	l.ctx = ctx
	return l
}

func (l *DummyLogger) Debug(msg string, v ...interface{}) {}

func (l *DummyLogger) Info(msg string, v ...interface{}) {}

func (l *DummyLogger) Warn(msg string, v ...interface{}) {}

func (l *DummyLogger) Error(msg string, v ...interface{}) {}

func (l *DummyLogger) Panic(msg string, v ...interface{}) {}

func (l *DummyLogger) Fatal(msg string, v ...interface{}) {}

func (l *DummyLogger) DebugF(format string, args ...interface{}) {}

func (l *DummyLogger) InfoF(format string, args ...interface{}) {}

func (l *DummyLogger) WarnF(format string, args ...interface{}) {}

func (l *DummyLogger) ErrorF(format string, args ...interface{}) {}

func (l *DummyLogger) PanicF(format string, args ...interface{}) {}

func (l *DummyLogger) FatalF(format string, args ...interface{}) {}
