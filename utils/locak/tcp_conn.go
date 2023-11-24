package lock

import (
	"net"
	"sync"
)

var TcpClientMap = &TcpClientDataMap{}

type TcpClientDataMap struct {
	TcpClientMap map[int32]*net.TCPConn
	Lock         *sync.RWMutex
}

func NewTcpClientMap() *TcpClientDataMap {
	return &TcpClientDataMap{
		TcpClientMap: make(map[int32]*net.TCPConn, 0),
		Lock:         &sync.RWMutex{},
	}
}

func (t TcpClientDataMap) Get(key int32) *net.TCPConn {
	t.Lock.RLock()
	defer t.Lock.RUnlock()
	return t.TcpClientMap[key]
}

func (t TcpClientDataMap) Set(key int32, conn *net.TCPConn) {
	t.Lock.Lock()
	defer t.Lock.Unlock()
	t.TcpClientMap[key] = conn
}

func (t TcpClientDataMap) Del(key int32) {
	t.Lock.Lock()
	defer t.Lock.Unlock()
	delete(t.TcpClientMap, key)
}
