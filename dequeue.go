package container

import (
	"errors"
	"sync"
)

// region - options

type queueConfig struct {
	sizeLimit  int
	preemption bool // Вытесняет элемент с противоположной стороны очереди при вставке нового элемента в заполненную очередь
}
type DequeOption func(*queueConfig)

// DequeWithSizeLimit - ограниченная в размере очередь
func DequeWithSizeLimit(sizeLimit int) DequeOption {
	return func(config *queueConfig) {
		config.sizeLimit = sizeLimit
	}
}

// DequeWithPreemption - разрешить вытеснение старого элемента при добавлении в заполненную очередь
func DequeWithPreemption() DequeOption {
	return func(config *queueConfig) {
		config.preemption = true
	}
}

// endregion
// region - errors

var (
	ErrDequeueEmpty = errors.New("empty dequeue")
	ErrDequeueFull  = errors.New("dequeue full")
)

// endregion

type Deque[T any] struct {
	items      []T  // Элементы очереди
	capacity   int  // Максимальный размер очереди (-1 для безлимитной очереди)
	preemption bool // Вытеснять крайний элемент с противоположной стороны при добавлении нового в заполненную очередь
	mutex      sync.RWMutex
}

func NewDeque[T any](opts ...DequeOption) *Deque[T] {

	qc := &queueConfig{
		sizeLimit:  -1,
		preemption: false,
	}
	for _, opt := range opts {
		opt(qc)
	}

	// для безлимитных очередей вытеснение не нужно
	if qc.sizeLimit <= 0 {
		qc.preemption = false
	}

	// задаём (начальный или постоянный) размер очереди
	initialCapacity := 8
	if qc.sizeLimit > 0 {
		initialCapacity = qc.sizeLimit
	}

	return &Deque[T]{
		items:      make([]T, 0, initialCapacity),
		capacity:   qc.sizeLimit,
		preemption: qc.preemption,
		mutex:      sync.RWMutex{},
	}

}

func (s *Deque[T]) PushHead(item T) error {

	if s.capacity > 0 && len(s.items) >= s.capacity && !s.preemption {
		return ErrDequeueFull
	}

	s.mutex.Lock()

	// добавляем новый элемент
	s.items = append([]T{item}, s.items...)

	// если нужно, вытесняем старый
	if len(s.items) > s.capacity && s.preemption {
		s.items = s.items[:len(s.items)-1]
	}

	s.mutex.Unlock()

	return nil

}
func (s *Deque[T]) PopHead() (T, error) {

	if len(s.items) == 0 {
		var empty T
		return empty, ErrDequeueEmpty
	}

	s.mutex.Lock()

	item := s.items[0]
	s.items = s.items[1:]

	s.mutex.Unlock()

	return item, nil

}
func (s *Deque[T]) PeekHead() (T, error) {

	if len(s.items) == 0 {
		var empty T
		return empty, ErrDequeueEmpty
	}

	s.mutex.RLock()
	item := s.items[0]
	s.mutex.RUnlock()

	return item, nil

}

func (s *Deque[T]) PushTail(item T) error {

	if s.capacity > 0 && len(s.items) >= s.capacity && !s.preemption {
		return ErrDequeueFull
	}

	s.mutex.Lock()

	// добавляем новый элемент
	s.items = append(s.items, item)

	// если нужно, вытесняем старый
	if len(s.items) > s.capacity && s.preemption {
		s.items = s.items[1:len(s.items)]
	}

	s.mutex.Unlock()

	return nil
}
func (s *Deque[T]) PopTail() (T, error) {

	if len(s.items) == 0 {
		var empty T
		return empty, ErrDequeueEmpty
	}

	s.mutex.Lock()

	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]

	s.mutex.Unlock()

	return item, nil

}
func (s *Deque[T]) PeekTail() (T, error) {

	if len(s.items) == 0 {
		var empty T
		return empty, ErrDequeueEmpty
	}

	s.mutex.RLock()

	item := s.items[len(s.items)-1]

	s.mutex.RUnlock()

	return item, nil
}

// Values - возвращает массив элементов в очереди
func (s *Deque[T]) Values() []T {
	s.mutex.RLock()
	var v []T
	for _, t := range s.items {
		v = append(v, t)
	}
	s.mutex.RUnlock()
	return v
}

// Size - возвращает количество элементов в очереди
func (s *Deque[T]) Size() int {
	return len(s.items)
}

// Capacity - возвращает максимальный размер очереди; -1 для неограниченной очереди
func (s *Deque[T]) Capacity() int {
	return s.capacity
}
