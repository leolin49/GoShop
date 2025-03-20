// Copyright (c) 2025, Yufeng Lin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package skiplist implements a skip list.
// Structure is thread concurrency safe.
// 
// SkipList, x present the dummy node.
// |
// x -> 1 ----------------> 8 -------------------> nil
// |    |                   |
// x -> 1 ------> 4 ------> 8 ----- -> 11 -------> nil
// |    |         |         |          |
// x -> 1 -> 3 -> 4 -> 6 -> 8 -> 10 -> 11 -> 15 -> nil
// 

package skiplist

import (
	"cmp"
	"errors"
	"fmt"
	"math/rand/v2"
	"sync"
)

const (
	MaxLevel = 128
	PFactor  = 0.25
)

// Comparator returns
//
//	-1 if x is less than y,
//	 0 if x equals y,
//	+1 if x is greater than y.
type Comparator[T any] func(x, y T) int

type SkipListNode[T comparable] struct {
	Val  T
	Span []int
	Next []*SkipListNode[T]
}

type Skiplist[T comparable] struct {
	Len   int
	Level int
	Head  *SkipListNode[T]
	Mu    sync.RWMutex
	Cmp   Comparator[T]
}

// NewSkiplist returns a new empty sorted list.
func NewSkiplist[T cmp.Ordered]() *Skiplist[T] {
	return &Skiplist[T]{
		Level: 0,
		Head: &SkipListNode[T]{
			Next: make([]*SkipListNode[T], MaxLevel),
			Span: make([]int, MaxLevel),
		},
		Len: 0,
		Cmp: cmp.Compare[T],
	}
}

func NewWithIntComparator() *Skiplist[int] {
	return NewSkiplist[int]()
}

func newErrorf(format string, a ...any) error {
	return errors.New(fmt.Sprintf("skiplist: " + format, a...))
}

func randomLevel() int {
	lv := 1
	for lv < MaxLevel && rand.Float64() < PFactor {
		lv++
	}
	return lv
}

// Exist reports whether the element is in the list.
// If the element not in, Exist returns false.
// Otherwise, Exist returns true.
func (s *Skiplist[T]) Exist(x T) bool {
	s.Mu.RLock()
	defer s.Mu.RUnlock()

	cur := s.Head
	for lv := s.Level - 1; lv >= 0; lv-- {
		for cur.Next[lv] != nil && s.greater(x, cur.Next[lv].Val) {
			cur = cur.Next[lv]
		}
	}
	cur = cur.Next[0]
	return cur != nil && cur.Val == x
}

// Insert inserts the element in the list.
func (s *Skiplist[T]) Insert(x T) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	update := make([]*SkipListNode[T], MaxLevel)
	rank := make([]int, MaxLevel)
	for i := range update {
		update[i] = s.Head
	}
	cur := s.Head
	for lv := s.Level - 1; lv >= 0; lv-- {
		if lv != s.Level-1 {
			rank[lv] = rank[lv+1]
		}
		for cur.Next[lv] != nil && s.greater(x, cur.Next[lv].Val) {
			rank[lv] += cur.Span[lv]
			cur = cur.Next[lv]
		}
		update[lv] = cur
	}
	level := randomLevel()
	s.Level = max(s.Level, level)
	newNode := &SkipListNode[T]{
		Val:  x,
		Next: make([]*SkipListNode[T], level),
		Span: make([]int, level),
	}
	for i, node := range update[:level] {
		newNode.Next[i] = node.Next[i]
		node.Next[i] = newNode

		newNode.Span[i] = node.Span[i] - (rank[0] - rank[i])
		node.Span[i] = rank[0] - rank[i] + 1
	}
	for i := level; i < s.Level; i++ {
		update[i].Span[i]++
	}
	s.Len++
}

// Erase reports whether the element is removed from the list.
// If the element is not existed, Erase returns an error.
// Otherwise, Erase will remove the element, and returns nil.
func (s *Skiplist[T]) Erase(x T) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	update := make([]*SkipListNode[T], MaxLevel)
	for i := range update {
		update[i] = s.Head
	}
	cur := s.Head
	for lv := s.Level - 1; lv >= 0; lv-- {
		for cur.Next[lv] != nil && s.greater(x, cur.Next[lv].Val) {
			cur = cur.Next[lv]
		}
		update[lv] = cur
	}
	cur = cur.Next[0]
	if cur == nil || !s.equal(cur.Val, x) {
		return newErrorf("not element [%v] in the list", x)
	}
	for lv := 0; lv < s.Level; lv++ {
		if update[lv].Next[lv] == cur {
			update[lv].Next[lv] = cur.Next[lv]
			update[lv].Span[lv] += cur.Span[lv] - 1
		} else {
			update[lv].Span[lv]--
		}
	}
	for s.Level > 1 && s.Head.Next[s.Level-1] == nil {
		s.Level--
	}
	s.Len--
	return nil
}

// Length returns the length of the list.
func (s *Skiplist[T]) Length() int {
	return s.Len
}

