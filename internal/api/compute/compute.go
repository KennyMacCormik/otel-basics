package compute

import (
	"jaeger/internal/api/cache"
	grpc2 "jaeger/internal/api/net/grpc"
	"log/slog"
)

type Compute interface {
	Get(key string, lg *slog.Logger) (string, bool, error)
	Set(key, value string, lg *slog.Logger) error
	Del(key string, lg *slog.Logger) error
}

type Comp struct {
	cache      cache.Storage
	grpcClient grpc2.StorageEndpointClient
}

func NewComp(cache cache.Storage, grpcClient grpc2.StorageEndpointClient) *Comp {
	return &Comp{cache: cache, grpcClient: grpcClient}
}

func (c Comp) Get(key string, lg *slog.Logger) (string, bool, error) {
	v, ok := c.cache.Get(key)
	if !ok {
		lg.Debug("cache miss")
		v, ok, err := c.grpcClient.Get(key)
		if err != nil {
			return "", false, err
		}
		c.cache.Set(key, v)
		return v, ok, nil
	}
	lg.Debug("cache hit")
	return v, ok, nil
}

func (c Comp) Set(key, value string, lg *slog.Logger) error {
	c.cache.Set(key, value)
	err := c.grpcClient.Set(key, value)
	if err != nil {
		return err
	}
	return nil
}

func (c Comp) Del(key string, lg *slog.Logger) error {
	c.cache.Del(key)
	err := c.grpcClient.Del(key)
	if err != nil {
		return err
	}
	return nil
}
