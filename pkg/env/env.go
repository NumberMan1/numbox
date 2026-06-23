package env

import (
	"os"
	"strconv"
	"strings"
)

type EnvValue string

func (e EnvValue) String() string {
	return string(e)
}

func (e EnvValue) Bool() bool {
	return strings.ToLower(e.String()) == "true" || e.Int() == 1
}

func (e EnvValue) Int() int {
	if e == "" {
		return 0
	}
	i, err := strconv.Atoi(e.String())
	if err != nil {
		panic(err)
	}
	return i
}

func GetEnv(key string, def ...EnvValue) EnvValue {
	val := os.Getenv(key)
	if val == "" && len(def) != 0 {
		return def[0]
	}
	return EnvValue(val)
}
