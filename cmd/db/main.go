package main

import (
	"context"
	"fmt"
	"jaeger/internal/cfg"
	mygrpc "jaeger/internal/db/net/grpc"
	mytcp "jaeger/internal/db/net/grpc/tcp"
	_map "jaeger/internal/db/repo/map"
	mylog "jaeger/internal/logging"
	"jaeger/internal/tracing"
	"log"
	"os"
	"os/signal"
	"syscall"
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
	// init storage
	st := _map.NewMap()
	// init tcp
	tcp, err := mytcp.NewTcpServer("0.0.0.0:8081")
	if err != nil {
		log.Fatalf("failed to listen 0.0.0.0:8081: %v", err)
	}
	defer func() { _ = tcp.Close() }()
	// init grpc
	gs := mygrpc.NewGrpcServer(st, lg)
	defer gs.GracefulStop()
	// run
	go func() {
		if err = gs.Serve(tcp); err != nil {
			lg.Error("grpc server failed", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
