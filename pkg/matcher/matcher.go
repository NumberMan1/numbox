// Adapted from github.com/NumberMan1/orbit (81cbe7b).
// Copyright (c) 2025 NumberMan1. See LICENSE in the repository root.

package matcher

import "unicode/utf8"

type matchNode struct {
	children map[rune]*matchNode
	terminal bool
}

// Match uses rune indices [Start,End). Matcher is immutable after construction,
// so searches are safe concurrently. Matching is exact, case-sensitive Unicode;
// normalization, dictionaries and application policy stay with the caller.
type Match struct{ Start, End int }
type Matcher struct{ root matchNode }

func NewMatcher(words []string) *Matcher {
	m := &Matcher{}
	for _, word := range words {
		if word == "" || !utf8.ValidString(word) {
			continue
		}
		n := &m.root
		for _, r := range word {
			if n.children == nil {
				n.children = make(map[rune]*matchNode)
			}
			if n.children[r] == nil {
				n.children[r] = &matchNode{}
			}
			n = n.children[r]
		}
		n.terminal = true
	}
	return m
}

// Find returns all overlapping matches, longest first for each start index.
func (m *Matcher) Find(text string) []Match {
	if !utf8.ValidString(text) {
		return nil
	}
	runes := []rune(text)
	var matches []Match
	for start := range runes {
		n := &m.root
		first := len(matches)
		for end := start; end < len(runes); end++ {
			n = n.children[runes[end]]
			if n == nil {
				break
			}
			if n.terminal {
				matches = append(matches, Match{start, end + 1})
			}
		}
		for a, b := first, len(matches)-1; a < b; a, b = a+1, b-1 {
			matches[a], matches[b] = matches[b], matches[a]
		}
	}
	return matches
}
func (m *Matcher) Contains(text string) bool {
	if !utf8.ValidString(text) {
		return false
	}
	runes := []rune(text)
	for start := range runes {
		n := &m.root
		for _, r := range runes[start:] {
			n = n.children[r]
			if n == nil {
				break
			}
			if n.terminal {
				return true
			}
		}
	}
	return false
}

// Replace masks the union of matched ranges, preserving other text/whitespace.
func (m *Matcher) Replace(text string, replacement rune) string {
	if !utf8.ValidString(text) {
		return text
	}
	runes := []rune(text)
	for _, match := range m.Find(text) {
		for i := match.Start; i < match.End; i++ {
			runes[i] = replacement
		}
	}
	return string(runes)
}
