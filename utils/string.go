package utils

import (
	"strconv"
	"strings"
	"unicode"
)

func ParseIntStringWithReceiver[T int | int8 | int32 | int64 | uint | uint8 | uint32 | uint64](str string, val *T) {
	inf, _ := strconv.ParseInt(str, 10, 64)
	*val = T(inf)
}

func ParseIntString[T int | int8 | int32 | int64 | uint | uint8 | uint32 | uint64](str string) T {
	inf, _ := strconv.ParseInt(str, 10, 64)
	return T(inf)
}

func GetNameStringLength(nameStr string) int {
	var strLength int
	for _, r := range nameStr {
		if unicode.Is(unicode.Han, r) { // 检查是否为中文字符
			strLength += 2
		} else if unicode.IsLetter(r) || unicode.IsDigit(r) { // 检查是否为英文字符或数字
			strLength += 1
		}
	}
	return strLength
}

func FormatIntString[T int | int8 | int32 | int64 | uint | uint8 | uint32 | uint64](intVal T) string {
	str := strconv.FormatInt(int64(intVal), 10)
	return str
}

func ParseIntSliceString[T int | int8 | int32 | int64 | uint | uint8 | uint32 | uint64](src string) []T {
	strArr := strings.Split(src, ",")
	resArr := make([]T, 0, len(strArr))
	for _, str := range strArr {
		resArr = append(resArr, ParseIntString[T](str))
	}
	return resArr
}

func FormatIntSliceString[T int | int8 | int32 | int64 | uint | uint8 | uint32 | uint64](srcArr []T) string {
	return strings.Join(FormatIntSliceToStringSlice(srcArr), ",")
}

func FormatIntSliceToStringSlice[T int | int8 | int32 | int64 | uint | uint8 | uint32 | uint64](srcArr []T) []string {
	var strArr []string
	for _, src := range srcArr {
		strArr = append(strArr, FormatIntString(src))
	}
	return strArr
}

func ParseStringSliceToIntSlice[T int | int8 | int32 | int64 | uint | uint8 | uint32 | uint64](src []string) []T {
	resArr := make([]T, 0, len(src))
	for _, str := range src {
		var v T
		ParseIntStringWithReceiver(str, &v)
		resArr = append(resArr, v)
	}
	return resArr
}

func FormatMapKeysToInt[T1 int | int8 | int32 | int64, T2 any](src map[string]T2) map[T1]T2 {
	res := make(map[T1]T2, len(src))
	for k, v := range src {
		res[ParseIntString[T1](k)] = v
	}
	return res
}
