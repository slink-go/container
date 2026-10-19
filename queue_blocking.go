// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (C) 2025-2026 Mikhail Kantur

package container

import (
	"sync"
)

// region - linked-list based queue

type node[T any] struct {
	data T
	next *node[T]
}
type BlockingQueue[T any] struct {
	size int
	head *node[T]
	tail *node[T]
	sync.RWMutex
}

func CreateBlockingQueue[T any]() *BlockingQueue[T] {
	return &BlockingQueue[T]{}
}

func (q *BlockingQueue[T]) Enqueue(value T) {
	q.Lock()
	node := &node[T]{data: value}
	if q.head == nil {
		q.head = node
		q.tail = node
		q.size++
		q.Unlock()
		return
	}
	q.tail.next = node
	q.tail = node
	q.size++
	q.Unlock()
}
func (q *BlockingQueue[T]) Dequeue() (T, bool) {
	var v T
	if q.IsEmpty() {
		return v, false //fmt.Errorf("empty queue")
	}
	q.Lock()
	v = q.head.data

	//q.head = q.head.next

	var empty T
	oldHead := q.head
	q.head = q.head.next
	oldHead.data = empty // Очистка данных
	oldHead.next = nil   // Очистка ссылки

	q.size--
	q.Unlock()
	return v, true
}

// Extract - достаёт из очереди (с удалением) элементы, соответствующие фильтру
func (q *BlockingQueue[T]) Extract(filter func(T) bool) []T {
	q.Lock()

	cur := q.head
	prev := q.head

	var res []T
	for ; cur != nil; cur = cur.next {
		if filter(cur.data) {
			res = append(res, cur.data)
			if cur == q.head {
				q.head = cur.next
			} else {
				prev.next = cur.next
			}
		}
		prev = cur
	}

	q.Unlock()
	return res
}

func (q *BlockingQueue[T]) IsEmpty() bool {
	return q.head == nil
}

// endregion
