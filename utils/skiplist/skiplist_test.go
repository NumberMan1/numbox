package skiplist

import (
	"fmt"
	"math/big"
	"testing"
)

func foreach(list *SkipList) {
	for e := list.Front(); e != nil; e = e.Next() {
		fmt.Println("for ea", e.Value())
	}
}

func TestSkipList(t *testing.T) {
	list := New()
	list.Set(big.NewFloat(123), "This string data is stored at key 123!")
	list.Set(big.NewFloat(12), "This string data is stored at key 12!")
	list.Set(big.NewFloat(1234), "This string data is stored at key 1234!")
	list.Set(big.NewFloat(123422), "This string data is stored at key 123422!")
	list.Set(big.NewFloat(2), "This string data is stored at key 2!")
	fmt.Println(list.Get(big.NewFloat(123)).Value())
	fmt.Println(list.Length) // prints 1
	fmt.Println(list.Get(big.NewFloat(1234)).Value())
	list.Remove(big.NewFloat(123))
	fmt.Println(list.Length) // prints 0
	list.Set(big.NewFloat(1234), "This string data is stored at key -1234!")
	// 插入一些元素
	for i := 0; i < 100; i++ {
		list.Set(big.NewFloat(float64(i)), i)
	}

	// 获取第2页的数据（假设每页25个元素）
	page := 2
	items := list.RangeVals((page-1)*25, page*25)
	fmt.Printf("Page %d items: %+v\n", page, items)
	// 获取第5页的数据（假设每页25个元素）
	page = 5
	items = list.RangeVals((page-1)*25, page*25)
	fmt.Printf("Page %d items: %+v\n", page, items)

	fmt.Println("RankIndex 1", list.RankIndex(big.NewFloat(1)))
	fmt.Println("RankIndex 25", list.RankIndex(big.NewFloat(25)))
	fmt.Println("RankIndex 1234", list.RankIndex(big.NewFloat(1234)))
	fmt.Println("RankIndex 123", list.RankIndex(big.NewFloat(123)))
	foreach(list)
}
