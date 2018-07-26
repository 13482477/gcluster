package monitor

import (
	"math"
	"time"
)

const (
	DistributionBucketGap  = 100
	DistributionBucketSize = 100
)

type CollectorMatrix struct {
	StartTime time.Time

	Seconds int64
	Value   int64
	Count   int64
	Max     int64
	Min     int64
	Avg     int64
	Qps     float32
	Dist50  int64
	Dist80  int64
	Dist90  int64
	Dist95  int64
	Dist99  int64
	Dist999 int64

	distribution [DistributionBucketSize]int64
}

func NewCollectorMatrix() *CollectorMatrix {
	return &CollectorMatrix{
		StartTime: time.Now(),
		Max:       0,
		Min:       math.MaxInt64,
	}
}

func (m *CollectorMatrix) Mark(value int64) {
	if m.Max < value {
		m.Max = value
	}

	if m.Min > value {
		m.Min = value
	}

	bucket := int(value / DistributionBucketGap)
	if bucket >= len(m.distribution) {
		bucket = len(m.distribution) - 1
	}

	m.Count++
	m.distribution[bucket]++
}

func (m *CollectorMatrix) Eval() {
	m.Seconds = int64(time.Since(m.StartTime).Seconds())

	if m.Count != 0 {
		m.Avg = m.Value / m.Count
	}

	if m.Seconds != 0 {
		m.Qps = float32(m.Count) / float32(m.Seconds)
	}

	m.Dist50 = m.getTimeDistribution(0.5)
	m.Dist80 = m.getTimeDistribution(0.8)
	m.Dist90 = m.getTimeDistribution(0.9)
	m.Dist95 = m.getTimeDistribution(0.95)
	m.Dist99 = m.getTimeDistribution(0.99)
	m.Dist999 = m.getTimeDistribution(0.999)
}

func (m *CollectorMatrix) getTimeDistribution(percentage float64) int64 {
	count := int64(float64(m.Count) * percentage)
	sumCount := int64(0)
	bucket := 0
	for bucket = 0; bucket != len(m.distribution) && sumCount < count; bucket++ {
		sumCount += m.distribution[bucket]
	}

	if bucket != 0 {
		bucket--
	}

	return int64(bucket) * int64(DistributionBucketGap)
}
