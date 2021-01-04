package window

import (
	"math"
	"sync"
	"time"
)

const (
	empty = -1
)

// SlidingWindowOption is for SlidingWindow.
type SlidingWindowOption struct {
	Capacity int
	Unit     time.Duration
}

// SlidingWindow contains a series of data.
type SlidingWindow struct {
	m       *sync.RWMutex
	buckets []bucket
	head    int
	tail    int
	size    int

	// option
	capacity int
	unit     time.Duration
}

// NewSlidingWindow creates an instance.
func NewSlidingWindow(opt *SlidingWindowOption) *SlidingWindow {
	if opt.Capacity <= 0 {
		return nil
	}

	return &SlidingWindow{
		m:        &sync.RWMutex{},
		buckets:  make([]bucket, opt.Capacity),
		head:     empty,
		tail:     empty,
		size:     0,
		capacity: opt.Capacity,
		unit:     opt.Unit,
	}
}

// Add the value into the latest window.
func (w *SlidingWindow) Add(value int64) {
	w.m.Lock()
	defer w.m.Unlock()
	align := Align(time.Now(), w.unit)
	w.add(value, align)
}

// GetStartTime returns the window start.
func (w *SlidingWindow) GetStartTime() time.Time {
	w.m.RLock()
	defer w.m.RUnlock()

	bound := w.getBoundByNow()
	current := w.head
	for i := 0; i < w.size; i++ {
		if b := w.buckets[current]; b.start.After(bound) {
			return b.start
		}
		current = (current + 1) % w.capacity
	}
	return zeroTime
}

func (w *SlidingWindow) add(value int64, align time.Time) {
	var b *bucket = w.getTail()
	if w.isEmpty() || !b.start.Equal(align) {
		w.enqueue()
		b = w.getTail()
		b.reset(align, w.unit)
	}
	b.Add(value)
	w.clean(align)
}

func (w *SlidingWindow) clean(align time.Time) {
	bound := w.getBound(align)
	for !w.isEmpty() {
		if w.getHead().start.After(bound) {
			break
		}
		w.dequeue()
	}
}

func (w *SlidingWindow) enqueue() {
	if w.isEmpty() {
		// Empty
		w.head = 0
		w.tail = 0
		w.size++
	} else if w.isFull() {
		// Full
		w.head = (w.head + 1) % w.capacity
		w.tail = (w.tail + 1) % w.capacity
	} else {
		w.tail = (w.tail + 1) % w.capacity
		w.size++
	}
}

func (w *SlidingWindow) dequeue() {
	if w.isEmpty() {
		return
	} else if w.size == 1 {
		w.head = empty
		w.tail = empty
	} else {
		w.head = (w.head + 1) % w.capacity
		w.size--
	}
}

func (w *SlidingWindow) getBound(align time.Time) time.Time {
	// bound = align - capacity
	// bound < valid window time
	return align.Add(w.unit * time.Duration(-w.capacity))
}

func (w *SlidingWindow) getBoundByNow() time.Time {
	align := Align(time.Now(), w.unit)
	return w.getBound(align)
}

func (w *SlidingWindow) getHead() *bucket {
	if w.isEmpty() {
		return nil
	}
	return &w.buckets[w.head]
}

func (w *SlidingWindow) getTail() *bucket {
	if w.isEmpty() {
		return nil
	}
	return &w.buckets[w.tail]
}

func (w *SlidingWindow) isEmpty() bool {
	return w.size == 0
}

func (w *SlidingWindow) isFull() bool {
	return w.size == w.capacity
}

/*
 * Statistic operations
 */

func (w *SlidingWindow) iter() <-chan int64 {
	c := make(chan int64, w.size)
	if w.isEmpty() {
		close(c)
	} else {
		go func() {
			w.m.RLock()
			defer w.m.RUnlock()
			bound := w.getBoundByNow()
			current := w.head
			for i := 0; i < w.size; i++ {
				if b := w.buckets[current]; b.start.After(bound) {
					c <- w.buckets[current].GetValue()
				}
				current = (current + 1) % w.capacity
			}
			close(c)
		}()
	}

	return c
}

// Avg returns average of values in window.
func (w *SlidingWindow) Avg() float64 {
	w.m.RLock()
	defer w.m.RUnlock()
	if w.isEmpty() {
		return 0
	}

	s := w.Sum()
	return float64(s) / float64(w.size)
}

// Max returns max of values in window.
func (w *SlidingWindow) Max() int64 {
	w.m.RLock()
	defer w.m.RUnlock()
	if w.isEmpty() {
		return 0
	}

	m := int64(math.MinInt64)
	for v := range w.iter() {
		if v > m {
			m = v
		}
	}
	return m
}

// Min returns min of values in window.
func (w *SlidingWindow) Min() int64 {
	w.m.RLock()
	defer w.m.RUnlock()
	if w.isEmpty() {
		return 0
	}

	m := int64(math.MaxInt64)
	for v := range w.iter() {
		if v < m {
			m = v
		}
	}
	return m
}

// Sum returns sum of values in window.
func (w *SlidingWindow) Sum() int64 {
	w.m.RLock()
	defer w.m.RUnlock()
	if w.isEmpty() {
		return 0
	}

	var s int64
	for v := range w.iter() {
		s += v
	}
	return s
}
