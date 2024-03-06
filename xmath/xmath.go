package xmath

import "cmp"

func AbsInt(a int64) int64 {
	return Abs(a)
}

func Abs[T ~int | ~int8 | ~int16 | ~int32 | ~int64 |
	~float32 | ~float64](val T) T {
	if val < 0 {
		return T(-1) * val
	}
	return val
}

func Sum[T ~int | ~int8 | ~int16 | ~int32 | ~int64 |
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
	~float32 | ~float64](val ...T) T {
	var result T
	for i := range val {
		result += val[i]
	}
	return result
}

func Min[T cmp.Ordered](val ...T) T {
	var result T
	cnt := len(val)
	if cnt == 0 {
		return result
	}
	result = val[0]
	for i := 1; i < cnt; i++ {
		if val[i] < result {
			result = val[i]
		}
	}
	return result
}

func Max[T cmp.Ordered](val ...T) T {
	var result T
	cnt := len(val)
	if cnt == 0 {
		return result
	}
	result = val[0]
	for i := 1; i < cnt; i++ {
		if val[i] > result {
			result = val[i]
		}
	}
	return result
}
