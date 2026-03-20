package container

import (
	"testing"
)

// BenchmarkUnboundedDequePushHead benchmarks PushHead operation on unbounded Deque
func BenchmarkUnboundedDequePushHead(b *testing.B) {
	deque := NewDeque[int]()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		deque.PushHead(i)
	}
}

// BenchmarkUnboundedDequePushTail benchmarks PushTail operation on unbounded Deque
func BenchmarkUnboundedDequePushTail(b *testing.B) {
	deque := NewDeque[int]()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		deque.PushTail(i)
	}
}

// BenchmarkUnboundedDequePopHead benchmarks PopHead operation on unbounded Deque
func BenchmarkUnboundedDequePopHead(b *testing.B) {
	deque := NewDeque[int]()

	// Pre-populate the deque
	for i := 0; i < 1000; i++ {
		deque.PushTail(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		deque.PopHead()
	}
}

// BenchmarkUnboundedDequePopTail benchmarks PopTail operation on unbounded Deque
func BenchmarkUnboundedDequePopTail(b *testing.B) {
	deque := NewDeque[int]()

	// Pre-populate the deque
	for i := 0; i < 1000; i++ {
		deque.PushTail(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		deque.PopTail()
	}
}

// BenchmarkBoundedDequePushHead benchmarks PushHead operation on bounded Deque
func BenchmarkBoundedDequePushHead(b *testing.B) {
	deque := NewDeque[int](DequeWithSizeLimit(1000))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		deque.PushHead(i)
	}
}

// BenchmarkBoundedDequePushTail benchmarks PushTail operation on bounded Deque
func BenchmarkBoundedDequePushTail(b *testing.B) {
	deque := NewDeque[int](DequeWithSizeLimit(1000))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		deque.PushTail(i)
	}
}

// BenchmarkBoundedDequePushHeadWithPreemption benchmarks PushHead operation with preemption on bounded Deque
func BenchmarkBoundedDequePushHeadWithPreemption(b *testing.B) {
	deque := NewDeque[int](DequeWithSizeLimit(100), DequeWithPreemption())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		deque.PushHead(i)
	}
}

// BenchmarkBoundedDequePushTailWithPreemption benchmarks PushTail operation with preemption on bounded Deque
func BenchmarkBoundedDequePushTailWithPreemption(b *testing.B) {
	deque := NewDeque[int](DequeWithSizeLimit(100), DequeWithPreemption())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		deque.PushTail(i)
	}
}

// BenchmarkUnboundedDequePushHeadAll benchmarks PushHeadAll operation on unbounded Deque
func BenchmarkUnboundedDequePushHeadAll(b *testing.B) {
	deque := NewDeque[int]()

	items := make([]int, 100)
	for i := range items {
		items[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		deque.PushHeadAll(items...)
	}
}

// BenchmarkUnboundedDequePushTailAll benchmarks PushTailAll operation on unbounded Deque
func BenchmarkUnboundedDequePushTailAll(b *testing.B) {
	deque := NewDeque[int]()

	items := make([]int, 100)
	for i := range items {
		items[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		deque.PushTailAll(items...)
	}
}

// BenchmarkDequeSize benchmarks Size operation on Deque
func BenchmarkDequeSize(b *testing.B) {
	deque := NewDeque[int]()

	// Pre-populate the deque
	for i := 0; i < 1000; i++ {
		deque.PushTail(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		deque.Size()
	}
}

// BenchmarkDequeValues benchmarks Values operation on Deque
func BenchmarkDequeValues(b *testing.B) {
	deque := NewDeque[int]()

	// Pre-populate the deque
	for i := 0; i < 1000; i++ {
		deque.PushTail(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		deque.Values()
	}
}

// BenchmarkDequeExpand benchmarks the Expand operation
func BenchmarkDequeExpand(b *testing.B) {
	deque := NewDeque[int](DequeWithSizeLimit(100))

	// Fill up the deque to trigger expansion
	for i := 0; i < 100; i++ {
		deque.PushTail(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		deque.Expand()
	}
}

// BenchmarkDequeShrink benchmarks the Shrink operation
func BenchmarkDequeShrink(b *testing.B) {
	deque := NewDeque[int](DequeWithSizeLimit(1000))

	// Fill up the deque and then shrink it
	for i := 0; i < 1000; i++ {
		deque.PushTail(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		deque.Shrink()
	}
}
