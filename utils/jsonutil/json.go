package jsonutil

import (
	"github.com/NumberMan1/numbox/utils/assertutil"
	"github.com/json-iterator/go"
)

// ToJSONString marshals v as JSON, returning an empty string for nil.
// A marshal error panics, matching the existing JSON helpers.
func ToJSONString(v interface{}) string {
	if v == nil {
		return ""
	}
	jsonBytes := ToJSONBytes(v)
	return string(jsonBytes)
}

// ToJSONBytes marshals v as JSON and panics on a marshal error.
func ToJSONBytes(v interface{}) []byte {
	jsonBytes, err := jsoniter.ConfigFastest.Marshal(v)
	if err != nil {
		panic(err)
	}
	return jsonBytes
}

func GetInt32SliceFromJSON(str string) []int32 {
	res := make([]int32, 0)
	if str == "" {
		return res
	}
	assertutil.Must(jsoniter.ConfigFastest.Unmarshal([]byte(str), &res))
	return res
}

func GetIntSliceFromJSON[T int | int8 | int32 | int64](str string) []T {
	res := make([]T, 0)
	if str == "" {
		return res
	}
	assertutil.Must(jsoniter.ConfigFastest.Unmarshal([]byte(str), &res))
	return res
}

func GetMapFromJSON[T int | int32 | int64](str string, res *map[T]T) {
	if str == "" {
		return
	}
	assertutil.Must(jsoniter.ConfigFastest.Unmarshal([]byte(str), res))
}

func UnmarshalFromJSON(src string, data interface{}) {
	assertutil.Must(jsoniter.ConfigFastest.Unmarshal([]byte(src), data))
}
