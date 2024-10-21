package ws

import (
	"IM/config"
	"IM/pkg/db"
	"IM/pkg/protocol/pb"
	"fmt"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"net/http"
	"sync"
	"time"
)

// 连接ID
var connId int64 = 0

type WebSocketConn struct {
	Conn              *websocket.Conn
	WscMgr            *WebSocketConnMgr
	ConnID            int64       // 连接编号，通过对编号取余，能够让 Conn 始终进入同一个 worker，保持有序性
	UserID            int64       // 连接所属用户id
	sendChannel       chan []byte // 写消息队列
	closeChannel      chan bool   // 监听channel是否关闭
	isClosed          bool        // 标识
	isCloseMutex      sync.Mutex  // 保护 isClose 的锁
	lastHeartBeatTime time.Time   // 最后活跃时间
	heartMutex        sync.Mutex  // 保护最后活跃时间的锁
	maxClientId       int64       // 该连接收到的最大 clientId，确保消息的可靠性
	maxClientIdMutex  sync.Mutex  // 保护 maxClientId 的锁
}

func NewWebSocketConn(w http.ResponseWriter, r *http.Request, uid int64) (wsc *WebSocketConn, err error) {
	//升级成websocket
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		zap.L().Error("建立websocket协议失败：", zap.Error(err))
		return nil, err
	}

	wsc = &WebSocketConn{
		ConnID:            connId,
		Conn:              conn,                    // websocket连接
		WscMgr:            WSCMgr,                  // 所属管理
		UserID:            uid,                     // 该连接的当前用户
		sendChannel:       make(chan []byte, 1024), // 发送消息管道
		closeChannel:      make(chan bool, 1),      // 关闭管道
		isClosed:          false,
		lastHeartBeatTime: time.Now(),
	}
	connId++
	return
}

func (wsc *WebSocketConn) Start() {
	// 开启读，读取客户端发送给服务端，再由服务端转发过来的消息
	go wsc.StartReader()
	// 开启写，将消息发送给服务端，再由服务端去转发
	go wsc.StartWriter()
}

// websocket后台读取消息
func (wsc *WebSocketConn) StartReader() {
	for {
		_, data, err := wsc.Conn.ReadMessage()
		if err != nil {
			fmt.Println("conn.ReadMessage error:", err)
			if !wsc.isClosed {
				wsc.closeChannel <- true
			}
			return
		}
		wsc.HandlerMessage(data)
	}
}

// websocket后台发送消息
func (wsc *WebSocketConn) StartWriter() {
	for {
		select {
		case data := <-wsc.sendChannel:
			err := wsc.Conn.WriteMessage(websocket.BinaryMessage, data)
			if err != nil {
				wsc.Close()
				return
			}
			wsc.KeepLive()
		case <-wsc.closeChannel:
			wsc.Close()
			return
		}
	}

}

func (wsc *WebSocketConn) Close() {
	wsc.isCloseMutex.Lock()
	defer wsc.isCloseMutex.Unlock()
	if wsc.isClosed {
		return
	}
	wsc.Conn.Close()

	wsc.closeChannel <- true
	WSCMgr.RemoveConn(wsc.UserID)
	db.DelUserOnline(wsc.UserID)

	wsc.isClosed = true
	close(wsc.sendChannel)
	close(wsc.closeChannel)
	fmt.Println("Conn Stop()... UserID = ", wsc.UserID)

}

// KeepLive 更新心跳
func (wsc *WebSocketConn) KeepLive() {
	now := time.Now()
	wsc.heartMutex.Lock()
	defer wsc.heartMutex.Unlock()
	wsc.lastHeartBeatTime = now
}

// IsAlive 是否存活
func (wsc *WebSocketConn) IsAlive() bool {
	now := time.Now()
	wsc.heartMutex.Lock()
	defer wsc.heartMutex.Unlock()
	if wsc.isClosed || now.Sub(wsc.lastHeartBeatTime) > time.Duration(config.Conf.HeartbeatTimeout)*time.Second {
		return false
	}
	return true
}

func (wsc *WebSocketConn) CompareAndIncrClientID(newMaxClientId int64) bool {
	wsc.maxClientIdMutex.Lock()
	defer wsc.maxClientIdMutex.Unlock()

	if wsc.maxClientId+1 == newMaxClientId {
		wsc.maxClientId++
		return true
	}
	fmt.Println("收到的 newMaxClientId 是：", newMaxClientId, "此时 c.maxClientId 是：", wsc.maxClientId)
	return false
}

func (wsc *WebSocketConn) HandlerMessage(data []byte) {
	// 拒收空消息
	if len(data) == 0 {
		fmt.Println("空消息")
		return
	}

	// 消息
	input := new(pb.Input)
	err := proto.Unmarshal(data, input)
	if err != nil {
		fmt.Println("HandlerMessage.proto.Unmarshal error:", err)
		return
	}
	req := &Req{
		conn: wsc,
		data: input.Data,
		f:    nil,
	}
	switch input.Type {
	case pb.CmdType_CT_LOGIN:
		fmt.Println("用户登录")
	case pb.CmdType_CT_MESSAGE: // 消息投递
		req.f = req.MessageHandler
	case pb.CmdType_CT_SYNC: // 拉取离线消息
		req.f = req.Sync
	case pb.CmdType_CT_HEARTBEAT:
		req.f = req.HeartBeat
	default:
		fmt.Println("未知消息类型")
	}
	if req.f == nil {
		return
	}
	// 更新心跳时间
	wsc.KeepLive()

	// 送入worker队列等待调度执行
	wsc.WscMgr.SendMsgToTaskQueue(req)
}
