package weight

import (
	"strconv"
	"strings"

	"github.com/NumberMan1/numbox/utils/convutil"
	"github.com/NumberMan1/numbox/utils/typeutil"
)

func ParseWeightItems[T1 Item, T2 typeutil.Number](configString string, fn func(T2, int) T1, defaultWeightArg ...int) (res []T1) {
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

func ParseWeightMultiItems[T Item, T2 typeutil.Number](configString string, fn func([]T2, int) T, defaultWeightArg ...int) (res []T) {
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

func ParseWeightMultiOptions[T Item, T2 typeutil.Number](configString string, fn func([]T2, int) T, defaultWeightArg ...int) (res []T) {
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

func ParseFixedTotalWeightPool[T1 Item, T2 typeutil.Number](src string, fn func(t T2, i int) T1) (totalWeight int, items []T1) {
	strArr := strings.Split(src, ",")
	if len(strArr) < 2 {
		return
	}
	totalWeight = convutil.ParseIntString[int](strArr[0])
	items = ParseWeightItems[T1, T2](strArr[1], fn)
	return
}
