package skiplist

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"sync"
)

const (
	MaxLevel = 49
	PFactor  = 0.25
)

type SkipListNode struct {
	Val  int
	Span []int
	Next []*SkipListNode
}

type Skiplist struct {
	Len   int
	Level int
	Head  *SkipListNode
	Mu	  sync.RWMutex
}

// NewSkiplist returns a new empty sorted list.
func NewSkiplist() *Skiplist {
	return &Skiplist{
		Level: 0,
		Head: &SkipListNode{
			Val:  -1,
			Next: make([]*SkipListNode, MaxLevel),
			Span: make([]int, MaxLevel),
		},
		Len: 0,
	}
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
func (s *Skiplist) Exist(num int) bool {
	s.Mu.RLock()
	defer s.Mu.RUnlock()

	cur := s.Head
	for lv := s.Level - 1; lv >= 0; lv-- {
		for cur.Next[lv] != nil && num > cur.Next[lv].Val {
			cur = cur.Next[lv]
		}
	}
	cur = cur.Next[0]
	return cur != nil && cur.Val == num
}

// Insert inserts the element in the list.
func (s *Skiplist) Insert(num int) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	update := make([]*SkipListNode, MaxLevel)
	rank := make([]int, MaxLevel)
	for i := range update {
		update[i] = s.Head
	}
	cur := s.Head
	for lv := s.Level - 1; lv >= 0; lv-- {
		if lv != s.Level-1 {
			rank[lv] = rank[lv+1]
		}
		for cur.Next[lv] != nil && num > cur.Next[lv].Val {
			rank[lv] += cur.Span[lv]
			cur = cur.Next[lv]
		}
		update[lv] = cur
	}
	level := randomLevel()
	s.Level = max(s.Level, level)
	newNode := &SkipListNode{
		Val:  num,
		Next: make([]*SkipListNode, level),
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
func (s *Skiplist) Erase(num int) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	
	update := make([]*SkipListNode, MaxLevel)
	for i := range update {
		update[i] = s.Head
	}
	cur := s.Head
	for lv := s.Level - 1; lv >= 0; lv-- {
		for cur.Next[lv] != nil && num > cur.Next[lv].Val {
			cur = cur.Next[lv]
		}
		update[lv] = cur
	}
	cur = cur.Next[0]
	if cur == nil || cur.Val != num {
		return errors.New(fmt.Sprintf("not element [%v] in the list", num))
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
func (s *Skiplist) Length() int {
	return s.Len
}

// Empty reports whether the list is empty.
// Return true if the list has not any element.
// Otherwise, Empty returns the false.
func (s *Skiplist) Empty() bool {
	return s.Len == 0
}

// GetKth returns the K th element of the list.
// If there are fewer than k elements in the list, an error is returned.
func (s *Skiplist) Kth(k int) (int, error) {
	s.Mu.RLock()
	defer s.Mu.RUnlock()

	if k < 1 || k > s.Len {
		return 0, errors.New("[Skiplist] list have no such number node")
	}
	node := s.getKthElement(k)
	return node.Val, nil
}

// GetRange returns the slices of all elements between
// the L-th and R-th elements. If the interval [l, r]
// is invalid, GetRange returns an empty slice and an error.
// NOTE: If `l` is greater than `r`, the elements slice
// of interval [r, l] will be also returned, just like [l, r].
func (s *Skiplist) Range(l, r int) ([]int, error) {
	s.Mu.RLock()
	defer s.Mu.RUnlock()

	if l < 1 || r < 1 || l > s.Len || r > s.Len {
		return nil, errors.New("[Skiplist] list have no such number node")
	}
	if l > r {
		l, r = r, l
	}
	lNode, rNode := s.getKthElement(l), s.getKthElement(r)
	elements := make([]int, r-l+1)
	for cur, idx := lNode, 0; cur != rNode.Next[0]; cur, idx = cur.Next[0], idx+1 {
		elements[idx] = cur.Val
	}
	return elements, nil
}

// Index returns the number of elements in the list(begin with 1).
func (s *Skiplist) Index(x int) (int, error) {
	s.Mu.RLock()
	defer s.Mu.RUnlock()

	idx := 0
	cur := s.Head
	for lv := s.Level - 1; lv >= 0; lv-- {
		for cur.Next[lv] != nil && x > cur.Next[lv].Val {
			idx += cur.Span[lv]
			cur = cur.Next[lv]
		}
	}
	cur = cur.Next[0]
	if cur == nil || cur.Val != x {
		return -1, errors.New(fmt.Sprintf("not element [%v] in the list", x))
	}
	return idx + 1, nil
}

// Rank returns the rank of elements in the list(begin with 1).
// The rank of element means the reciprocal number of the element.
// Rank and Index should satisfy Index(x) + Rank(x) is Len() + 2.
func (s *Skiplist) Rank(x int) (int, error) {
	idx, err := s.Index(x)
	if err != nil {
		return -1, err
	}
	return s.Len - idx + 1, nil
}

// First returns the first element in the list.
func (s *Skiplist) First() (int, error) {
	return s.Kth(1)
}

// Last returns the last element in the list.
func (s *Skiplist) Last() (int, error) {
	return s.Kth(s.Len)
}

// func (s *Skiplist) Lower() (int, error) {
// }

// func (s *Skiplist) Upper() (int, error) {
// }

// Remove removes the element which Index() is index.
func (s *Skiplist) Remove(index int) error {
	if index > s.Len {
		return errors.New(fmt.Sprintf("index [%v] out of range in the list", index))
	}
	v, _ := s.Kth(index)
	return s.Erase(v)
}

func (s *Skiplist) getKthElement(k int) *SkipListNode {
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
