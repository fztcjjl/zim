package ztimer

type TimerTask struct {
	delayMs        int64
	f              func()
	timerTaskEntry *TimerTaskEntry
}

func (t *TimerTask) Cancel() {
	if t.timerTaskEntry != nil {
		t.timerTaskEntry.Remove()
		t.timerTaskEntry = nil
	}
}

func (t *TimerTask) SetTimerTaskEntry(entry *TimerTaskEntry) {
	if t.timerTaskEntry != nil && t.timerTaskEntry != entry {
		t.timerTaskEntry.Remove()
	}
	t.timerTaskEntry = entry
}

func (t *TimerTask) GetTimerTaskEntry() *TimerTaskEntry {
	return t.timerTaskEntry
}

func (t *TimerTask) Run() {
	t.f()
}
