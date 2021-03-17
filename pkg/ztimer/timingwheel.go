package ztimer

import (
	"github.com/fztcjjl/zim/pkg/ztimer/dqueue"
	"github.com/fztcjjl/zim/pkg/ztimer/util"
	"sync/atomic"
)

type TimingWheel struct {
	tickMs        int64
	wheelSize     int
	interval      int64
	currentTime   int64
	overflowWheel *TimingWheel
	taskCounter   *util.AtomicInteger
	queue         *dqueue.DelayQueue
	buckets       []*TimerTaskList
}

func NewTimingWheel(tickMs int64, wheelSize int, startMs int64, taskCounter *util.AtomicInteger, queue *dqueue.DelayQueue) *TimingWheel {
	buckets := make([]*TimerTaskList, wheelSize)
	for i := range buckets {
		buckets[i] = NewTimerTaskList(taskCounter)
	}
	return &TimingWheel{
		tickMs:        tickMs,
		wheelSize:     wheelSize,
		interval:      tickMs * int64(wheelSize),
		currentTime:   startMs - (startMs % tickMs),
		overflowWheel: nil,
		taskCounter:   taskCounter,
		queue:         queue,
		buckets:       buckets,
	}
}

func (tw *TimingWheel) Add(entry *TimerTaskEntry) bool {
	expiration := entry.expirationMs
	if entry.Cancelled() {
		// Canceled
		return false
	} else if expiration < tw.currentTime+tw.tickMs {
		// Already expired
		return false
	} else if expiration < tw.currentTime+tw.interval {
		// Put in its own bucket
		virtualId := expiration / tw.tickMs
		bucket := tw.buckets[int(virtualId)%tw.wheelSize]
		bucket.Add(entry)

		// Set the bucket expiration time
		if bucket.SetExpiration(virtualId * tw.tickMs) {
			// The bucket needs to be enqueued because it was an expired bucket
			// We only need to enqueue the bucket when its expiration time has changed, i.e. the wheel has advanced
			// and the previous buckets gets reused; further calls to set the expiration within the same wheel cycle
			// will pass in the same value and hence return false, thus the bucket with the same expiration will not
			// be enqueued multiple times.

			tw.queue.Offer(bucket)
		}
		return true
	} else {
		if tw.overflowWheel == nil {
			tw.addOverflowWheel()
		}

		return tw.overflowWheel.Add(entry)
	}
}

func (tw *TimingWheel) addOverflowWheel() {
	if tw.overflowWheel == nil {
		tw.overflowWheel = NewTimingWheel(tw.interval, tw.wheelSize, tw.currentTime, tw.taskCounter, tw.queue)
	}
}

func (tw *TimingWheel) AdvanceClock(timeMs int64) {
	currentTime := atomic.LoadInt64(&tw.currentTime)
	if timeMs >= tw.currentTime+tw.tickMs {
		currentTime = timeMs - (timeMs % tw.tickMs)
		atomic.StoreInt64(&tw.currentTime, currentTime)
		if tw.overflowWheel != nil {
			tw.overflowWheel.AdvanceClock(currentTime)
		}
	}
}
