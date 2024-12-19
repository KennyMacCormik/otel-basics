package repo

import "context"

type Storage interface {
	Get(ctx context.Context, key string) (string, bool, error)
	Set(ctx context.Context, key string, value string) error
	Del(ctx context.Context, key string) error
}
