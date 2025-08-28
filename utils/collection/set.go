package collection

import (
	"github.com/NumberMan1/numbox/utils"
)

func NewSet[T comparable](items ...T) Set[T] {
	var res = make(Set[T], len(items))
	res.Add(items...)
	return res
}

type Set[T comparable] map[T]struct{}

func (set Set[T]) Add(items ...T) {
	for _, item := range items {
		set[item] = struct{}{}
	}
}

func (set Set[T]) Remove(item T) {
	delete(set, item)
}

func (set Set[T]) Contains(item T) bool {
	_, ok := set[item]
	return ok
}

func (set Set[T]) ContainsAny(items ...T) bool {
	for _, item := range items {
		if set.Contains(item) {
			return true
		}
	}
	return false
}

func (set Set[T]) All() []T {
	return utils.MapKeys(set)
}
