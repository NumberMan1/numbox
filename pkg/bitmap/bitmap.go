package bitmap

import "fmt"

// BitMap 是一个位图数据结构。
// 它使用一个字节切片来存储大量的布尔标记，非常节省内存。
// 例如，需要标记数字 1, 2, 5 是否存在，它会在底层字节的相应比特位上标记为1。
// 注意：此实现不是线程安全的。
type BitMap struct {
	bits   []byte // 底层存储的字节切片
	maxVal uint32 // 可以存储的最大数值，决定了位图的容量
	count  int32  // 当前已添加元素的数量
}

// NewBitMap 创建一个指定容量的 BitMap。
// max 参数定义了这个位图能处理的最大数值。例如，如果 max 为 1000，则可以处理 0-1000 的数字。
func NewBitMap(max uint32) *BitMap {
	// 计算需要的字节数。max >> 3 等价于 max / 8。
	// +1 是为了确保即使 max 是8的倍数，也能容纳该数值。
	b := make([]byte, (max>>3)+1)
	return &BitMap{bits: b, maxVal: max}
}

// Add 将一个数字添加到 BitMap 中。
// 这会在对应的比特位上标记为 1。
func (bm *BitMap) Add(num uint32) error {
	if num > bm.maxVal {
		return fmt.Errorf("number %d exceeds the max value %d", num, bm.maxVal)
	}
	// 找到 num 对应的字节索引
	byteIndex := num >> 3
	// 找到 num 在该字节中的比特位索引
	bitPosition := num & 0x07 // 等价于 num % 8
	// 使用位掩码检查该位是否已经为1
	if bm.bits[byteIndex]&(1<<bitPosition) == 0 {
		// 使用 OR 操作将该位置为 1
		bm.bits[byteIndex] |= 1 << bitPosition
		bm.count++
	}
	return nil
}

// IsExist 检查一个数字是否存在于 BitMap 中。
func (bm *BitMap) IsExist(num uint32) bool {
	if num > bm.maxVal {
		return false
	}
	byteIndex := num >> 3
	bitPosition := num & 0x07
	// 检查对应比特位是否为 1
	return bm.bits[byteIndex]&(1<<bitPosition) != 0
}

// Remove 从 BitMap 中移除一个数字。
// 这会将对应的比特位清零。
func (bm *BitMap) Remove(num uint32) error {
	if num > bm.maxVal {
		return fmt.Errorf("number %d exceeds the max value %d", num, bm.maxVal)
	}
	byteIndex := num >> 3
	bitPosition := num & 0x07
	// 检查该位是否原本就存在
	if bm.bits[byteIndex]&(1<<bitPosition) != 0 {
		// 使用 AND 和 NOT 操作将该位清零
		bm.bits[byteIndex] &= ^(1 << bitPosition)
		bm.count--
	}
	return nil
}

// Max 返回位图支持的最大数值。
func (bm *BitMap) Max() uint32 {
	return bm.maxVal
}

// Bits 返回位图底层的字节切片。
// 这可以用于序列化和持久化存储。
func (bm *BitMap) Bits() []byte {
	return bm.bits
}

// SetBits 使用一个已有的字节切片重置位图。
// 这通常用于从持久化存储中恢复位图状态。
func (bm *BitMap) SetBits(b []byte) {
	bm.bits = b
	// 重置后需要重新计算元素数量
	bm.recount()
}

// String 返回位图内容的字符串表示，方便调试。
func (bm *BitMap) String() string {
	return fmt.Sprint(bm.bits)
}

// Count 返回当前位图中已添加元素的数量。
func (bm *BitMap) Count() int32 {
	return bm.count
}

// recount 重新计算位图中为1的比特位总数。
func (bm *BitMap) recount() {
	bm.count = 0
	for _, b := range bm.bits {
		// 遍历每个字节，计算其中1的个数
		for b > 0 {
			b &= b - 1 // 这个操作可以高效地移除最低位的1
			bm.count++
		}
	}
}
