package lock

import (
	log "github.com/sirupsen/logrus"
	"sync"
)

var TerminalReadMap = &TerminalReadDataMap{}

type TerminalReadDataMap struct {
	TerminalReadMap map[uint32]chan string
	Lock            *sync.RWMutex
}

func NewTerminalReadMap() *TerminalReadDataMap {
	// 初始化通道并添加到map中
	terminalReadMap := make(map[uint32]chan string)
	return &TerminalReadDataMap{
		TerminalReadMap: terminalReadMap,
		Lock:            &sync.RWMutex{},
	}
}

func (t TerminalReadDataMap) Get(key uint32) chan string {
	t.Lock.RLock()
	defer t.Lock.RUnlock()
	log.Infof("TerminalReadDataMap Get key: %v, value: %v", key, t.TerminalReadMap[key])
	return t.TerminalReadMap[key]
}

func (t TerminalReadDataMap) Set(key uint32, data chan string) {
	t.Lock.Lock()
	defer t.Lock.Unlock()
	if data == nil {
		return
	}
	t.TerminalReadMap[key] = data
	log.Infof("TerminalReadDataMap Set key: %v, value: %v", key, t.TerminalReadMap[key])
}

func (t TerminalReadDataMap) Del(key uint32) {
	t.Lock.Lock()
	defer t.Lock.Unlock()
	delete(t.TerminalReadMap, key)
}
