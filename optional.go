package container

import (
	"fmt"
	"time"
)

type Value[T any] struct {
	value   T
	isEmpty bool
}

func (o *Value[T]) Get() T {
	return o.value
}
func (o *Value[T]) Set(value T) {
	o.value = value
	o.isEmpty = false
}
func (o *Value[T]) Empty() bool {
	return o.isEmpty
}
func (o *Value[T]) OrElse(value T) T {
	if o.isEmpty {
		return value
	}
	return o.value
}
func (o *Value[T]) OrElseString(value string) string {
	if o.isEmpty {
		return value
	}
	return fmt.Sprintf("%v", o.value)
}
func (o *Value[T]) OrElseFormatted(format, value string) string {
	if o.isEmpty {
		return value
	}
	var v interface{}
	v = o.value
	switch x := v.(type) {
	case time.Time:
		return x.Format(format)
	default:
		return fmt.Sprintf(format, o.value)
	}

}

func NewFloatOptional(value ...float64) Value[float64] {
	if len(value) == 0 {
		return Value[float64]{
			value:   0,
			isEmpty: true,
		}
	}
	return Value[float64]{
		value:   value[0],
		isEmpty: false,
	}
}
func NewIntOptional(value ...int) Value[int] {
	if len(value) == 0 {
		return Value[int]{
			value:   0,
			isEmpty: true,
		}
	}
	return Value[int]{
		value:   value[0],
		isEmpty: false,
	}
}
func NewDateOptional(value ...time.Time) Value[time.Time] {
	if len(value) == 0 {
		return Value[time.Time]{
			value:   time.Time{},
			isEmpty: true,
		}
	}
	return Value[time.Time]{
		value:   value[0],
		isEmpty: false,
	}
}
func NewOptional[T any](value ...T) Value[T] {
	if len(value) == 0 {
		return Value[T]{
			isEmpty: true,
		}
	}
	return Value[T]{
		value:   value[0],
		isEmpty: false,
	}
}
