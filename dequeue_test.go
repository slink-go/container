package container

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"
)

func TestUnboundedDequePushAll(t *testing.T) {

	tests := []struct {
		name             string
		method           string
		items            []int
		expectedSize     int
		expectedCapacity int
		expectedHead     int
		expectedTail     int
	}{
		{
			name:             "unbounded: head push all",
			method:           "PushHeadAll",
			items:            []int{1, 2, 3, 4, 5},
			expectedSize:     5,
			expectedCapacity: -1,
			expectedHead:     5,
			expectedTail:     1,
		},
		{
			name:             "unbounded: head push all reversed",
			method:           "PushHeadAllReversed",
			items:            []int{1, 2, 3, 4, 5},
			expectedSize:     5,
			expectedCapacity: -1,
			expectedHead:     1,
			expectedTail:     5,
		},
		{
			name:             "unbounded: tail push all",
			method:           "PushTailAll",
			items:            []int{1, 2, 3, 4, 5},
			expectedSize:     5,
			expectedCapacity: -1,
			expectedHead:     1,
			expectedTail:     5,
		},
		{
			name:             "unbounded: tail push all reversed",
			method:           "PushTailAllReversed",
			items:            []int{1, 2, 3, 4, 5},
			expectedSize:     5,
			expectedCapacity: -1,
			expectedHead:     5,
			expectedTail:     1,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			d := NewDeque[int]()
			assert.Equal(t, 0, d.Size())
			assert.Equal(t, test.expectedCapacity, d.Capacity())

			arr := make([]interface{}, len(test.items))
			for i := 0; i < len(test.items); i++ {
				arr[i] = test.items[i]
			}
			val, err := invoke(d, test.method, arr...)
			if err != nil {
				t.Fatal(err)
			}
			assert.True(t, val.IsZero()) // means no error

			v, err := d.PeekTail()
			assert.NoError(t, err)
			assert.Equal(t, test.expectedTail, v)

			v, err = d.PeekHead()
			assert.NoError(t, err)
			assert.Equal(t, test.expectedHead, v)
		})
	}

}
func TestBoundedNoPreemptionDequePushAll(t *testing.T) {

	tests := []struct {
		name             string
		method           string
		items            []int
		expectedSize     int
		expectedCapacity int
		expectedHead     int
		expectedTail     int
	}{
		{
			name:             "bounded: head push all",
			method:           "PushHeadAll",
			items:            []int{1, 2, 3, 4, 5},
			expectedSize:     3,
			expectedCapacity: 3,
			expectedHead:     3,
			expectedTail:     1,
		},
		{
			name:             "bounded: head push all reversed",
			method:           "PushHeadAllReversed",
			items:            []int{1, 2, 3, 4, 5},
			expectedSize:     3,
			expectedCapacity: 3,
			expectedHead:     3,
			expectedTail:     5,
		},
		{
			name:             "bounded: tail push all",
			method:           "PushTailAll",
			items:            []int{1, 2, 3, 4, 5},
			expectedSize:     3,
			expectedCapacity: 3,
			expectedHead:     1,
			expectedTail:     3,
		},
		{
			name:             "bounded: tail push all reversed",
			method:           "PushTailAllReversed",
			items:            []int{1, 2, 3, 4, 5},
			expectedSize:     3,
			expectedCapacity: 3,
			expectedHead:     5,
			expectedTail:     3,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			d := NewDeque[int](
				DequeWithSizeLimit(test.expectedCapacity),
			)
			assert.Equal(t, 0, d.Size())
			assert.Equal(t, test.expectedCapacity, d.Capacity())

			arr := make([]interface{}, len(test.items))
			for i := 0; i < len(test.items); i++ {
				arr[i] = test.items[i]
			}
			val, err := invoke(d, test.method, arr...)
			if err != nil {
				t.Fatal(err)
			}
			vv := reflect.ValueOf(val).Interface().(reflect.Value)
			vvv := vv.Interface().(error)
			assert.ErrorIs(t, vvv, ErrDequeueFull)

			v, err := d.PeekTail()
			assert.NoError(t, err)
			assert.Equal(t, test.expectedTail, v)

			v, err = d.PeekHead()
			assert.NoError(t, err)
			assert.Equal(t, test.expectedHead, v)
		})
	}

}
func TestBoundedPreemptionDequePushAll(t *testing.T) {

	tests := []struct {
		name             string
		method           string
		items            []int
		expectedSize     int
		expectedCapacity int
		expectedHead     int
		expectedTail     int
	}{
		{
			name:             "bounded preempted: head push all",
			method:           "PushHeadAll",
			items:            []int{1, 2, 3, 4, 5},
			expectedSize:     3,
			expectedCapacity: 3,
			expectedHead:     5,
			expectedTail:     3,
		},
		{
			name:             "bounded preempted: head push all reversed",
			method:           "PushHeadAllReversed",
			items:            []int{1, 2, 3, 4, 5},
			expectedSize:     3,
			expectedCapacity: 3,
			expectedHead:     1,
			expectedTail:     3,
		},
		{
			name:             "bounded preempted: tail push all",
			method:           "PushTailAll",
			items:            []int{1, 2, 3, 4, 5},
			expectedSize:     3,
			expectedCapacity: 3,
			expectedHead:     3,
			expectedTail:     5,
		},
		{
			name:             "bounded preempted: tail push all reversed",
			method:           "PushTailAllReversed",
			items:            []int{1, 2, 3, 4, 5},
			expectedSize:     3,
			expectedCapacity: 3,
			expectedHead:     3,
			expectedTail:     1,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			d := NewDeque[int](
				DequeWithPreemption(),
				DequeWithSizeLimit(test.expectedCapacity),
			)
			assert.Equal(t, 0, d.Size())
			assert.Equal(t, test.expectedCapacity, d.Capacity())

			arr := make([]interface{}, len(test.items))
			for i := 0; i < len(test.items); i++ {
				arr[i] = test.items[i]
			}
			val, err := invoke(d, test.method, arr...)
			if err != nil {
				t.Fatal(err)
			}
			assert.True(t, val.IsZero()) // means no error

			v, err := d.PeekTail()
			assert.NoError(t, err)
			assert.Equal(t, test.expectedTail, v)

			v, err = d.PeekHead()
			assert.NoError(t, err)
			assert.Equal(t, test.expectedHead, v)
		})
	}

}
func invoke(any interface{}, name string, args ...interface{}) (reflect.Value, error) {
	method := reflect.ValueOf(any).MethodByName(name)
	methodType := method.Type()
	numIn := methodType.NumIn()
	if numIn > len(args) {
		return reflect.ValueOf(nil), fmt.Errorf("method %s must have minimum %d params; has %d", name, numIn, len(args))
	}
	if numIn != len(args) && !methodType.IsVariadic() {
		return reflect.ValueOf(nil), fmt.Errorf("method %s must have %d params; has %d", name, numIn, len(args))
	}
	in := make([]reflect.Value, len(args))
	for i := 0; i < len(args); i++ {
		var inType reflect.Type
		if methodType.IsVariadic() && i >= numIn-1 {
			inType = methodType.In(numIn - 1).Elem()
		} else {
			inType = methodType.In(i)
		}
		argValue := reflect.ValueOf(args[i])
		if !argValue.IsValid() {
			return reflect.ValueOf(nil), fmt.Errorf("method %s. Param[%d] must be %s; has %s", name, i, inType, argValue.String())
		}
		argType := argValue.Type()
		if argType.ConvertibleTo(inType) {
			in[i] = argValue.Convert(inType)
		} else {
			return reflect.ValueOf(nil), fmt.Errorf("method %s. Param[%d] must be %s; has %s", name, i, inType, argType)
		}
	}
	return method.Call(in)[0], nil
}

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

