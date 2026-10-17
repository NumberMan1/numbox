package convutil

import (
	"strings"
)

func ParseItemsString[T any](itemsStr string, splitStr string, fn func(s string) (T, bool)) (res []T) {
	res = make([]T, 0)
	if itemsStr == "" {
		return
	}

	var arr = strings.Split(itemsStr, splitStr)
	for _, item := range arr {
		item = strings.TrimSpace(item)
		val, ok := fn(item)
		if !ok {
			continue
		}
		res = append(res, val)
	}
	return
}
