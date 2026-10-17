package convutil

func FormatIntTypeMap[T1 ~int32 | ~int64, T2 int | int8 | int32 | int64 | uint | uint8 | uint32 | uint64](src map[T1]T2) map[int32]int64 {
	var res = make(map[int32]int64, len(src))
	for k, v := range src {
		res[int32(k)] = int64(v)
	}
	return res
}

func ParseIntTypeMap[T1 ~int32 | ~int64, T2 int | int8 | int32 | int64 | uint | uint8 | uint32 | uint64](src map[int32]int64) map[T1]T2 {
	var res = make(map[T1]T2, len(src))
	for k, v := range src {
		res[T1(k)] = T2(v)
	}
	return res
}
