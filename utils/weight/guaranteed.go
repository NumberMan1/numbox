// Filename: guaranteed.go
// This file should be placed in the same 'weight' package as your existing 'weight.go'.
// It depends on the 'Item' interface and 'Shuffle' function from that file.

package weight

import (
	"sort"
	"time"

	"golang.org/x/exp/rand"
)

// GuaranteedItem is an interface for items that can be placed in a GuaranteedPool.
// It extends the base Item interface with methods to define the guarantee threshold
// and an initial miss count.
type GuaranteedItem interface {
	Item // Embeds Weight() int method.
	// GuaranteedTimes returns the number of misses after which this item is guaranteed to be picked.
	// If GuaranteedTimes() returns 0 or a negative number, the item has no guarantee.
	GuaranteedTimes() int
	// MissTimes returns the initial number of times the item has already been missed.
	// This is used to initialize the item's state when it's added to the pool.
	MissTimes() int
}

// guaranteedItemState is an internal struct used by the GuaranteedPool to track the state
// of an item within that specific pool instance. This is crucial for managing the 'miss' count
// without requiring the user's item object to be mutable, which is a better design practice.
type guaranteedItemState struct {
	item      GuaranteedItem
	missTimes int
}

// GuaranteedPool is a random-weighted pool that supports guaranteed picks.
// After an item has been passed over (missed) for a certain number of times,
// it is guaranteed to be selected on a subsequent pick.
type GuaranteedPool struct {
	rand        *rand.Rand
	randSeed    int64
	items       []*guaranteedItemState // A slice of pointers to the internal state representation of items.
	totalWeight int
}

// NewGuaranteedPool creates a new, empty GuaranteedPool.
// You can optionally provide a seed for the random number generator.
func NewGuaranteedPool(randSeeds ...int64) *GuaranteedPool {
	var randSeed = time.Now().UnixNano()
	if len(randSeeds) > 0 {
		randSeed = randSeeds[0]
	}
	return &GuaranteedPool{
		randSeed:    randSeed,
		rand:        rand.New(rand.NewSource(uint64(randSeed))),
		items:       make([]*guaranteedItemState, 0),
		totalWeight: 0,
	}
}

// Add adds a new GuaranteedItem to the pool.
// The initial miss count is taken from the item's MissTimes() method.
func (p *GuaranteedPool) Add(item GuaranteedItem) {
	state := &guaranteedItemState{
		item:      item,
		missTimes: item.MissTimes(), // Use the initial miss count from the item.
	}
	p.items = append(p.items, state)
	p.totalWeight += item.Weight()
}

// Length returns the number of items currently in the pool.
func (p *GuaranteedPool) Length() int {
	return len(p.items)
}

// SetTotalWeight explicitly sets the total weight for the pool. This is used for "miss" calculations.
// If a random number is greater than the sum of weights of items in the pool but less than
// or equal to this totalWeight, it's considered a "miss", and no item is picked.
func (p *GuaranteedPool) SetTotalWeight(totalWeight int) {
	p.totalWeight = totalWeight
}

// SetRandomSeed sets a new random seed for the pool's random number generator.
func (p *GuaranteedPool) SetRandomSeed(seed int64) {
	p.randSeed = seed
	p.rand = rand.New(rand.NewSource(uint64(seed)))
}

