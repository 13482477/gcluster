package monitor

import (
	"time"
)

type MatrixItem struct {
	Name  string
	Value int
}

type Collector interface {
	Run()
	Push(item *MatrixItem) error
	Get() *CollectorData
}

var globalCollector Collector = nil

func GetGlobalCollector() Collector {
	return globalCollector
}

func StartGlobalCollector(collector Collector) {
	globalCollector = collector
	GetGlobalCollector().Run()
}

func ScopeMonitor(name string, start time.Time) {
	if nil != GetGlobalCollector() {
		GetGlobalCollector().Push(&MatrixItem{
			Name:  name,
			Value: int(time.Since(start).Nanoseconds() / 1000),
		})
	}
}

func Push(name string, value int) {
	if nil != GetGlobalCollector() {
		GetGlobalCollector().Push(&MatrixItem{
			Name:  name,
			Value: value,
		})
	}
}
