package ztimer

import (
	"context"
	"github.com/fztcjjl/zim/pkg/ztimer/dqueue"
	"github.com/fztcjjl/zim/pkg/ztimer/util"
	"time"
)

type Timer struct {
	tickMs      int64
	wheelSize   int
	taskCounter *util.AtomicInteger
	delayQueue  *dqueue.DelayQueue
	timingWheel *TimingWheel
	ctx         context.Context
	cancel      context.CancelFunc
}

func NewTimer(tickMs int64, wheelSize int) *Timer {
	t := new(Timer)
	t.tickMs = tickMs
	t.wheelSize = wheelSize
	t.taskCounter = util.NewAtomicInteger()
	t.delayQueue = dqueue.NewDelayQueue()
	startMs := time.Now().UnixNano() / int64(time.Millisecond)
	t.timingWheel = NewTimingWheel(tickMs, wheelSize, startMs, t.taskCounter, t.delayQueue)
	t.ctx, t.cancel = context.WithCancel(context.Background())

	return t
}

func (t *Timer) AfterFunc(d time.Duration, f func()) *TimerTask {
	delayMs := int64(d / time.Millisecond)

	task := &TimerTask{
		delayMs: delayMs,
		f:       f,
	}

	entry := NewTimerTaskEntry(task, task.delayMs+util.GetTimeMs())

	t.addTimerTaskEntry(entry)

	return task
}

func (t *Timer) addTimerTaskEntry(entry *TimerTaskEntry) {
	if !t.timingWheel.Add(entry) {
		// Already expired or cancelled
		if !entry.Cancelled() {
			// TODO: goroutine pool
			go func() {
				entry.timerTask.Run()
			}()
		}
	}
}

func (t *Timer) reinsert(entry *TimerTaskEntry) {
	t.addTimerTaskEntry(entry)
}

func (t *Timer) Start() {
	go func() {
		for {
			d := t.delayQueue.Take(t.ctx)
			if d == nil {
				break
			}
			bucket := d.(*TimerTaskList)
			t.timingWheel.AdvanceClock(bucket.GetExpiration())
			bucket.Flush(t.reinsert)
		}
	}()
}

func (t *Timer) Stop() {
	if t.cancel != nil {
		t.cancel()
	}
}

func (t *Timer) Size() int64 {
	return t.taskCounter.Get()
}
