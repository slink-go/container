package container

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOptionalFloatGetSet(t *testing.T) {
	opt := NewFloatOptional()
	assert.True(t, opt.IsEmpty)
	assert.Equal(t, 0.0, opt.Get())
	opt.Set(1.0)
	assert.False(t, opt.IsEmpty)
	assert.Equal(t, 1.0, opt.Get())
}
func TestOptionalIntGetSet(t *testing.T) {
	opt := NewIntOptional(2)
	assert.False(t, opt.IsEmpty)
	assert.Equal(t, 2, opt.Get())
	opt.Set(1)
	assert.False(t, opt.IsEmpty)
	assert.Equal(t, 1, opt.Get())
}
func TestOptionalGenericGetSet(t *testing.T) {
	opt := NewOptional[int64](2)
	assert.False(t, opt.IsEmpty)
	v := any(opt.Get())
	switch v.(type) {
	case int64:
	default:
		assert.Fail(t, "Optional is not int64")
	}
	assert.Equal(t, int64(2), opt.Get())
	opt.Set(1)
	assert.False(t, opt.IsEmpty)
	assert.IsType(t, int64(1), any(opt.Get()))
	assert.Equal(t, int64(1), opt.Get())
}
