package ws

import (
	"IM/common"
	"IM/model"
	"IM/pkg/protocol/pb"
	"fmt"
	"google.golang.org/protobuf/proto"
)

// Handler 路由函数
type Handler func()

// Req 请求
type Req struct {
	conn *WebSocketConn // 连接
	data []byte         // 客户端发送的请求数据
	f    Handler        // 该请求需要执行的路由函数
}

func (r *Req) MessageHandler() {
	// 消息解析 proto string -> struct
	msg := new(pb.UpMsg)
	err := proto.Unmarshal(r.data, msg)
	if err != nil {
		fmt.Println("MessageHandler.proto.Unmarshal error:", err)
		return
	}

	// 实现消息可靠性
	if !r.conn.CompareAndIncrClientID(msg.ClientId) {
		fmt.Println("不是想要收到的 clientID，不进行处理, msg:", msg)
		return
	}

	if msg.Msg.SenderId != r.conn.UserID {
		fmt.Println("发送者有误,", msg.Msg.SenderId, "    ", r.conn.UserID)
		return
	}

	// 单聊不能发给自己
	if msg.Msg.SessionType == pb.SessionType_ST_SINGLE && msg.Msg.ReceiverId == r.conn.UserID {
		fmt.Println("接收者有误")
		return
	}

	// 给自己发一份，消息落库但是不推送
	seq, err := SendToUser(msg.Msg.SenderId, msg.Msg)
	if err != nil {
		fmt.Println("[消息处理] send to 自己出错, err:", err)
		return
	}

	// 单聊、群聊
	switch msg.Msg.SessionType {
	case pb.SessionType_ST_SINGLE:
		_, err = SendToUser(msg.Msg.ReceiverId, msg.Msg)
	case pb.SessionType_ST_GROUP:
		err = SendToGroup(msg.Msg.ReceiverId, msg.Msg)
	default:
		fmt.Println("会话类型错误")
		return
	}
	if err != nil {
		fmt.Println("[消息处理] 系统错误")
		return
	}

	// 回复发送上行消息的客户端 ACK
	ackBytes, err := GetOutputMsg(pb.CmdType_CT_ACK, common.OK, &pb.ACKMsg{
		Type:     pb.ACKType_AT_UP,
		ClientId: msg.ClientId, // 回复客户端，当前已 ACK 的消息
		Seq:      seq,          // 回复客户端当前其 seq
	})
	if err != nil {
		fmt.Println("[消息处理] proto.Marshal err:", err)
		return
	}
	//回复发送 Message 请求的客户端 A
	r.conn.SendMsg(ackBytes)

}

func (r *Req) HeartBeat() {
	// TODO 更新当前用户状态，不做回复
}

func (r *Req) Sync() {
	msg := new(pb.SyncInputMsg)
	err := proto.Unmarshal(r.data, msg)
	if err != nil {
		fmt.Println("Sync.proto.Unmarshal error:", err)
		return
	}
	// 根据seq查询，得到比 seq 大的用户消息
	messages, hasMore, err := model.MessagesListByUserIdAndSeq(r.conn.UserID, msg.Seq, model.MessageLimit)
	if err != nil {
		fmt.Println("Sync.model.MessagesListByUserIdAndSeq error:", err)
		return
	}
	pbMessage := model.MessagesToPB(messages)

	ackBytes, err := GetOutputMsg(pb.CmdType_CT_SYNC, common.OK, &pb.SyncOutputMsg{
		Messages: pbMessage,
		HasMore:  hasMore,
	})
	if err != nil {
		fmt.Println("[离线消息] proto.Marshal err:", err)
		return
	}
	r.conn.SendMsg(ackBytes)

}

// 下行消息确认
func (r *Req) ACKMsg() {
	ackMsg := new(pb.ACKMsg)
	err := proto.Unmarshal(r.data, ackMsg)
	if err != nil {
		fmt.Println("AckMsg.proto.Unmarshal error:", err)
		return
	}
	r.conn.maxClientIdMutex.Lock()
	defer r.conn.maxClientIdMutex.Unlock()
	if ackMsg.ClientId != r.conn.maxClientId {
		fmt.Printf("ackMsg.ClientId: %d != r.conn.maxClientId: %d\n", ackMsg.ClientId, r.conn.maxClientId)
		return
	}
	//fmt.Println("下行消息确认")
}
