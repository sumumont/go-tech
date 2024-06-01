package common

func Min[T int | int32 | int64 | float64 | float32](x, y T) T {
	if x < y {
		return x
	}
	return y
}

func Max[T int | int32 | int64 | float64 | float32](x, y T) T {
	if x > y {
		return x
	}
	return y
}
