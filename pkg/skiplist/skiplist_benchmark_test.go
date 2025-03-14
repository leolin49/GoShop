// Copyright (c) 2025, Yufeng Lin. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package skiplist implements a skip list.
//
// Run the benchmark by `go test -bench=.` in package path.

package skiplist

import (
	"math/rand/v2"
	"testing"
)

const ElementNum = 100000

func BenchmarkSkipListKth(b *testing.B) {
	sl := NewSkiplist[int]()
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
	sl := NewSkiplist[int]()
	for range ElementNum {
		x := rand.Int()
		sl.Insert(x)
	}
}

func BenchmarkErase(b *testing.B) {
	sl := NewSkiplist[int]()
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
