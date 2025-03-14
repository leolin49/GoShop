package skiplist

type SkipListIterator struct {
	p *SkipListNode
}

func (s *Skiplist) NewIterator() *SkipListIterator {
	return &SkipListIterator{
		p: s.Head.Next[0],
	}
}

func (it *SkipListIterator) Next() *SkipListIterator {
	if it.p == nil {
		return nil
	}
	return &SkipListIterator{
		p: it.p.Next[0],
	}
}

func (it *SkipListIterator) Value() int {
	if it.p == nil {
		panic("skiplist iterator is out of bounds")
	}
	return it.p.Val
}

func (it *SkipListIterator) HasNext() bool {
	return it.p != nil
}

