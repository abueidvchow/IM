package ws

import (
	"fmt"
	"time"
)

type HeartBeatChecker struct {
	interval time.Duration     // 心跳检测时间间隔
	quit     chan struct{}     // 退出信号
	wsc      *WebSocketConnMgr // 所属服务器
}

func NewHeartBeatChecker(interval time.Duration, wsc *WebSocketConnMgr) *HeartBeatChecker {
	return &HeartBeatChecker{
		interval: interval,
		quit:     make(chan struct{}, 1),
		wsc:      wsc,
	}
}

func (h *HeartBeatChecker) Start() {
	fmt.Println("HeartBeatChecker start")

	ticker := time.NewTicker(h.interval)
	for {
		select {
		case <-ticker.C:
			h.Check()
		case <-h.quit:
			ticker.Stop()
			return
		}
	}
}

func (h *HeartBeatChecker) Check() {
	fmt.Println("HeartBeatChecker check", time.Now().Format("2006-01-02 15:04:05"))
	conns := h.wsc.GetAllConn()
	for _, conn := range conns {
		if !conn.IsAlive() {
			conn.Close()
		}
	}
}

func (h *HeartBeatChecker) Stop() {
	h.quit <- struct{}{}
}
