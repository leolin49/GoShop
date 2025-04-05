package mutex

import (
	"sync/atomic"
	"time"
)

type RWMutex struct {
	state int32
}

func (m *RWMutex) RLock() {
	for {
		if v := atomic.LoadInt32(&m.state); v <= 0 {
			if atomic.CompareAndSwapInt32(&m.state, v, v-1) {
				return
			}
		}
		time.Sleep(10 * time.Nanosecond)
	}
}

func (m *RWMutex) RUnLock() {
	atomic.AddInt32(&m.state, 1)
}

func (m *RWMutex) WLock() {
	for {
		if v := atomic.LoadInt32(&m.state); v >= 0 {
			if atomic.CompareAndSwapInt32(&m.state, v, v+1) {
				return
			}
		}
		time.Sleep(10 * time.Nanosecond)
	}
}

func (m *RWMutex) WUnLock() {
	atomic.AddInt32(&m.state, -1)
}
