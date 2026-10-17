package percent

import (
	"testing"

	"golang.org/x/exp/rand"
)

func TestHitBoundsAndExactThreshold(t *testing.T) {
	for seed := int64(0); seed < 10000; seed++ {
		draw := rand.New(rand.NewSource(uint64(seed))).Intn(1000)
		for _, p := range []Percent{-1, 0, 1, 500, 999, 1000, 1001} {
			want := draw < int(p)
			if got := p.HitBySeed(seed); got != want {
				t.Fatalf("seed=%d draw=%d p=%d: %v, want %v", seed, draw, p, got, want)
			}
		}
		if ZeroPercent.Hit() || !FullPercent.Hit() {
			t.Fatal("invalid probability bound")
		}
	}
}
func TestDivideRounding(t *testing.T) {
	if got := GetPercentByDivide(1, 3); got != 333 {
		t.Fatal(got)
	}
	if got := GetPercentByDivideCeil(1, 3); got != 334 {
		t.Fatal(got)
	}
	if got := GetPercentByDivide(1, 2); got != 500 {
		t.Fatal(got)
	}
	defer func() {
		if recover() == nil {
			t.Fatal("zero divisor accepted")
		}
	}()
	GetPercentByDivide(1, 0)
}
