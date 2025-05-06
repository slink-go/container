package container

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestUnboundedDequeHead(t *testing.T) {

	d := NewDeque[int]()
	assert.Equal(t, 0, d.Size())
	assert.Equal(t, -1, d.Capacity())

	d.PushHead(1)
	assert.Equal(t, 1, d.Size())
	assert.Equal(t, -1, d.Capacity())

	d.PushHead(2)
	assert.Equal(t, 2, d.Size())
	assert.Equal(t, -1, d.Capacity())

	v, err := d.PeekHead()
	assert.NoError(t, err)
	assert.Equal(t, 2, v)

	v, err = d.PopHead()
	assert.NoError(t, err)
	assert.Equal(t, 2, v)

	v, err = d.PeekHead()
	assert.NoError(t, err)
	assert.Equal(t, 1, v)

	v, err = d.PopHead()
	assert.NoError(t, err)
	assert.Equal(t, 1, v)

	v, err = d.PopHead()
	assert.ErrorIs(t, err, ErrDequeueEmpty)
	assert.Equal(t, 0, d.Size())
	assert.Equal(t, -1, d.Capacity())

}
func TestUnboundedDequeTail(t *testing.T) {

	d := NewDeque[int]()
	assert.Equal(t, 0, d.Size())
	assert.Equal(t, -1, d.Capacity())

	d.PushTail(1)
	assert.Equal(t, 1, d.Size())
	assert.Equal(t, -1, d.Capacity())

	d.PushTail(2)
	assert.Equal(t, 2, d.Size())
	assert.Equal(t, -1, d.Capacity())

	v, err := d.PeekTail()
	assert.NoError(t, err)
	assert.Equal(t, 2, v)

	v, err = d.PopTail()
	assert.NoError(t, err)
	assert.Equal(t, 2, v)

	v, err = d.PeekTail()
	assert.NoError(t, err)
	assert.Equal(t, 1, v)

	v, err = d.PopTail()
	assert.NoError(t, err)
	assert.Equal(t, 1, v)

	v, err = d.PopTail()
	assert.ErrorIs(t, err, ErrDequeueEmpty)
	assert.Equal(t, 0, d.Size())
	assert.Equal(t, -1, d.Capacity())

}
func TestUnboundedDequeBoth(t *testing.T) {

	d := NewDeque[int]()
	assert.Equal(t, 0, d.Size())
	assert.Equal(t, -1, d.Capacity())

	d.PushHead(1)
	assert.Equal(t, 1, d.Size())
	assert.Equal(t, -1, d.Capacity())

	d.PushHead(2)
	assert.Equal(t, 2, d.Size())
	assert.Equal(t, -1, d.Capacity())

	d.PushTail(3)
	assert.Equal(t, 3, d.Size())
	assert.Equal(t, -1, d.Capacity())

	d.PushTail(4)
	assert.Equal(t, 4, d.Size())
	assert.Equal(t, -1, d.Capacity())

	v, err := d.PeekHead()
	assert.NoError(t, err)
	assert.Equal(t, 2, v)

	v, err = d.PeekTail()
	assert.NoError(t, err)
	assert.Equal(t, 4, v)

	v, err = d.PopHead()
	assert.NoError(t, err)
	assert.Equal(t, 2, v)

	v, err = d.PopHead()
	assert.NoError(t, err)
	assert.Equal(t, 1, v)

	v, err = d.PeekHead()
	assert.NoError(t, err)
	assert.Equal(t, 3, v)

	v, err = d.PopHead()
	assert.NoError(t, err)
	assert.Equal(t, 3, v)

	v, err = d.PeekHead()
	assert.NoError(t, err)
	assert.Equal(t, 4, v)

	v, err = d.PeekTail()
	assert.NoError(t, err)
	assert.Equal(t, 4, v)

}

