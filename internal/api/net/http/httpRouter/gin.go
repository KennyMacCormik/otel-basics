package httpRouter

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
