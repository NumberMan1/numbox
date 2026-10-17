package zset

import (
	"cmp"
	"fmt"
	"math"
	"math/rand/v2"
	"slices"
	"strconv"
	"sync"
	"testing"
)

func modelOrder(model map[string]float64) []string {
	members := make([]string, 0, len(model))
	for member := range model {
		members = append(members, member)
	}
	slices.SortFunc(members, func(a, b string) int {
		if order := cmp.Compare(model[a], model[b]); order != 0 {
			return order
		}
		return cmp.Compare(a, b)
	})
	return members
}

// Check public behavior against an independently sorted map, and every tower's
// spans/backlinks after updates and trims. Call only when writers are stopped.
func verifySet(t testing.TB, z *ZSet, model map[string]float64) {
	t.Helper()
	members := modelOrder(model)
	if z.Count() != len(members) {
		t.Fatalf("count %d, model %d", z.Count(), len(members))
	}
	if got := z.Range(1, math.MaxInt); !slices.Equal(got, members) {
		t.Fatalf("range %q, model %q", got, members)
	}
	if got := z.RangeByScore(math.Inf(-1), math.Inf(1)); !slices.Equal(got, members) {
		t.Fatalf("score range %q, model %q", got, members)
	}
	reversed := slices.Clone(members)
	slices.Reverse(reversed)
	if got := z.RevRange(1, math.MaxInt); !slices.Equal(got, reversed) {
		t.Fatalf("reverse range %q, model %q", got, reversed)
	}
	if got := z.RangeByScore(math.Inf(1), math.Inf(-1)); !slices.Equal(got, reversed) {
		t.Fatalf("reverse score range %q, model %q", got, reversed)
	}
	for i, member := range members {
		if rank := z.Rank(member); rank != i+1 {
			t.Fatalf("%q rank %d, want %d", member, rank, i+1)
		}
		if rank := z.RevRank(member); rank != len(members)-i {
			t.Fatalf("%q reverse rank %d", member, rank)
		}
		if score, ok := z.Score(member); !ok || score != model[member] {
			t.Fatalf("%q score %v, %v", member, score, ok)
		}
	}
	if z.head == nil {
		if len(model) != 0 {
			t.Fatal("nil head with members")
		}
		return
	}
	positions := map[*node]int{z.head: 0}
	previous := (*node)(nil)
	index, tallest := 0, 1
	for item := z.head.links[0].next; item != nil; item = item.links[0].next {
		if index >= len(members) || item.member != members[index] || item.score != model[item.member] || item.previous != previous {
			t.Fatalf("invalid node at %d", index)
		}
		index++
		positions[item] = index
		previous = item
		tallest = max(tallest, len(item.links))
	}
	if index != z.length || previous != z.tail || z.level != tallest || len(z.tbl) != z.length {
		t.Fatal("invalid length, tail or height")
	}
	for level := 0; level < z.level; level++ {
		for current := z.head; ; {
			edge := current.links[level]
			expected := z.length - positions[current]
			if edge.next != nil {
				nextRank, ok := positions[edge.next]
				if !ok || nextRank <= positions[current] {
					t.Fatal("invalid tower link")
				}
				expected = nextRank - positions[current]
			}
			if edge.span != expected {
				t.Fatalf("level %d rank %d span %d, want %d", level, positions[current], edge.span, expected)
			}
			if edge.next == nil {
				break
			}
			current = edge.next
		}
	}
}

func trimModel(model map[string]float64, keep int, reverse bool) int {
	members := modelOrder(model)
	keep = max(0, keep)
	if keep >= len(members) {
		return 0
	}
	removed := members[keep:]
	if reverse {
		removed = members[:len(members)-keep]
	}
	for _, member := range removed {
		delete(model, member)
	}
	return len(removed)
}

