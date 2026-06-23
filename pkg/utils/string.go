package utils

import (
	"strconv"
	"strings"
	"unicode"

	"github.com/NumberMan1/numbox/pkg/percent"
)

func ParseIntTypString[T ~int32 | ~int64](str string) T {
	return T(ParseIntString[int64](str))
}

func ParseIntStringWithReceiver[T ~int | ~int8 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint32 | ~uint64](str string, val *T) {
	inf, _ := strconv.ParseInt(str, 10, 64)
	*val = T(inf)
}

func ParseIntString[T ~int | ~int8 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint32 | ~uint64](str string) T {
	inf, _ := strconv.ParseInt(str, 10, 64)
	return T(inf)
}

func ParseIntStrings[T ~int | ~int8 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint32 | ~uint64](strVals []string) []T {
	var res = make([]T, 0, len(strVals))
	for _, strVal := range strVals {
		res = append(res, ParseIntString[T](strVal))
	}
	return res
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

func FormatIntString[T ~int | ~int8 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint32 | ~uint64](intVal T) string {
	return strconv.FormatInt(int64(intVal), 10)
}

func FormatIntTypString[T ~int32 | ~int64](src T) string {
	return FormatIntString(int64(src))
}

func FormatIntStrings[T ~int | ~int8 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint32 | ~uint64](intVals []T) []string {
	var res = make([]string, 0, len(intVals))
	for _, intVal := range intVals {
		res = append(res, FormatIntString[T](intVal))
	}
	return res
}

func ParseIntSliceString[T ~int | ~int8 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint32 | ~uint64](src string, splitStrArg ...string) []T {
	if src == "" {
		return []T{}
	}
	var splitStr = ","
	if len(splitStrArg) > 0 {
		splitStr = splitStrArg[0]
	}
	strArr := strings.Split(src, splitStr)
	resArr := make([]T, 0, len(strArr))
	for _, str := range strArr {
		resArr = append(resArr, ParseIntString[T](str))
	}
	return resArr
}

func FormatIntSliceString[T ~int | ~int8 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint32 | ~uint64](srcArr []T) string {
	return strings.Join(FormatIntSliceToStringSlice(srcArr), ",")
}

type stringTyp interface {
	ToString() string
}

func FormatAnySliceToString[T stringTyp](srcArr []T) string {
	var res = make([]string, 0, len(srcArr))
	for _, str := range srcArr {
		res = append(res, str.ToString())
	}
	return strings.Join(res, ",")
}

func FormatIntSliceToStringSlice[T ~int | ~int8 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint32 | ~uint64](srcArr []T) []string {
	var strArr []string
	for _, src := range srcArr {
		strArr = append(strArr, FormatIntString(src))
	}
	return strArr
}

func ParseStringSliceToIntSlice[T ~int | ~int8 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint32 | ~uint64](src []string) []T {
	resArr := make([]T, 0, len(src))
	for _, str := range src {
		var v T
		ParseIntStringWithReceiver(str, &v)
		resArr = append(resArr, v)
	}
	return resArr
}

func FormatMapKeysToInt[T1 ~int | ~int8 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint32 | ~uint64, T2 any](src map[string]T2) map[T1]T2 {
	res := make(map[T1]T2, len(src))
	for k, v := range src {
		res[ParseIntString[T1](k)] = v
	}
	return res
}

func FormatPercentString(percent percent.Percent) string {
	return FormatIntString(percent.Int32())
}

func ParsePercentString(src string) percent.Percent {
	return percent.Percent(ParseIntString[int32](src))
}

func FormatToIntSlice[T1 ~int32 | ~int64, T2 ~int | ~int8 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint32 | ~uint64](slice []T1) []T2 {
	res := make([]T2, 0, len(slice))
	for _, t := range slice {
		res = append(res, T2(t))
	}
	return res
}

func ParseIntValueRange[T ~int | ~int8 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint32 | ~uint64](src string) (min, max T, ok bool) {
	if src == "" {
		return
	}
	strArr := strings.Split(src, "~")
	if len(strArr) < 1 {
		return
	}
	min = ParseIntString[T](strArr[0])
	max = min
	if len(strArr) > 1 {
		max = ParseIntString[T](strArr[1])
	}
	ok = true
	return
}
