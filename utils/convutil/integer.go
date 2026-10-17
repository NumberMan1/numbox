package convutil

import (
	"errors"
	"strconv"
	"strings"

	"github.com/NumberMan1/numbox/utils/typeutil"
)

// ParseIntTypeString parses a signed integer, including named int32/int64 types.
func ParseIntTypeString[T ~int32 | ~int64](str string) T {
	return T(ParseIntString[int64](str))
}

func ParseIntStringWithReceiver[T typeutil.Integer](str string, val *T) {
	*val = ParseIntString[T](str)
}

// ParseIntString preserves loose parsing: errors are ignored and narrow values
// may truncate. Use ParseInteger to validate input and the destination range.
func ParseIntString[T typeutil.Integer](str string) T {
	var zero T
	// ^0 is positive only for unsigned types, including named integer types.
	if ^zero > zero {
		value, _ := strconv.ParseUint(str, 10, 64)
		return T(value)
	}
	value, _ := strconv.ParseInt(str, 10, 64)
	return T(value)
}

func ParseIntStrings[T typeutil.Integer](strVals []string) []T {
	res := make([]T, len(strVals))
	for i, strVal := range strVals {
		res[i] = ParseIntString[T](strVal)
	}
	return res
}

// ParseInteger parses a base-10 integer, rejecting malformed input and values
// outside T's range. Whitespace is not trimmed. Errors return T's zero value.
func ParseInteger[T typeutil.Integer](str string) (T, error) {
	var zero T
	if ^zero > zero {
		value, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			return zero, err
		}
		result := T(value)
		if uint64(result) != value {
			return zero, &strconv.NumError{Func: "ParseInteger", Num: str, Err: strconv.ErrRange}
		}
		return result, nil
	}
	value, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return zero, err
	}
	result := T(value)
	if int64(result) != value {
		return zero, &strconv.NumError{Func: "ParseInteger", Num: str, Err: strconv.ErrRange}
	}
	return result, nil
}

// ParseIntegerSlice parses every separated element strictly. Empty input returns
// an empty non-nil slice. Empty elements/separators are errors; errors return nil.
func ParseIntegerSlice[T typeutil.Integer](src, separator string) ([]T, error) {
	if separator == "" {
		return nil, errors.New("utils: integer separator must not be empty")
	}
	if src == "" {
		return []T{}, nil
	}
	result := make([]T, 0, strings.Count(src, separator)+1)
	for part := range strings.SplitSeq(src, separator) {
		value, err := ParseInteger[T](part)
		if err != nil {
			return nil, err
		}
		result = append(result, value)
	}
	return result, nil
}

func FormatIntString[T typeutil.Integer](intVal T) string {
	if intVal < 0 {
		return strconv.FormatInt(int64(intVal), 10)
	}
	return strconv.FormatUint(uint64(intVal), 10)
}

// FormatIntTypeString formats a signed integer, including named int32/int64 types.
func FormatIntTypeString[T ~int32 | ~int64](src T) string {
	return FormatIntString(int64(src))
}

func FormatIntStrings[T typeutil.Integer](intVals []T) []string {
	res := make([]string, len(intVals))
	for i, intVal := range intVals {
		res[i] = FormatIntString(intVal)
	}
	return res
}

func ParseIntSliceString[T typeutil.Integer](src string, splitStrArg ...string) []T {
	if src == "" {
		return []T{}
	}
	var splitStr = ","
	if len(splitStrArg) > 0 {
		splitStr = splitStrArg[0]
	}
	resArr := make([]T, 0, strings.Count(src, splitStr)+1)
	for str := range strings.SplitSeq(src, splitStr) {
		resArr = append(resArr, ParseIntString[T](str))
	}
	return resArr
}

func FormatIntSliceString[T typeutil.Integer](srcArr []T) string {
	if len(srcArr) == 0 {
		return ""
	}
	var builder strings.Builder
	builder.Grow(len(srcArr))
	var buffer [20]byte // max decimal width of uint64 and int64
	for i, value := range srcArr {
		if i > 0 {
			builder.WriteByte(',')
		}
		var digits []byte
		if value < 0 {
			digits = strconv.AppendInt(buffer[:0], int64(value), 10)
		} else {
			digits = strconv.AppendUint(buffer[:0], uint64(value), 10)
		}
		builder.Write(digits)
	}
	return builder.String()
}

type stringType interface {
	ToString() string
}

func FormatAnySliceToString[T stringType](srcArr []T) string {
	var res = make([]string, 0, len(srcArr))
	for _, str := range srcArr {
		res = append(res, str.ToString())
	}
	return strings.Join(res, ",")
}

func FormatIntSliceToStringSlice[T typeutil.Integer](srcArr []T) []string {
	// Preserve the old nil result for empty input.
	if len(srcArr) == 0 {
		return nil
	}
	strArr := make([]string, len(srcArr))
	for i, src := range srcArr {
		strArr[i] = FormatIntString(src)
	}
	return strArr
}

func ParseStringSliceToIntSlice[T typeutil.Integer](src []string) []T {
	return ParseIntStrings[T](src)
}

func FormatMapKeysToInt[T1 typeutil.Integer, T2 any](src map[string]T2) map[T1]T2 {
	res := make(map[T1]T2, len(src))
	for k, v := range src {
		res[ParseIntString[T1](k)] = v
	}
	return res
}

func FormatToIntSlice[T1 typeutil.Integer, T2 typeutil.Integer](slice []T1) []T2 {
	res := make([]T2, len(slice))
	for i, t := range slice {
		res[i] = T2(t)
	}
	return res
}

func ParseIntValueRange[T typeutil.Integer](src string) (min, max T, ok bool) {
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
