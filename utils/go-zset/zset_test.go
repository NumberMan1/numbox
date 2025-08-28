package zset

import (
	"fmt"
	"github.com/NumberMan1/numbox/utils"
	"strconv"
	"testing"
)

func equal(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func assert(t *testing.T, ok bool, s string) {
	if !ok {
		t.Error(s)
	}
}

func TestBase(t *testing.T) {
	z := New()
	assert(t, z.Count() == 0, "empty Count error")
	z.Add(1, "12")
	z.Add(1, "32")
	assert(t, z.Count() == 2, "not empty Count error")
	var score float64
	var ex bool
	score, ex = z.Score("12")
	assert(t, score == 1, "Score error")
	z.Add(2, "12")
	assert(t, z.Count() == 2, "after add duplicate Count error")
	score, ex = z.Score("12")
	assert(t, score == 2, "after add Score error")
	z.Rem("12")
	assert(t, z.Count() == 1, "after rem Count error")
	score, ex = z.Score("12")
	assert(t, ex == false, "not exist Score error")
	fmt.Println("")
}

func TestRangeByScore(t *testing.T) {
	z := New()
	z.Add(2, "22")
	z.Add(1, "11")
	z.Add(3, "33")
	s := "TestRangeByScore error"
	assert(t, equal(z.RangeByScore(2, 3), []string{"22", "33"}), s)
	assert(t, equal(z.RangeByScore(0, 5), []string{"11", "22", "33"}), s)
	assert(t, equal(z.RangeByScore(10, 5), []string{}), s)
	assert(t, equal(z.RangeByScore(10, 0), []string{"33", "22", "11"}), s)
}

func TestRange(t *testing.T) {
	z := New()
	z.Add(100.1, "1")
	z.Add(100.9, "9")
	z.Add(100.5, "5")
	assert(t, equal(z.Range(1, 3), []string{"1", "5", "9"}), "Range1 error")
	assert(t, equal(z.Range(3, 1), []string{"9", "5", "1"}), "Range2 error")
	assert(t, equal(z.Range(2, 6), []string{"5", "9"}), "Range3 error")
	assert(t, equal(z.Range(6, 2), nil), "Range4 error")

	assert(t, equal(z.RevRange(1, 2), []string{"9", "5"}), "RevRange1 error")
	assert(t, equal(z.RevRange(3, 2), []string{"1", "5"}), "RevRange2 error")
	assert(t, equal(z.RevRange(6, 2), nil), "RevRange3 error")
	assert(t, equal(z.RevRange(2, 6), []string{"5", "1"}), "RevRange4 error")
}

func TestRank(t *testing.T) {
	z := New()
	assert(t, z.Rank("kehan") == 0, "Rank empty error")
	z.Add(1111.1111, "kehan")
	assert(t, z.Rank("kehan") == 1, "Rank error")
	z.Add(222.2222, "lwy")
	assert(t, z.Rank("kehan") == 2, "Rank 2 error")
	assert(t, z.RevRank("kehan") == 1, "RevRank error")
}

func TestLimit(t *testing.T) {
	z := New()
	z.Add(1, "1")
	z.Add(2, "2")
	z.Add(3, "3")
	z.Limit(1)
	assert(t, z.Count() == 1, "Limit error")
	assert(t, z.Rank("3") == 0, "Limit Rank error")
	z.Add(4.4, "4")
	z.Add(5.5, "5")
	z.Add(0.5, "0.5")
	z.Dump()
	assert(t, z.RevLimit(4) == 0, "RevLimit error")
	assert(t, z.RevLimit(0) == 4, "RevLimit2 error")
}

func BenchmarkZset(b *testing.B) {
	z := New()

	b.Run("add", func(b *testing.B) {
		b.ResetTimer()
		// for i := range 10000 {
		// 	z.Add(float64(util.RandInt(1, 10000)), strconv.Itoa(i))
		// 	z.RevLimit(30)
		// }
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				z.Add(float64(utils.RandInt(1, 100000)), strconv.Itoa(int(utils.RandInt(1, 100000))))
				z.RevLimit(10000)
			}
		})
	})

	b.Run("incr", func(b *testing.B) {
		b.ResetTimer()
		// for i := range 10000 {
		// 	z.Add(float64(util.RandInt(1, 10000)), strconv.Itoa(i))
		// 	z.RevLimit(30)
		// }
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				z.Add(float64(utils.RandInt(1, 100000)), strconv.Itoa(int(utils.RandInt(1, 100000))))
				z.RevLimit(10000)
			}
		})
	})

	a, _ := z.Score("30")
	fmt.Println(z.Count(), a)
	for _, member := range z.RevRange(1, 30) {
		score, _ := z.Score(member)
		fmt.Println(member, score)
	}
}
