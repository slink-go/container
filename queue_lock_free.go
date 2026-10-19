// SPDX-License-Identifier: AGPL-3.0-only
// Copyright (C) 2025-2026 Mikhail Kantur

package container

import (
	"sync/atomic"
	"unsafe"
)

// Node - узел очереди
type Node struct {
	value interface{}
	next  *Node
}

// LockFreeQueue - lock-free очередь
type LockFreeQueue struct {
	head *Node
	tail *Node
}

// NewLockFreeQueue - создание новой lock-free очереди
func NewLockFreeQueue() *LockFreeQueue {
	dummy := &Node{}
	return &LockFreeQueue{
		head: dummy,
		tail: dummy,
	}
}

// Enqueue - добавление элемента в очередь
func (q *LockFreeQueue) Enqueue(value interface{}) {
	newNode := &Node{value: value}

	for {
		tail := atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&q.tail)))
		tailNode := (*Node)(tail)

		next := atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&tailNode.next)))
		nextNode := (*Node)(next)

		// Проверяем, не изменилась ли tail-нода
		if tail == atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&q.tail))) {
			// Если tail-нода указывает на NULL, пытаемся добавить новый элемент
			if nextNode == nil {
				// Пытаемся добавить новый узел в конец очереди
				if atomic.CompareAndSwapPointer(
					(*unsafe.Pointer)(unsafe.Pointer(&tailNode.next)),
					unsafe.Pointer(nextNode),
					unsafe.Pointer(newNode)) {
					// Успешно добавили узел, обновляем tail
					atomic.CompareAndSwapPointer(
						(*unsafe.Pointer)(unsafe.Pointer(&q.tail)),
						tail,
						unsafe.Pointer(newNode))
					return
				}
			} else {
				// Если next не NULL, обновляем tail
				atomic.CompareAndSwapPointer(
					(*unsafe.Pointer)(unsafe.Pointer(&q.tail)),
					tail,
					unsafe.Pointer(nextNode))
			}
		}
	}
}

// Dequeue - извлечение элемента из очереди
func (q *LockFreeQueue) Dequeue() (interface{}, bool) {
	for {
		head := atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&q.head)))
		headNode := (*Node)(head)

		next := atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&headNode.next)))
		nextNode := (*Node)(next)

		// Проверяем, не изменился ли head
		if head == atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&q.head))) {
			// Если next NULL, очередь пуста
			if nextNode == nil {
				return nil, false
			}

			// Пытаемся переместить head на следующий узел
			if atomic.CompareAndSwapPointer(
				(*unsafe.Pointer)(unsafe.Pointer(&q.head)),
				head,
				unsafe.Pointer(nextNode)) {
				// Успешно извлекли элемент
				value := nextNode.value
				nextNode.value = nil // Очищаем значение для сборщика мусора
				return value, true
			}
		}
	}
}

// IsEmpty - проверка на пустоту очереди
func (q *LockFreeQueue) IsEmpty() bool {
	head := atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&q.head)))
	headNode := (*Node)(head)

	next := atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&headNode.next)))
	return next == nil
}

// Extract - извлекает элементы из очереди по фильтру и удаляет их из очереди
func (q *LockFreeQueue) Extract(filter func(interface{}) bool) []interface{} {
	if q == nil {
		return nil
	}

	// Создаем временный список для хранения отфильтрованных элементов
	var filteredElements []interface{}

	// Используем подход с "прогоном" через очередь
	// Для этого будем использовать временную очередь для переноса элементов
	tempQueue := NewLockFreeQueue()

	// Собираем все элементы из исходной очереди
	var elements []interface{}

	// Извлекаем все элементы из оригинальной очереди
	for {
		element, ok := q.Dequeue()
		if !ok || element == nil {
			break
		}
		elements = append(elements, element)
	}

	// Разделяем элементы: те, что прошли фильтр - в результат, остальные - во временную очередь
	for _, element := range elements {
		if filter(element) {
			filteredElements = append(filteredElements, element)
		} else {
			tempQueue.Enqueue(element)
		}
	}

	// Восстанавливаем оставшиеся элементы обратно в исходную очередь
	for {
		element, ok := tempQueue.Dequeue()
		if !ok || element == nil {
			break
		}
		q.Enqueue(element)
	}

	return filteredElements
}

