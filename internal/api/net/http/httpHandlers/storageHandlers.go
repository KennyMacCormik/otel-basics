package httpHandlers

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"jaeger/internal/api/compute"
	"log/slog"
	"net/http"
)

type body struct {
	Key string `json:"key"`
	Val string `json:"val,omitempty"`
}

type errorMsg struct {
	Err string `json:"err"`
}

func NewStorageHandlers(comp compute.Compute, logger func(c *gin.Context) *slog.Logger) func(*gin.Engine) {
	return func(router *gin.Engine) {
		router.GET("/st", get(comp, logger))
		router.PUT("/st", set(comp, logger))
		router.POST("/st", set(comp, logger))
		router.DELETE("/st", del(comp, logger))
	}
}

func get(comp compute.Compute, logger func(c *gin.Context) *slog.Logger) func(c *gin.Context) {
	return func(c *gin.Context) {
		// child span
		tracer := otel.Tracer("gin-handler")
		ctx := c.Request.Context()

		spnCtx, span := tracer.Start(ctx, "get-handler")
		defer span.End()

		span.AddEvent("your message", trace.WithAttributes(attribute.String("key", "value")))
		//init attrs
		span.SetAttributes(
			semconv.ClientAddressKey.String(c.Request.RemoteAddr),
			semconv.HTTPRequestMethodKey.String(c.Request.Method),
			semconv.NetworkProtocolVersionKey.String(c.Request.Proto),
			semconv.URLSchemeKey.String(c.Request.URL.Scheme),
			semconv.URLPathKey.String(c.Request.URL.Path),
			semconv.HostNameKey.String(c.Request.Host),
		)
		// logger
		lg := logger(c)
		b := &body{}
		err := c.ShouldBindJSON(&b)
		if err != nil {
			lg.Error("failed read body", "error", err)
			c.JSON(http.StatusBadRequest, errorMsg{Err: err.Error()})
			return
		}
		lg.Debug("request body", "body", b)

		val, _, err := comp.Get(spnCtx, b.Key, lg)
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.NotFound {
			lg.Error("key does not exist", "key", b.Key)
			c.Status(http.StatusNotFound)
			return
		}
		if err != nil {
			lg.Error("db connection error", "error", err)
			c.JSON(http.StatusInternalServerError, errorMsg{Err: err.Error()})
			return
		}

		result := body{Key: b.Key, Val: val}
		lg.Debug("response body", "body", result)
		c.JSON(http.StatusOK, result)
		lg.Info("success")
	}
}

func set(comp compute.Compute, logger func(c *gin.Context) *slog.Logger) func(c *gin.Context) {
	return func(c *gin.Context) {
		lg := logger(c)

		tracer := otel.Tracer("gin-handler")
		ctx := c.Request.Context()

		spnCtx, span := tracer.Start(ctx, "get-handler")
		defer span.End()

		span.AddEvent("your message", trace.WithAttributes(attribute.String("key", "value")))
		//init attrs
		span.SetAttributes(
			semconv.ClientAddressKey.String(c.Request.RemoteAddr),
			semconv.HTTPRequestMethodKey.String(c.Request.Method),
			semconv.NetworkProtocolVersionKey.String(c.Request.Proto),
			semconv.URLSchemeKey.String(c.Request.URL.Scheme),
			semconv.URLPathKey.String(c.Request.URL.Path),
			semconv.HostNameKey.String(c.Request.Host),
		)

		b := &body{}
		err := c.ShouldBindJSON(&b)
		if err != nil {
			lg.Error("failed read body", "error", err)
			c.JSON(http.StatusBadRequest, errorMsg{Err: err.Error()})
			return
		}
		lg.Debug("request body", "body", b)

		err = comp.Set(spnCtx, b.Key, b.Val, lg)
		if err != nil {
			lg.Error("db connection error", "error", err)
			c.JSON(http.StatusInternalServerError, errorMsg{Err: err.Error()})
			return
		}

		c.Status(http.StatusOK)
		lg.Info("success")
	}
}

func del(comp compute.Compute, logger func(c *gin.Context) *slog.Logger) func(c *gin.Context) {
	return func(c *gin.Context) {
		lg := logger(c)

		tracer := otel.Tracer("gin-handler")
		ctx := c.Request.Context()

		spnCtx, span := tracer.Start(ctx, "get-handler")
		defer span.End()

		span.AddEvent("your message", trace.WithAttributes(attribute.String("key", "value")))
		//init attrs
		span.SetAttributes(
			semconv.ClientAddressKey.String(c.Request.RemoteAddr),
			semconv.HTTPRequestMethodKey.String(c.Request.Method),
			semconv.NetworkProtocolVersionKey.String(c.Request.Proto),
			semconv.URLSchemeKey.String(c.Request.URL.Scheme),
			semconv.URLPathKey.String(c.Request.URL.Path),
			semconv.HostNameKey.String(c.Request.Host),
		)

		b := &body{}
		err := c.ShouldBindJSON(&b)
		if err != nil {
			c.JSON(http.StatusBadRequest, errorMsg{Err: err.Error()})
			return
		}
		lg.Debug("request body", "body", b)

		err = comp.Del(spnCtx, b.Key, lg)
		if err != nil {
			lg.Error("db connection error", "error", err)
			c.JSON(http.StatusInternalServerError, errorMsg{Err: err.Error()})
			return
		}

		c.Status(http.StatusOK)
		lg.Info("success")
	}
}
