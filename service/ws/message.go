package ws

import (
	"IM/common"
	"IM/model"
	"IM/pkg/db"
	"IM/pkg/mq"
	"IM/pkg/protocol/pb"
	"context"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/proto"
	"time"
)

// GetOutputMsg 组装出下行消息
func GetOutputMsg(cmdType pb.CmdType, code int32, messages ...proto.Message) ([]byte, error) {
	outputBatch := &pb.OutputBatch{}
	if messages != nil {
		for _, msg := range messages {
			// 组装output消息
			output := &pb.Output{
				Type:    cmdType,
				Code:    code,
				CodeMsg: common.HttpStatusCode(code).GetMsg(),
				Data:    nil,
			}

			msgBytes, err := proto.Marshal(msg)
			if err != nil {
				fmt.Println("GetOutputMsg.proto.Marshal(msg) error:", err)
				return nil, err
			}
			output.Data = msgBytes

			outputBytes, err := proto.Marshal(output)
			if err != nil {
				fmt.Println("GetOutputMsg.proto.Marshal(output) error:", err)
				return nil, err
			}
			// 加入outputBatch
			outputBatch.Outputs = append(outputBatch.Outputs, outputBytes)
		}
	}

	// 最终要发送的消息批次
	bytes, err := proto.Marshal(outputBatch)
	if err != nil {
		fmt.Println("[GetOutputMsg] output marshal err:", err)
		return nil, err
	}
	return bytes, nil
}

func SendToUser(rid int64, msg *pb.Message) (int64, error) {
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
		return 0, err
	}

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
			//ContentType: "application/x-protobuf", // POSTMAN
			ContentType: "application/json",

			Body: mData,
		})
	if err != nil {
		fmt.Println("发送消息给消息队列失败:", err)
		return 0, err
	}

	// 如果发给自己的，只落库不进行发送
	if rid == msg.SenderId {
		return seq, nil
	}

	// 检验对方是否在线
	online, err := db.GetUserOnline(rid)
	if err != nil {
		fmt.Println("SendToUser.db.GetUserOnline error:", err)
		return 0, err
	}
	// 对方离线
	if online == "" {
		fmt.Printf("用户ID：%d 不在线 \n", rid)
		return 0, err
	}

	// 对方在线
	wsc := WSCMgr.GetConn(rid)

	// 组装下行消息
	bytes, err := GetOutputMsg(pb.CmdType_CT_MESSAGE, common.OK, &pb.PushMsg{Msg: msg})
	if err != nil {
		fmt.Println("[消息处理] GetOutputMsg Marshal error,err:", err)
		return 0, err
	}

	//wsc.sendChannel <- msg.Content // POSTMAN接收
	wsc.sendChannel <- bytes
	return 0, nil
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

	// 获取当前群的群成员
	groupUsers, err := model.GetGroupUserList(groupId)
	if err != nil {
		fmt.Println("SendToGroup.service.GroupUserListService Error: ", err)
		return

	}
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

	// receiverIds 接收消息的群成员
	receiverIds := make([]int64, 0, len(m))
	for userId := range m {
		receiverIds = append(receiverIds, userId)
	}

	// 获取群用户的seqID
	seqIDs, err := db.GetUserNextSeqBatch(db.SeqObjectTypeUser, receiverIds)
	if err != nil {
		fmt.Println("SendToGroup.db.GetUserNextSeqBatch error:", err)
		return
	}

	//  k:userid v:该userId的seq
	receiverId2Seq := make(map[int64]int64, len(seqIDs))
	for i, userId := range receiverIds {
		receiverId2Seq[userId] = seqIDs[i]
	}

	messages := make([]*model.Message, 0, len(m))
	for userID, seq := range receiverId2Seq {
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
			//ContentType: "application/x-protobuf", // POSTMAN
			ContentType: "application/json",
			Body:        model.MessageToProtoMarshal(messages...),
		})
	if err != nil {
		fmt.Println("发送消息给消息队列失败:", err)
		return
	}

	// 组装消息推送
	for UserID, seq := range receiverId2Seq {
		wsc := WSCMgr.GetConn(UserID)
		msg.Seq = seq
		// 组装下行消息
		bytes, err := GetOutputMsg(pb.CmdType_CT_MESSAGE, common.OK, &pb.PushMsg{Msg: msg})
		if err != nil {
			fmt.Println("[消息处理] GetOutputMsg Marshal error,err:", err)
			return err
		}

		wsc.sendChannel <- bytes
		//wsc.sendChannel <- msg.Content // POSTMAN发送
	}
	fmt.Println("群消息发送完成 ")

	return nil
}
