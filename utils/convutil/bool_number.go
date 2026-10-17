package convutil

func ParseBoolFromInt[T int | int8 | int32 | int64](intVal T) bool {
	return intVal > 0
}
