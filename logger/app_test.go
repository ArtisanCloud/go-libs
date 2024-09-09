package logger

import (
	"context"
	"github.com/ArtisanCloud/PowerLibs/v3/object"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
	"net/http"
	"os"
	"testing"
	"time"
)

var strArtisanCloudPath = "./"
var strOutputPath = strArtisanCloudPath + "/output.log"
var strErrorPath = strArtisanCloudPath + "/errors.log"

func init() {
	err := initLogPath(strArtisanCloudPath, strOutputPath, strErrorPath)
	if err != nil {
		panic(err)
	}

	initTracer()
}

func initTracer() {
	tp := trace.NewTracerProvider()
	// Set Global Tracer Provider
	otel.SetTracerProvider(tp)
}

func Test_Log_Info(t *testing.T) {
	logger, err := NewLogger(nil, &object.HashMap{
		"env":        "test",
		"outputPath": strOutputPath,
		"errorPath":  strErrorPath,
		"stdout":     true,
	})
	if err != nil {
		t.Error(err)
	}

	// without context
	logger.Info("test info", "app response", &http.Response{})

	// log with context，will append traceId and spanId if ctx hash trace info
	tracer := otel.Tracer("example-tracer")
	ctx, span := tracer.Start(context.Background(), "test")
	defer span.End()

	logger.WithContext(ctx).Info("test info with context")
	logger.WithContext(ctx).InfoF("current time %s", time.Now().Format("2006-01-02 15:04:05"))

}

func initLogPath(path string, files ...string) (err error) {
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
