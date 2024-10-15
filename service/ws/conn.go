package ws

import (
	"IM/config"
	"IM/model"
	"IM/pkg/mq"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"net/http"
	"sync"
	"time"
)

type WebSocketConn struct {
	Conn              *websocket.Conn
	UserID            int64
	inChannel         chan []byte //读消息队列
	outChannel        chan []byte //写消息队列
	closeChannel      chan bool   //监听channel是否关闭
	isClosed          bool        // 标识
	isCloseMutex      sync.Mutex  // 保护 isClose 的锁
	lastHeartBeatTime time.Time   // 最后活跃时间
	heartMutex        sync.Mutex  // 保护最后活跃时间的锁
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
		Conn:         conn,                    // websocket连接
		UserID:       uid,                     // 该连接的当前用户
		inChannel:    make(chan []byte, 1024), // 接收消息管道
		outChannel:   make(chan []byte, 1024), // 发送消息管道
		closeChannel: make(chan bool, 1),      // 关闭管道
		isClosed:     false,
	}
	// 开启读，读取客户端发送给服务端，再由服务端转发过来的消息
	go wsc.StartReader()
	// 开启写，将消息发送给服务端，再由服务端去转发
	go wsc.StartWriter()
	return
}

// 从管道里接收消息
//func (wsc *WebSocketConn) ReadMessage() ([]byte, error) {
//	for {
//		select {
//		case data := <-wsc.inChannel:
//			return data, nil
//		case <-wsc.closeChannel:
//			wsc.Close()
//			return nil, errors.New("conn is closed")
//		}
//	}
//}

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
		case data := <-wsc.outChannel:
			err := wsc.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				wsc.Close()
				return
			}
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

	WSCMgr.RemoveConn(wsc.UserID)
	wsc.isClosed = true
	wsc.closeChannel <- true
	// 停止读取和写入
	close(wsc.inChannel)
	close(wsc.outChannel)
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

func (wsc *WebSocketConn) HandlerMessage(data []byte) {
	if len(data) == 0 {
		return
	}
	msg := &model.Message{}
	err := json.Unmarshal(data, msg)
	if err != nil {
		fmt.Println("HandlerMessage.json.Unmarshal 错误：", err)
		return
	}
	if msg.SessionType == 1 && msg.ReceiverId == wsc.UserID { //不给自己发
		fmt.Println("不给自己发送")
		return
	}
	msg.SenderID = wsc.UserID
	msg.UserID = msg.ReceiverId
	switch msg.SessionType {
	case 1: //私聊
		SendToUser(msg.ReceiverId, msg)
	case 2: //群聊
		SendToGroup(msg.ReceiverId, msg)
	}
}

func SendToUser(rid int64, msg *model.Message) {
	/* POSTMAN 私聊消息发送格式
	URL:127.0.0.1:8080/api/ws?receiver_id=2433797081530368
	{
	    "message_type":1,
	    "session_type":1,
	    "receiver_id":2433797081530368,
	    "sender_id":2433910990438400,
	    "Content":"你好"
	}
	*/
	//先假设都在线
	wsc := WSCMgr.GetConn(rid)
	msg.SendTime = time.Now()
	msg.CreateTime = time.Now()
	msg.UpdateTime = time.Now()
	messages := make([]model.Message, 1)
	messages[0] = *msg
	mData, _ := json.Marshal(messages) // 序列化数据

	// 将数据发送给服务器的消息队列
	err := mq.RabbitMQ.Producer.PublishWithContext(context.Background(), mq.RabbitMQ.ExchangeName, mq.RabbitMQ.RouteKey, false, false,
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        mData,
		})
	if err != nil {
		fmt.Println("发送消息给消息队列失败:", err)
		return
	}
	wsc.outChannel <- mData
}

func SendToGroup(groupId int64, msg *model.Message) {
	// 假设都在线
	// 得到群成员列表
	/* POSTMAN 发送群聊格式
	{
	    "message_type":1,
	    "session_type":2,
	    "receiver_id":1,
	    "Content":"你好"
	}
	*/

	groupUsers, err := model.GetGroupUserList(groupId)
	if err != nil {
		fmt.Println("SendToGroup.service.GroupUserListService Error: ", err)
		return

	}
	// 当前群的群成员
	m := make(map[int64]struct{}, len(groupUsers))
	for _, user := range groupUsers {
		m[user.UserID] = struct{}{}
	}
	// 检验当前用户是否属于当前群
	if _, ok := m[msg.SenderID]; !ok {
		fmt.Println("用户不属于该群组")
		return
	}

	// 自己不再进行推送
	delete(m, msg.SenderID)
	i := 0
	messages := make([]model.Message, len(m))
	for k, _ := range m {
		messages[i] = model.Message{
			UserID:      k,
			SenderID:    msg.SenderID,
			SessionType: 2,
			ReceiverId:  groupId,
			MessageType: 1,
			Content:     msg.Content,
			//Seq:         0,
			SendTime:   time.Now(),
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
		}
		i++
	}
	mData, err := json.Marshal(messages)
	if err != nil {
		fmt.Println("SendToGroup.json.Marshal Error: ", err)
		return
	}

	// 写入消息队列
	err = mq.RabbitMQ.Producer.PublishWithContext(context.Background(), mq.RabbitMQ.ExchangeName, mq.RabbitMQ.RouteKey, false, false,
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        mData,
		})
	if err != nil {
		fmt.Println("发送消息给消息队列失败:", err)
		return
	}

	for UserID, _ := range m {
		wsc := WSCMgr.GetConn(UserID)
		mData, _ := json.Marshal(msg.Content)
		wsc.outChannel <- mData
	}
	fmt.Println("群消息发送完成 ")
}