// PickRandom selects one item from the pool and removes it.
// It first checks if any items have met their guaranteed pick threshold.
// If so, it randomly picks one from the guaranteed items.
// If not, it performs a standard weighted random selection.
// On every pick (or miss), the miss counter for unpicked items is incremented.
func (p *GuaranteedPool) PickRandom() (GuaranteedItem, bool) {
	if len(p.items) == 0 {
		return nil, false
	}

	// 1. Identify all items that are guaranteed to be picked.
	guaranteedIndices := make([]int, 0)
	for i, state := range p.items {
		// An item is guaranteed if its guarantee threshold is positive and the miss count has reached it.
		if state.item.GuaranteedTimes() > 0 && state.missTimes >= state.item.GuaranteedTimes() {
			guaranteedIndices = append(guaranteedIndices, i)
		}
	}

	var pickedIndex = -1

	if len(guaranteedIndices) > 0 {
		// 2. We have guaranteed items. Pick one of them randomly.
		randGuaranteedIndex := p.rand.Intn(len(guaranteedIndices))
		pickedIndex = guaranteedIndices[randGuaranteedIndex]
	} else {
		// 3. No guaranteed items, perform a normal weighted random pick.
		// Build the cumulative weight slice for selection.
		sum := 0
		cum := make([]int, len(p.items))
		for i, s := range p.items {
			sum += s.item.Weight()
			cum[i] = sum
		}

		// The effective total weight determines the chance of a "miss".
		effectiveTotalWeight := p.totalWeight
		if effectiveTotalWeight < sum {
			effectiveTotalWeight = sum
		}

		// Draw a random number.
		r := p.rand.Intn(effectiveTotalWeight) + 1

		if r > sum {
			// This is a "miss". No item is picked.
			// Increment miss times for ALL items and return.
			for _, state := range p.items {
				state.missTimes++
			}
			return nil, false
		}

		// Find the item corresponding to the random number `r` using binary search.
		idx := sort.Search(len(cum), func(i int) bool { return cum[i] >= r })
		if idx < len(p.items) {
			pickedIndex = idx
		}
	}

	if pickedIndex == -1 {
		// This should not be reached if the pool has items, but acts as a safeguard.
		return nil, false
	}

	// 4. An item has been selected.
	pickedItem := p.items[pickedIndex].item

	// Update miss times: increment for all non-picked items.
	for i, state := range p.items {
		if i != pickedIndex {
			state.missTimes++
		}
	}
	// Reset the miss counter for the picked item.
	// We do not reset the state in the p.items slice because the item is about to be removed.
	// This logic is sound because the state is removed along with the item.

	// 5. Remove the picked item from the pool.
	p.totalWeight -= pickedItem.Weight()
	p.items = append(p.items[:pickedIndex], p.items[pickedIndex+1:]...)

	return pickedItem, true
}

// PickOneFromGuaranteedItems is a convenience function to pick a single item from a list.
func PickOneFromGuaranteedItems[T GuaranteedItem](randomSeed int64, items ...T) (res T, ok bool) {
	pool := NewGuaranteedPool(randomSeed)
	for _, item := range items {
		pool.Add(item)
	}
	picked, success := pool.PickRandom()
	if !success {
		return
	}
	res, ok = picked.(T)
	return
}

// PickManyFromGuaranteedItems creates a temporary pool and picks multiple items from it.
// The state of the pool (like miss times) is contained within this function call.
// For stateful multi-picking, manage the GuaranteedPool object directly.
func PickManyFromGuaranteedItems[T GuaranteedItem](pickCount int, randomSeed int64, items ...T) []T {
	if len(items) == 0 || pickCount <= 0 {
		return []T{}
	}

	pool := NewGuaranteedPool(randomSeed)
	for _, item := range items {
		pool.Add(item)
	}

	// If we need to pick all or more than available, shuffle and return all items.
	if pool.Length() <= pickCount {
		res := make([]T, len(items))
		for i, it := range items {
			res[i] = it
		}
		Shuffle(res) // Using Shuffle from the original `weight.go` file.
		return res
	}

	// Pick items one by one from the pool.
	pickedItems := make([]T, 0, pickCount)
	for i := 0; i < pickCount; i++ {
		item, ok := pool.PickRandom()
		if ok {
			pickedItems = append(pickedItems, item.(T))
		}
		// If !ok (a miss), we simply pick one less item, matching the behavior
		// of PickManyRandom from the original file.
	}

	return pickedItems
}
