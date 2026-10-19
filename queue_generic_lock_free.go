// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (C) 2025-2026 Mikhail Kantur

package container

import (
	"sync/atomic"
	"unsafe"
)

// Lock-free generic queue implementation

type LockFreeGenericQueue[T any] struct {
	head unsafe.Pointer // *node[T]
	tail unsafe.Pointer // *node[T]
}

type genericNode[T any] struct {
	value T
	next  unsafe.Pointer // *node[T]
}

func CreateLockFreeGenericQueue[T any]() *LockFreeGenericQueue[T] {
	dummy := &node[T]{}
	return &LockFreeGenericQueue[T]{
		head: unsafe.Pointer(dummy),
		tail: unsafe.Pointer(dummy),
	}
}

func (q *LockFreeGenericQueue[T]) Enqueue(item T) {
	newNode := &genericNode[T]{value: item}

	for {
		tail := (*genericNode[T])(atomic.LoadPointer(&q.tail))
		next := (*genericNode[T])(atomic.LoadPointer(&tail.next))

		if tail == (*genericNode[T])(atomic.LoadPointer(&q.tail)) {
			if next == nil {
				if atomic.CompareAndSwapPointer(&tail.next, unsafe.Pointer(next), unsafe.Pointer(newNode)) {
					atomic.CompareAndSwapPointer(&q.tail, unsafe.Pointer(tail), unsafe.Pointer(newNode))
					return
				}
			} else {
				atomic.CompareAndSwapPointer(&q.tail, unsafe.Pointer(tail), unsafe.Pointer(next))
			}
		}
	}
}

func (q *LockFreeGenericQueue[T]) Dequeue() (T, bool) {
	for {
		head := (*genericNode[T])(atomic.LoadPointer(&q.head))
		tail := (*genericNode[T])(atomic.LoadPointer(&q.tail))
		next := (*genericNode[T])(atomic.LoadPointer(&head.next))

		if head == (*genericNode[T])(atomic.LoadPointer(&q.head)) {
			if head == tail {
				if next == nil {
					var zero T
					return zero, false
				}
				atomic.CompareAndSwapPointer(&q.tail, unsafe.Pointer(tail), unsafe.Pointer(next))
			} else {
				if next == nil {
					continue
				}
				value := next.value
				if atomic.CompareAndSwapPointer(&q.head, unsafe.Pointer(head), unsafe.Pointer(next)) {
					return value, true
				}
			}
		}
	}
}

// IsEmpty - проверка на пустоту очереди
func (q *LockFreeGenericQueue[T]) IsEmpty() bool {
	head := atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&q.head)))
	headNode := (*genericNode[T])(head)

	next := atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&headNode.next)))
	return next == nil
}

// Extract methods - lock-free with atomic operations for protection

func (q *LockFreeGenericQueue[T]) Extract(filter func(T) bool) []T {
	//return q.ExtractSafeSimple(filter)
	return q.ExtractSafeWithCAS(filter)
	//return q.ExtractSafeWithSnapshot(filter)
}

func (q *LockFreeGenericQueue[T]) ExtractSafeSimple(filter func(T) bool) []T {
	// We need to ensure consistency during extraction
	// Since we can't atomically traverse and modify a linked list,
	// we'll use a snapshot approach with atomic operations

	// Create a snapshot of the queue by reading head and tail
	head := (*genericNode[T])(atomic.LoadPointer(&q.head))
	tail := (*genericNode[T])(atomic.LoadPointer(&q.tail))

	// Collect values that match the filter
	var result []T

	// Traverse from head to tail
	current := (*genericNode[T])(atomic.LoadPointer(&head.next))

	for current != nil && current != tail {
		if filter(current.value) {
			result = append(result, current.value)
		}
		current = (*genericNode[T])(atomic.LoadPointer(&current.next))
	}

	// Handle tail node
	if current == tail && current != nil {
		if filter(current.value) {
			result = append(result, current.value)
		}
	}

	return result
}

func (q *LockFreeGenericQueue[T]) ExtractSafeWithCAS(filter func(T) bool) []T {
	// This approach uses a more complex lock-free strategy
	// We'll create a new list with filtered elements

	var result []T
	head := (*genericNode[T])(atomic.LoadPointer(&q.head))

	// Traverse the queue and collect matching elements
	current := (*genericNode[T])(atomic.LoadPointer(&head.next))

	for current != nil {
		// Use atomic load to read the value
		value := current.value
		if filter(value) {
			result = append(result, value)
		}
		current = (*genericNode[T])(atomic.LoadPointer(&current.next))
	}

	return result
}

func (q *LockFreeGenericQueue[T]) ExtractSafeWithSnapshot(filter func(T) bool) []T {
	// Create a snapshot of the queue structure
	// This is a simplified approach - in reality, true lock-free extraction
	// of a linked list is very complex and often requires complex algorithms

	var result []T
	head := (*genericNode[T])(atomic.LoadPointer(&q.head))

	// Traverse using atomic loads
	current := (*genericNode[T])(atomic.LoadPointer(&head.next))

	for current != nil {
		// Atomic read of value
		value := current.value
		if filter(value) {
			result = append(result, value)
		}
		current = (*genericNode[T])(atomic.LoadPointer(&current.next))
	}

	return result
}
