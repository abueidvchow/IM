package ws

import (
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
	msg := new(pb.Message)
	err := proto.Unmarshal(r.data, msg)
	if err != nil {
		fmt.Println("MessageHandler.proto.Unmarshal error:", err)
		return
	}

	if msg.SenderId != r.conn.UserID {
		fmt.Println("发送者有误,", msg.SenderId, "    ", r.conn.UserID)
		return
	}

	// 单聊不能发给自己
	if msg.SessionType == pb.SessionType_ST_SINGLE && msg.ReceiverId == r.conn.UserID {
		fmt.Println("接收者有误")
		return
	}

	// 给自己发一份，消息落库但是不推送
	//seq, err := SendToUser(msg.Msg, msg.Msg.SenderId)
	//if err != nil {
	//	fmt.Println("[消息处理] send to 自己出错, err:", err)
	//	return
	//}

	// 单聊、群聊
	switch msg.SessionType {
	case pb.SessionType_ST_SINGLE:
		err = SendToUser(msg.ReceiverId, msg)
	case pb.SessionType_ST_GROUP:
		err = SendToGroup(msg.ReceiverId, msg)
	default:
		fmt.Println("会话类型错误")
		return
	}
	if err != nil {
		fmt.Println("[消息处理] 系统错误")
		return
	}
	//
	//// 回复发送上行消息的客户端 ACK
	//ackBytes, err := GetOutputMsg(pb.CmdType_CT_ACK, common.OK, &pb.ACKMsg{
	//	Type:     pb.ACKType_AT_Up,
	//	ClientId: msg.ClientId, // 回复客户端，当前已 ACK 的消息
	//	Seq:      seq,          // 回复客户端当前其 seq
	//})
	//if err != nil {
	//	fmt.Println("[消息处理] proto.Marshal err:", err)
	//	return
	//}
	// 回复发送 Message 请求的客户端 A
	//r.conn.SendMsg(msg.Msg.SenderId, ackBytes)
}

func (r *Req) Sync() {

}
