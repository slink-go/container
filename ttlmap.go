package container

import (
	"sync"
	"time"
)

type Item[T any] struct {
	value      T
	lastAccess int64
	created    int64
}

type TTLMap[T any] struct {
	m   map[string]*Item[T]
	l   sync.RWMutex
	ttl int64
}

func NewTTLMap[T any](ln int, maxTTL int) *TTLMap[T] {
	m := &TTLMap[T]{
		m:   make(map[string]*Item[T], ln),
		ttl: int64(maxTTL),
	}
	//go func() {
	//	for now := range time.Tick(5 * time.Second) {
	//		m.l.Lock()
	//		for k, v := range m.m {
	//			if now.Unix()-v.lastAccess > int64(maxTTL) {
	//				deletinere(m.m, k)
	//			}
	//		}
	//		m.l.Unlock()
	//	}
	//}()
	return m
}

func (m *TTLMap[T]) Len() int {
	return len(m.m)
}
func (m *TTLMap[T]) Put(k string, v T) {
	m.l.Lock()
	defer m.l.Unlock()
	it, ok := m.m[k]
	if !ok {
		it = &Item[T]{value: v}
		m.m[k] = it
	}
	it.lastAccess = time.Now().Unix()
	it.created = time.Now().Unix()
}
func (m *TTLMap[T]) Get(k string) (v T, ok bool) {
	m.l.RLock()
	defer m.l.RUnlock()
	it, ok := m.m[k]
	if ok {
		v = it.value
		it.lastAccess = time.Now().Unix()
		if time.Now().Unix()-it.created > m.ttl {
			delete(m.m, k)
			return v, false
		}
	}
	return v, ok
}
func (m *TTLMap[T]) Clear() {
	m.l.Lock()
	for k, _ := range m.m {
		delete(m.m, k)
	}
	m.l.Unlock()
}
