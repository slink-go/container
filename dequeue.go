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

const (
	expandFactor      = 3.0 / 4.0
	expandCoefficient = 1.25
	shrinkFactor      = 2.0
	shrinkCoefficient = 1.25
)

type Deque[T any] struct {
	items           []T  // Элементы очереди
	initialCapacity int  // Начальая ёмкость
	capacity        int  // Максимальный размер очереди (-1 для безлимитной очереди)
	preemption      bool // Вытеснять крайний элемент с противоположной стороны при добавлении нового в заполненную очередь
	mutex           sync.RWMutex
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
		// TODO: м.б. надо это пересмотреть - сделать макс. размер в несколько раз больше
		//       заданного настройкой, чтобы снизить количество "реаллокаций" массива
		//       (можно это задавать отдельным конфигурационным параметром, чтобы можно было
		//        регулировать размер потребляемой памяти vs. число реаллокаций)
		initialCapacity = qc.sizeLimit
	}

	return &Deque[T]{
		items:           make([]T, 0, initialCapacity),
		initialCapacity: initialCapacity,
		capacity:        qc.sizeLimit,
		preemption:      qc.preemption,
		mutex:           sync.RWMutex{},
	}

}

//func (s *Deque[T]) Lock() {
//	s.mutex.Lock()
//}
//func (s *Deque[T]) Unlock() {
//	s.mutex.Unlock()
//}

func (s *Deque[T]) Flush() {
	s.mutex.Lock()
	s.flush()
	s.mutex.Unlock()
}
func (s *Deque[T]) flush() {
	s.items = make([]T, 0, s.initialCapacity)
	s.capacity = s.initialCapacity
}

func (s *Deque[T]) Replace(items ...T) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.flush()
	return s.processAll(true, s.pushHeadInternal, items...)
}

func (s *Deque[T]) PushHeadAll(items ...T) (err error) {
	s.mutex.Lock()
	err = s.processAll(false, s.pushHeadInternal, items...)
	s.mutex.Unlock()
	return err
}
func (s *Deque[T]) PushHeadAllReversed(items ...T) (err error) {
	s.mutex.Lock()
	err = s.processAll(true, s.pushHeadInternal, items...)
	s.mutex.Unlock()
	return err
}
func (s *Deque[T]) PushHead(item T) error {

	// TODO: возможно, лучше это всё обернуть в один WriteLock

	s.mutex.RLock()
	err := s.checkSize()
	s.mutex.RUnlock()
	if err != nil {
		return err
	}

	s.mutex.Lock()
	s.pushHeadInternal(item)
	s.mutex.Unlock()

	return nil

}
func (s *Deque[T]) PopHead() (T, error) {

	if s == nil {
		var empty T
		return empty, ErrDequeueEmpty
	}

	s.mutex.RLock()

	if len(s.items) == 0 {
		s.mutex.RUnlock()
		var empty T
		return empty, ErrDequeueEmpty
	}

	s.mutex.RUnlock()

	s.mutex.Lock()

	item := s.items[0]
	s.items = s.items[1:]

	s.mutex.Unlock()

	return item, nil

}
func (s *Deque[T]) PeekHead() (T, error) {

	if s == nil {
		var empty T
		return empty, ErrDequeueEmpty
	}

	s.mutex.RLock()

	if len(s.items) == 0 {
		s.mutex.RUnlock()
		var empty T
		return empty, ErrDequeueEmpty
	}

	item := s.items[0]

	s.mutex.RUnlock()

	return item, nil

}

