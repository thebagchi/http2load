package main

import (
	"math"
	"sync"
)

// Gauge structure for holding gauge data
type Gauge struct {
	lock sync.Mutex
	sum  float64

	Count int
	Min   float64
	Max   float64
}

// NewGauge Creates a new gauge.
func NewGauge() *Gauge {
	gauge := &Gauge{}
	gauge.Reset()
	return gauge
}

// Reset Resets the gauge with inital values.
func (gauge *Gauge) Reset() {
	gauge.sum, gauge.Count, gauge.Min, gauge.Max = 0, 0, math.MaxFloat64, 0
}

// Mean get mean value from gauge.
func (gauge *Gauge) Mean() float64 {
	return gauge.sum / float64(gauge.Count)
}

// Add add new value to gauge.
func (gauge *Gauge) Add(n float64) {
	gauge.lock.Lock()
	defer gauge.lock.Unlock()
	if n < gauge.Min {
		gauge.Min = n
	}
	if n > gauge.Max {
		gauge.Max = n
	}
	gauge.Count++
	gauge.sum = gauge.sum + n
}