func TestBoundedNoPreemptionDequeHead(t *testing.T) {

	d := NewDeque[int](DequeWithSizeLimit(3))
	assert.Equal(t, 0, d.Size())
	assert.Equal(t, 3, d.Capacity())

	// 1
	err := d.PushHead(1)
	assert.NoError(t, err)
	assert.Equal(t, 1, d.Size())
	assert.Equal(t, 3, d.Capacity())

	// 2 1
	err = d.PushHead(2)
	assert.NoError(t, err)
	assert.Equal(t, 2, d.Size())
	assert.Equal(t, 3, d.Capacity())

	// 3 2 1
	err = d.PushHead(3)
	assert.NoError(t, err)
	assert.Equal(t, 3, d.Size())
	assert.Equal(t, 3, d.Capacity())

	// 3 2 1
	err = d.PushHead(4)
	assert.ErrorIs(t, err, ErrDequeueFull)
	assert.Equal(t, 3, d.Size())
	assert.Equal(t, 3, d.Capacity())

	// 3 2 1
	err = d.PushHead(5)
	assert.ErrorIs(t, err, ErrDequeueFull)
	assert.Equal(t, 3, d.Size())
	assert.Equal(t, 3, d.Capacity())

	// 3 2 1
	v, err := d.PeekHead()
	assert.NoError(t, err)
	assert.Equal(t, 3, v)
	assert.Equal(t, 3, d.Size())
	assert.Equal(t, 3, d.Capacity())

	// 2 1
	v, err = d.PopHead()
	assert.NoError(t, err)
	assert.Equal(t, 3, v)
	assert.Equal(t, 2, d.Size())
	assert.Equal(t, 3, d.Capacity())

	// 2 1
	v, err = d.PeekHead()
	assert.NoError(t, err)
	assert.Equal(t, 2, v)
	assert.Equal(t, 2, d.Size())
	assert.Equal(t, 3, d.Capacity())

	// 1
	v, err = d.PopHead()
	assert.NoError(t, err)
	assert.Equal(t, 2, v)
	assert.Equal(t, 1, d.Size())
	assert.Equal(t, 3, d.Capacity())

	// 1
	v, err = d.PeekHead()
	assert.NoError(t, err)
	assert.Equal(t, 1, v)
	assert.Equal(t, 1, d.Size())
	assert.Equal(t, 3, d.Capacity())

	// -
	v, err = d.PopHead()
	assert.NoError(t, err)
	assert.Equal(t, 1, v)
	assert.Equal(t, 0, d.Size())
	assert.Equal(t, 3, d.Capacity())

	v, err = d.PopHead()
	assert.ErrorIs(t, err, ErrDequeueEmpty)
	assert.Equal(t, 0, d.Size())
	assert.Equal(t, 3, d.Capacity())

}
func TestBoundedNoPreemptionDequeTail(t *testing.T) {

	d := NewDeque[int](DequeWithSizeLimit(3))
	assert.Equal(t, 0, d.Size())
	assert.Equal(t, 3, d.Capacity())

	// 1
	err := d.PushTail(1)
	assert.NoError(t, err)
	assert.Equal(t, 1, d.Size())
	assert.Equal(t, 3, d.Capacity())

	// 1 2
	err = d.PushTail(2)
	assert.NoError(t, err)
	assert.Equal(t, 2, d.Size())
	assert.Equal(t, 3, d.Capacity())

	// 1 2 3
	err = d.PushTail(3)
	assert.NoError(t, err)
	assert.Equal(t, 3, d.Size())
	assert.Equal(t, 3, d.Capacity())

	// 1 2 3
	err = d.PushTail(4)
	assert.ErrorIs(t, err, ErrDequeueFull)
	assert.Equal(t, 3, d.Size())
	assert.Equal(t, 3, d.Capacity())

	// 1 2 3
	err = d.PushTail(5)
	assert.ErrorIs(t, err, ErrDequeueFull)
	assert.Equal(t, 3, d.Size())
	assert.Equal(t, 3, d.Capacity())

	// 1 2 3
	v, err := d.PeekTail()
	assert.NoError(t, err)
	assert.Equal(t, 3, v)
	assert.Equal(t, 3, d.Size())
	assert.Equal(t, 3, d.Capacity())

	// 1 2
	v, err = d.PopTail()
	assert.NoError(t, err)
	assert.Equal(t, 3, v)
	assert.Equal(t, 2, d.Size())
	assert.Equal(t, 3, d.Capacity())

	// 1 2
	v, err = d.PeekTail()
	assert.NoError(t, err)
	assert.Equal(t, 2, v)
	assert.Equal(t, 2, d.Size())
	assert.Equal(t, 3, d.Capacity())

	// 1
	v, err = d.PopTail()
	assert.NoError(t, err)
	assert.Equal(t, 2, v)
	assert.Equal(t, 1, d.Size())
	assert.Equal(t, 3, d.Capacity())

	// 1
	v, err = d.PeekTail()
	assert.NoError(t, err)
	assert.Equal(t, 1, v)
	assert.Equal(t, 1, d.Size())
	assert.Equal(t, 3, d.Capacity())

	// -
	v, err = d.PopTail()
	assert.NoError(t, err)
	assert.Equal(t, 1, v)
	assert.Equal(t, 0, d.Size())
	assert.Equal(t, 3, d.Capacity())

	v, err = d.PopTail()
	assert.ErrorIs(t, err, ErrDequeueEmpty)
	assert.Equal(t, 0, d.Size())
	assert.Equal(t, 3, d.Capacity())

}
func TestBoundedNoPreemptionDequeBoth(t *testing.T) {

	d := NewDeque[int](DequeWithSizeLimit(3))
	assert.Equal(t, 0, d.Size())
	assert.Equal(t, 3, d.Capacity())

	// 1
	err := d.PushTail(1)
	assert.NoError(t, err)
	assert.Equal(t, 1, d.Size())
	assert.Equal(t, 3, d.Capacity())

	// 2 1
	err = d.PushHead(2)
	assert.NoError(t, err)
	assert.Equal(t, 2, d.Size())
	assert.Equal(t, 3, d.Capacity())

	// 2 1 3
	err = d.PushTail(3)
	assert.NoError(t, err)
	assert.Equal(t, 3, d.Size())
	assert.Equal(t, 3, d.Capacity())

	// 2 1 3
	err = d.PushTail(4)
	assert.ErrorIs(t, err, ErrDequeueFull)
	assert.Equal(t, 3, d.Size())
	assert.Equal(t, 3, d.Capacity())

	// 2 1 3
	err = d.PushTail(5)
	assert.ErrorIs(t, err, ErrDequeueFull)
	assert.Equal(t, 3, d.Size())
	assert.Equal(t, 3, d.Capacity())

	// 2 1 3
	v, err := d.PeekTail()
	assert.NoError(t, err)
	assert.Equal(t, 3, v)
	assert.Equal(t, 3, d.Size())
	assert.Equal(t, 3, d.Capacity())

	// 2 1 3
	v, err = d.PeekHead()
	assert.NoError(t, err)
	assert.Equal(t, 2, v)
	assert.Equal(t, 3, d.Size())
	assert.Equal(t, 3, d.Capacity())

	// 2 1
	v, err = d.PopTail()
	assert.NoError(t, err)
	assert.Equal(t, 3, v)
	assert.Equal(t, 2, d.Size())
	assert.Equal(t, 3, d.Capacity())

	// 1
	v, err = d.PopHead()
	assert.NoError(t, err)
	assert.Equal(t, 2, v)
	assert.Equal(t, 1, d.Size())
	assert.Equal(t, 3, d.Capacity())

	// 1
	v, err = d.PeekTail()
	assert.NoError(t, err)
	assert.Equal(t, 1, v)
	assert.Equal(t, 1, d.Size())
	assert.Equal(t, 3, d.Capacity())

	// 1
	v, err = d.PeekHead()
	assert.NoError(t, err)
	assert.Equal(t, 1, v)
	assert.Equal(t, 1, d.Size())
	assert.Equal(t, 3, d.Capacity())

	// -
	v, err = d.PopTail()
	assert.NoError(t, err)
	assert.Equal(t, 1, v)
	assert.Equal(t, 0, d.Size())
	assert.Equal(t, 3, d.Capacity())

	v, err = d.PopTail()
	assert.ErrorIs(t, err, ErrDequeueEmpty)
	assert.Equal(t, 0, d.Size())
	assert.Equal(t, 3, d.Capacity())

}

