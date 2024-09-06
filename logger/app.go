package logger

import (
	"context"
	"errors"
	"github.com/ArtisanCloud/PowerLibs/v3/logger/contract"
	"github.com/ArtisanCloud/PowerLibs/v3/logger/drivers/zap"
	"github.com/ArtisanCloud/PowerLibs/v3/object"
	"os"
)

type Logger struct {
	Driver contract.LoggerInterface
	ctx    context.Context
}

func NewLogger(driver interface{}, config *object.HashMap) (logger *Logger, err error) {

	var driverLogger contract.LoggerInterface
	if driver != nil {
		d, ok := driver.(contract.LoggerInterface)
		if !ok {
			return nil, errors.New("driver is not of type contract.LoggerInterface")
		}
		driverLogger = d
	} else {
		driverLogger, err = zap.NewLogger(config)
	}

	logger = &Logger{
		Driver: driverLogger,
	}

	return logger, err

}

func (log *Logger) WithContext(ctx context.Context) *Logger {
	log.ctx = ctx
	return log
}

func (log *Logger) Debug(msg string, v ...interface{}) {
	log.Driver.WithContext(log.ctx).WithContext(log.ctx).Debug(msg, v...)
}
func (log *Logger) Info(msg string, v ...interface{}) {
	log.Driver.WithContext(log.ctx).Info(msg, v...)
}
func (log *Logger) Warn(msg string, v ...interface{}) {
	log.Driver.WithContext(log.ctx).Warn(msg, v...)
}
func (log *Logger) Error(msg string, v ...interface{}) {
	log.Driver.WithContext(log.ctx).Error(msg, v...)
}
func (log *Logger) Panic(msg string, v ...interface{}) {
	log.Driver.WithContext(log.ctx).Panic(msg, v...)
}
func (log *Logger) Fatal(msg string, v ...interface{}) {
	log.Driver.WithContext(log.ctx).Fatal(msg, v...)
}

func (log *Logger) DebugF(format string, args ...interface{}) {
	log.Driver.WithContext(log.ctx).DebugF(format, args)
}
func (log *Logger) InfoF(format string, args ...interface{}) {
	log.Driver.WithContext(log.ctx).InfoF(format, args)
}
func (log *Logger) WarnF(format string, args ...interface{}) {
	log.Driver.WithContext(log.ctx).WarnF(format, args)
}
func (log *Logger) ErrorF(format string, args ...interface{}) {
	log.Driver.WithContext(log.ctx).ErrorF(format, args)
}
func (log *Logger) PanicF(format string, args ...interface{}) {
	log.Driver.WithContext(log.ctx).PanicF(format, args)
}
func (log *Logger) FatalF(format string, args ...interface{}) {
	log.Driver.WithContext(log.ctx).FatalF(format, args)
}

func InitLogPath(path string, files ...string) (err error) {
	if _, err = os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
	} else if os.IsPermission(err) {
		return err
	}

	for _, fileName := range files {
		if _, err = os.Stat(fileName); os.IsNotExist(err) {
			_, err = os.Create(fileName)
			if err != nil {
				return err
			}
		}
	}

	return err

}
