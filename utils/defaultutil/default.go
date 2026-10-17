package defaultutil

func GetParamWithDefaultValue[T1 ~int | ~int32 | ~int64, T2 ~int64 | ~int32](defVal T1, args ...T2) T1 {
	if len(args) <= 0 {
		return defVal
	}
	return T1(args[0])
}

func GetBoolParamWithDefault(defVal bool, args ...bool) bool {
	if len(args) <= 0 {
		return defVal
	}
	return args[0]
}

func GetStringParamWithDefault(defVal string, args ...string) string {
	if len(args) <= 0 {
		return defVal
	}
	return args[0]
}
