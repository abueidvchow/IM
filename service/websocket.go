package service

import (
	"IM/model"
	"fmt"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"sync"
)

type WebSocketConn struct {
	Conn   *websocket.Conn
	UserID int64
}

var wc map[int64]*websocket.Conn = make(map[int64]*websocket.Conn, 0)
var lock sync.RWMutex

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
	lock.Lock()
	wc[uid] = conn
	lock.Unlock()

	return &WebSocketConn{
		Conn:   conn,
		UserID: uid,
	}, nil
}

func StartReader(wsc *WebSocketConn) {
	//读取和发送消息
	for {
		ms := &model.MessageStruct{}
		//读取数据
		err := wsc.Conn.ReadJSON(ms)
		if err != nil {
			fmt.Println("wsc.Conn.ReadJSON error:", err)
		}

		//处理消息
		//HandlerMessage(ms)
		//判断当前用户是否属于发送消息体的聊天室
		if !model.CheckUserFromUserRoom(wsc.UserID, strconv.FormatInt(ms.RoomID, 10)) {
			continue
		}

		//保存消息
		//msg := &model.Message{
		//	UserID:     uid,
		//	RoomID:     ms.RoomID,
		//	Content:    ms.Content,
		//	CreateTime: time.Now(),
		//	UpdateTime: time.Now(),
		//}
		//model.SaveMessage(msg)

		//返回该消息体的聊天室的在线用户
		//users, err := model.GetUserFromUserRoom(ms.RoomID)
		//if err != nil {
		//	zap.L().Error(err.Error())
		//	return
		//}
		//
		//for _, u := range users {
		//	//不给自己发消息
		//	if u.UserID == uid {
		//		continue
		//	}
		//	if cc, ok := wc[uid]; ok { //这里要加判断是否可以得到，不然当只有一个用户的时候，wc里面只有一个，但是得到的users（这里会存在没有在线的用户）
		//		err = cc.WriteMessage(websocket.TextMessage, []byte(ms.Content))
		//		if err != nil {
		//			zap.L().Error(err.Error())
		//			return
		//		}
		//	}
		//}
	}
}

func StartWriter(wsc *WebSocketConn) {

}
