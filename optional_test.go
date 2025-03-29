package container

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOptionalFloatGetSet(t *testing.T) {
	opt := NewFloatOptional()
	assert.True(t, opt.isEmpty)
	assert.Equal(t, 0.0, opt.Get())
	opt.Set(1.0)
	assert.False(t, opt.isEmpty)
	assert.Equal(t, 1.0, opt.Get())
}
func TestOptionalIntGetSet(t *testing.T) {
	opt := NewIntOptional(2)
	assert.False(t, opt.isEmpty)
	assert.Equal(t, 2, opt.Get())
	opt.Set(1)
	assert.False(t, opt.isEmpty)
	assert.Equal(t, 1, opt.Get())
}
func TestOptionalGenericGetSetV1(t *testing.T) {
	opt := NewOptional[int64](2)
	assert.False(t, opt.isEmpty)
	v := any(opt.Get())
	switch v.(type) {
	case int64:
	default:
		assert.Fail(t, "value is not int64")
	}
	assert.Equal(t, int64(2), opt.Get())
	opt.Set(1)
	assert.False(t, opt.isEmpty)
	assert.IsType(t, int64(1), any(opt.Get()))
	assert.Equal(t, int64(1), opt.Get())
}