func TestDequeExpandShrink(t *testing.T) {

	d := NewDeque[int](
		DequeWithSizeLimit(8),
	)
	assert.Equal(t, 0, d.Size())
	assert.Equal(t, 8, d.Capacity())

	d.PushHeadAll(1, 2, 3, 4, 5, 6)
	assert.Equal(t, 6, d.Size())
	assert.Equal(t, 8, d.Capacity())

	v, err := d.PeekTail()
	assert.NoError(t, err)
	assert.Equal(t, 1, v)

	v, err = d.PeekHead()
	assert.NoError(t, err)
	assert.Equal(t, 6, v)

	err = d.Expand()
	assert.NoError(t, err)

	assert.Equal(t, 6, d.Size())
	assert.Equal(t, 10, d.Capacity())

	v, err = d.PeekTail()
	assert.NoError(t, err)
	assert.Equal(t, 1, v)

	v, err = d.PeekHead()
	assert.NoError(t, err)
	assert.Equal(t, 6, v)

	d.PushHeadAll(7, 8, 9, 10)

	v, err = d.PeekTail()
	assert.NoError(t, err)
	assert.Equal(t, 1, v)

	v, err = d.PeekHead()
	assert.NoError(t, err)
	assert.Equal(t, 10, v)

	v, err = d.PopHead() // 10
	v, err = d.PopHead() // 9
	v, err = d.PopHead() // 8
	v, err = d.PopTail() // 1
	v, err = d.PopTail() // 2
	v, err = d.PopTail() // 3
	v, err = d.PopTail() // 4

	err = d.Shrink()
	assert.NoError(t, err)
	assert.Equal(t, 3, d.Size())
	assert.Equal(t, 8, d.Capacity())

	v, err = d.PeekTail()
	assert.NoError(t, err)
	assert.Equal(t, 5, v)

	v, err = d.PeekHead()
	assert.NoError(t, err)
	assert.Equal(t, 7, v)

}
