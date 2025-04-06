package skiplist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func NewSkiplistWithTest[T int]() *Skiplist[int] {
	sl := NewSkiplist[int]()
	sl.Insert(3)
	sl.Insert(1)
	sl.Insert(10)
	sl.Insert(8)
	sl.Insert(4)
	sl.Insert(6)
	sl.Insert(11)
	sl.Insert(15)
	return sl
}

func TestSkipListEmpty(t *testing.T) {
	sl := NewSkiplist[int]()
	assert.True(t, sl.Empty())
	sl.Insert(1)
	assert.False(t, sl.Empty())
}

func TestSkipListLength(t *testing.T) {
	sl := NewSkiplist[int]()
	assert.Equal(t, 0, sl.Length())
	sl.Insert(1)
	assert.Equal(t, 1, sl.Length())
	sl.Insert(1)
	sl.Insert(1)
	assert.Equal(t, 3, sl.Length())
}

func TestSkipListGetKth(t *testing.T) {
	sl := NewSkiplistWithTest()
	// the list should be: 1->3->4->6->8->10->11->15
	val, _ := sl.Kth(4)
	assert.Equal(t, 6, val)
	val, _ = sl.Kth(7)
	assert.Equal(t, 11, val)

	val, _ = sl.Min()
	assert.Equal(t, 1, val)
	val, _ = sl.Max()
	assert.Equal(t, 15, val)
}

func TestSkipListGetRange(t *testing.T) {
	sl := NewSkiplistWithTest()
	// the list should be: 1->3->4->6->8->10->11->15
	elements, _ := sl.Range(2, 6)
	assert.Equal(t, []int{3, 4, 6, 8, 10}, elements)
	elements, _ = sl.Range(8, 4)
	assert.Equal(t, []int{6, 8, 10, 11, 15}, elements)
}

func TestSkipListIndexRank(t *testing.T) {
	sl := NewSkiplistWithTest()
	// the list should be: 1->3->4->6->8->10->11->15
	idx, _ := sl.Index(3)
	assert.Equal(t, 2, idx)
	idx, _ = sl.Index(8)
	assert.Equal(t, 5, idx)
	idx, _ = sl.Index(15)
	assert.Equal(t, 8, idx)

	rank, _ := sl.Rank(11)
	assert.Equal(t, 2, rank)
	rank, _ = sl.Rank(1)
	assert.Equal(t, 8, rank)
	rank, _ = sl.Rank(10)
	assert.Equal(t, 3, rank)
}

func TestSkipListLowerUpper(t *testing.T) {
	sl := NewSkiplistWithTest()

	x, _ := sl.Lower(4)
	assert.Equal(t, 4, x)
	x, _ = sl.Lower(5)
	assert.Equal(t, 6, x)
	x, _ = sl.Upper(9)
	assert.Equal(t, 10, x)
}

func TestSkipListAll(t *testing.T) {
	sl := NewSkiplist[int]()

	sl.Insert(3)
	sl.Insert(1)
	sl.Insert(10)
	sl.Insert(8)
	sl.Insert(4)
	sl.Insert(6)
	sl.Insert(11)
	sl.Insert(15)

	assert.Equal(t, 8, sl.Length(), "[SkipList] get len error")

	// 测试查找功能
	var found bool
	found = sl.Exist(11)
	assert.True(t, found, "[SkipList] find element failed")

	found = sl.Exist(6)
	assert.True(t, found, "[SkipList] find element failed")

	// 测试删除功能
	err := sl.Erase(6)
	assert.Nil(t, err, "[SkipList] erase element failed")
	assert.Equal(t, 7, sl.Length(), "[SkipList] get len error")

	found = sl.Exist(6)
	assert.False(t, found, "[SkipList] find element failed")

	// 测试删除不存在的节点
	err = sl.Erase(100)
	assert.NotNil(t, err)
	assert.Equal(t, 7, sl.Length(), "[SkipList] get len error")

	// 测试查找不存在的节点
	found = sl.Exist(100)
	assert.False(t, found, "[SkipList] find non-existent element failed")

	// 查找第k个的节点
	val, err := sl.Kth(1)
	if err != nil {
		t.Error("[SkipList] Get Kth node failed")
	}
	assert.Equal(t, 1, val, "[SkipList] Get Kth node failed")

	val, err = sl.Kth(6)
	if err != nil {
		t.Error("[SkipList] Get Kth node failed")
	}
	assert.Equal(t, 11, val, "[SkipList] Get Kth node failed")
}

type Person struct {
	Name  string
	Age   int
	Score float64
}

func CompareByAge(a, b Person) int {
	if a.Age < b.Age {
		return -1
	}
	if a.Age > b.Age {
		return 1
	}
	return 0
}

func TestSkipListStruct(t *testing.T) {
	sl := NewSkiplistWithCmp(CompareByAge)
	people := []Person{
		{"Alice", 25, 88.5},
		{"Bob", 30, 92.0},
		{"Charlie", 22, 76.5},
	}
	for _, p := range people {
		sl.Insert(p)
	}
	// 查询年龄最小的人
	youngest, _ := sl.Min()
	assert.Equal(t, 22, youngest.Age, "[SkipList] get element error")

	// 查找年龄>=25的人
	p, _ := sl.Lower(Person{Age: 25})
	assert.Equal(t, 25, p.Age, "[SkipList] get element error")

	// 查找年龄>25的人
	p, _ = sl.Upper(Person{Age: 25})
	assert.Equal(t, 30, p.Age, "[SkipList] get element error")
}
