package dqueue

import (
	"context"
	"github.com/fztcjjl/zim/pkg/ztimer/pqueue"
	"sync"
	"sync/atomic"
	"time"
)

type Delayed interface {
	pqueue.Comparable
	GetDelay() int64
}

type DelayQueue struct {
	sync.Mutex
	q        *pqueue.PriorityQueue
	sleeping int32
	wakeupC  chan struct{}
}

func NewDelayQueue() *DelayQueue {
	return &DelayQueue{
		q:       pqueue.NewPriorityQueue(1),
		wakeupC: make(chan struct{}),
	}
}

func (dq *DelayQueue) Offer(d Delayed) {
	dq.Lock()
	dq.q.Push(d)
	first := dq.q.Front()
	dq.Unlock()

	if first.(Delayed) == d {
		if atomic.CompareAndSwapInt32(&dq.sleeping, 1, 0) {
			dq.wakeupC <- struct{}{}
		}
	}
}

func (dq *DelayQueue) Take(ctx context.Context) Delayed {
	for {
		dq.Lock()

		var first Delayed
		var delay int64
		if e := dq.q.Front(); e != nil {
			first = e.(Delayed)
			delay = first.GetDelay()
			if delay <= 0 {
				dq.q.Remove(0)
				dq.Unlock()
				return first
			}
		}

		atomic.StoreInt32(&dq.sleeping, 1)

		dq.Unlock()

		if first == nil {
			select {
			case <-dq.wakeupC:
				continue
			case <-ctx.Done():
				return nil
			}
		}

		if first != nil {
			select {
			case <-dq.wakeupC:
				continue
			case <-ctx.Done():
				return nil
			case <-time.After(time.Duration(delay) * time.Millisecond):
				if atomic.SwapInt32(&dq.sleeping, 0) == 0 {
					<-dq.wakeupC
				}
				continue
			}
		}
	}
}

func (dq *DelayQueue) Poll() Delayed {
	dq.Lock()
	defer dq.Unlock()

	if e := dq.q.Front(); e != nil {
		first := e.(Delayed)
		if first.GetDelay() <= 0 {
			dq.q.Remove(0)
			return first
		}
	}

	return nil
}
