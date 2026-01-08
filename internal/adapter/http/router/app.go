package router

import (
	"github.com/gin-contrib/logger"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

type AppHandlers struct {
	Health gin.HandlerFunc
}

func NewAppRouter(handlers AppHandlers) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	router.Use(gin.Recovery())
	router.Use(
		logger.SetLogger(
			logger.WithLogger(func(c *gin.Context, l zerolog.Logger) zerolog.Logger {
				if trace.SpanFromContext(c.Request.Context()).SpanContext().IsValid() {
					l = l.With().
						Str("trace_id", trace.SpanFromContext(c.Request.Context()).SpanContext().TraceID().String()).
						Str("span_id", trace.SpanFromContext(c.Request.Context()).SpanContext().SpanID().String()).
						Logger()
				}
				return l.With().
					Str("request_id", requestid.Get(c)).
					Str("method", c.Request.Method).
					Str("path", c.Request.URL.Path).
					Logger()
			}),
		),
	)

	router.Any("/system/healthz", handlers.Health)

	return router
}
