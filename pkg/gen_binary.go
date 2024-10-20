package pkg

import (
	"IM/pkg/protocol/pb"
	"fmt"
	"google.golang.org/protobuf/proto"
	"io/ioutil"
	"log"
)

// 序列化消息生成二进制文件
func GenBin() {
	// 私聊
	msg1 := pb.Message{
		SessionType: 1,
		ReceiverId:  2658626632155136,
		SenderId:    2433910990438400,
		MessageType: 1,
		Content:     []byte("你好"),
	}

	marshal, err := proto.Marshal(&msg1)
	if err != nil {
		return
	}
	input := &pb.Input{
		Type: 1,
		Data: marshal,
	}
	data, err := proto.Marshal(input)
	if err != nil {
		return
	}
	// 写入文件
	err = ioutil.WriteFile("siliao.bin", data, 0644)
	if err != nil {
		log.Fatalf("Failed to write file: %v", err)
	}

	fmt.Println("Data has been written to siliao.bin")

	// 群聊
	msg2 := pb.Message{
		SessionType: 2,
		ReceiverId:  1,
		SenderId:    2433910990438400,
		MessageType: 1,
		Content:     []byte("你们好"),
	}
	marshal, err = proto.Marshal(&msg2)
	if err != nil {
		return
	}
	input = &pb.Input{
		Type: 1,
		Data: marshal,
	}
	data, err = proto.Marshal(input)
	if err != nil {
		return
	}
	// 写入文件
	err = ioutil.WriteFile("group.bin", data, 0644)
	if err != nil {
		log.Fatalf("Failed to write file: %v", err)
	}

	fmt.Println("Data has been written to group.bin")
}
