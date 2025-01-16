package httpRouter

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
	"jaeger/internal/api/compute"
	"jaeger/internal/api/net/http/httpHandlers"
	"log/slog"
)

func NewGinRouter(comp compute.Compute, lg *slog.Logger) *gin.Engine {
	router := gin.New()
	_ = router.SetTrustedProxies(nil)

	addMiddleware(router)
	addHandlers(router, comp, lg)

	return router
}

func addMiddleware(router *gin.Engine) {
	router.Use(gin.Recovery())
	router.Use(TracingMiddleware())
}

func addHandlers(router *gin.Engine, comp compute.Compute, lg *slog.Logger) {
	logger := func(c *gin.Context) *slog.Logger {
		uid := uuid.New()
		newLg := lg.With("ID", uid)
		newLg.Info("connection accepted",
			"ClientIP", c.ClientIP(),
			"Method", c.Request.Method,
			"Path", c.Request.URL.Path,
			"Proto", c.Request.Proto,
			"Headers", c.Request.Header)
		return newLg
	}
	httpHandlers.NewStorageHandlers(comp, logger)(router)
}

func TracingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tracer := otel.Tracer("")
		ctx, span := tracer.Start(c.Request.Context(), c.Request.URL.Path)
		defer span.End()

		span.SetAttributes(
			semconv.HTTPMethodKey.String(c.Request.Method),
			semconv.HTTPURLKey.String(c.Request.URL.String()),
		)

		c.Request = c.Request.WithContext(ctx)
		c.Next()

		span.SetAttributes(
			semconv.HTTPStatusCodeKey.Int(c.Writer.Status()),
		)
	}
}
