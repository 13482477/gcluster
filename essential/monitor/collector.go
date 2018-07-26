package monitor

import (
	"encoding/json"
	"math/rand"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type CollectorData struct {
	Data map[string]*CollectorMatrix
}

type ChannelCollector struct {
	channel         []chan *MatrixItem
	collectorMatrix *CollectorData
	lock            *sync.Mutex
	freeze          *CollectorData

	logTimeFre time.Duration
}

func NewCollectorData() *CollectorData {
	return &CollectorData{
		Data: make(map[string]*CollectorMatrix),
	}
}

func (data *CollectorData) GetMatrix(name string) *CollectorMatrix {
	m, ok := data.Data[name]
	if ok {
		return m
	}
	m = NewCollectorMatrix()
	data.Data[name] = m
	return m
}

func NewChannelCollector(channelCount int, queueSize int, logTimeGap time.Duration) *ChannelCollector {
	c := &ChannelCollector{
		channel:         make([]chan *MatrixItem, channelCount),
		collectorMatrix: NewCollectorData(),
		lock:            new(sync.Mutex),
		logTimeFre:      logTimeGap,
	}

	for i := 0; i != len(c.channel); i++ {
		c.channel[i] = make(chan *MatrixItem, queueSize)
	}

	return c
}

func (c *ChannelCollector) Run() {
	for _, mc := range c.channel {
		go func(mc chan *MatrixItem) {
			for m := range mc {
				c.innerPush(m)
			}
		}(mc)
	}

	ticker := time.NewTicker(c.logTimeFre * time.Duration(1))
	go func() {
		for range ticker.C {
			c.freeze = c.innerGet()

			data, _ := json.Marshal(c.freeze)
			log.Info(string(data))
		}
	}()
}

func (c *ChannelCollector) Push(item *MatrixItem) error {
	c.channel[rand.Int31n(int32(len(c.channel)))] <- item
	return nil
}

func (c *ChannelCollector) Get() *CollectorData {
	return c.freeze
}

func (c *ChannelCollector) innerGet() *CollectorData {
	c.lock.Lock()
	freeze := c.collectorMatrix
	c.collectorMatrix = NewCollectorData()
	c.lock.Unlock()

	for _, m := range freeze.Data {
		m.Eval()
	}

	return freeze
}

func (c *ChannelCollector) innerPush(item *MatrixItem) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.collectorMatrix.GetMatrix(item.Name).Mark(int64(item.Value))
}
