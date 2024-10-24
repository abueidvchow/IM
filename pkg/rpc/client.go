package rpc

import (
	"IM/pkg/protocol/pb"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const maxMsgSize = 4 * 1024 * 1024 // 4MB

var (
	ConnServerClient pb.ConnectClient
)

// GetServerClient 获取 grpc 连接
func GetServerClient(addr string) pb.ConnectClient {
	client, err := grpc.Dial(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxMsgSize)),
		grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(maxMsgSize)))
	if err != nil {
		fmt.Println("grpc client Dial err, err:", err)
		panic(err)
	}
	ConnServerClient = pb.NewConnectClient(client)
	return ConnServerClient
}
