// Copyright (c) 2025, Yufeng Lin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// 
// Package skiplist implements a skip list.
// 
package skiplist

type SkipListIterator[T comparable] struct {
	p *SkipListNode[T]
}

func (s *Skiplist[T]) NewIterator() *SkipListIterator[T] {
	return &SkipListIterator[T]{
		p: s.Head.Next[0],
	}
}

func (it *SkipListIterator[T]) Next() *SkipListIterator[T] {
	if it.p == nil {
		return nil
	}
	return &SkipListIterator[T]{
		p: it.p.Next[0],
	}
}

func (it *SkipListIterator[T]) Value() T {
	if it.p == nil {
		panic("skiplist iterator is out of bounds")
	}
	return it.p.Val
}

func (it *SkipListIterator[T]) HasNext() bool {
	return it.p != nil
}
