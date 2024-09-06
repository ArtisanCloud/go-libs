package zap

import (
	"context"
	"github.com/ArtisanCloud/PowerLibs/v3/object"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
	"net/http"
	"testing"
	"time"
)

var strArtisanCloudPath = "./"
var strOutputPath = strArtisanCloudPath + "/output.log"
var strErrorPath = strArtisanCloudPath + "/errors.log"

func init() {
	initTracer()
}

func initTracer() {
	tp := trace.NewTracerProvider()
	// Set Global Tracer Provider
	otel.SetTracerProvider(tp)
}

func Test_Log_Info(t *testing.T) {

	logger, err := NewLogger(&object.HashMap{
		"env":        "test",
		"outputPath": strOutputPath,
		"errorPath":  strErrorPath,
	})
	if err != nil {
		t.Error(err)
	}

	logger.Info("test info", "response", &http.Response{})

	tracer := otel.Tracer("example-tracer")
	ctx, span := tracer.Start(context.Background(), "test")
	defer span.End()

	logger = logger.WithContext(ctx)

	logger.Info("test info with context")
	logger.InfoF("current time %s", time.Now().Format("2006-01-02 15:04:05"))
	logger.Error("test info with context", "1234567", "abcdefg")
	logger.ErrorF("test info with context %s", "1234567")

}