func TestBoundedPreemptionDequeHead(t *testing.T) {

	d := NewDeque[int](
		DequeWithSizeLimit(3),
		DequeWithPreemption(),
	)
	assert.Equal(t, 0, d.Size())
	assert.Equal(t, 3, d.Capacity())

	// 1
	err := d.PushHead(1)
	assert.NoError(t, err)
	assert.Equal(t, 1, d.Size())
	assert.Equal(t, 3, d.Capacity())

	// 2 1
	err = d.PushHead(2)
	assert.NoError(t, err)
	assert.Equal(t, 2, d.Size())
	assert.Equal(t, 3, d.Capacity())

	// 3 2 1
	err = d.PushHead(3)
	assert.NoError(t, err)
	assert.Equal(t, 3, d.Size())
	assert.Equal(t, 3, d.Capacity())

	// 4 3 2
	err = d.PushHead(4)
	assert.NoError(t, err)
	assert.Equal(t, 3, d.Size())
	assert.Equal(t, 3, d.Capacity())

	// 5 4 3
	err = d.PushHead(5)
	assert.NoError(t, err)
	assert.Equal(t, 3, d.Size())
	assert.Equal(t, 3, d.Capacity())

	// 5 4 3
	v, err := d.PeekHead()
	assert.NoError(t, err)
	assert.Equal(t, 5, v)
	assert.Equal(t, 3, d.Size())
	assert.Equal(t, 3, d.Capacity())

	// 4 3
	v, err = d.PopHead()
	assert.NoError(t, err)
	assert.Equal(t, 5, v)
	assert.Equal(t, 2, d.Size())
	assert.Equal(t, 3, d.Capacity())

	// 4 3
	v, err = d.PeekHead()
	assert.NoError(t, err)
	assert.Equal(t, 4, v)
	assert.Equal(t, 2, d.Size())
	assert.Equal(t, 3, d.Capacity())

	// 3
	v, err = d.PopHead()
	assert.NoError(t, err)
	assert.Equal(t, 4, v)
	assert.Equal(t, 1, d.Size())
	assert.Equal(t, 3, d.Capacity())

	// 3
	v, err = d.PeekHead()
	assert.NoError(t, err)
	assert.Equal(t, 3, v)
	assert.Equal(t, 1, d.Size())
	assert.Equal(t, 3, d.Capacity())

	// -
	v, err = d.PopHead()
	assert.NoError(t, err)
	assert.Equal(t, 3, v)
	assert.Equal(t, 0, d.Size())
	assert.Equal(t, 3, d.Capacity())

	v, err = d.PopHead()
	assert.ErrorIs(t, err, ErrDequeueEmpty)
	assert.Equal(t, 0, d.Size())
	assert.Equal(t, 3, d.Capacity())

}
func TestBoundedPreemptionDequeTail(t *testing.T) {

	d := NewDeque[int](
		DequeWithSizeLimit(3),
		DequeWithPreemption(),
	)
	assert.Equal(t, 0, d.Size())
	assert.Equal(t, 3, d.Capacity())

	// 1
	err := d.PushTail(1)
	assert.NoError(t, err)
	assert.Equal(t, 1, d.Size())
	assert.Equal(t, 3, d.Capacity())

	// 1 2
	err = d.PushTail(2)
	assert.NoError(t, err)
	assert.Equal(t, 2, d.Size())
	assert.Equal(t, 3, d.Capacity())

	// 1 2 3
	err = d.PushTail(3)
	assert.NoError(t, err)
	assert.Equal(t, 3, d.Size())
	assert.Equal(t, 3, d.Capacity())

	// 2 3 4
	err = d.PushTail(4)
	assert.NoError(t, err)
	assert.Equal(t, 3, d.Size())
	assert.Equal(t, 3, d.Capacity())

	// 3 4 5
	err = d.PushTail(5)
	assert.NoError(t, err)
	assert.Equal(t, 3, d.Size())
	assert.Equal(t, 3, d.Capacity())

	// 3 4 5
	v, err := d.PeekTail()
	assert.NoError(t, err)
	assert.Equal(t, 5, v)
	assert.Equal(t, 3, d.Size())
	assert.Equal(t, 3, d.Capacity())

	// 3 4
	v, err = d.PopTail()
	assert.NoError(t, err)
	assert.Equal(t, 5, v)
	assert.Equal(t, 2, d.Size())
	assert.Equal(t, 3, d.Capacity())

	// 3 4
	v, err = d.PeekTail()
	assert.NoError(t, err)
	assert.Equal(t, 4, v)
	assert.Equal(t, 2, d.Size())
	assert.Equal(t, 3, d.Capacity())

	// 3
	v, err = d.PopTail()
	assert.NoError(t, err)
	assert.Equal(t, 4, v)
	assert.Equal(t, 1, d.Size())
	assert.Equal(t, 3, d.Capacity())

	// 3
	v, err = d.PeekTail()
	assert.NoError(t, err)
	assert.Equal(t, 3, v)
	assert.Equal(t, 1, d.Size())
	assert.Equal(t, 3, d.Capacity())

	// -
	v, err = d.PopTail()
	assert.NoError(t, err)
	assert.Equal(t, 3, v)
	assert.Equal(t, 0, d.Size())
	assert.Equal(t, 3, d.Capacity())

	v, err = d.PopTail()
	assert.ErrorIs(t, err, ErrDequeueEmpty)
	assert.Equal(t, 0, d.Size())
	assert.Equal(t, 3, d.Capacity())

}
func TestBoundedPreemptionDequeBoth(t *testing.T) {

	d := NewDeque[int](
		DequeWithSizeLimit(3),
		DequeWithPreemption(),
	)
	assert.Equal(t, 0, d.Size())
	assert.Equal(t, 3, d.Capacity())

	// 1
	err := d.PushTail(1)
	assert.NoError(t, err)
	assert.Equal(t, 1, d.Size())
	assert.Equal(t, 3, d.Capacity())

	// 2 1
	err = d.PushHead(2)
	assert.NoError(t, err)
	assert.Equal(t, 2, d.Size())
	assert.Equal(t, 3, d.Capacity())

	// 2 1 3
	err = d.PushTail(3)
	assert.NoError(t, err)
	assert.Equal(t, 3, d.Size())
	assert.Equal(t, 3, d.Capacity())

	// 4 2 1
	err = d.PushHead(4)
	assert.NoError(t, err)
	assert.Equal(t, 3, d.Size())
	assert.Equal(t, 3, d.Capacity())

	// 2 1 5
	err = d.PushTail(5)
	assert.NoError(t, err)
	assert.Equal(t, 3, d.Size())
	assert.Equal(t, 3, d.Capacity())

	// 2 1 5
	v, err := d.PeekTail()
	assert.NoError(t, err)
	assert.Equal(t, 5, v)
	assert.Equal(t, 3, d.Size())
	assert.Equal(t, 3, d.Capacity())

	// 2 1 5
	v, err = d.PeekHead()
	assert.NoError(t, err)
	assert.Equal(t, 2, v)
	assert.Equal(t, 3, d.Size())
	assert.Equal(t, 3, d.Capacity())

	// 2 1
	v, err = d.PopTail()
	assert.NoError(t, err)
	assert.Equal(t, 5, v)
	assert.Equal(t, 2, d.Size())
	assert.Equal(t, 3, d.Capacity())

	// 1
	v, err = d.PopHead()
	assert.NoError(t, err)
	assert.Equal(t, 2, v)
	assert.Equal(t, 1, d.Size())
	assert.Equal(t, 3, d.Capacity())

	// 1
	v, err = d.PeekTail()
	assert.NoError(t, err)
	assert.Equal(t, 1, v)
	assert.Equal(t, 1, d.Size())
	assert.Equal(t, 3, d.Capacity())

	// 1
	v, err = d.PeekHead()
	assert.NoError(t, err)
	assert.Equal(t, 1, v)
	assert.Equal(t, 1, d.Size())
	assert.Equal(t, 3, d.Capacity())

	// -
	v, err = d.PopTail()
	assert.NoError(t, err)
	assert.Equal(t, 1, v)
	assert.Equal(t, 0, d.Size())
	assert.Equal(t, 3, d.Capacity())

	v, err = d.PopTail()
	assert.ErrorIs(t, err, ErrDequeueEmpty)
	assert.Equal(t, 0, d.Size())
	assert.Equal(t, 3, d.Capacity())

}

func TestConcurrentAccess(t *testing.T) {

	d := NewDeque[int](
		DequeWithSizeLimit(3),
		DequeWithPreemption(),
	)

	go func() {
		err := d.PushTail(1)
		assert.NoError(t, err)
		v, err := d.PeekTail()
		assert.NoError(t, err)
		assert.Equal(t, 1, v)
	}()

	go func() {
		err := d.PushHead(2)
		assert.NoError(t, err)
		v, err := d.PeekHead()
		assert.NoError(t, err)
		assert.Equal(t, 2, v)
	}()

	time.Sleep(3 * time.Second)

}