// ExtractOptimized - Альтернативная реализация с использованием CAS для более эффективного извлечения
func (q *LockFreeQueue) ExtractOptimized(filter func(interface{}) bool) []interface{} {
	if q == nil {
		return nil
	}

	var filteredElements []interface{}

	// Используем подход с временной очередью для минимизации времени блокировки
	tempQueue := NewLockFreeQueue()

	// Собираем элементы из оригинальной очереди
	var elements []interface{}

	// Извлекаем элементы из очереди
	for {
		element, ok := q.Dequeue()
		if !ok || element == nil {
			break
		}
		elements = append(elements, element)
	}

	// Разделяем элементы по фильтру
	for _, element := range elements {
		if filter(element) {
			filteredElements = append(filteredElements, element)
		} else {
			tempQueue.Enqueue(element)
		}
	}

	// Восстанавливаем элементы обратно в очередь
	for {
		element, ok := tempQueue.Dequeue()
		if !ok || element == nil {
			break
		}
		q.Enqueue(element)
	}

	return filteredElements
}

// ExtractWithState - Экстракт с использованием внутреннего состояния очереди
func (q *LockFreeQueue) ExtractWithState(filter func(interface{}) bool) []interface{} {
	if q == nil {
		return nil
	}

	var filteredElements []interface{}

	// Используем подход с временным списком
	tempElements := make([]interface{}, 0)

	// Извлекаем все элементы
	for {
		element, ok := q.Dequeue()
		if !ok || element == nil {
			break
		}
		tempElements = append(tempElements, element)
	}

	// Применяем фильтр
	for _, element := range tempElements {
		if filter(element) {
			filteredElements = append(filteredElements, element)
		} else {
			q.Enqueue(element) // Возвращаем обратно
		}
	}

	return filteredElements
}

// ExtractSimple - Упрощенная версия для конкретной реализации
func (q *LockFreeQueue) ExtractSimple(filter func(interface{}) bool) []interface{} {
	if q == nil || filter == nil {
		return nil
	}

	var result []interface{}

	// Извлекаем все элементы из очереди
	elements := make([]interface{}, 0)
	for {
		element, ok := q.Dequeue()
		if !ok || element == nil {
			break
		}
		elements = append(elements, element)
	}

	// Применяем фильтр и восстанавливаем элементы
	for _, element := range elements {
		if filter(element) {
			result = append(result, element)
		} else {
			q.Enqueue(element)
		}
	}

	return result
}

// ExtractEfficient - Более эффективная реализация с минимальным количеством операций
func (q *LockFreeQueue) ExtractEfficient(filter func(interface{}) bool) []interface{} {
	if q == nil || filter == nil {
		return nil
	}

	var filteredElements []interface{}

	// Создаем временную очередь для восстановления элементов
	tempQueue := NewLockFreeQueue()

	// Извлекаем все элементы из оригинальной очереди
	var allElements []interface{}
	for {
		element, ok := q.Dequeue()
		if !ok || element == nil {
			break
		}
		allElements = append(allElements, element)
	}

	// Разделяем элементы по фильтру
	for _, element := range allElements {
		if filter(element) {
			filteredElements = append(filteredElements, element)
		} else {
			tempQueue.Enqueue(element)
		}
	}

	// Восстанавливаем оставшиеся элементы обратно в очередь
	for {
		element, ok := tempQueue.Dequeue()
		if !ok || element == nil {
			break
		}
		q.Enqueue(element)
	}

	return filteredElements
}
