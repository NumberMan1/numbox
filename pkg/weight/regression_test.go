package weight

import (
	"reflect"
	"testing"
)

func mustPanic(t *testing.T, fn func()) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Error("expected validation panic")
		}
	}()
	fn()
}
func TestPoolBatchPreservesIndexAndOwnership(t *testing.T) {
	p := NewWeightPool(7)
	p.Add(testWeightItem{1, 10})
	p.Add(testWeightItem{2, 20})
	p.Add(testWeightItem{3, 30})
	before := append([]Item(nil), p.Items...)
	got := p.PickManyRandom(100)
	if len(got) != 3 || !reflect.DeepEqual(p.Items, before) || p.GetTotalWeight() != 60 {
		t.Fatal(got, p)
	}
	got[0] = nil
	if !reflect.DeepEqual(p.Items, before) {
		t.Fatal("batch aliases original pool")
	}
	for range 100 {
		if _, ok := p.PickRandomAndPutBack(); !ok {
			t.Fatal("index corrupted")
		}
	}
	if len(p.PickManyRandom(-1)) != 0 {
		t.Fatal("negative count")
	}
}
func TestPoolZeroWeightsMissesAndReuse(t *testing.T) {
	p := NewWeightPool(7)
	p.Add(testWeightItem{1, 0})
	if _, ok := p.PickRandom(); ok {
		t.Fatal("zero-weight item picked")
	}
	if _, ok := p.PickRandomAndPutBack(); ok {
		t.Fatal("zero-weight item picked")
	}
	if len(p.PickRandomManyAndPutBack(3)) != 0 || len(p.PickManyRandom(3)) != 0 {
		t.Fatal("zero-weight batch")
	}
	p.Add(testWeightItem{2, 10})
	p.SetTotalWeight(11)
	p.Add(testWeightItem{3, 10})
	// The miss range must not be inserted into the new item's cumulative weight.
	if !reflect.DeepEqual(p.cum, []int{0, 10, 20}) {
		t.Fatal(p.cum)
	}
	for p.Length() > 1 {
		if item, ok := p.PickRandom(); ok && item.(testWeightItem).weightVal == 0 {
			t.Fatal("picked zero weight")
		}
	}
	if p.GetTotalWeight() != 1 {
		t.Fatal(p.GetTotalWeight())
	}
	p.Add(testWeightItem{4, 10})
	if !reflect.DeepEqual(p.cum, []int{0, 10}) {
		t.Fatal(p.cum)
	}
}
func TestPoolSingletonRefillAndValidation(t *testing.T) {
	p := NewWeightPool(1)
	p.Add(testWeightItem{1, 10})
	p.PickRandom()
	p.Add(testWeightItem{2, 10})
	if item, ok := p.PickRandom(); !ok || item.(testWeightItem).value != 2 || len(p.cum) != 0 {
		t.Fatal(item, ok, p.cum)
	}
	mustPanic(t, func() { p.Add(testWeightItem{3, -1}) })
	p.Add(testWeightItem{4, 10})
	mustPanic(t, func() { p.SetTotalWeight(9) })
	mustPanic(t, func() { p.Add(testWeightItem{5, int(^uint(0) >> 1)}) })
	if p.Length() != 1 || p.GetTotalWeight() != 10 {
		t.Fatal("invalid add mutated pool")
	}
}

type guaranteedTestItem struct {
	testWeightItem
	threshold, misses int
}

func (i guaranteedTestItem) GuaranteedTimes() int { return i.threshold }
func (i guaranteedTestItem) MissTimes() int       { return i.misses }
func TestGuaranteedZeroWeightAndBatchThreshold(t *testing.T) {
	item := guaranteedTestItem{testWeightItem{1, 0}, 2, 0}
	p := NewGuaranteedPool(1)
	p.Add(item)
	for range 2 {
		if _, ok := p.PickRandom(); ok {
			t.Fatal("premature guarantee")
		}
	}
	if _, ok := p.PickRandom(); !ok {
		t.Fatal("guarantee failed after two misses")
	}
	if got := PickManyFromGuaranteedItems(2, 1, item); len(got) != 0 {
		t.Fatal("batch bypassed threshold", got)
	}
	if got := PickManyFromGuaranteedItems(3, 1, item); len(got) != 1 {
		t.Fatal("batch lost guarantee", got)
	}
	mustPanic(t, func() { p.Add(guaranteedTestItem{testWeightItem{2, -1}, 0, 0}) })
}
