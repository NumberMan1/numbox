package utils

import (
	"strconv"
	"strings"

	"github.com/NumberMan1/numbox/pkg/weight"
)

func ParseWeightItems[T1 weight.Item, T2 Numeric](configString string, fn func(T2, int) T1, defaultWeightArg ...int) (res []T1) {
	var (
		arr           = strings.Split(configString, "|")
		defaultWeight = 100
	)
	if len(defaultWeightArg) > 0 {
		defaultWeight = defaultWeightArg[0]
	}
	for _, item := range arr {
		item = strings.TrimSpace(item)
		itemArr := strings.Split(item, ":")
		id, _ := strconv.Atoi(itemArr[0])
		if id == 0 {
			continue
		}
		if len(itemArr) == 1 {
			res = append(res, fn(T2(id), defaultWeight))
			continue
		}
		weightVal, _ := strconv.Atoi(itemArr[1])
		res = append(res, fn(T2(id), weightVal))
	}
	return
}

func ParseWeightMultiItems[T weight.Item, T2 Numeric](configString string, fn func([]T2, int) T, defaultWeightArg ...int) (res []T) {
	var (
		arr           = strings.Split(configString, "|")
		defaultWeight = 100
	)
	if len(defaultWeightArg) > 0 {
		defaultWeight = defaultWeightArg[0]
	}
	for _, item := range arr {
		item = strings.TrimSpace(item)
		itemArr := strings.Split(item, ":")
		var ids []T2
		idsArr := strings.Split(itemArr[0], "&")
		for _, idStr := range idsArr {
			id, _ := strconv.Atoi(idStr)
			if id == 0 {
				continue
			}
			ids = append(ids, T2(id))
		}
		if len(itemArr) == 1 {
			res = append(res, fn(ids, defaultWeight))
			continue
		}
		weightVal, _ := strconv.Atoi(itemArr[1])
		res = append(res, fn(ids, weightVal))
	}
	return
}

func ParseWeightMultiOptions[T weight.Item, T2 Numeric](configString string, fn func([]T2, int) T, defaultWeightArg ...int) (res []T) {
	var (
		arr           = strings.Split(configString, "|")
		defaultWeight = 100
	)
	if len(defaultWeightArg) > 0 {
		defaultWeight = defaultWeightArg[0]
	}
	for _, item := range arr {
		item = strings.TrimSpace(item)
		itemArr := strings.Split(item, ":")
		var ids []T2
		idsArr := strings.Split(itemArr[0], "/")
		for _, idStr := range idsArr {
			id, _ := strconv.Atoi(idStr)
			if id == 0 {
				continue
			}
			ids = append(ids, T2(id))
		}
		if len(itemArr) == 1 {
			res = append(res, fn(ids, defaultWeight))
			continue
		}
		weightVal, _ := strconv.Atoi(itemArr[1])
		res = append(res, fn(ids, weightVal))
	}
	return
}

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

func ParseBoolFromInt[T int | int8 | int32 | int64](intVal T) bool {
	return intVal > 0
}

func ParseFixedTotalWeightPool[T1 weight.Item, T2 Numeric](src string, fn func(t T2, i int) T1) (totalWeight int, items []T1) {
	strArr := strings.Split(src, ",")
	if len(strArr) < 2 {
		return
	}
	totalWeight = ParseIntString[int](strArr[0])
	items = ParseWeightItems[T1, T2](strArr[1], fn)
	return
}
