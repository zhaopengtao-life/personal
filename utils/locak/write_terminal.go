package lock

import (
	log "github.com/sirupsen/logrus"
	"sync"
)

var TerminalWriteMap = &TerminalWriteDataMap{}

type TerminalWriteDataMap struct {
	TerminalWriteMap map[uint32]chan string
	Lock             *sync.RWMutex
}

func NewTerminalWriteMap() *TerminalWriteDataMap {
	// 初始化通道并添加到map中
	terminalWriteMap := make(map[uint32]chan string)
	return &TerminalWriteDataMap{
		TerminalWriteMap: terminalWriteMap,
		Lock:             &sync.RWMutex{},
	}
}

func (t TerminalWriteDataMap) Get(key uint32) chan string {
	t.Lock.RLock()
	defer t.Lock.RUnlock()
	log.Infof("TerminalWriteDataMap Get key: %v, value: %v", key, t.TerminalWriteMap[key])
	return t.TerminalWriteMap[key]
}

func (t TerminalWriteDataMap) Set(key uint32, data chan string) {
	t.Lock.Lock()
	defer t.Lock.Unlock()
	if data == nil {
		return
	}
	t.TerminalWriteMap[key] = data
	log.Infof("TerminalWriteDataMap Set key: %v, value: %v", key, t.TerminalWriteMap[key])
}

func (t TerminalWriteDataMap) Del(key uint32) {
	t.Lock.Lock()
	defer t.Lock.Unlock()
	delete(t.TerminalWriteMap, key)
}
