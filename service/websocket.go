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

	// TODO
	select {}
}