func TestPureGoBoundariesAndUpdates(t *testing.T) {
	var z ZSet
	model := map[string]float64{}
	verifySet(t, &z, model)
	for _, member := range []string{"", "a", "a\x00", "a\x00b", "aa", "z", "中文", "\xff"} {
		z.Add(1, member)
		model[member] = 1
	}
	z.Add(math.Inf(-1), "minimum")
	model["minimum"] = math.Inf(-1)
	z.Add(math.Inf(1), "maximum")
	model["maximum"] = math.Inf(1)
	verifySet(t, &z, model)
	for _, score := range []float64{math.Inf(-1), 2, -3, math.Inf(1), 1, 1} {
		z.Add(score, "a")
		model["a"] = score
		verifySet(t, &z, model)
	}
	z.Rem("missing")
	z.Rem("")
	delete(model, "")
	verifySet(t, &z, model)
	if z.Rank("missing") != 0 || z.RevRank("missing") != 0 {
		t.Fatal("missing rank")
	}
	if _, ok := z.Score("missing"); ok {
		t.Fatal("missing score")
	}
	for _, member := range modelOrder(model) {
		z.Rem(member)
		delete(model, member)
		verifySet(t, &z, model)
	}
	z.Add(42, "reused")
	model["reused"] = 42
	verifySet(t, &z, model)
}

func TestPureGoRankAndScoreRanges(t *testing.T) {
	z := New()
	for _, value := range []struct {
		score  float64
		member string
	}{{-1, "a"}, {0, "b"}, {0, "c"}, {1, "d"}, {2, "e"}} {
		z.Add(value.score, value.member)
	}
	for _, c := range []struct {
		first, last int
		reverse     bool
		want        []string
	}{
		{1, math.MaxInt, false, []string{"a", "b", "c", "d", "e"}},
		{5, 2, false, []string{"e", "d", "c", "b"}},
		{0, 1, false, nil}, {1, 0, false, nil}, {math.MinInt, math.MaxInt, false, nil}, {math.MaxInt, 1, false, []string{}},
		{1, math.MaxInt, true, []string{"e", "d", "c", "b", "a"}},
		{5, 2, true, []string{"a", "b", "c", "d"}},
		{2, math.MaxInt, true, []string{"d", "c", "b", "a"}},
		{0, 1, true, nil}, {1, 0, true, nil}, {math.MaxInt, 1, true, nil},
	} {
		got := z.Range(c.first, c.last)
		if c.reverse {
			got = z.RevRange(c.first, c.last)
		}
		if !slices.Equal(got, c.want) {
			t.Fatalf("%+v: got %q", c, got)
		}
	}
	for _, c := range []struct {
		first, last float64
		want        []string
	}{
		{0, 0, []string{"b", "c"}}, {1, 0, []string{"d", "c", "b"}}, {-1, 1, []string{"a", "b", "c", "d"}},
		{3, 4, []string{}}, {math.NaN(), 1, []string{}}, {0, math.NaN(), []string{}},
	} {
		if got := z.RangeByScore(c.first, c.last); !slices.Equal(got, c.want) {
			t.Fatalf("score %+v: %q", c, got)
		}
	}
}

func TestPureGoNaNDoesNotChangeSet(t *testing.T) {
	z := New()
	z.Add(1, "a")
	func() {
		defer func() {
			if recover() == nil {
				t.Error("NaN accepted")
			}
		}()
		z.Add(math.NaN(), "a")
	}()
	verifySet(t, z, map[string]float64{"a": 1})
}

func TestPureGoTrimBoundaries(t *testing.T) {
	for _, reverse := range []bool{false, true} {
		for _, keep := range []int{math.MinInt, -1, 0, 1, 3, 5, math.MaxInt} {
			z := New()
			model := map[string]float64{}
			for i := 0; i < 5; i++ {
				member := strconv.Itoa(i)
				z.Add(float64(i%2), member)
				model[member] = float64(i % 2)
			}
			want := trimModel(model, keep, reverse)
			var got int
			if reverse {
				got = z.RevLimit(keep)
			} else {
				got = z.Limit(keep)
			}
			if got != want {
				t.Fatalf("keep %d reverse %v: %d, want %d", keep, reverse, got, want)
			}
			verifySet(t, z, model)
		}
	}
}

