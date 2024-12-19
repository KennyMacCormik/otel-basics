package httpHandlers

import (
	"github.com/gin-gonic/gin"
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
		lg := logger(c)
		b := &body{}
		err := c.ShouldBindJSON(&b)
		if err != nil {
			lg.Error("failed read body", "error", err)
			c.JSON(http.StatusBadRequest, errorMsg{Err: err.Error()})
			return
		}
		lg.Debug("request body", "body", b)

		val, _, err := comp.Get(b.Key, lg)
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
		b := &body{}
		err := c.ShouldBindJSON(&b)
		if err != nil {
			lg.Error("failed read body", "error", err)
			c.JSON(http.StatusBadRequest, errorMsg{Err: err.Error()})
			return
		}
		lg.Debug("request body", "body", b)

		err = comp.Set(b.Key, b.Val, lg)
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
		b := &body{}
		err := c.ShouldBindJSON(&b)
		if err != nil {
			c.JSON(http.StatusBadRequest, errorMsg{Err: err.Error()})
			return
		}
		lg.Debug("request body", "body", b)

		err = comp.Del(b.Key, lg)
		if err != nil {
			lg.Error("db connection error", "error", err)
			c.JSON(http.StatusInternalServerError, errorMsg{Err: err.Error()})
			return
		}

		c.Status(http.StatusOK)
		lg.Info("success")
	}
}
