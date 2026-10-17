// Copyright (c) 2017 sean
// SPDX-License-Identifier: MIT
// Derived from github.com/sean-public/fast-skiplist and modified for numbox.
// See LICENSE and README.md in this directory for the license and local changes.

package skiplist

import (
	"math/big"
	"math/rand/v2"
	"sync"
)

type elementNode struct {
	next []*Element
}

type Element struct {
	elementNode
	key   *big.Float
	value interface{}
}

// Key allows retrieval of the key for a given Element
func (e *Element) Key() *big.Float {
	return e.key
}

// Value allows retrieval of the value for a given Element
func (e *Element) Value() interface{} {
	return e.value
}

// Next returns the following Element or nil if we're at the end of the list.
// Only operates on the bottom level of the skip list (a fully linked list).
func (element *Element) Next() *Element {
	return element.next[0]
}

type SkipList struct {
	elementNode
	maxLevel       int
	Length         int
	randSource     rand.Source
	probability    float64
	probTable      []float64
	mutex          sync.RWMutex
	prevNodesCache []*elementNode
}
