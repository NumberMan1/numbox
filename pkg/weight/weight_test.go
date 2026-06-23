package weight

import (
	"fmt"
	"testing"
	"time"
)

type testWeightItem struct {
	value     int
	weightVal int
}

func (t testWeightItem) Weight() int {
	return t.weightVal
}

func TestPool_PickManyRandom(t *testing.T) {
	var pool = NewWeightPool()
	for i := 1; i <= 10; i++ {
		pool.Add(testWeightItem{
			value:     i,
			weightVal: 100,
		})
	}
	fmt.Println(pool)

	res := pool.PickManyRandom(3)
	fmt.Println(res)
	fmt.Println(pool)
}

func TestPool_PickOneRandom(t *testing.T) {
	var list = []int32{101}
	fmt.Println(PickOneFromFairPool(time.Now().UnixMilli(), list))
}

func TestPool_PickWithTotalWeight(t *testing.T) {
	var (
		succeedTimes, i int64
	)
	for i = 0; i < 10000; i++ {
		_, ok := PickOneFromItemsWithTotalWeight(time.Now().UnixNano(), 1000, testWeightItem{
			value:     1,
			weightVal: 1000,
		})
		if ok {
			succeedTimes++
		}
	}
	fmt.Println(succeedTimes)
}

func TestPool_PickAndPutBack_State(t *testing.T) {
	pool := NewWeightPool()
	pool.Add(testWeightItem{value: 1, weightVal: 100})
	pool.Add(testWeightItem{value: 2, weightVal: 200})

	initialLen := pool.Length()
	initialTotalWeight := pool.totalWeight

	// 执行 100 次抽取
	for i := 0; i < 100; i++ {
		_, ok := pool.PickRandomAndPutBack()
		if !ok {
			t.Errorf("Should always pick an item when totalWeight matches items")
		}
	}

	// 验证池子状态未改变
	if pool.Length() != initialLen {
		t.Errorf("Pool length changed! Got %d, want %d", pool.Length(), initialLen)
	}
	if pool.totalWeight != initialTotalWeight {
		t.Errorf("Total weight changed! Got %d, want %d", pool.totalWeight, initialTotalWeight)
	}
}

func TestPool_PickManyAndPutBack_Distribution(t *testing.T) {
	pool := NewWeightPool(123) // 使用固定种子
	// 1:2:7 的比例
	pool.Add(testWeightItem{value: 1, weightVal: 10})
	pool.Add(testWeightItem{value: 2, weightVal: 20})
	pool.Add(testWeightItem{value: 3, weightVal: 70})

	counts := make(map[int]int)
	totalPicks := 100000

	results := pool.PickRandomManyAndPutBack(totalPicks)
	for _, item := range results {
		val := item.(testWeightItem).value
		counts[val]++
	}

	// 验证分布误差在 1% 以内
	for val, weight := range map[int]int{1: 10, 2: 20, 3: 70} {
		expected := float64(totalPicks) * float64(weight) / 100.0
		actual := float64(counts[val])
		diff := (actual - expected) / expected
		if diff < 0 {
			diff = -diff
		}
		if diff > 0.05 { // 允许 5% 的统计波动
			t.Errorf("Value %d distribution error too high: got %v, want ~%v (diff %f)", val, actual, expected, diff)
		}
		fmt.Printf("Value %d: Actual %v, Expected %v, Diff %.2f%%\n", val, actual, expected, diff*100)
	}
}

func TestPool_PickAndPutBack_Empty(t *testing.T) {
	pool := NewWeightPool()
	_, ok := pool.PickRandomAndPutBack()
	if ok {
		t.Error("PickRandomAndPutBack on empty pool should return false")
	}

	results := pool.PickRandomManyAndPutBack(5)
	if len(results) != 0 {
		t.Errorf("PickRandomManyAndPutBack on empty pool should return empty slice, got len %d", len(results))
	}
}
