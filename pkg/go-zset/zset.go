// Package zset provides a concurrent, ranked sorted set implemented in pure Go.
package zset

import (
	"fmt"
	"math"
	"math/rand/v2"
	"sync"
)

// DEFAULT_TBL_LEN is the initial member-index capacity, retained for compatibility.
const DEFAULT_TBL_LEN = 16
const maxLevel = 32

type link struct {
	next *node
	span int // number of level-zero members crossed, or remaining to the end
}

type node struct {
	member   string
	score    float64
	previous *node
	links    []link
}

// ZSet sorts by score ascending, then member bytes ascending. Methods are safe
// concurrently; each query and trim observes one locked state. The zero value
// works. Do not copy a ZSet after first use.
type ZSet struct {
	mu            sync.RWMutex
	head, tail    *node
	level, length int
	tbl           map[string]float64
}

// New creates an empty set. No C compiler, finalizer or explicit Close is needed.
func New() *ZSet {
	z := &ZSet{}
	z.initLocked()
	return z
}

func (z *ZSet) initLocked() {
	if z.head == nil {
		z.head = &node{links: make([]link, maxLevel)}
		z.level = 1
		z.tbl = make(map[string]float64, DEFAULT_TBL_LEN)
	}
}

func nodeBefore(item *node, score float64, member string) bool {
	return item.score < score || (item.score == score && item.member < member)
}

// predecessorsLocked records the path and zero-based rank at each level.
func (z *ZSet) predecessorsLocked(score float64, member string, path *[maxLevel]*node, ranks *[maxLevel]int) {
	current := z.head
	rank := 0
	for i := z.level - 1; i >= 0; i-- {
		for next := current.links[i].next; next != nil && nodeBefore(next, score, member); next = current.links[i].next {
			rank += current.links[i].span
			current = next
		}
		path[i], ranks[i] = current, rank
	}
}

func randomLevel() int {
	level := 1
	bits := rand.Uint64()
	for level < maxLevel && bits&3 == 0 {
		level++
		bits >>= 2
	}
	return level
}

// Add inserts a member or updates its score. NaN panics before changing the set;
// infinities are supported. Updating a member reuses its allocated node/tower.
func (z *ZSet) Add(score float64, member string) {
	if math.IsNaN(score) {
		panic("zset: NaN score")
	}
	z.mu.Lock()
	defer z.mu.Unlock()
	z.initLocked()
	var path [maxLevel]*node
	var ranks [maxLevel]int
	var item *node
	if old, ok := z.tbl[member]; ok {
		if old == score {
			return
		}
		z.predecessorsLocked(old, member, &path, &ranks)
		item = path[0].links[0].next
		z.unlinkLocked(item, &path)
		item.score = score
	} else {
		item = &node{member: member, score: score, links: make([]link, randomLevel())}
	}
	z.predecessorsLocked(score, member, &path, &ranks)
	if height := len(item.links); height > z.level {
		for i := z.level; i < height; i++ {
			path[i], ranks[i] = z.head, 0
			z.head.links[i].span = z.length
		}
		z.level = height
	}
	for i := range item.links {
		distance := ranks[0] - ranks[i]
		item.links[i] = link{next: path[i].links[i].next, span: path[i].links[i].span - distance}
		path[i].links[i] = link{next: item, span: distance + 1}
	}
	for i := len(item.links); i < z.level; i++ {
		path[i].links[i].span++
	}
	item.previous = path[0]
	if item.previous == z.head {
		item.previous = nil
	}
	if next := item.links[0].next; next != nil {
		next.previous = item
	} else {
		z.tail = item
	}
	z.length++
	z.tbl[member] = score
}

// unlinkLocked removes a known node. The caller updates the member index.
func (z *ZSet) unlinkLocked(item *node, path *[maxLevel]*node) {
	for i := 0; i < z.level; i++ {
		edge := &path[i].links[i]
		if edge.next == item {
			edge.span += item.links[i].span - 1
			edge.next = item.links[i].next
		} else {
			edge.span--
		}
	}
	if next := item.links[0].next; next != nil {
		next.previous = item.previous
	} else {
		z.tail = item.previous
	}
	z.length--
	for z.level > 1 && z.head.links[z.level-1].next == nil {
		z.level--
	}
}

// Rem removes a member, doing nothing if it is absent.
func (z *ZSet) Rem(member string) {
	z.mu.Lock()
	defer z.mu.Unlock()
	score, ok := z.tbl[member]
	if !ok {
		return
	}
	var path [maxLevel]*node
	var ranks [maxLevel]int
	z.predecessorsLocked(score, member, &path, &ranks)
	z.unlinkLocked(path[0].links[0].next, &path)
	delete(z.tbl, member)
}

// Count returns the number of members.
func (z *ZSet) Count() int {
	z.mu.RLock()
	defer z.mu.RUnlock()
	return z.length
}

// Score returns a member's score and whether it exists.
func (z *ZSet) Score(member string) (float64, bool) {
	z.mu.RLock()
	defer z.mu.RUnlock()
	score, ok := z.tbl[member]
	return score, ok
}

func (z *ZSet) nodeByRankLocked(rank int) *node {
	if rank < 1 || rank > z.length {
		return nil
	}
	current := z.head
	traversed := 0
	for i := z.level - 1; i >= 0; i-- {
		for current.links[i].next != nil && current.links[i].span <= rank-traversed {
			traversed += current.links[i].span
			current = current.links[i].next
		}
		if traversed == rank {
			return current
		}
	}
	return nil
}

