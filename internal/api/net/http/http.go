package http

import (
	"jaeger/internal/api/compute"
	"jaeger/internal/api/net/http/httpRouter"
	"log/slog"
	"net/http"
	"time"
)

func NewHttpServer(addr string, timeout time.Duration, comp compute.Compute, lg *slog.Logger) *http.Server {
	return &http.Server{
		Addr:         addr,
		Handler:      httpRouter.NewGinRouter(comp, lg),
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
		IdleTimeout:  timeout,
	}
}
