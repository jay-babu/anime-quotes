//go:build wireinject
// +build wireinject

package web

import (
	"bytes"
	"io"
	"os"
	"sync"
	"time"

	"github.com/gin-contrib/requestid"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/jayp0521/anime-quotes/utils"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	ginEngine *gin.Engine
	ginInit   sync.Once
)

var superset = wire.NewSet(
	utils.SuperSet,
)

var SuperSet = wire.NewSet(
	ProvideRouter,
	InjectMain,
)

func injectRouter(log *zap.SugaredLogger) *gin.Engine {
	ginInit.Do(func() {
		if os.Getenv("ENV") == "LOCAL" {
		} else {
			gin.SetMode(gin.ReleaseMode)
		}
		ginEngine = gin.New()
		config := &ginzap.Config{
			TraceID:    true,
			UTC:        true,
			TimeFormat: time.RFC3339,
			Context: ginzap.Fn(func(c *gin.Context) []zapcore.Field {
				fields := []zapcore.Field{}
				// log request ID
				if requestID := c.Writer.Header().Get("X-Request-Id"); requestID != "" {
					fields = append(fields, zap.String("request_id", requestID))
				}

				// log trace and span ID
				if trace.SpanFromContext(c.Request.Context()).SpanContext().IsValid() {
					fields = append(fields, zap.String("trace_id", trace.SpanFromContext(c.Request.Context()).SpanContext().TraceID().String()))
					fields = append(fields, zap.String("span_id", trace.SpanFromContext(c.Request.Context()).SpanContext().SpanID().String()))
				}

				// log request body
				var body []byte
				var buf bytes.Buffer
				tee := io.TeeReader(c.Request.Body, &buf)
				body, _ = io.ReadAll(tee)
				c.Request.Body = io.NopCloser(&buf)
				fields = append(fields, zap.String("body", string(body)))

				return fields
			}),
		}
		ginEngine.Use(requestid.New())
		ginEngine.Use(ginzap.GinzapWithConfig(log.Desugar(), config))
		ginEngine.Use(ginzap.RecoveryWithZap(log.Desugar(), true))
	})
	return ginEngine
}

func ProvideRouter() *gin.Engine {
	panic(wire.Build(
		superset,
		injectRouter,
	))
}

func InjectMain() Main {
	panic(wire.Build(
		superset,
		ProvideRouter,
		wire.Struct(new(Main), "*"),
	))
}
