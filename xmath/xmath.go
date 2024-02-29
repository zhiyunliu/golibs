package xmath

func AbsInt(a int64) int64 {
	if a < 0 {
		return -1 * a
	}
	return a
}

func Sum[T ~int | ~int32 | ~int64 | ~uint | ~uint32 | ~uint64 | ~float32 | ~float64](val ...T) T {
	var result T
	for i := range val {
		result += val[i]
	}
	return result
}

func Min[T ~int | ~int32 | ~int64 | ~uint | ~uint32 | ~uint64 | ~float32 | ~float64](val ...T) T {
	var result T
	cnt := len(val)
	if cnt == 0 {
		return result
	}
	result = val[0]
	cnt = cnt - 1
	for i := 1; i < cnt; i++ {
		if val[i] < result {
			result = val[i]
		}
	}
	return result
}

func Max[T int | ~int32 | ~int64 | ~uint | ~uint32 | ~uint64 | ~float32 | ~float64](val ...T) T {
	var result T
	cnt := len(val)
	if cnt == 0 {
		return result
	}
	result = val[0]
	cnt = cnt - 1
	for i := 1; i < cnt; i++ {
		if val[i] > result {
			result = val[i]
		}
	}
	return result
}
