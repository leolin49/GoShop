package skiplist

import (
	"math/rand/v2"
	"testing"
)

const ElementNum = 100000

func BenchmarkSkipListKth(b *testing.B) {
	sl := NewSkiplist()
	for range ElementNum {
		x := rand.Int()
		sl.Insert(x)
	}
	b.ResetTimer()
	for range ElementNum {
		i := rand.IntN(ElementNum) + 1
		_, _ = sl.Kth(i)
	}
}

func BenchmarkInsert(b *testing.B) {
	sl := NewSkiplist()
	for range ElementNum {
		x := rand.Int()
		sl.Insert(x)
	}
}

func BenchmarkErase(b *testing.B) {
	sl := NewSkiplist()
	for range ElementNum {
		x := rand.Int()
		sl.Insert(x)
	}
	b.ResetTimer()
	for range ElementNum {
		i := rand.IntN(sl.Length()) + 1
		_ = sl.Remove(i)
	}
}
