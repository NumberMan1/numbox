package weight

import (
	"sort"
	"time"

	"golang.org/x/exp/rand"
)

const (
	DefaultWeight = 100
)

// Shuffle 打乱切片中的元素
func Shuffle[T any](seed int64, slice []T) {
	r := rand.New(rand.NewSource(uint64(seed)))
	// Fisher-Yates 洗牌算法
	for i := len(slice) - 1; i > 0; i-- {
		// 生成一个 [0, i] 范围内的随机索引
		j := r.Intn(i + 1)
		// 交换 slice[i] 和 slice[j] 的元素
		slice[i], slice[j] = slice[j], slice[i]
	}
}

type fairWeightItem[T any] struct {
	val T
}

func (f fairWeightItem[T]) Weight() int {
	return 100
}

func newFairWeightPool[T any](randomSeed int64, slice []T) Pool {
	pool := NewWeightPool(randomSeed)
	for _, v := range slice {
		pool.Add(fairWeightItem[T]{v})
	}
	return pool
}

func pickFairItem[T any](pool *Pool) (res T, ok bool) {
	resItem, ok := pool.PickRandom()
	if !ok {
		return
	}
	res = resItem.(fairWeightItem[T]).val
	return
}

func PickOneFromFairPool[T any](randomSeed int64, slice []T) (res T, ok bool) {
	pool := newFairWeightPool(randomSeed, slice)
	return pickFairItem[T](&pool)
}

func PickOneFromFairPoolWithoutSeed[T any](slice []T) (res T, ok bool) {
	return PickOneFromFairPool(time.Now().UnixNano(), slice)
}

func newItemWeightPool[T Item](randomSeed int64, weightItems []T) Pool {
	pool := NewWeightPool(randomSeed)
	for _, weightItem := range weightItems {
		pool.Add(weightItem)
	}
	return pool
}

func PickOneFromItems[T Item](randomSeed int64, weightItems ...T) (res T, ok bool) {
	pool := newItemWeightPool(randomSeed, weightItems)
	resItem, ok := pool.PickRandom()
	if !ok {
		return
	}
	res, ok = resItem.(T)
	return
}

func PickManyFromFairPool[T any](pickCount int, randomSeed int64, slice []T) (res []T) {
	pool := newFairWeightPool(randomSeed, slice)
	for _, pickedItem := range pool.PickManyRandom(pickCount) {
		res = append(res, pickedItem.(fairWeightItem[T]).val)
	}
	return
}

func PickManyFromItems[T Item](pickCount int, randomSeed int64, weightItems ...T) (res []T) {
	pool := newItemWeightPool(randomSeed, weightItems)
	for _, pickedItem := range pool.PickManyRandom(pickCount) {
		res = append(res, pickedItem.(T))
	}
	return
}

func PickOneFromItemsWithTotalWeight[T Item](randomSeed int64, totalWeight int, weightItems ...T) (res T, ok bool) {
	pool := newItemWeightPool(randomSeed, weightItems)
	pool.SetTotalWeight(totalWeight)
	pickItem, ok := pool.PickRandom()
	if !ok {
		return
	}
	res = pickItem.(T)
	return
}

func NewWeightPool(randSeeds ...int64) Pool {
	var randSeed = time.Now().UnixNano()
	if len(randSeeds) > 0 {
		randSeed = randSeeds[0]
	}
	return Pool{
		randSeed:    randSeed,
		rand:        rand.New(rand.NewSource(uint64(randSeed))),
		Items:       make([]Item, 0, 8), // 预分配一点容量，减少初期扩容
		cum:         make([]int, 0, 8),
		totalWeight: 0,
	}
}

type Item interface {
	Weight() int
}

type Pool struct {
	randSeed int64
	rand     *rand.Rand
	// Items is exposed for compatibility. Do not mutate it or item weights after Add.
	Items       []Item
	totalWeight int
	cum         []int // 内部维护的累加权重，对外部完全透明
}

func (p *Pool) Length() int {
	return len(p.Items)
}

