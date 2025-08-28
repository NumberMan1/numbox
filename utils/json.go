package utils

import (
	jsoniter "github.com/json-iterator/go"
)

func ToJsonString(v interface{}) string {
	if v == nil {
		return ""
	}
	jsonBytes := ToJsonBytes(v)
	return string(jsonBytes)
}

func ToJsonBytes(v interface{}) []byte {
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
	Must(jsoniter.ConfigFastest.Unmarshal([]byte(str), &res))
	return res
}

func GetIntSliceFromJSON[T int | int8 | int32 | int64](str string) []T {
	res := make([]T, 0)
	if str == "" {
		return res
	}
	Must(jsoniter.ConfigFastest.Unmarshal([]byte(str), &res))
	return res
}

func GetMapFromJSON[T int | int32 | int64](str string, res *map[T]T) {
	if str == "" {
		return
	}
	Must(jsoniter.ConfigFastest.Unmarshal([]byte(str), res))
}

func UnmarshalFromJSON(src string, data interface{}) {
	Must(jsoniter.ConfigFastest.Unmarshal([]byte(src), data))
}
