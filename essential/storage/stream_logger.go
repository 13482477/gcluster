package storage

import (
	"time"
	//"sync/atomic"
	"os"
	//"unsafe"
	"bytes"
	"container/list"
	"errors"
	log "github.com/sirupsen/logrus"
	"io"
	"path/filepath"
	"poseidon/essential/utils"
	"sync"
)

/*
	usage:
	in global app, or app main, call the register your logger topic,
	then in other place, just call GetStreamLogger().CollectionLog();

eg:
	// In main.go setup logger

	logger := storage.GetStreamLogger()

	// Follow One Need Call Once
	logger.SetFallbackDir(".")
	logger.RegisterTopic(kTopic, ".", 60, "win")
	//logger.RegisterTopic("xxx", ".", 60, "click")
	//register other
	logger.StartLogging()


	// in Other Module
	import storage

	logger := storage.GetStreamLogger()

	logger.CollectionLog(topic, data)
*/

const (
	kfallback_topic = "fallback_topic"
)

var (
	logInstance *StreamLogger = newStreamLogger()
)

type LogItem struct {
	logTopic      string
	path          string
	namePrefix    string
	repeator      *time.Ticker
	flushInterval int64

	//need load/set by atomic opts
	queuelist        *list.List
	locker           *sync.Mutex
	fallbackFile     *os.File
	inProgress       bool
	wg               *sync.WaitGroup
	exitChannel      chan bool
	skipEmptyContent bool
}

func NewLogItem(path string, prefix string, interval int64) *LogItem {
	n := new(LogItem)

	n.path = path
	n.namePrefix = prefix
	n.flushInterval = interval
	n.queuelist = list.New()
	n.inProgress = false
	n.wg = new(sync.WaitGroup)
	n.exitChannel = make(chan bool)
	n.skipEmptyContent = true
	n.locker = new(sync.Mutex)
	return n
}

func (logitem *LogItem) AppendData(data string) {
	logitem.locker.Lock()
	logitem.queuelist.PushBack(data)
	logitem.locker.Unlock()
}

func (logitem *LogItem) LogHandle() {
	//atomic.StoreInt32(&logitem.inProgress, 1)
	logitem.inProgress = true
	logitem.wg.Add(1)

	for logitem.inProgress {
		select {
		case <-logitem.repeator.C:

			logitem.flushToFile()

		case <-logitem.exitChannel:

			logitem.inProgress = false
			logitem.repeator.Stop()
			logitem.flushToFile()
			//atomic.StoreInt32(&logitem.inProgress, 0)
			logitem.wg.Done()
			//default:
			//log.Debugf("Not Should Reached Topic: %s \n", logitem.logTopic)
		}
	}
}

func (logitem *LogItem) flushToFile() {

	newList := list.New()
	//oldQueue := atomic.SwapPointer((*unsafe.Pointer)(unsafe.Pointer(&logitem.queuelist)), unsafe.Pointer(newList)))
	logitem.locker.Lock()
	oldQueue := logitem.queuelist
	logitem.queuelist = newList
	logitem.locker.Unlock()

	if oldQueue == nil {
		return
	}

	var outstream bytes.Buffer
	for elem := oldQueue.Front(); elem != nil; elem = elem.Next() {
		outstream.WriteString(elem.Value.(string))
	}

	fallback_file := false
	if outstream.Len() == 0 && logitem.skipEmptyContent {
		return
	}
	file, err := logitem.CreateLogFile()
	if err != nil {
		log.Warnf("Topic %s Create LogFile Failed, error %s", logitem.logTopic, err.Error())
		file = logitem.fallbackFile
		fallback_file = true
	}
	file_name := file.Name()

	io.WriteString(file, outstream.String())
	file.Close()
	if !fallback_file {
		//TODO: rename file ?
		go func(name string) {
			os.Rename(name, name+".cp")
		}(file_name)
	}
}

