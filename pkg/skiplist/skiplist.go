package skiplist

import "math/rand/v2"

const (
	maxLevel = 32
	pFactor  = 0.25
)

type SkipListNode struct {
	Val  int
	Next []*SkipListNode
	Span []int
}

type Skiplist struct {
	Level int
	Head  *SkipListNode
}

func NewSkiplist() *Skiplist {
	return &Skiplist{
		Level: 0,
		Head: &SkipListNode{
			Val:  -1,
			Next: make([]*SkipListNode, maxLevel),
			Span: make([]int, maxLevel),
		},
	}
}

func (Skiplist) randomLevel() int {
	lv := 1
	for lv < maxLevel && rand.Float64() < pFactor {
		lv++
	}
	return lv
}

func (s *Skiplist) Search(target int) bool {
	cur := s.Head
	for lv := s.Level - 1; lv >= 0; lv-- {
		for cur.Next[lv] != nil && target < cur.Next[lv].Val {
			cur = cur.Next[lv]
		}
	}
	cur = cur.Next[0] 
	return cur != nil && cur.Val == target
}

func (s *Skiplist) Add(num int) {
	update := make([]*SkipListNode, maxLevel)
	rank := make([]int, maxLevel)
	for i := range update {
		update[i] = s.Head
	}
	cur := s.Head
	for lv := s.Level - 1; lv >= 0; lv-- {
		if lv != s.Level - 1 {
			rank[lv] = rank[lv+1]
		}
		for cur.Next[lv] != nil && num < cur.Next[lv].Val {
			rank[lv] += cur.Span[lv]
			cur = cur.Next[lv]
		}
		update[lv] = cur
	}
	level := s.randomLevel()
	s.Level = max(s.Level, level)
	newNode := &SkipListNode{
		Val:  num,
		Next: make([]*SkipListNode, level),
	}
	for i, node := range update[:level] {
		newNode.Next[i] = node.Next[i]
		node.Next[i] = newNode

		// NOTE: need test
		newNode.Span[i] = node.Span[i] - (rank[0] - rank[i])
		node.Span[i] = rank[0] - rank[i] + 1
	}
	for i := level; i < s.Level; i++ {
		update[i].Span[i]++
	}
}

func (s *Skiplist) Erase(num int) bool {
	update := make([]*SkipListNode, maxLevel)
	for i := range update {
		update[i] = s.Head
	}
	cur := s.Head
	for lv := maxLevel - 1; lv >= 0; lv-- {
		for cur.Next[lv] != nil && num < cur.Next[lv].Val {
			cur = cur.Next[lv]
		}
		update[lv] = cur
	}
	cur = cur.Next[0]
	if cur == nil || cur.Val != num {
		return false
	}
	for lv := 0; lv < maxLevel && update[lv].Next[lv] == cur; lv++ {
		update[lv].Next[lv] = cur.Next[lv]
	}
	for s.Level > 1 && s.Head.Next[s.Level-1] == nil {
		s.Level--
	}
	return true
}

/**
 * Your Skiplist object will be instantiated and called as such:
 * obj := Constructor();
 * param_1 := obj.Search(target);
 * obj.Add(num);
 * param_3 := obj.Erase(num);
 */
