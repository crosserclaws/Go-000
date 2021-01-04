package window

import (
	"time"
)

var (
	zeroTime   = time.Time{}
	zeroBucket = bucket{}
)

type bucket struct {
	start time.Time
	end   time.Time
	value int64
}

func newBucket(t time.Time, unit time.Duration) *bucket {
	start := Align(t, unit)
	return &bucket{
		start: start,
		end:   start.Add(unit),
	}
}

func (b *bucket) Add(n int64) {
	b.value += n
}

func (b *bucket) GetValue() int64 {
	return b.value
}

func (b *bucket) reset(t time.Time, unit time.Duration) {
	start := Align(t, unit)
	b.start = start
	b.end = start.Add(unit)
	b.value = 0
}

// Align returns a new aligned time value.
func Align(t time.Time, unit time.Duration) time.Time {
	return t.Truncate(unit)
}
