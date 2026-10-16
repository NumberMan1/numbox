// Adapted from github.com/NumberMan1/orbit (81cbe7b).
// Copyright (c) 2025 NumberMan1. See LICENSE in the repository root.

package matcher

import (
	"reflect"
	"sync"
	"testing"
)

func TestMatcherOverlapsUnicodeAndPreservesText(t *testing.T) {
	m := NewMatcher([]string{"aba", "ba", "坏词", "坏", ""})
	if got := m.Find("ababa"); !reflect.DeepEqual(got, []Match{{0, 3}, {1, 3}, {2, 5}, {3, 5}}) {
		t.Fatal(got)
	}
	if got := m.Replace("  坏词 & ababa!  ", '*'); got != "  ** & *****!  " {
		t.Fatal(got)
	}
	if m.Contains("not present") || !m.Contains("末尾坏词") {
		t.Fatal("match detection")
	}
	var wg sync.WaitGroup
	for range 8 {
		wg.Go(func() {
			for range 100 {
				if m.Replace("坏词", '*') != "**" {
					t.Error("concurrent search")
				}
			}
		})
	}
	wg.Wait()
}
