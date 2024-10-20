package ws

import (
	"IM/config"
	"IM/model"
	"IM/pkg/db"
	"IM/pkg/mq"
	"IM/pkg/protocol/pb"
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"net/http"
	"sync"
	"time"
)

type WebSocketConn struct {
	Conn              *websocket.Conn
	WscMgr            *WebSocketConnMgr
	UserID            int64
	sendChannel       chan []byte //写消息队列
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
		WscMgr:       WSCMgr,                  // 所属管理
		UserID:       uid,                     // 该连接的当前用户
		sendChannel:  make(chan []byte, 1024), // 发送消息管道
		closeChannel: make(chan bool, 1),      // 关闭管道
		isClosed:     false,
	}
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
	case pb.CmdType_CT_MESSAGE: // 消息投递
		req.f = req.MessageHandler
	case pb.CmdType_CT_SYNC: // 拉取离线消息
		req.f = req.Sync
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

func SendToUser(rid int64, msg *pb.Message) (err error) {
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

	// 获取接受者的seqID
	seq, err := db.GetUserNextSeq(db.SeqObjectTypeUser, rid)
	if err != nil {
		fmt.Println("SendToUser.db.GetUserNextSeq error:", err)
		return
	}

	// 检验对方是否在线
	online, err := db.GetUserOnline(rid)
	if err != nil {
		fmt.Println("SendToUser.db.GetUserOnline error:", err)
		return
	}
	// 对方离线
	if online == "" {
		fmt.Printf("用户ID：%d 不在线 \n", rid)
		return
	}

	// 对方在线
	wsc := WSCMgr.GetConn(rid)
	msg.Seq = seq
	msg.SendTime = time.Now().UnixMilli()

	mData := model.MessageToProtoMarshal(&model.Message{
		UserID:      rid,
		SenderID:    msg.SenderId,
		SessionType: int8(msg.SessionType),
		ReceiverId:  msg.ReceiverId,
		MessageType: int8(msg.MessageType),
		Content:     msg.Content,
		Seq:         seq,
		SendTime:    time.UnixMilli(msg.SendTime),
	})

	// 将数据发送给服务器的消息队列
	err = mq.RabbitMQ.Producer.PublishWithContext(context.Background(), mq.RabbitMQ.ExchangeName, mq.RabbitMQ.RouteKey, false, false,
		amqp091.Publishing{
			ContentType: "application/x-protobuf",
			Body:        mData,
		})
	if err != nil {
		fmt.Println("发送消息给消息队列失败:", err)
		return
	}
	wsc.sendChannel <- msg.Content
	return nil
}

func SendToGroup(groupId int64, msg *pb.Message) (err error) {
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
	if _, ok := m[msg.SenderId]; !ok {
		fmt.Println("用户不属于该群组")
		return
	}

	// 自己不再进行推送
	delete(m, msg.SenderId)

	sendUserIds := make([]int64, 0, len(m))
	for userId := range m {
		sendUserIds = append(sendUserIds, userId)
	}

	// 获取群用户的seqID
	seqIDs, err := db.GetUserNextSeqBatch(db.SeqObjectTypeUser, sendUserIds)
	if err != nil {
		fmt.Println("SendToGroup.db.GetUserNextSeqBatch error:", err)
		return
	}

	//  k:userid v:该userId的seq
	sendUserSet := make(map[int64]int64, len(seqIDs))
	for i, userId := range sendUserIds {
		sendUserSet[userId] = seqIDs[i]
	}

	messages := make([]*model.Message, 0, len(m))
	for userID, seq := range sendUserSet {
		messages = append(messages, &model.Message{
			UserID:      userID,
			SenderID:    msg.SenderId,
			SessionType: int8(msg.SessionType),
			ReceiverId:  groupId,
			MessageType: int8(msg.MessageType),
			Content:     msg.Content,
			Seq:         seq,
			SendTime:    time.UnixMilli(msg.SendTime),
		})
	}

	// 写入消息队列
	err = mq.RabbitMQ.Producer.PublishWithContext(context.Background(), mq.RabbitMQ.ExchangeName, mq.RabbitMQ.RouteKey, false, false,
		amqp091.Publishing{
			ContentType: "application/x-protobuf",
			Body:        model.MessageToProtoMarshal(messages...),
		})
	if err != nil {
		fmt.Println("发送消息给消息队列失败:", err)
		return
	}

	for UserID, _ := range m {
		wsc := WSCMgr.GetConn(UserID)
		//mData, _ := json.Marshal(msg.Content)
		wsc.sendChannel <- msg.Content
	}
	fmt.Println("群消息发送完成 ")
	return nil
}
