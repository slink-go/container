// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (C) 2025-2026 Mikhail Kantur

package container

import (
	"context"
	"sync"
	"time"
)

type Item[T any] struct {
	value      T
	lastAccess time.Time
	created    time.Time
}

type TTLMap[K comparable, V any] struct {
	m   map[K]*Item[V]
	l   sync.RWMutex
	ttl time.Duration
}

func NewTTLMap[K comparable, V any](initialSize int, ttl time.Duration) *TTLMap[K, V] {
	// Validate TTL
	if ttl <= 0 {
		ttl = time.Hour // Default to 1 hour if invalid
	}

	m := &TTLMap[K, V]{
		m:   make(map[K]*Item[V], initialSize),
		ttl: ttl,
	}
	return m
}

func (m *TTLMap[K, V]) Start(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(m.ttl / 2)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				m.Cleanup()
			}
		}
	}()
}

func (m *TTLMap[K, V]) Len() int {
	m.l.RLock()
	defer m.l.RUnlock()
	return len(m.m)
}

func (m *TTLMap[K, V]) Put(k K, v V) {
	m.l.Lock()
	it, ok := m.m[k]
	if !ok {
		it = &Item[V]{}
		m.m[k] = it
	}

	now := time.Now()
	it.value = v
	it.lastAccess = now
	it.created = now

	m.l.Unlock()
}

func (m *TTLMap[K, V]) Get(k K) (v V, ok bool) {
	m.l.RLock()
	it, ok := m.m[k]
	if !ok {
		m.l.RUnlock()
		var empty V
		return empty, false
	}

	now := time.Now()

	// Check if item has expired before accessing it
	if now.Sub(it.created) > m.ttl {
		m.l.RUnlock()

		// Remove expired item
		m.l.Lock()
		if it2, exists := m.m[k]; exists && it2 == it {
			delete(m.m, k)
		}
		m.l.Unlock()

		var empty V
		return empty, false
	}

	// Update access time atomically with the read
	it.lastAccess = now
	v = it.value

	m.l.RUnlock()
	return v, true
}

func (m *TTLMap[K, V]) Has(k K) bool {
	m.l.RLock()
	it, ok := m.m[k]
	if !ok {
		m.l.RUnlock()
		return false
	}
	m.l.RUnlock()
	// Check if item has expired
	return time.Since(it.created) <= m.ttl
}

func (m *TTLMap[K, V]) GetWithTTL(k K) (v V, ttl time.Duration, ok bool) {
	m.l.RLock()
	it, ok := m.m[k]
	if !ok {
		m.l.RUnlock()
		var empty V
		return empty, 0, false
	}

	now := time.Now()

	// Check if item has expired before accessing it
	if now.Sub(it.created) > m.ttl {
		m.l.RUnlock()

		// Remove expired item
		m.l.Lock()
		if it2, exists := m.m[k]; exists && it2 == it {
			delete(m.m, k)
		}
		m.l.Unlock()

		var empty V
		return empty, 0, false
	}

	// Calculate remaining TTL
	remaining := m.ttl - now.Sub(it.created)
	v = it.value
	m.l.RUnlock()

	return v, remaining, true
}

func (m *TTLMap[K, V]) Clear() {
	m.l.Lock()
	for k := range m.m {
		delete(m.m, k)
	}
	m.l.Unlock()
}

func (m *TTLMap[K, V]) Cleanup() {
	now := time.Now()
	m.l.Lock()

	for k, v := range m.m {
		if now.Sub(v.created) > m.ttl {
			delete(m.m, k)
		}
	}
	m.l.Unlock()
}

func (m *TTLMap[K, V]) ExtendTTL(k K) bool {
	m.l.Lock()
	it, ok := m.m[k]
	if !ok {
		m.l.Unlock()
		return false
	}
	// Extend the TTL by updating the created time to now
	it.created = time.Now()
	m.l.Unlock()
	return true
}

func (m *TTLMap[K, V]) Keys() []K {
	m.l.RLock()
	keys := make([]K, 0, len(m.m))
	for k := range m.m {
		keys = append(keys, k)
	}
	m.l.RUnlock()
	return keys
}

func (m *TTLMap[K, V]) Values() []V {
	m.l.RLock()
	values := make([]V, 0, len(m.m))
	for _, item := range m.m {
		values = append(values, item.value)
	}
	m.l.RUnlock()
	return values
}