func (s *Deque[T]) PushTailAll(items ...T) (err error) {
	s.mutex.Lock()
	err = s.processAll(false, s.pushTailInternal, items...)
	s.mutex.Unlock()
	return err
}
func (s *Deque[T]) PushTailAllReversed(items ...T) (err error) {
	s.mutex.Lock()
	err = s.processAll(true, s.pushTailInternal, items...)
	s.mutex.Unlock()
	return err
}
func (s *Deque[T]) PushTail(item T) error {

	// TODO: возможно, лучше это всё обернуть в один WriteLock

	s.mutex.RLock()
	err := s.checkSize()
	s.mutex.RUnlock()
	if err != nil {
		return err
	}

	s.mutex.Lock()
	s.pushTailInternal(item)
	s.mutex.Unlock()

	return nil
}
func (s *Deque[T]) PopTail() (T, error) {

	if s == nil {
		var empty T
		return empty, ErrDequeueEmpty
	}

	s.mutex.RLock()

	if len(s.items) == 0 {
		s.mutex.RUnlock()
		var empty T
		return empty, ErrDequeueEmpty
	}

	s.mutex.RUnlock()

	s.mutex.Lock()

	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]

	s.mutex.Unlock()

	return item, nil

}
func (s *Deque[T]) PeekTail() (T, error) {

	if s == nil {
		var empty T
		return empty, ErrDequeueEmpty
	}

	s.mutex.RLock()

	if len(s.items) == 0 {
		s.mutex.RUnlock()
		var empty T
		return empty, ErrDequeueEmpty
	}

	item := s.items[len(s.items)-1]

	s.mutex.RUnlock()

	return item, nil
}

func (s *Deque[T]) Expand() (err error) {

	if s.capacity <= 0 {
		return // для unbounded deque ничего не делаем
	}

	s.mutex.Lock()

	if float64(s.size()) >= float64(s.capacity)*expandFactor {
		arr := s.values()
		s.initialCapacity = int(float64(s.capacity) * expandCoefficient)
		s.flush()
		err = s.processAll(false, s.pushTailInternal, arr...)
	}

	s.mutex.Unlock()

	return

}
func (s *Deque[T]) Shrink() (err error) {

	// TODO: сжимать без учёта текущего размера, отрезая хвост (конфигурабельно)

	if s.capacity <= 0 {
		return // для unbounded deque ничего не делаем
	}

	s.mutex.Lock()

	if float64(s.size()) <= float64(s.capacity)*shrinkFactor {
		arr := s.values()
		s.initialCapacity = int(float64(s.capacity) / shrinkCoefficient)
		s.flush()
		err = s.processAll(false, s.pushTailInternal, arr...)
	}

	s.mutex.Unlock()

	return
}

// Values - возвращает массив элементов в очереди
func (s *Deque[T]) Values() []T {
	s.mutex.RLock()
	v := s.values()
	s.mutex.RUnlock()
	return v
}
func (s *Deque[T]) values() []T {
	var v []T
	for _, t := range s.items {
		v = append(v, t)
	}
	return v
}

// Size - возвращает количество элементов в очереди
func (s *Deque[T]) Size() int {
	s.mutex.RLock()
	sz := s.size()
	s.mutex.RUnlock()
	return sz
}
func (s *Deque[T]) size() int {
	return len(s.items)
}

// Capacity - возвращает максимальный размер очереди; -1 для неограниченной очереди
func (s *Deque[T]) Capacity() int {
	return s.capacity
}

func (s *Deque[T]) pushHeadInternal(item T) {

	// добавляем новый элемент
	s.items = append([]T{item}, s.items...)

	// если нужно, вытесняем старый
	if len(s.items) > s.capacity && s.preemption {
		s.items = s.items[:len(s.items)-1]
	}

}
func (s *Deque[T]) pushTailInternal(item T) {

	// добавляем новый элемент
	s.items = append(s.items, item)

	// если нужно, вытесняем старый
	if len(s.items) > s.capacity && s.preemption {
		s.items = s.items[1:len(s.items)]
	}

}
func (s *Deque[T]) checkSize() error {
	if s.capacity > 0 && len(s.items) >= s.capacity && !s.preemption {
		return ErrDequeueFull
	}
	return nil
}
func (s *Deque[T]) processAll(reversed bool, receiver func(T), items ...T) (err error) {

	if len(items) == 0 {
		return nil
	}

	if reversed {
		for i := len(items) - 1; i >= 0; i-- {
			err = s.checkSize()
			if err != nil {
				break
			}
			receiver(items[i])
		}
	} else {
		for i := 0; i < len(items); i++ {
			err = s.checkSize()
			if err != nil {
				break
			}
			receiver(items[i])
		}
	}

	return err

}
