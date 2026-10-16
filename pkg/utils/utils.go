package utils

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/exp/rand"
)

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

func Assert(isOk bool, err error) {
	if !isOk {
		panic(err)
	}
}

func AssertStringError(isOk bool, str string) {
	if !isOk {
		panic(fmt.Errorf("%s", str))
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

// durationMap 定义中文时间单位对应的时长。
var durationMap = map[string]time.Duration{
	"毫秒": time.Millisecond,
	"秒钟": time.Second,
	"秒":  time.Second,
	"分钟": time.Minute,
	"分":  time.Minute,
	"小时": time.Hour,
	"天":  24 * time.Hour,
	"周":  7 * 24 * time.Hour,
}

var cnDurationPattern = regexp.MustCompile(`(\d+)(毫秒|秒钟|秒|分钟|分|小时|天|周)`)

// ParseCNDuration parses complete, non-negative Chinese duration expressions.
// Unknown text, partial matches and time.Duration overflow are errors.
func ParseCNDuration(input string) (time.Duration, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return 0, fmt.Errorf("empty duration")
	}
	matches := cnDurationPattern.FindAllStringSubmatchIndex(input, -1)
	var total time.Duration
	end := 0
	for _, match := range matches {
		if strings.TrimSpace(input[end:match[0]]) != "" {
			return 0, fmt.Errorf("invalid duration: %q", input)
		}
		value, err := strconv.ParseUint(input[match[2]:match[3]], 10, 64)
		if err != nil {
			return 0, err
		}
		unit := durationMap[input[match[4]:match[5]]]
		if value > uint64((time.Duration(1<<63-1)-total)/unit) {
			return 0, fmt.Errorf("duration overflow: %q", input)
		}
		total += time.Duration(value) * unit
		end = match[1]
	}
	if len(matches) == 0 || strings.TrimSpace(input[end:]) != "" {
		return 0, fmt.Errorf("invalid duration: %q", input)
	}
	return total, nil
}

func ParsePercentageToPermille(input string) (int, error) {
	// 去掉空格并检查是否包含百分号
	input = strings.TrimSpace(input)
	if !strings.HasSuffix(input, "%") {
		return 0, fmt.Errorf("输入格式错误，必须包含百分号")
	}

	// 去掉百分号并解析浮点数
	numStr := strings.TrimSuffix(input, "%")
	value, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return 0, err
	}

	// 计算千分比并四舍五入
	permille := math.Round(value * 10)
	if math.IsNaN(permille) || math.IsInf(permille, 0) || permille >= float64(int(^uint(0)>>1)) || permille < -float64(int(^uint(0)>>1))-1 {
		return 0, fmt.Errorf("percentage outside integer range")
	}
	return int(permille), nil
}

func ShuffleRandom[T any](slice []T) {
	rand.New(rand.NewSource(uint64(time.Now().UnixNano()))).Shuffle(len(slice), func(i, j int) {
		slice[i], slice[j] = slice[j], slice[i]
	})
}

func Shuffle[T any](seed int64, slice []T) {
	rand.New(rand.NewSource(uint64(seed))).Shuffle(len(slice), func(i, j int) {
		slice[i], slice[j] = slice[j], slice[i]
	})
}

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
