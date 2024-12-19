package _map

import (
	"context"
)

type Map struct {
	m    map[string]string
	lock chan struct{}
}

func NewMap() *Map {
	return &Map{m: make(map[string]string), lock: make(chan struct{}, 1)}
}

// expire context if too long waiting
func (m *Map) acquireLock(ctx context.Context) error {
	select {
	case m.lock <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (m *Map) releaseLock() {
	<-m.lock
}

func (m *Map) Get(ctx context.Context, key string) (string, bool, error) {
	if err := m.acquireLock(ctx); err != nil {
		return "", false, err
	}
	defer m.releaseLock()
	v, ok := m.m[key]
	return v, ok, nil
}

func (m *Map) Set(ctx context.Context, key string, value string) error {
	if err := m.acquireLock(ctx); err != nil {
		return err
	}
	defer m.releaseLock()
	m.m[key] = value
	return nil
}

func (m *Map) Del(ctx context.Context, key string) error {
	if err := m.acquireLock(ctx); err != nil {
		return err
	}
	defer m.releaseLock()
	delete(m.m, key)
	return nil
}
