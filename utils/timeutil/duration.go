package timeutil

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

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
