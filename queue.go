package container

type Queue[T any] interface {
	Enqueue(value T)
	Dequeue() (T, bool)
	Extract(filter func(T) bool) []T
	IsEmpty() bool
}

func CreateGenericQueue[T any]() Queue[T] {
	return CreateBlockingQueue[T]()
	//return CreateLockFreeGenericQueue[T]()
}
