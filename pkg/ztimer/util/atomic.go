package util

import "sync/atomic"

type AtomicInteger struct {
	val int64
}

func NewAtomicInteger() *AtomicInteger {
	return &AtomicInteger{0}
}

func (t *AtomicInteger) IncrementAndGet() int64 {
	return atomic.AddInt64(&t.val, 1)
}

func (t *AtomicInteger) DecrementAndGet() int64 {
	return atomic.AddInt64(&t.val, -1)
}

func (t *AtomicInteger) GetAndSet(newVal int64) int64 {
	return atomic.SwapInt64(&t.val, newVal)
}

func (t *AtomicInteger) Get() int64 {
	return atomic.LoadInt64(&t.val)
}
