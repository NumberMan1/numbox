package utils

import "fmt"

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

func Asset(isOk bool, err error) {
	if !isOk {
		panic(err)
	}
}

func WithRecover(work func() error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Recovered from panic:", err)
		}
	}()

	err := work()
	if err != nil {
		fmt.Println(err)
	}
}

func ParseBoolString(src string) bool {
	switch src {
	case "t":
		return true
	default:
		return false
	}
}

func FormatBoolToString(ok bool) string {
	if ok {
		return "t"
	}
	return "f"
}

func FormatBoolToCNStr(ok bool) string {
	if ok {
		return "是"
	}
	return "否"
}

func FormatBoolToInt[T int | int8 | int32 | int64 | uint | uint8 | uint32 | uint64](ok bool) T {
	if ok {
		return T(1)
	}
	return T(0)
}
