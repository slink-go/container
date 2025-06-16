package container

import (
	"sync"
)

type Set[T comparable] struct {
	values map[T]struct{}
	mutex  sync.RWMutex
}

func NewSet[T comparable]() *Set[T] {
	values := make(map[T]struct{})
	set := &Set[T]{
		values: values,
		mutex:  sync.RWMutex{},
	}
	return set
}

func (s *Set[T]) Add(value ...T) *Set[T] {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, v := range value {
		s.values[v] = struct{}{}
	}
	return s
}
func (s *Set[T]) Remove(value T) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.values, value)
}
func (s *Set[T]) Contains(value T) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	_, ok := s.values[value]
	return ok
}
func (s *Set[T]) Len() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return len(s.values)
}
func (s *Set[T]) Clear() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.values = make(map[T]struct{})
}
func (s *Set[T]) IsEmpty() bool {
	return s.Len() == 0
}
func (s *Set[T]) Size() int {
	return s.Len()
}
func (s *Set[T]) Values() []T {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	values := make([]T, 0, len(s.values))
	for v := range s.values {
		values = append(values, v)
	}
	return values
}

func (s *Set[T]) RemoveAll(other *Set[T]) *Set[T] {
	result := NewSet[T]()
	for v := range s.values {
		if !other.Contains(v) {
			result.Add(v)
		}
	}
	return result
}
func (s *Set[T]) AddAll(other *Set[T]) *Set[T] {
	result := NewSet[T]()
	for v := range s.values {
		result.Add(v)
	}
	for v := range other.values {
		result.Add(v)
	}
	return result
}