func TestPureGoRandomizedModel(t *testing.T) {
	z := New()
	model := map[string]float64{}
	for i := 0; i < 2000; i++ {
		key := strconv.Itoa(i)
		z.Add(float64(i%11), key)
		model[key] = float64(i % 11)
	}
	verifySet(t, z, model)
	rng := rand.New(rand.NewPCG(12, 34))
	for step := 0; step < 5000; step++ {
		member := strconv.Itoa(rng.IntN(2500))
		switch rng.IntN(10) {
		case 0, 1, 2, 3, 4, 5:
			score := float64(rng.IntN(41) - 20)
			z.Add(score, member)
			model[member] = score
		case 6:
			z.Rem(member)
			delete(model, member)
		case 7:
			keep := rng.IntN(2600)
			if got, want := z.Limit(keep), trimModel(model, keep, false); got != want {
				t.Fatal(got, want)
			}
		case 8:
			keep := rng.IntN(2600)
			if got, want := z.RevLimit(keep), trimModel(model, keep, true); got != want {
				t.Fatal(got, want)
			}
		case 9:
			first, last := float64(rng.IntN(51)-25), float64(rng.IntN(51)-25)
			want := make([]string, 0)
			ordered := modelOrder(model)
			if first > last {
				slices.Reverse(ordered)
			}
			for _, key := range ordered {
				if model[key] >= min(first, last) && model[key] <= max(first, last) {
					want = append(want, key)
				}
			}
			if got := z.RangeByScore(first, last); !slices.Equal(got, want) {
				t.Fatal(got, want)
			}
		}
		if step%29 == 0 {
			verifySet(t, z, model)
		}
	}
	verifySet(t, z, model)
}

func TestPureGoConcurrentReadersWritersAndTrims(t *testing.T) {
	z := New()
	var workers sync.WaitGroup
	for worker := 0; worker < 8; worker++ {
		workers.Add(1)
		go func(worker int) {
			defer workers.Done()
			for step := 0; step < 500; step++ {
				member := strconv.Itoa((worker*37 + step) % 300)
				z.Add(float64(step%17), member)
				z.Count()
				z.Score(member)
				z.Rank(member)
				z.RevRank(member)
				z.Range(1, 10)
				z.RevRange(1, 10)
				z.RangeByScore(5, 12)
				if step%7 == 0 {
					z.Rem(member)
				}
				if step%23 == 0 {
					z.Limit(180)
				}
				if step%31 == 0 {
					z.RevLimit(150)
				}
			}
		}(worker)
	}
	workers.Wait()
	model := map[string]float64{}
	for _, key := range z.Range(1, math.MaxInt) {
		score, _ := z.Score(key)
		model[key] = score
	}
	verifySet(t, z, model)
}

func FuzzPureGoOperations(f *testing.F) {
	f.Add([]byte{0, 1, 2, 0, 2, 2, 0, 1, 3, 1, 2, 0, 2, 1, 0, 3, 0, 0})
	f.Add([]byte{0, 0, 255, 0, 255, 0, 0, 0, 0, 2, 255, 0})
	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) > 1536 {
			data = data[:1536]
		}
		z := New()
		model := map[string]float64{}
		for i := 0; i+2 < len(data); i += 3 {
			key := strconv.Itoa(int(data[i+1]))
			score := float64(int8(data[i+2]))
			switch data[i] % 4 {
			case 0:
				z.Add(score, key)
				model[key] = score
			case 1:
				z.Rem(key)
				delete(model, key)
			case 2:
				keep := int(data[i+1])
				if z.Limit(keep) != trimModel(model, keep, false) {
					t.Fatal("limit count")
				}
			case 3:
				keep := int(data[i+1])
				if z.RevLimit(keep) != trimModel(model, keep, true) {
					t.Fatal("reverse limit count")
				}
			}
			if i%96 == 0 {
				verifySet(t, z, model)
			}
		}
		verifySet(t, z, model)
	})
}

func BenchmarkPureGoZSet(b *testing.B) {
	keys := make([]string, 10000)
	for i := range keys {
		keys[i] = fmt.Sprintf("member-%05d", i)
	}
	for _, operation := range []string{"rank", "update", "range20"} {
		b.Run(operation, func(b *testing.B) {
			z := New()
			for i, key := range keys {
				z.Add(float64(i), key)
			}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				key := keys[i%len(keys)]
				switch operation {
				case "rank":
					z.Rank(key)
				case "update":
					z.Add(float64((i+1)%10001), key)
				case "range20":
					z.Range(5000, 5019)
				}
			}
		})
	}
}