func (p *Pool) Add(item Item) {
	if item == nil {
		panic("weight: nil item")
	}
	weight := item.Weight()
	if weight < 0 || weight > int(^uint(0)>>1)-p.totalWeight {
		panic("weight: negative or overflowing weight")
	}
	sum := weight
	if len(p.cum) > 0 {
		sum += p.cum[len(p.cum)-1]
	}
	p.Items = append(p.Items, item)
	p.totalWeight += weight
	p.cum = append(p.cum, sum)
}

func (p *Pool) SetTotalWeight(totalWeight int) {
	sum := 0
	if len(p.cum) > 0 {
		sum = p.cum[len(p.cum)-1]
	}
	if totalWeight < sum {
		panic("weight: total weight is less than item weights")
	}
	p.totalWeight = totalWeight
}

func (p *Pool) GetTotalWeight() int {
	return p.totalWeight
}

func (p *Pool) SetRandomSeed(seed int64) {
	p.randSeed = seed
	p.rand = rand.New(rand.NewSource(uint64(seed)))
}

func (p *Pool) Copy() Pool {
	newPool := NewWeightPool(p.randSeed)
	newPool.Items = append([]Item{}, p.Items...)
	newPool.cum = append([]int(nil), p.cum...)
	newPool.totalWeight = p.totalWeight
	return newPool
}

// PickManyRandom picks up to pickCount items, skipping misses if r exceeds sum of weights.
func (p *Pool) PickManyRandom(pickCount int) []Item {
	if pickCount <= 0 {
		return nil
	}
	newPool := p.Copy()
	res := make([]Item, 0, min(pickCount, len(p.Items)))
	for i := 0; i < pickCount && newPool.Length() > 0; i++ {
		item, ok := newPool.PickRandom()
		if ok {
			res = append(res, item)
		}
		// if not ok, skip and continue
	}
	return res
}

// PickRandom returns a random item based on p.totalWeight; if r > sum of item weights, returns nil,false.
func (p *Pool) PickRandom() (Item, bool) {
	n := len(p.Items)
	if n == 0 || p.totalWeight <= 0 {
		return nil, false
	}
	r := p.rand.Intn(p.totalWeight) + 1
	if len(p.cum) > 0 && r > p.cum[len(p.cum)-1] {
		return nil, false
	}

	idx := sort.Search(n, func(i int) bool { return p.cum[i] >= r })
	if idx < 0 || idx >= n {
		return nil, false
	}
	// Remove using the indexed weight, keeping zero weights and any miss range.
	item := p.Items[idx]
	weight := p.cum[idx]
	if idx > 0 {
		weight -= p.cum[idx-1]
	}
	copy(p.Items[idx:], p.Items[idx+1:])
	p.Items[n-1] = nil
	p.Items = p.Items[:n-1]
	p.totalWeight -= weight
	for i := idx; i < n-1; i++ {
		p.cum[i] = p.cum[i+1] - weight
	}
	p.cum = p.cum[:n-1]
	return item, true
}

// PickRandomAndPutBack 随机获取一个物品但不从池子中移除 (放回抽样)
func (p *Pool) PickRandomAndPutBack() (Item, bool) {
	n := len(p.Items)
	if n == 0 || p.totalWeight <= 0 {
		return nil, false
	}

	r := p.rand.Intn(p.totalWeight) + 1
	if len(p.cum) > 0 && r > p.cum[len(p.cum)-1] {
		return nil, false
	}
	// 二分查找命中区间 (O(log N))
	idx := sort.Search(n, func(i int) bool { return p.cum[i] >= r })
	if idx < n {
		return p.Items[idx], true
	}

	return nil, false
}

// PickRandomManyAndPutBack 随机获取多个物品但不移除 (多次放回抽样)
func (p *Pool) PickRandomManyAndPutBack(pickCount int) []Item {
	n := len(p.Items)
	if n == 0 || pickCount <= 0 || p.totalWeight <= 0 {
		return nil
	}

	res := make([]Item, 0, pickCount)
	for i := 0; i < pickCount; i++ {
		r := p.rand.Intn(p.totalWeight) + 1
		if len(p.cum) > 0 && r > p.cum[len(p.cum)-1] {
			continue // 命中空位，模拟 PickManyRandom 跳过逻辑
		}

		// 每次抽取仅需 O(log N)
		idx := sort.Search(n, func(i int) bool { return p.cum[i] >= r })
		if idx < n {
			res = append(res, p.Items[idx])
		}
	}
	return res
}
