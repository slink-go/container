package container

type Tuple[T1, T2 any] struct {
	First  T1
	Second T2
}

func NewTuple[T1, T2 any](f T1, s T2) Tuple[T1, T2] {
	return Tuple[T1, T2]{
		First:  f,
		Second: s,
	}
}

func (t Tuple[T1, T2]) GetFirstAsString() string {
	var x any = t.First
	return x.(string)
}
func (t Tuple[T1, T2]) GetSecondAsString() string {
	var x any = t.Second
	return x.(string)
}

type StringTuple Tuple[string, string]

func (t StringTuple) GetFirst() string {
	return t.First
}
func (t StringTuple) GetSecond() string {
	return t.Second
}
