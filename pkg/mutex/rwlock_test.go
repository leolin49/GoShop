package mutex

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestRWMutex(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		var m RWMutex
		var counter int32

		// 测试读锁不互斥
		var wg sync.WaitGroup
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				m.RLock()
				atomic.AddInt32(&counter, 1)
				m.RUnLock()
				wg.Done()
			}()
		}
		wg.Wait()

		if counter != 10 {
			t.Errorf("Expected counter=10, got %d", counter)
		}
	})

	t.Run("WriteExclusive", func(t *testing.T) {
		var m RWMutex
		var counter int32
		var writeHappened bool

		// 先获取读锁
		m.RLock()

		// 启动写goroutine
		go func() {
			m.WLock()
			writeHappened = true
			counter = 100
			m.WUnLock()
		}()

		// 确保写操作被阻塞
		time.Sleep(100 * time.Millisecond)
		if writeHappened {
			t.Error("Write happened while read lock was held")
		}

		// 释放读锁后写操作应该能继续
		m.RUnLock()
		time.Sleep(100 * time.Millisecond)
		if !writeHappened || counter != 100 {
			t.Error("Write lock didn't work as expected")
		}
	})

	t.Run("ReadBlocksWrite", func(t *testing.T) {
		var m RWMutex
		var readFinished bool

		// 先获取写锁
		m.WLock()

		// 启动读goroutine
		go func() {
			m.RLock()
			readFinished = true
			m.RUnLock()
		}()

		// 确保读操作被阻塞
		time.Sleep(100 * time.Millisecond)
		if readFinished {
			t.Error("Read happened while write lock was held")
		}

		// 释放写锁后读操作应该能继续
		m.WUnLock()
		time.Sleep(100 * time.Millisecond)
		if !readFinished {
			t.Error("Read lock didn't work as expected")
		}
	})

	t.Run("MultipleWrite", func(t *testing.T) {
		var m RWMutex
		var counter int32

		// 测试写锁互斥
		for i := 0; i < 10; i++ {
			go func() {
				m.WLock()
				atomic.AddInt32(&counter, 1)
				time.Sleep(10 * time.Millisecond)
				atomic.AddInt32(&counter, -1)
				m.WUnLock()
			}()
		}

		time.Sleep(500 * time.Millisecond)
		if counter != 0 {
			t.Errorf("Write locks not exclusive, counter=%d", counter)
		}
	})

	t.Run("ConcurrentReadWrite", func(t *testing.T) {
		var m RWMutex
		var data int
		var wg sync.WaitGroup

		// 启动多个读写goroutine
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				m.RLock()
				_ = data // 读操作
				m.RUnLock()
			}()
		}

		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				m.WLock()
				data++ // 写操作
				m.WUnLock()
			}()
		}

		wg.Wait()
		if data != 10 {
			t.Errorf("Expected data=10, got %d", data)
		}
	})
}

func BenchmarkRWMutex(b *testing.B) {
	var m RWMutex
	var data int

	b.Run("Read", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				m.RLock()
				_ = data
				m.RUnLock()
			}
		})
	})

	b.Run("Write", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				m.WLock()
				data++
				m.WUnLock()
			}
		})
	})

	// b.Run("Mixed", func(b *testing.B) {
	// 	b.RunParallel(func(pb *testing.PB) {
	// 		for pb.Next() {
	// 			if pb.id()%10 == 0 { // 10%写操作
	// 				m.WLock()
	// 				data++
	// 				m.WUnLock()
	// 			} else { // 90%读操作
	// 				m.RLock()
	// 				_ = data
	// 				m.RUnLock()
	// 			}
	// 		}
	// 	})
	// })
}
