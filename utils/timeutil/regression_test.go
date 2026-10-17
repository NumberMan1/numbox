package timeutil

import (
	"testing"
	"time"
)

func TestChineseDurationValidation(t *testing.T) {
	for _, s := range []string{"", "bad", "1天junk", "-1秒", "1.5秒", "1000000000000000000天"} {
		if _, err := ParseCNDuration(s); err == nil {
			t.Fatalf("accepted %q", s)
		}
	}
	if got, err := ParseCNDuration("1天 2小时3分钟4秒5毫秒"); err != nil || got != 26*time.Hour+3*time.Minute+4*time.Second+5*time.Millisecond {
		t.Fatal(got, err)
	}
}