func (logitem *LogItem) CreateFallbackFile() (*os.File, error) {
	absPath, _ := filepath.Abs(logitem.path)
	fullDir := filepath.Join(absPath, logitem.logTopic)
	if exists, _ := utils.PathExists(fullDir); exists == false {
		os.MkdirAll(fullDir, os.ModePerm)
	}

	fileName := filepath.Join(fullDir, logitem.namePrefix+"_fallback"+".log")
	var f *os.File
	var err error

	if f, err = os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err != nil {
		log.Errorf("Create LogFile %s Failed For Topic %s\n", fileName, logitem.logTopic)
	}
	return f, err
}

func (logitem *LogItem) CreateLogFile() (*os.File, error) {
	now := time.Now()
	absPath, _ := filepath.Abs(logitem.path)
	fullDir := filepath.Join(absPath, logitem.logTopic, now.Format("20060102"))
	if exists, _ := utils.PathExists(fullDir); exists == false {
		os.MkdirAll(fullDir, os.ModePerm)
	}

	fileName := filepath.Join(fullDir, logitem.namePrefix+now.Format("-150405")+".log")

	var f *os.File
	var err error

	if f, err = os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err != nil {
		log.Errorf("Create LogFile %s Failed For Topic %s\n", fileName, logitem.logTopic)
	}

	return f, err
}

func (logitem *LogItem) Start() {
	logitem.repeator = time.NewTicker(time.Duration(logitem.flushInterval) * time.Second)

	var err error
	logitem.fallbackFile, err = logitem.CreateFallbackFile()
	if err != nil {
		log.Panicf("Topic %s Create Fallback LogFile Faild, error: %s", err.Error())
	}

	go logitem.LogHandle()
}

func (logitem *LogItem) Stop() {
	if !logitem.inProgress {
		return
	}

	logitem.repeator.Stop()
	logitem.exitChannel <- false
	logitem.wg.Wait()

	log.Infof("Topic %s Exited \n", logitem.logTopic)
}

type StreamLogger struct {
	inLogging   bool
	fallbackDir string
	logItems    map[string]*LogItem
}

func GetStreamLogger() *StreamLogger {
	return logInstance
}

func newStreamLogger() *StreamLogger {
	logger := StreamLogger{
		inLogging:   false,
		fallbackDir: ".",
		logItems:    make(map[string]*LogItem),
	}
	return &logger
}

func (logger *StreamLogger) SetFallbackDir(dir string) {
	logger.fallbackDir = dir
}

func (logger *StreamLogger) RegisterTopic(topic string, path string, flushInterval uint32, namePrefix string) error {

	if nil != logger.logItems[topic] {
		var err error = errors.New("This Topic Has Been Registered")
		return err
	}
	newLogItem := NewLogItem(path, namePrefix, (int64)(flushInterval))
	newLogItem.logTopic = topic
	logger.logItems[topic] = newLogItem
	return nil
}

func (logger *StreamLogger) StartLogging() {
	if logger.inLogging {
		log.Errorln("StreamLogger Has Ben Started. Can't Start Twice")
		return
	}
	logger.RegisterTopic(kfallback_topic, logger.fallbackDir, 60, "fallback")

	for topic, logitem := range logger.logItems {
		log.Infof("Stream Logger Start Logging %s", topic)
		logitem.Start()
	}
	logger.inLogging = true
}

func (logger *StreamLogger) StopLogging() {
	if !logger.inLogging {
		return
	}

	for _, logitem := range logger.logItems {
		logitem.Stop()
	}
}

func (logger *StreamLogger) CollectionLog(topic string, data string) error {

	logItem, ok := logger.logItems[topic]
	if !ok {
		logItem = logger.logItems[kfallback_topic]
	}

	if logItem == nil {
		return errors.New("Bad Topic, And Fallback also Failed")
	}

	logItem.AppendData(data)
	return nil
}
