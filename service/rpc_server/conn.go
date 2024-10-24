package rpc_server

import (
	"IM/pkg/protocol/pb"
	"IM/service/ws"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"net"
)

type ConnectServer struct {
	pb.UnimplementedConnectServer
}

// 从另一个服务器发送过来给本地服务器的用户
func (R ConnectServer) DeliverMessage(ctx context.Context, req *pb.DeliverMessageReq) (*emptypb.Empty, error) {
	resp := &emptypb.Empty{}

	// 获取本地连接
	conn := ws.WSCMgr.GetConn(req.ReceiverId)
	if conn == nil || conn.UserID != req.ReceiverId {
		fmt.Println("连接不存在，userID:", req.ReceiverId)
		return resp, nil
	}
	// 消息发送
	conn.SendMsg(req.Data)

	return resp, nil
}

func (R ConnectServer) DeliverGroupMessage(ctx context.Context, req *pb.DeliverMessageGroupReq) (*emptypb.Empty, error) {
	resp := &emptypb.Empty{}

	//// 进行本地推送
	//ws.WSCMgr
	//ws.GetServer().SendMessageAll(req.GetReceiverId_2Data())

	return resp, nil
}

func InitRPCServer(port int) {
	grpcServer := grpc.NewServer()
	pb.RegisterConnectServer(grpcServer, &ConnectServer{})

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Println("listen err:", err)
		return
	}
	defer listener.Close()
	fmt.Println("RPC服务启动端口:", port)
	for {
		err := grpcServer.Serve(listener)
		if err != nil {
			fmt.Println("grpcServer.Serve err:", err)
			return
		}
	}
}