// Empty reports whether the list is empty.
// Return true if the list has not any element.
// Otherwise, Empty returns the false.
func (s *Skiplist[T]) Empty() bool {
	return s.Len == 0
}

// GetKth returns the K th element of the list.
// If there are fewer than k elements in the list, an error is returned.
func (s *Skiplist[T]) Kth(k int) (T, error) {
	s.Mu.RLock()
	defer s.Mu.RUnlock()

	if k < 1 || k > s.Len {
		var zero T
		return zero, newErrorf("have no such number [%v] node", k)
	}
	node := s.getKthElement(k)
	return node.Val, nil
}

// GetRange returns the slices of all elements between
// the L-th and R-th elements. If the interval [l, r]
// is invalid, GetRange returns an empty slice and an error.
// NOTE: If `l` is greater than `r`, the elements slice
// of interval [r, l] will be also returned, just like [l, r].
func (s *Skiplist[T]) Range(l, r int) ([]T, error) {
	s.Mu.RLock()
	defer s.Mu.RUnlock()

	if l < 1 || r < 1 || l > s.Len || r > s.Len {
		return nil, newErrorf("have no such number [%v] - [%v] node", l, r)
	}
	if l > r {
		l, r = r, l
	}
	lNode, rNode := s.getKthElement(l), s.getKthElement(r)
	elements := make([]T, r-l+1)
	for cur, idx := lNode, 0; cur != rNode.Next[0]; cur, idx = cur.Next[0], idx+1 {
		elements[idx] = cur.Val
	}
	return elements, nil
}

// Index returns the number of elements in the list(begin with 1).
func (s *Skiplist[T]) Index(x T) (int, error) {
	s.Mu.RLock()
	defer s.Mu.RUnlock()

	idx := 0
	cur := s.Head
	for lv := s.Level - 1; lv >= 0; lv-- {
		for cur.Next[lv] != nil && s.greater(x, cur.Next[lv].Val) {
			idx += cur.Span[lv]
			cur = cur.Next[lv]
		}
	}
	cur = cur.Next[0]
	if cur == nil || !s.equal(cur.Val, x) {
		return -1, newErrorf("not element [%v] in the list", x)
	}
	return idx + 1, nil
}

// Rank returns the rank of elements in the list(begin with 1).
// The rank of element means the reciprocal number of the element.
// Rank and Index should satisfy Index(x) + Rank(x) is Len() + 2.
func (s *Skiplist[T]) Rank(x T) (int, error) {
	s.Mu.RLock()
	defer s.Mu.RUnlock()

	idx, err := s.Index(x)
	if err != nil {
		return -1, err
	}
	return s.Len - idx + 1, nil
}

// Min returns the minimum element in the list.
func (s *Skiplist[T]) Min() (T, error) {
	return s.Kth(1)
}

// Max returns the maximum element in the list.
func (s *Skiplist[T]) Max() (T, error) {
	return s.Kth(s.Len)
}

// Lower returns the first element in the list
// which is greater than or equal to element x.
func (s *Skiplist[T]) Lower(x T) (T, error) {
	s.Mu.RLock()
	defer s.Mu.RUnlock()

	cur := s.Head
	for lv := s.Level - 1; lv >= 0; lv-- {
		for cur.Next[lv] != nil && s.greater(x, cur.Next[lv].Val) {
			cur = cur.Next[lv]
		}
	}
	cur = cur.Next[0]
	if cur == nil {
		var zero T
		return zero, newErrorf("no element greater than or equal to [%v]", x)
	}
	return cur.Val, nil
}

// Upper returns the first element in the list
// which is greater to element x.
func (s *Skiplist[T]) Upper(x T) (T, error) {
	s.Mu.RLock()
	defer s.Mu.RUnlock()

	cur := s.Head
	for lv := s.Level - 1; lv >= 0; lv-- {
		for cur.Next[lv] != nil && !s.less(x, cur.Next[lv].Val) {
			cur = cur.Next[lv]
		}
	}
	cur = cur.Next[0]
	if cur == nil {
		var zero T
		return zero, newErrorf("no element greater to [%v]", x)
	}
	return cur.Val, nil
}

// Remove removes the element which Index() is index.
func (s *Skiplist[T]) Remove(index int) error {
	s.Mu.RLock()
	if index > s.Len {
		return newErrorf("index [%v] out of range in the list", index)
	}
	v, _ := s.Kth(index)
	s.Mu.RUnlock()
	return s.Erase(v)
}

func (s *Skiplist[T]) getKthElement(k int) *SkipListNode[T] {
	n := 0
	cur := s.Head
	for lv := s.Level - 1; lv >= 0; lv-- {
		for cur.Next[lv] != nil && n+cur.Span[lv] <= k {
			n += cur.Span[lv]
			cur = cur.Next[lv]
		}
	}
	return cur
}

func (s *Skiplist[T]) greater(x, y T) bool {
	return s.Cmp(x, y) == 1
}

func (s *Skiplist[T]) less(x, y T) bool {
	return s.Cmp(x, y) == -1
}

func (s *Skiplist[T]) equal(x, y T) bool {
	return s.Cmp(x, y) == 0
}
