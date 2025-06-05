package container

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRingBufferCreate(t *testing.T) {
	rb := NewRingBuffer[int](10)
	assert.NotNil(t, rb)
	assert.Equal(t, 10, rb.Cap())
}
func TestRingBufferPushPeekPopBasic(t *testing.T) {

	rb := NewRingBuffer[int](10)
	assert.NotNil(t, rb)
	assert.Equal(t, 10, rb.Cap())

	rb.Push(1)
	v, err := rb.Peek()
	assert.NoError(t, err)
	assert.Equal(t, v, 1)
	assert.Equal(t, 1, rb.Len())

	rb.Push(2)
	v, err = rb.Peek()
	assert.NoError(t, err)
	assert.Equal(t, v, 2)
	assert.Equal(t, 2, rb.Len())

	v, err = rb.Pop()
	assert.NoError(t, err)
	assert.Equal(t, v, 2)
	assert.Equal(t, 1, rb.Len())

}
func TestRingBufferPushPeekPopOverflow(t *testing.T) {

	rb := NewRingBuffer[int](3)
	assert.NotNil(t, rb)
	assert.Equal(t, 3, rb.Cap())
	assert.Equal(t, 0, rb.Len())

	rb.Push(1)
	v, err := rb.Peek()
	assert.NoError(t, err)
	assert.Equal(t, 1, v)
	assert.Equal(t, 1, rb.Len())

	rb.Push(2)
	v, err = rb.Peek()
	assert.NoError(t, err)
	assert.Equal(t, 2, v)
	assert.Equal(t, 2, rb.Len())

	rb.Push(3)
	v, err = rb.Peek()
	assert.NoError(t, err)
	assert.Equal(t, 3, v)
	assert.Equal(t, 3, rb.Len())

	rb.Push(4)
	v, err = rb.Peek()
	assert.NoError(t, err)
	assert.Equal(t, 4, v)
	assert.Equal(t, 3, rb.Len())

	rb.Push(5)
	v, err = rb.Peek()
	assert.NoError(t, err)
	assert.Equal(t, 5, v)
	assert.Equal(t, 3, rb.Len())

	rb.Push(6)
	v, err = rb.Peek()
	assert.NoError(t, err)
	assert.Equal(t, 6, v)
	assert.Equal(t, 3, rb.Len())

	rb.Push(7)
	v, err = rb.Peek()
	assert.NoError(t, err)
	assert.Equal(t, 7, v)
	assert.Equal(t, 3, rb.Len())

	for _, v := range rb.Values() {
		fmt.Println(v)
	}

	v, err = rb.Pop()
	assert.NoError(t, err)
	assert.Equal(t, 7, v)
	assert.Equal(t, 2, rb.Len())

	v, err = rb.Pop()
	assert.NoError(t, err)
	assert.Equal(t, 6, v)
	assert.Equal(t, 1, rb.Len())

	v, err = rb.Pop()
	assert.NoError(t, err)
	assert.Equal(t, 5, v)
	assert.Equal(t, 0, rb.Len())

}
