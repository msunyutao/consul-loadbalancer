package util

import (
	"bytes"
	"sort"
	"strconv"
	"sync"
	"time"
)

type AvgWatch struct {
	count int
	total float64
}

type Watch struct {
	stopChan    chan bool
	mutex       *sync.Mutex
	avg_mutex   *sync.Mutex
	content     map[string]int
	avg_content map[string]*AvgWatch
	logger      Logger
}

func (watcher *Watch) AddWatchValue(key string, value int) {
	watcher.mutex.Lock()
	defer watcher.mutex.Unlock()
	v, in := watcher.content[key]
	if in {
		v = v + value
	} else {
		v = value
	}

	watcher.content[key] = v
}

func (watcher *Watch) SetWatchValue(key string, value int) {
	watcher.mutex.Lock()
	defer watcher.mutex.Unlock()
	watcher.content[key] = value
}

func (watcher *Watch) AddAvgWatchValue(key string, value float64) {
	watcher.avg_mutex.Lock()
	defer watcher.avg_mutex.Unlock()
	v, in := watcher.avg_content[key]
	if !in {
		v = &AvgWatch{}
		watcher.avg_content[key] = v
	}

	v.total += value
	v.count++
}

func (watcher *Watch) RunWatch() {
	go func() {
		for {
			select {
			case <-time.Tick(1 * time.Minute):
				WriteDisk(watcher)
			case <-watcher.stopChan:
				return
			}
		}
	}()
}

func (watcher *Watch) Stop() {
	watcher.stopChan <- true
	WriteDisk(watcher)
}

func WriteDisk(watcher *Watch) {
	content := getContent(watcher)
	avg_content := getAvgContent(watcher)
	watcher.logger.Infof(content + avg_content)
}

func getContent(watcher *Watch) string {
	var buffer bytes.Buffer
	watcher.mutex.Lock()
	defer watcher.mutex.Unlock()

	keys := make([]string, len(watcher.content))
	i := 0
	for k := range watcher.content {
		keys[i] = k
		i++
	}
	sort.Strings(keys)

	for _, key := range keys {
		buffer.WriteString("\t")
		buffer.WriteString(key)
		buffer.WriteString(":")
		buffer.WriteString(strconv.Itoa(watcher.content[key]))
		watcher.content[key] = 0.0
	}

	return buffer.String()
}

func getAvgContent(watcher *Watch) string {
	var buffer bytes.Buffer
	watcher.avg_mutex.Lock()
	defer watcher.avg_mutex.Unlock()
	// avg
	keys := make([]string, len(watcher.avg_content))
	i := 0
	for k := range watcher.avg_content {
		keys[i] = k
		i++
	}
	sort.Strings(keys)

	for _, key := range keys {
		avg := 0
		if watcher.avg_content[key].count > 0 {
			avg = int(watcher.avg_content[key].total) / watcher.avg_content[key].count
		}
		buffer.WriteString("\t")
		buffer.WriteString(key)
		buffer.WriteString(":")
		buffer.WriteString(strconv.Itoa(avg))
		buffer.WriteString("\t")
		buffer.WriteString(key + "_count")
		buffer.WriteString(":")
		buffer.WriteString(strconv.Itoa(watcher.avg_content[key].count))
		watcher.avg_content[key].count = 0.0
		watcher.avg_content[key].total = 0.0
	}

	return buffer.String()
}

func NewWatch(log Logger) *Watch {
	watcher := Watch{}
	watcher.stopChan = make(chan bool)
	watcher.mutex = new(sync.Mutex)
	watcher.avg_mutex = new(sync.Mutex)
	watcher.content = make(map[string]int)
	watcher.avg_content = make(map[string]*AvgWatch)
	watcher.logger = log
	return &watcher
}
