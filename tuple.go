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

//func (t Tuple[T1, T2]) FirstAsString() string {
//	var x any = t.First
//	return x.(string)
//}
//func (t Tuple[T1, T2]) SecondAsString() string {
//	var x any = t.Second
//	return x.(string)
//}

type IntTuple Tuple[int, int]
type Int8Tuple Tuple[int8, int8]
type Int16Tuple Tuple[int16, int16]
type Int32Tuple Tuple[int32, int32]
type Int64Tuple Tuple[int64, int64]

type UintTuple Tuple[uint, uint]
type Uint8Tuple Tuple[uint8, uint8]
type Uint16Tuple Tuple[uint16, uint16]
type Uint32Tuple Tuple[uint32, uint32]
type Uint64Tuple Tuple[uint64, uint64]

type Float32Tuple Tuple[float32, float32]
type Float64Tuple Tuple[float64, float64]

type StringTuple Tuple[string, string]
