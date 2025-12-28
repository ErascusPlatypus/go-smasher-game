package helpers

import "time"

type Timer struct {
    duration time.Duration
    start    time.Time
    active   bool
}

func NewTimer(d time.Duration) *Timer {
    return &Timer{
        duration: d,
        active:   false,
    }
}

func (t *Timer) Start() {
    t.start = time.Now()
    t.active = true
}

func (t *Timer) Reset() {
    t.start = time.Now()
}

func (t *Timer) Stop() {
    t.active = false
}


func (t *Timer) IsReady() bool {
    if !t.active {
        return false
    }

    return time.Since(t.start) >= t.duration
}

func (t *Timer) IsActive() bool {
    return t.active
}