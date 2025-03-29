package container

import (
	"fmt"
	"time"
)

type Value[T any] struct {
	Value   T
	IsEmpty bool
}

func (o *Value[T]) Get() T {
	return o.Value
}
func (o *Value[T]) Set(value T) {
	o.Value = value
	o.IsEmpty = false
}
func (o *Value[T]) Empty() bool {
	return o.IsEmpty
}
func (o *Value[T]) OrElse(value T) T {
	if o.IsEmpty {
		return value
	}
	return o.Value
}
func (o *Value[T]) OrElseString(value string) string {
	if o.IsEmpty {
		return value
	}
	return fmt.Sprintf("%v", o.Value)
}
func (o *Value[T]) OrElseFormatted(format, value string) string {
	if o.IsEmpty {
		return value
	}
	var v interface{}
	v = o.Value
	switch x := v.(type) {
	case time.Time:
		return x.Format(format)
	default:
		return fmt.Sprintf(format, o.Value)
	}

}

func NewFloatOptional(value ...float64) Value[float64] {
	if len(value) == 0 {
		return Value[float64]{
			Value:   0,
			IsEmpty: true,
		}
	}
	return Value[float64]{
		Value:   value[0],
		IsEmpty: false,
	}
}
func NewIntOptional(value ...int) Value[int] {
	if len(value) == 0 {
		return Value[int]{
			Value:   0,
			IsEmpty: true,
		}
	}
	return Value[int]{
		Value:   value[0],
		IsEmpty: false,
	}
}
func NewDateOptional(value ...time.Time) Value[time.Time] {
	if len(value) == 0 {
		return Value[time.Time]{
			Value:   time.Time{},
			IsEmpty: true,
		}
	}
	return Value[time.Time]{
		Value:   value[0],
		IsEmpty: false,
	}
}
func NewOptional[T any](value ...T) Value[T] {
	if len(value) == 0 {
		return Value[T]{
			IsEmpty: true,
		}
	}
	return Value[T]{
		Value:   value[0],
		IsEmpty: false,
	}
}
