package time

import (
	"testing"
	"time"
)

func TestTimestampIsValid(t *testing.T) {
	// Timestamp 是毫秒 int64,零值 0 对应 1970-01-01,time.UnixMilli(0).IsZero()
	// 恒为 false,因此 IsValid 必须直接判 t > 0,不能落到 time.Time。
	cases := []struct {
		name string
		ts   Timestamp
		want bool
	}{
		{name: "zero", ts: 0, want: false},
		{name: "negative", ts: -1, want: false},
		{name: "unix-epoch-milli", ts: 1, want: true},
		{name: "now", ts: NewTimestampNow(), want: true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.ts.IsValid(); got != c.want {
				t.Fatalf("Timestamp(%d).IsValid() = %v, want %v", c.ts, got, c.want)
			}
		})
	}
}

func TestDurationToNextInterval(t *testing.T) {
	// 零值 Timestamp 视为“未设定”，应返回 0（走调用方当前时间兜底），
	// 而不能被 time.UnixMilli(0).IsZero()==false 误判为 1970 起点。
	if got := DurationToNextInterval(0, time.Minute); got != 0 {
		t.Fatalf("DurationToNextInterval(0, 1m) = %v, want 0", got)
	}
	if got := DurationToNextInterval(NewTimestampNow(), 0); got != 0 {
		t.Fatalf("DurationToNextInterval(now, 0) = %v, want 0", got)
	}
	// 非零时间戳应落在下一个整分钟边界。
	base := NewTimestampWithUnixMilli(1_000_000_000)
	got := DurationToNextInterval(base, time.Minute)
	if got <= 0 || got > time.Minute {
		t.Fatalf("DurationToNextInterval(%d, 1m) = %v, want (0, 1m]", base, got)
	}
}

func TestOffsetClock(t *testing.T) {
	defer SetOffsetMillis(0)

	SetOffsetMillis(0)
	if got := OffsetMillis(); got != 0 {
		t.Fatalf("OffsetMillis() = %d, want 0", got)
	}

	// 偏移 +24h 后 Now() 应比真实时间前移约一天。
	const dayMillis = int64(24 * time.Hour / time.Millisecond)
	SetOffsetMillis(dayMillis)
	delta := Now().Sub(time.Now())
	if delta < 23*time.Hour || delta > 25*time.Hour {
		t.Fatalf("Now() delta = %v, want ~24h", delta)
	}

	// NewTimestampNow 走同一时钟，应同步前移。
	tsDelta := NewTimestampNow().Time().Sub(time.Now())
	if tsDelta < 23*time.Hour || tsDelta > 25*time.Hour {
		t.Fatalf("NewTimestampNow delta = %v, want ~24h", tsDelta)
	}

	// 重置回真实时间。
	SetOffsetMillis(0)
	if delta := Now().Sub(time.Now()); delta > time.Second || delta < -time.Second {
		t.Fatalf("after reset Now() delta = %v, want ~0", delta)
	}
}

func TestBackwardOffset(t *testing.T) {
	defer SetOffsetMillis(0)

	// 负偏移（回拨）：Now() 应比真实时间落后约一天。
	const dayMillis = int64(24 * time.Hour / time.Millisecond)
	SetOffsetMillis(-dayMillis)
	delta := Now().Sub(time.Now())
	if delta > -23*time.Hour || delta < -25*time.Hour {
		t.Fatalf("backward Now() delta = %v, want ~-24h", delta)
	}

	// 回拨后，对一个刚记录的“服务器时间”时间戳求 Since 应为负——
	// 这正是回拨会让冷却/每日刷新落到未来的根因，锁定该语义。
	ts := NewTimestampNow()
	SetOffsetMillis(-2 * dayMillis)
	if since := Since(ts); since >= 0 {
		t.Fatalf("Since after further backward jump = %v, want negative", since)
	}
}

func TestCalendarBoundariesAndLocalOwnership(t *testing.T) {
	saved := time.Local
	t.Cleanup(func() { time.Local = saved })
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatal(err)
	}
	time.Local = loc
	for _, d := range []struct {
		month time.Month
		day   int
		want  time.Duration
	}{{time.March, 8, 23 * time.Hour}, {time.November, 1, 25 * time.Hour}} {
		now := NewTimestamp(time.Date(2026, d.month, d.day, 12, 0, 0, 0, loc))
		start, end := GetTodayStartEnd(now)
		if end.SinceFrom(start) != d.want || end != GetTodayEnd(now) {
			t.Fatal(d, start, end)
		}
		ws, we := GetWeekStartEnd(now)
		if we != GetWeekEnd(now) || we.Time().Hour() != 0 || ws.Time().Weekday() != time.Monday || we.Time().Weekday() != time.Monday {
			t.Fatal(ws, we)
		}
	}
	time.Local = time.FixedZone("UTC+5:30", 19800)
	if GetZoneOffset() != 5.5 {
		t.Fatal(GetZoneOffset())
	}
}
func TestIsTodayHalfOpenBoundary(t *testing.T) {
	old := OffsetMillis()
	t.Cleanup(func() { SetOffsetMillis(old) })
	SetOffsetMillis(0)
	start, end := GetTodayStartEnd()
	if !IsToday(start) || IsToday(end) {
		t.Fatal("invalid [start,end) boundary")
	}
}
