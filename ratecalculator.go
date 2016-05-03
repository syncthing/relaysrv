package main

import (
	"sync/atomic"
	"time"
)

var rc *rateCalculator


type rateCalculator struct {
	rates     []int64
	prev      int64
	counter   *int64
	startTime time.Time
}

func newRateCalculator(keepIntervals int, interval time.Duration, counter *int64) *rateCalculator {
	r := &rateCalculator{
		rates:     make([]int64, keepIntervals),
		counter:   counter,
		startTime: time.Now(),
	}

	go r.updateRates(interval)

	return r
}

func (r *rateCalculator) updateRates(interval time.Duration) {
	for {
		now := time.Now()
		next := now.Truncate(interval).Add(interval)
		time.Sleep(next.Sub(now))

		cur := atomic.LoadInt64(r.counter)
		rate := int64(float64(cur-r.prev) / interval.Seconds())
		copy(r.rates[1:], r.rates)
		r.rates[0] = rate
		r.prev = cur
	}
}

func (r *rateCalculator) rate(periods int) int64 {
	var tot int64
	for i := 0; i < periods; i++ {
		tot += r.rates[i]
	}
	return tot / int64(periods)
}
