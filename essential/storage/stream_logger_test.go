package storage

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"
)

const (
	kTopic = "test_topic"
)

var wg sync.WaitGroup = sync.WaitGroup{}

func Send100WLog(logger *StreamLogger) {
	items := 1000000
	index := 0

	for index < items {
		logger.CollectionLog(kTopic, "this is a test log message\n")
		index++
	}
	wg.Done()
}

func TestNewStreamLogger(t *testing.T) {

	logger := GetStreamLogger()

	// Follow One Need Call Once
	logger.SetFallbackDir(".")
	logger.RegisterTopic(kTopic, ".", 3, "win")
	//logger.RegisterTopic("xxx", ".", 60, "click")
	//register other
	logger.StartLogging()
	// Note End

	start := time.Now()

	for i := 1; i <= runtime.NumCPU()*4; i++ {
		wg.Add(1)
		go Send100WLog(logger)
	}
	wg.Wait()

	stop := time.Now()

	fmt.Printf("Start at %d, Stop at %d, Spend %d Sencond, lThreads: %d\n", start.Unix(), stop.Unix(), stop.Unix()-start.Unix(), runtime.NumCPU()*4)
	logger.StopLogging()
	fmt.Printf("Programing Exit\n")
}
