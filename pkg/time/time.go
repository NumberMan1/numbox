package time

import (
	"time"
)

func init() {
	var err error
	time.Local, err = time.LoadLocation("Asia/Shanghai")
	if err != nil {
		time.Local = time.FixedZone("CST", 8*3600)
	}
}

const (
	Nanosecond  = time.Nanosecond
	Microsecond = time.Microsecond
	Millisecond = time.Millisecond
	Second      = time.Second
	Minute      = time.Minute
	Hour        = time.Hour
	Day         = time.Hour * 24
)

type Duration = time.Duration

type Time = time.Time

func NewTimestamp(t time.Time) Timestamp {
	return Timestamp(t.UnixMilli())
}

func NewTimestampNow() Timestamp {
	return Timestamp(time.Now().UnixMilli())
}

func GetTodayStartEnd(nowArg ...Timestamp) (Timestamp, Timestamp) {
	var now = NewTimestamp(time.Now())
	if len(nowArg) > 0 {
		now = nowArg[0]
	}
	todayStart := GetTodayStart(now)
	todayEnd := todayStart.AddDuration(24 * time.Hour)
	return todayStart, todayEnd
}

func GetTodayStart(nowArg ...Timestamp) Timestamp {
	var now = time.Now()
	if len(nowArg) > 0 {
		now = nowArg[0].Time()
	}
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return NewTimestamp(startOfDay)
}

func GetTodayEnd(nowArg ...Timestamp) Timestamp {
	var now = time.Now()
	if len(nowArg) > 0 {
		now = nowArg[0].Time()
	}
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Add(24 * time.Hour)
	return NewTimestamp(endOfDay)
}

func GetWeekStartEnd(nowArg ...Timestamp) (Timestamp, Timestamp) {
	var now = NewTimestamp(time.Now())
	if len(nowArg) > 0 {
		now = nowArg[0]
	}
	weekStart := GetWeekStart(now)
	weekEnd := weekStart.AddDuration(7 * 24 * time.Hour)
	return weekStart, weekEnd
}

func GetWeekStart(nowArg ...Timestamp) Timestamp {
	var now = time.Now()
	if len(nowArg) > 0 {
		now = nowArg[0].Time()
	}
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	monday := now.AddDate(0, 0, -weekday+1)
	weekStart := time.Date(monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, now.Location())
	return NewTimestamp(weekStart)
}

func GetWeekEnd(nowArg ...Timestamp) Timestamp {
	var now = time.Now()
	if len(nowArg) > 0 {
		now = nowArg[0].Time()
	}
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	monday := now.AddDate(0, 0, -weekday+1)
	weekEnd := time.Date(monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, now.Location()).Add(7 * 24 * time.Hour)
	return NewTimestamp(weekEnd)
}

func NewTimestampWithUnixMilli(unixMilli int64) Timestamp {
	return Timestamp(unixMilli)
}

func NewTimestampWithUnix(unix int64) Timestamp {
	return Timestamp(unix * 1000)
}

type Timestamp int64

func (t Timestamp) Time() time.Time {
	return time.UnixMilli(int64(t))
}

func (t Timestamp) UnixNano() int64 {
	return t.Time().UnixNano()
}

func (t Timestamp) Unix() int64 {
	return t.Time().Unix()
}

func (t Timestamp) Int64() int64 {
	return int64(t)
}

func (t Timestamp) AddDuration(d time.Duration) Timestamp {
	return Timestamp(t.Time().Add(d).UnixMilli())
}

func (t Timestamp) SinceFrom(fromTimestamp Timestamp) time.Duration {
	duration := t - fromTimestamp
	return time.Duration(duration) * time.Millisecond
}

func (t Timestamp) Before(fromTimestamp Timestamp) bool {
	return t.Time().Before(fromTimestamp.Time())
}

func (t Timestamp) After(fromTimestamp Timestamp) bool {
	return t.Time().After(fromTimestamp.Time())
}

func ParseDuration(duration string) (Duration, error) {
	timeDuration, err := time.ParseDuration(duration)
	if err != nil {
		return 0, err
	}
	return timeDuration, nil
}

func ParseDateString(layout string, dateStr string) (Timestamp, error) {
	resTime, err := time.ParseInLocation(layout, dateStr, time.Local)
	if err != nil {
		return 0, err
	}
	return NewTimestamp(resTime), nil
}

func GetZoneOffset() float64 {
	_, zoneOffset := time.Now().Zone()
	return float64(zoneOffset / 3600)
}

func FormatTimeToDateString(dateline Timestamp) string {
	if dateline <= 0 {
		return "无"
	}
	return dateline.Time().Format("2006-01-02 15:04:05")
}

func Since(timestamp Timestamp) Duration {
	return time.Since(timestamp.Time())
}

func SinceDays(timestamp Timestamp) int {
	return int(time.Since(timestamp.Time()).Hours()/24) + 1
}

// DurationToNextInterval returns the time.Duration until the next multiple of the given interval.
// If t is zero value, uses current time. Interval must be >0.
func DurationToNextInterval(timestamp Timestamp, interval time.Duration) time.Duration {
	var t = timestamp.Time()
	if interval <= 0 || t.IsZero() {
		return 0
	}
	// Truncate to the interval and add one interval to get next boundary
	next := t.Truncate(interval).Add(interval)
	return next.Sub(t)
}

func IsToday(checkTimestamp Timestamp) bool {
	if checkTimestamp <= 0 {
		return false
	}
	now := NewTimestampNow()
	todayStart, todayEnd := GetTodayStartEnd(now)
	return checkTimestamp >= todayStart && checkTimestamp <= todayEnd
}
