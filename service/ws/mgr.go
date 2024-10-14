package ws

import (
	"fmt"
	"github.com/gorilla/websocket"
	"sync"
)

type WebSocketConnMgr struct {
	Wscs map[int64]*WebSocketConn
	lock sync.RWMutex
}

var WSCMgr *WebSocketConnMgr = &WebSocketConnMgr{Wscs: make(map[int64]*WebSocketConn)}

func (this *WebSocketConnMgr) AddConn(uid int64, wsc *WebSocketConn) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.Wscs[uid] = wsc
}

func (this *WebSocketConnMgr) RemoveConn(uid int64) {
	this.lock.Lock()
	defer this.lock.Unlock()
	delete(this.Wscs, uid)
	return
}

func (this *WebSocketConnMgr) GetConn(uid int64) *WebSocketConn {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.Wscs[uid]
}

// 广播消息
func (this *WebSocketConnMgr) BroadCast(data []byte) {
	this.lock.RLock()
	defer this.lock.RUnlock()
	for _, v := range this.Wscs {
		err := v.Conn.WriteMessage(websocket.BinaryMessage, data)
		if err != nil {
			fmt.Println("WriteMessage:", err)
			return
		}
	}
}

// 转发消息
func (this *WebSocketConnMgr) Transfer(rid int64, data []byte) {
	this.lock.RLock()
	defer this.lock.RUnlock()
	err := this.Wscs[rid].Conn.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		fmt.Println("Transfer.this.wscs[rid].Conn.WriteMessage:", err)
		return
	}
}
