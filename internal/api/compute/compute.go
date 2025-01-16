package compute

import (
	"context"
	"jaeger/internal/api/cache"
	grpc2 "jaeger/internal/api/net/grpc"
	"log/slog"
)

type Compute interface {
	Get(ctx context.Context, key string, lg *slog.Logger) (string, bool, error)
	Set(ctx context.Context, key, value string, lg *slog.Logger) error
	Del(ctx context.Context, key string, lg *slog.Logger) error
}

type Comp struct {
	cache      cache.Storage
	grpcClient grpc2.StorageEndpointClient
}

func NewComp(cache cache.Storage, grpcClient grpc2.StorageEndpointClient) *Comp {
	return &Comp{cache: cache, grpcClient: grpcClient}
}

func (c Comp) Get(ctx context.Context, key string, lg *slog.Logger) (string, bool, error) {
	v, ok := c.cache.Get(key)
	if !ok {
		lg.Debug("cache miss")
		v, ok, err := c.grpcClient.Get(ctx, key)
		if err != nil {
			return "", false, err
		}
		c.cache.Set(key, v)
		return v, ok, nil
	}
	lg.Debug("cache hit")
	return v, ok, nil
}

func (c Comp) Set(ctx context.Context, key, value string, lg *slog.Logger) error {
	c.cache.Set(key, value)
	err := c.grpcClient.Set(ctx, key, value)
	if err != nil {
		return err
	}
	return nil
}

func (c Comp) Del(ctx context.Context, key string, lg *slog.Logger) error {
	c.cache.Del(key)
	err := c.grpcClient.Del(ctx, key)
	if err != nil {
		return err
	}
	return nil
}
