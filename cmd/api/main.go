package main

import (
	"context"
	"fmt"
	"jaeger/internal/api/cache/kvCleanup"
	"jaeger/internal/api/compute"
	"jaeger/internal/api/net/grpc"
	"jaeger/internal/api/net/http"
	"jaeger/internal/cfg"
	mylog "jaeger/internal/logging"
	"jaeger/internal/tracing"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// load app config
	conf, err := cfg.NewConfig()
	if err != nil {
		log.Fatal(fmt.Errorf("error loading config: %w", err))
	}
	// init logger
	lg := mylog.NewLogger(conf.Log)
	// init tracer
	tp, err := tracing.NewTraceProvider(context.Background(), conf.Tracing, true)
	if err != nil {
		log.Fatal(fmt.Errorf("error trace init: %w", err))
	}
	defer func() {
		_ = tp.Shutdown(context.Background())
	}()
	// init cache
	c := kvCleanup.NewKVCleanup(60*time.Second, 60*time.Second)
	// init grpcClient
	grpcCl, err := grpc.NewGrpcClient("localhost:8081", 100*time.Millisecond)
	if err != nil {
		log.Fatal(fmt.Errorf("error creating grpc client: %w", err))
	}
	defer func() {
		_ = grpcCl.Close()
	}()
	// init compute
	comp := compute.NewComp(c, grpcCl)
	// init rest api server
	server := http.NewHttpServer("0.0.0.0:8080", 2*time.Second, comp, lg)
	go func() {
		_ = server.ListenAndServe()
	}()
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = server.Shutdown(ctx)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
