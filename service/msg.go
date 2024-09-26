package service

import "C"
import (
	"IM/controller"
	"IM/db/mysql"
	"IM/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var wc map[string]*websocket.Conn = make(map[string]*websocket.Conn, 0)
var lock sync.RWMutex

func SendMsgService(w http.ResponseWriter, r *http.Request, uid int64) {
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
		return
	}
	defer conn.Close()

	user_id := strconv.FormatInt(uid, 10)
	lock.Lock()
	wc[user_id] = conn
	lock.Unlock()

	//读取和发送消息
	for {
		ms := &model.MessageStruct{}
		//读取数据
		err = conn.ReadJSON(ms)
		if err != nil {
			fmt.Println("err:", err)
			return
		}
		//判断当前用户是否属于发送消息体的聊天室
		if !mysql.CheckUserFromUserRoom(user_id, strconv.FormatInt(ms.RoomID, 10)) {
			continue
		}

		//保存消息
		msg := &model.Message{
			UserID:     uid,
			RoomID:     ms.RoomID,
			Content:    ms.Content,
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
		}
		mysql.SaveMessage(msg)

		//返回该消息体的聊天室的在线用户
		users, err := mysql.GetUserFromUserRoom(ms.RoomID)
		if err != nil {
			zap.L().Error(err.Error())
			return
		}

		for _, u := range users {
			//不给自己发消息
			if u.UserID == uid {
				continue
			}
			uidStr := strconv.FormatInt(u.UserID, 10)
			if cc, ok := wc[uidStr]; ok { //这里要加判断是否可以得到，不然当只有一个用户的时候，wc里面只有一个，但是得到的users（这里会存在没有在线的用户）
				err = cc.WriteMessage(websocket.TextMessage, []byte(ms.Content))
				if err != nil {
					zap.L().Error(err.Error())
					return
				}
			}
		}
	}

}

func ChatListService(c *gin.Context) {
	room_id := c.Query("room_id")
	if room_id == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "房间号不能为空",
		})
		return
	}
	//判断用户是否属于该房间
	uid, _ := c.Get(controller.CtxUserIdKey)
	if !mysql.CheckUserFromUserRoom(uid.(string), room_id) {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "非法访问",
		})
		return
	}
	page, _ := strconv.Atoi(c.Query("page"))
	size, _ := strconv.Atoi(c.Query("size"))

	list, err := mysql.GetChatList(room_id, page, size)
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
		"data": list,
	})
	return
}
