package kvCleanup

import (
	"math/rand"
	"sync"
	"time"
)

const RandTimeout = 10

type item struct {
	Value   string
	Created time.Time
	Timeout time.Duration
	mtx     sync.RWMutex
}

type KVCleanup struct {
	m           sync.Map
	tick        *time.Ticker
	itemTimeout time.Duration
	close       chan struct{}
}

func NewKVCleanup(itemTimeout, cleanupTimeout time.Duration) *KVCleanup {
	c := &KVCleanup{close: make(chan struct{}), tick: time.NewTicker(cleanupTimeout), itemTimeout: itemTimeout}
	go c.cleanup()
	return c
}

func (K *KVCleanup) Close() {
	close(K.close)
}

// TODO: fix this nonsense with sync.Map

func (K *KVCleanup) Get(key string) (string, bool) {
	v, ok := K.m.Load(key)
	if !ok {
		return "", ok
	}
	val := v.(*item)
	val.mtx.RLock()
	defer val.mtx.RUnlock()
	return val.Value, true
}

func (K *KVCleanup) Set(key, value string) {
	timeout := K.itemTimeout + time.Duration(rand.Intn(RandTimeout))*time.Second
	K.m.Store(key, &item{Value: value, Created: time.Now(), Timeout: timeout})
}

func (K *KVCleanup) Del(key string) {
	K.m.Delete(key)
}

func (K *KVCleanup) cleanup() {
	defer K.tick.Stop()
	for {
		select {
		case <-K.close:
			return
		case <-K.tick.C:
			K.m.Range(func(k, v interface{}) bool {
				val := v.(*item) // Type assertion
				val.mtx.Lock()
				defer val.mtx.Unlock()
				if time.Since(val.Created) > val.Timeout {
					K.m.Delete(k)
				}
				return true
			})
		}
	}
}
