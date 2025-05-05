package container

import (
	"encoding/json"
	"fmt"
	"go.slink.ws/types"
	"time"
)

type Optional[T comparable] struct {
	Value   T
	IsEmpty bool
}

func (o *Optional[T]) Equals(other Optional[T]) bool {
	return o.IsEmpty && other.IsEmpty ||
		o.Value == other.Value
}

func (o *Optional[T]) MarshalJSON() ([]byte, error) {
	if o.IsEmpty {
		return []byte("null"), nil
	}
	return json.Marshal(o.Value)
}

func (o *Optional[T]) Get() T {
	return o.Value
}
func (o *Optional[T]) Set(value T) {
	o.Value = value
	o.IsEmpty = false
}
func (o *Optional[T]) Empty() bool {
	return o.IsEmpty
}
func (o *Optional[T]) OrElse(value T) T {
	if o.IsEmpty {
		return value
	}
	return o.Value
}
func (o *Optional[T]) OrElseString(value string) string {
	if o.IsEmpty {
		return value
	}
	return fmt.Sprintf("%v", o.Value)
}
func (o *Optional[T]) OrElseFormatted(format, value string) string {
	if o.IsEmpty {
		return value
	}
	var v interface{}
	v = o.Value
	switch x := v.(type) {
	case time.Time:
		return x.Format(format)
	case types.Date:
		return x.Format(format)
	default:
		return fmt.Sprintf(format, o.Value)
	}

}

func NewFloatOptional(value ...float64) Optional[float64] {
	if len(value) == 0 {
		return Optional[float64]{
			Value:   0,
			IsEmpty: true,
		}
	}
	return Optional[float64]{
		Value:   value[0],
		IsEmpty: false,
	}
}
func NewIntOptional(value ...int) Optional[int] {
	if len(value) == 0 {
		return Optional[int]{
			Value:   0,
			IsEmpty: true,
		}
	}
	return Optional[int]{
		Value:   value[0],
		IsEmpty: false,
	}
}
func NewDateOptional(value ...time.Time) Optional[time.Time] {
	if len(value) == 0 {
		return Optional[time.Time]{
			Value:   time.Time{},
			IsEmpty: true,
		}
	}
	return Optional[time.Time]{
		Value:   value[0],
		IsEmpty: false,
	}
}
func NewOptional[T comparable](value ...T) Optional[T] {
	if len(value) == 0 {
		return Optional[T]{
			IsEmpty: true,
		}
	}
	return Optional[T]{
		Value:   value[0],
		IsEmpty: false,
	}
}
