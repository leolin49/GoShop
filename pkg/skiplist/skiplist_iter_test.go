package skiplist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSkipListIter(t *testing.T) {
	sl := NewSkiplist[int]()

	sl.Insert(3)
	sl.Insert(1)
	sl.Insert(10)
	sl.Insert(8)
	sl.Insert(4)
	sl.Insert(6)
	sl.Insert(11)
	sl.Insert(15)

	res := []int{}
	for it := sl.NewIterator(); !it.End(); it = it.Next() {
		res = append(res, it.Value())
	}
	assert.Equal(t, []int{1, 3, 4, 6, 8, 10, 11, 15}, res)
}
