package contract

import "context"

const (
	DebugLevel   int8 = 0
	InfoLevel    int8 = 1
	WarningLevel int8 = 2
	ErrorLevel   int8 = 3
	PanicLevel   int8 = 4
	FatalLevel   int8 = 5
)

type LoggerInterface interface {
	Debug(msg string, v ...interface{})
	Info(msg string, v ...interface{})
	Warn(msg string, v ...interface{})
	Error(msg string, v ...interface{})
	Panic(msg string, v ...interface{})
	Fatal(msg string, v ...interface{})

	DebugF(format string, args ...interface{})
	InfoF(format string, args ...interface{})
	WarnF(format string, args ...interface{})
	ErrorF(format string, args ...interface{})
	PanicF(format string, args ...interface{})
	FatalF(format string, args ...interface{})
	WithContext(ctx context.Context) LoggerInterface

	//Log( level int8, msg string, objs ...interface{})
	//Logf( format string, v ...interface{})
}