func (z *ZSet) rangeLocked(first, last int) []string {
	if first < 1 || last < 1 {
		return nil
	}
	if first > z.length {
		return []string{}
	}
	forward := first <= last
	if forward {
		last = min(last, z.length)
	}
	count := first - last + 1
	if forward {
		count = last - first + 1
	}
	result := make([]string, 0, count)
	for item := z.nodeByRankLocked(first); item != nil && len(result) < count; {
		result = append(result, item.member)
		if forward {
			item = item.links[0].next
		} else {
			item = item.previous
		}
	}
	return result
}

// Range returns an inclusive, 1-based rank range. first>last walks backwards
// from first. Nonpositive ranks return nil; an out-of-set first returns empty.
// Endpoints beyond the end are bounded without allocating for the excess range.
func (z *ZSet) Range(first, last int) []string {
	z.mu.RLock()
	defer z.mu.RUnlock()
	return z.rangeLocked(first, last)
}

// RevRange uses 1-based ranks counted from the largest member. An end beyond
// Count is clamped. Reversing the endpoints reverses traversal direction.
func (z *ZSet) RevRange(first, last int) []string {
	z.mu.RLock()
	defer z.mu.RUnlock()
	if first < 1 || last < 1 || first > z.length {
		return nil
	}
	last = min(last, z.length)
	return z.rangeLocked(z.length-first+1, z.length-last+1)
}

// RangeByScore returns the inclusive score range, descending when first>last.
// NaN bounds return an empty slice. Equal scores are ordered by member bytes.
func (z *ZSet) RangeByScore(first, last float64) []string {
	z.mu.RLock()
	defer z.mu.RUnlock()
	result := make([]string, 0)
	if z.length == 0 || math.IsNaN(first) || math.IsNaN(last) {
		return result
	}
	current := z.head
	if first <= last {
		for i := z.level - 1; i >= 0; i-- {
			for next := current.links[i].next; next != nil && next.score < first; next = current.links[i].next {
				current = next
			}
		}
		for current = current.links[0].next; current != nil && current.score <= last; current = current.links[0].next {
			result = append(result, current.member)
		}
	} else {
		for i := z.level - 1; i >= 0; i-- {
			for next := current.links[i].next; next != nil && next.score <= first; next = current.links[i].next {
				current = next
			}
		}
		for current != nil && current != z.head && current.score >= last {
			result = append(result, current.member)
			current = current.previous
		}
	}
	return result
}

func (z *ZSet) rankLocked(member string) int {
	score, ok := z.tbl[member]
	if !ok {
		return 0
	}
	current := z.head
	rank := 0
	for i := z.level - 1; i >= 0; i-- {
		for next := current.links[i].next; next != nil && (nodeBefore(next, score, member) || (next.score == score && next.member == member)); next = current.links[i].next {
			rank += current.links[i].span
			current = next
		}
		if current != z.head && current.member == member {
			return rank
		}
	}
	return 0
}

// Rank returns the member's 1-based ascending rank, or zero when absent.
func (z *ZSet) Rank(member string) int {
	z.mu.RLock()
	defer z.mu.RUnlock()
	return z.rankLocked(member)
}

// RevRank returns the member's 1-based descending rank, or zero when absent.
func (z *ZSet) RevRank(member string) int {
	z.mu.RLock()
	defer z.mu.RUnlock()
	rank := z.rankLocked(member)
	if rank == 0 {
		return 0
	}
	return z.length - rank + 1
}

// deleteRanksLocked removes a contiguous range using one predecessor search.
func (z *ZSet) deleteRanksLocked(first, last int) int {
	var path [maxLevel]*node
	current := z.head
	traversed := 0
	for i := z.level - 1; i >= 0; i-- {
		for current.links[i].next != nil && current.links[i].span < first-traversed {
			traversed += current.links[i].span
			current = current.links[i].next
		}
		path[i] = current
	}
	removed := 0
	for item := path[0].links[0].next; item != nil && removed < last-first+1; {
		next := item.links[0].next
		z.unlinkLocked(item, &path)
		delete(z.tbl, item.member)
		removed++
		item = next
	}
	return removed
}

// Limit retains the lowest count members, returning the number removed.
// A negative count is treated as zero. The entire trim is atomic.
func (z *ZSet) Limit(count int) int {
	z.mu.Lock()
	defer z.mu.Unlock()
	count = max(0, count)
	if count >= z.length {
		return 0
	}
	return z.deleteRanksLocked(count+1, z.length)
}

// RevLimit retains the highest count members, returning the number removed.
// A negative count is treated as zero. The entire trim is atomic.
func (z *ZSet) RevLimit(count int) int {
	z.mu.Lock()
	defer z.mu.Unlock()
	count = max(0, count)
	if count >= z.length {
		return 0
	}
	return z.deleteRanksLocked(1, z.length-count)
}

// Dump prints a diagnostic snapshot to stdout. The text format is not an API.
func (z *ZSet) Dump() {
	z.mu.RLock()
	defer z.mu.RUnlock()
	if z.head == nil {
		return
	}
	rank := 1
	for item := z.head.links[0].next; item != nil; item = item.links[0].next {
		fmt.Printf("node %d: score:%f, member:%s\n", rank, item.score, item.member)
		rank++
	}
}
