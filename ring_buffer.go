package container

import (
	"errors"
	"sync"
)

var (
	ErrEmptyBuffer = errors.New("empty buffer")
)

type RingBuffer[T any] struct {
	items        []T
	capacity     int
	size         int
	readPointer  int
	writePointer int
	mutex        sync.RWMutex
}

func NewRingBuffer[T any](capacity int) *RingBuffer[T] {
	return &RingBuffer[T]{
		items:        make([]T, capacity),
		capacity:     capacity,
		size:         0,
		readPointer:  0,
		writePointer: 0,
		mutex:        sync.RWMutex{},
	}
}

func (rb *RingBuffer[T]) Cap() int {
	rb.mutex.RLock()
	v := rb.capacity
	rb.mutex.RUnlock()
	return v
}
func (rb *RingBuffer[T]) Len() int {
	rb.mutex.RLock()
	v := rb.size
	rb.mutex.RUnlock()
	return v
}

func (rb *RingBuffer[T]) Push(item T) {
	if rb.size < rb.capacity {
		rb.size++
	}
	rb.items[rb.writePointer] = item
	rb.writePointer = rb.stepUp(rb.writePointer)
	rb.readPointer = rb.stepDown(rb.writePointer)
}
func (rb *RingBuffer[T]) Peek() (T, error) {
	rb.mutex.RLock()
	if rb.size == 0 {
		rb.mutex.RUnlock()
		var v T
		return v, ErrEmptyBuffer
	}
	v := rb.items[rb.readPointer]
	rb.mutex.RUnlock()
	return v, nil
}
func (rb *RingBuffer[T]) Pop() (T, error) {
	rb.mutex.Lock()
	if rb.size == 0 {
		rb.mutex.Unlock()
		var v T
		return v, ErrEmptyBuffer
	}
	item := rb.items[rb.readPointer]
	rb.size--
	rb.writePointer = rb.stepDown(rb.writePointer)
	rb.readPointer = rb.stepDown(rb.writePointer)
	rb.mutex.Unlock()
	return item, nil
}
func (rb *RingBuffer[T]) Values() []T {
	rb.mutex.RLock()
	s := rb.readPointer
	cnt := 0
	res := make([]T, 0, rb.size)
	for idx := s; cnt < rb.size; idx = rb.stepDown(idx) {
		res = append(res, rb.items[idx])
		cnt++
	}
	rb.mutex.RUnlock()
	return res
}

func (rb *RingBuffer[T]) Flush() {
	rb.mutex.Lock()
	rb.size = 0
	rb.readPointer = 0
	rb.writePointer = 0
	rb.mutex.Unlock()
}

func (rb *RingBuffer[T]) stepDown(idx int) int {
	idx = idx - 1
	if idx < 0 {
		idx = rb.capacity - 1
	}
	return idx
}
func (rb *RingBuffer[T]) stepUp(idx int) int {
	return (idx + 1) % rb.capacity
}
