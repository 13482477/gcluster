package monitor

import (
	"math/rand"
	"testing"
	"time"
)

func TestStartGloablCollector(t *testing.T) {
	StartGlobalCollector(NewChannelCollector(8, 1024, time.Duration(1)))
}

func testFunc() {
	ScopeMonitor("test", time.Now())
	time.Sleep(time.Millisecond * time.Duration(rand.Int31n(10)))
}

func TestMonitorScope(t *testing.T) {
	var c chan int

	go func() {
		for {
			testFunc()
		}
	}()

	ticker := time.NewTicker(time.Second * time.Duration(30))
	go func() {
		for range ticker.C {
			data := GetGlobalCollector().Get()
			if data != nil {
				t.Log(data)
			}
		}
	}()

	<-c
}
