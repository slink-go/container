package container

import (
	"context"
	"sync"
	"time"
)

type Item[T any] struct {
	value      T
	lastAccess int64
	created    int64
}

type TTLMap[K comparable, V any] struct {
	m   map[K]*Item[V]
	l   sync.RWMutex
	ttl time.Duration
}

func NewTTLMap[K comparable, V any](ln int, maxTTL time.Duration) *TTLMap[K, V] {
	m := &TTLMap[K, V]{
		m:   make(map[K]*Item[V], ln),
		ttl: maxTTL,
	}
	return m
}

func (m *TTLMap[K, V]) Start(ctx context.Context) {
	go func() {
		ticker := time.NewTimer(m.ttl / 2)
		defer ticker.Stop()
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			m.Cleanup()
		}
	}()
}

func (m *TTLMap[K, V]) Len() int {
	return len(m.m)
}
func (m *TTLMap[K, V]) Put(k K, v V) {
	m.l.Lock()
	it, ok := m.m[k]
	if !ok {
		it = &Item[V]{value: v}
		m.m[k] = it
	}
	it.lastAccess = time.Now().Unix()
	it.created = time.Now().Unix()
	m.l.Unlock()
}
func (m *TTLMap[K, V]) Get(k K) (v V, ok bool) {
	m.l.RLock()
	it, ok := m.m[k]
	if ok {
		v = it.value
		it.lastAccess = time.Now().Unix()
		if time.Now().Unix()-it.created > int64(m.ttl.Seconds()) {
			delete(m.m, k)
			m.l.RUnlock()
			return v, false
		}
	}
	m.l.RUnlock()
	return v, ok
}
func (m *TTLMap[K, V]) Clear() {
	m.l.Lock()
	for k, _ := range m.m {
		delete(m.m, k)
	}
	m.l.Unlock()
}
func (m *TTLMap[K, V]) Cleanup() {
	m.l.Lock()
	for k, v := range m.m {
		if time.Now().Unix()-v.created > int64(m.ttl.Seconds()) {
			delete(m.m, k)
		}
	}
	m.l.Unlock()
}
