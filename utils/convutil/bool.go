package convutil

func ParseBoolString(src string) bool {
	switch src {
	case "t", "T", "true", "TRUE", "True":
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

func ParseBoolIntString(str string) bool {
	if str == "" || str == "0" {
		return false
	}
	return true
}

func FormatBoolIntString(ok bool) string {
	if ok {
		return "1"
	}
	return "0"
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
