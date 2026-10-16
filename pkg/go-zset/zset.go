package zset

//#include "skiplist.h"
import "C"

import (
	"runtime"
	"sync"
)

// DEFAULT_TBL_LEN 默认表size
const DEFAULT_TBL_LEN = 16

// ZSet 有序集合
type ZSet struct {
	sl  *C.skiplist
	tbl map[string]float64
	mu  sync.RWMutex
}

// tocstring 将 Go 字符串转换为 C 字符串，并返回指针和长度
func tocstring(s string) (*C.char, C.size_t) {
	// 将 Go 字符串转换为 C 字符串
	cstr := C.CString(s)
	// 计算字符串长度
	length := C.size_t(len(s))
	return cstr, length
}

// newslobj 创建 skiplist 的对象节点
func newslobj(s string) *C.slobj {
	p, l := tocstring(s)
	return C.slCreateObj(p, l)
}

// New 创建一个新的 ZSet 实例
func New() *ZSet {
	z := &ZSet{sl: C.slCreate(), tbl: make(map[string]float64, DEFAULT_TBL_LEN)}
	runtime.SetFinalizer(z, func(z *ZSet) {
		C.slFree(z.sl)
		z.tbl = nil
	})
	return z
}

// Add 向有序集合中添加成员，若成员已存在则更新分数
func (z *ZSet) Add(score float64, member string) {
	z.mu.Lock()
	defer z.mu.Unlock()

	if old, ok := z.tbl[member]; ok {
		if old == score {
			return
		}
		C.slDelete(z.sl, C.double(old), newslobj(member))
	}
	C.slInsert(z.sl, C.double(score), newslobj(member))
	z.tbl[member] = score
}

// Rem 从有序集合中移除指定成员
func (z *ZSet) Rem(member string) {
	z.mu.Lock()
	defer z.mu.Unlock()

	if score, ok := z.tbl[member]; ok {
		C.slDelete(z.sl, C.double(score), newslobj(member))
		delete(z.tbl, member)
	}
}

// Count 返回有序集合中的成员数量
func (z *ZSet) Count() int {
	return int(z.sl.length)
}

// Score 获取指定成员的分数，若不存在返回 false
func (z *ZSet) Score(member string) (float64, bool) {
	z.mu.RLock()
	defer z.mu.RUnlock()

	score, ex := z.tbl[member]
	return score, ex
}

// Range 按排名范围返回成员列表（1-based），支持正序和逆序
func (z *ZSet) Range(r1, r2 int) []string {
	z.mu.RLock()
	defer z.mu.RUnlock()

	// 从1开始，不支持倒数查询
	if r1 < 1 || r2 < 1 {
		return nil
	}

	var reverse, rangelen int
	if r1 <= r2 {
		reverse = 0
		rangelen = r2 - r1 + 1
	} else {
		reverse = 1
		rangelen = r1 - r2 + 1
	}
	node := C.slGetNodeByRank(z.sl, C.ulong(r1))
	result := make([]string, 0, rangelen)
	rr := C.int(reverse)
	for n := 0; node != nil && n < rangelen; {
		result = append(result, C.GoStringN(node.obj.ptr, C.int(node.obj.length)))
		node = C.getNextNode(node, rr)
		n++
	}
	return result
}

// reverseRank 计算逆序排名
func (z *ZSet) reverseRank(r int) int {
	return z.Count() - r + 1
}

// RevRange 按逆序排名范围返回成员列表
func (z *ZSet) RevRange(r1, r2 int) []string {
	// 支持末尾超出查询 如：实际只有3条，查询2-6，实际返回2-3的数据
	if r2 > z.Count() {
		r2 = z.Count()
	}
	return z.Range(z.reverseRank(r1), z.reverseRank(r2))
}

// RangeByScore 按分数范围返回成员列表，支持正序和逆序
func (z *ZSet) RangeByScore(s1, s2 float64) []string {
	z.mu.RLock()
	defer z.mu.RUnlock()

	var reverse int
	var node *C.skiplistNode
	cs1, cs2 := C.double(s1), C.double(s2)
	if s1 <= s2 {
		reverse = 0
		node = C.slFirstInRange(z.sl, cs1, cs2)
	} else {
		reverse = 1
		node = C.slLastInRange(z.sl, cs2, cs1)
	}

	result := make([]string, 0)
	rr := C.int(reverse)
	for node != nil {
		if reverse == 1 {
			if node.score < cs2 {
				break
			}
		} else {
			if node.score > cs2 {
				break
			}
		}
		result = append(result, C.GoStringN(node.obj.ptr, C.int(node.obj.length)))
		node = C.getNextNode(node, rr)
	}
	return result
}

// Rank 返回指定成员的排名（1-based），不存在返回 0
func (z *ZSet) Rank(member string) int {
	z.mu.RLock()
	defer z.mu.RUnlock()

	score, ex := z.tbl[member]
	if !ex {
		return 0
	}
	rank := C.slGetRank(z.sl, C.double(score), newslobj(member))
	return int(rank)
}

// RevRank 返回指定成员的逆序排名（1-based），不存在返回 0
func (z *ZSet) RevRank(member string) int {
	rank := z.Rank(member)
	if rank != 0 {
		rank = z.reverseRank(rank)
	}
	return rank
}

// deleteByRank 删除指定排名范围内的成员，返回删除数量
func (z *ZSet) deleteByRank(from, to int) int {
	if from > to {
		from, to = to, from
	}
	members := z.Range(from, to)
	for _, member := range members {
		z.Rem(member)
	}
	return len(members)
}

// func (z *zset) deleteByRank(from, to int) int {
// 	if from > to {
// 		from, to = to, from
// 	}

// 	// 将 slice 中的 data 指针指向的对象 Pin 住
// 	var pinner runtime.Pinner
// 	// defer pinner.Unpin()
// 	pinner.Pin(&z)

// 	return int(C.slDeleteByRank(z.sl, C.uint(from), C.uint(to), unsafe.Pointer(z)))
// }

// //export delCb
// func delCb(p unsafe.Pointer, obj *C.slobj) {
// 	z := (*zset)(p)
// 	member := C.GoStringN(obj.ptr, C.int(obj.length))
// 	delete(z.tbl, member)
// }

// Limit 保留前 count 个成员，其余删除，返回删除数量
func (z *ZSet) Limit(count int) int {
	total := z.Count()
	if total <= count {
		return 0
	}
	return z.deleteByRank(count+1, total)
}

// RevLimit 保留后 count 个成员，其余删除，返回删除数量
func (z *ZSet) RevLimit(count int) int {
	total := z.Count()
	if total <= count {
		return 0
	}
	from := z.reverseRank(count + 1)
	to := z.reverseRank(total)
	return z.deleteByRank(from, to)
}

// Dump 打印当前有序集合的内容（调试用）
func (z *ZSet) Dump() {
	z.mu.RLock()
	defer z.mu.RUnlock()

	C.slDump(z.sl)
}
