package service

import (
	"IM/service/ws"
	"fmt"
	"net/http"
)

func ChatService(uid int64, w http.ResponseWriter, r *http.Request) {
	var wsc *ws.WebSocketConn

	// 给当前用户创立WebSocket连接
	wsc, err := ws.NewWebSocketConn(w, r, uid)

	// 把当前用户加入websocket conn管理
	ws.WSCMgr.AddConn(uid, wsc)
	if err != nil {
		fmt.Println("SingleChatService.ws.NewWebSocketConn err: ", err)
		return
	}

	// 开启一个协程从管道读取消息和发送消息
	go func(c *ws.WebSocketConn) {
		for {
			fmt.Println("等待读取对方消息")
			data, err := c.ReadMessage()
			if err != nil {
				fmt.Println("ChatService.c.ReadMessage err: ", err)
				return
			}
			c.HandlerMessage(data)
		}
	}(wsc)

	// TODO
	select {}
}
