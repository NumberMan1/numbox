package time

import (
	"time"
)

var (
	Second = time.Second
	Minute = time.Minute
	Hour   = time.Hour
	Day    = time.Hour * 24
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

func ParseDuration(duration string) (Duration, error) {
	return time.ParseDuration(duration)
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
