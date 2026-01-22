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
func Shuffle[T any](slice []T) {
	// Fisher-Yates 洗牌算法
	for i := len(slice) - 1; i > 0; i-- {
		// 生成一个 [0, i] 范围内的随机索引
		j := rand.Intn(i + 1)
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

func PickOneFromFairPool[T any](randomSeed int64, slice []T) (res T, ok bool) {
	pool := NewWeightPool(randomSeed)
	for _, v := range slice {
		pool.Add(fairWeightItem[T]{v})
	}
	resItem, ok := pool.PickRandom()
	if !ok {
		return
	}
	res = resItem.(fairWeightItem[T]).val
	return
}

func PickOneFromFairPoolWithoutSeed[T any](slice []T) (res T, ok bool) {
	pool := NewWeightPool(time.Now().UnixNano())
	for _, v := range slice {
		pool.Add(fairWeightItem[T]{v})
	}
	resItem, ok := pool.PickRandom()
	if !ok {
		return
	}
	res = resItem.(fairWeightItem[T]).val
	return
}

func PickOneFromItems[T Item](randomSeed int64, weightItems ...T) (res T, ok bool) {
	pool := NewWeightPool(randomSeed)
	for _, weightItem := range weightItems {
		pool.Add(weightItem)
	}
	resItem, ok := pool.PickRandom()
	if !ok {
		return
	}
	res, ok = resItem.(T)
	return
}

func PickManyFromFairPool[T any](pickCount int, randomSeed int64, slice []T) (res []T) {
	pool := NewWeightPool(randomSeed)
	for _, v := range slice {
		pool.Add(fairWeightItem[T]{v})
	}
	for _, pickedItem := range pool.PickManyRandom(pickCount) {
		res = append(res, pickedItem.(fairWeightItem[T]).val)
	}
	return
}

func PickManyFromItems[T Item](pickCount int, randomSeed int64, weightItems ...T) (res []T) {
	pool := NewWeightPool(randomSeed)
	for _, weightItem := range weightItems {
		pool.Add(weightItem)
	}
	for _, pickedItem := range pool.PickManyRandom(pickCount) {
		res = append(res, pickedItem.(T))
	}
	return
}

func PickOneFromItemsWithTotalWeight[T Item](randomSeed int64, totalWeight int, weightItems ...T) (res T, ok bool) {
	pool := NewWeightPool(randomSeed)
	for _, weightItem := range weightItems {
		pool.Add(weightItem)
	}
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
		Items:       nil,
		totalWeight: 0,
	}
}

type Item interface {
	Weight() int
}

type Pool struct {
	randSeed    int64
	rand        *rand.Rand
	Items       []Item
	totalWeight int
}

func (p *Pool) Length() int {
	return len(p.Items)
}

func (p *Pool) Add(item Item) {
	p.Items = append(p.Items, item)
	p.totalWeight += item.Weight()
}

func (p *Pool) SetTotalWeight(totalWeight int) {
	p.totalWeight = totalWeight
}

func (p *Pool) SetRandomSeed(seed int64) {
	p.randSeed = seed
	p.rand = rand.New(rand.NewSource(uint64(seed)))
}

func (p *Pool) Copy() Pool {
	newPool := NewWeightPool(p.randSeed)
	newPool.Items = append([]Item{}, p.Items...)
	newPool.totalWeight = p.totalWeight
	return newPool
}

// PickManyRandom picks up to pickCount items, skipping misses if r exceeds sum of weights.
func (p *Pool) PickManyRandom(pickCount int) []Item {
	if len(p.Items) <= pickCount {
		Shuffle(p.Items)
		return p.Items
	}
	newPool := p.Copy()
	res := make([]Item, 0, pickCount)
	for i := 0; i < pickCount; i++ {
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
	if n == 0 {
		return nil, false
	}
	if n == 1 && p.Items[0].Weight() >= p.totalWeight {
		item := p.Items[0]
		p.Items = nil
		p.totalWeight = 0
		return item, true
	}
	// build cumulative weights
	sum := 0
	cum := make([]int, n)
	for i, it := range p.Items {
		sum += it.Weight()
		cum[i] = sum
	}
	// draw
	r := p.rand.Intn(p.totalWeight) + 1
	if r > sum {
		// miss
		return nil, false
	}
	// find slot
	idx := sort.Search(n, func(i int) bool { return cum[i] >= r })
	if idx < 0 || idx >= n {
		return nil, false
	}
	// pick and remove
	item := p.Items[idx]
	p.Items = append(p.Items[:idx], p.Items[idx+1:]...)
	p.totalWeight -= item.Weight()
	return item, true
}

// PickRandomAndPutBack 随机获取一个物品但不从池子中移除 (放回抽样)
func (p *Pool) PickRandomAndPutBack() (Item, bool) {
	n := len(p.Items)
	if n == 0 {
		return nil, false
	}

	// 1. 构建临时累加权重 (如果需要极致性能且池子静态，建议在 Pool 结构体中缓存此切片)
	sum := 0
	cum := make([]int, n)
	for i, it := range p.Items {
		sum += it.Weight()
		cum[i] = sum
	}

	// 2. 使用内部随机生成器 draw
	r := p.rand.Intn(p.totalWeight) + 1
	if r > sum {
		return nil, false
	}

	// 3. 二分查找命中区间 (O(log N))
	idx := sort.Search(n, func(i int) bool { return cum[i] >= r })
	if idx < n {
		return p.Items[idx], true
	}

	return nil, false
}

// PickRandomManyAndPutBack 随机获取多个物品但不移除 (多次放回抽样)
func (p *Pool) PickRandomManyAndPutBack(pickCount int) []Item {
	n := len(p.Items)
	if n == 0 || pickCount <= 0 {
		return nil
	}

	// 1. 预计算累加权重 (仅计算一次)
	sum := 0
	cum := make([]int, n)
	for i, it := range p.Items {
		sum += it.Weight()
		cum[i] = sum
	}

	res := make([]Item, 0, pickCount)
	for i := 0; i < pickCount; i++ {
		r := p.rand.Intn(p.totalWeight) + 1
		if r > sum {
			continue // 命中空位，模拟 PickManyRandom 跳过逻辑
		}

		// 2. 每次抽取仅需 O(log N)
		idx := sort.Search(n, func(i int) bool { return cum[i] >= r })
		if idx < n {
			res = append(res, p.Items[idx])
		}
	}
	return res
}
