package grpc

import (
	"google.golang.org/grpc"
	"jaeger/internal/db/net/grpc/StorageEndpoint"
	"jaeger/internal/db/repo"
	"jaeger/internal/proto/db"
	"log/slog"
)

func NewGrpcServer(st repo.Storage, lg *slog.Logger) *grpc.Server {
	gs := grpc.NewServer()
	db.RegisterStorageEndpointServer(gs, StorageEndpoint.NewStorageEndpoint(st, lg))
	return gs
}
