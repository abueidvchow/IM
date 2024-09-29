package service

import (
	"IM/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
)

func SendMessageService(c *gin.Context, uid int64) error {
	wsc, err := NewWebSocketConn(c.Writer, c.Request, uid)
	if err != nil {
		return err
	}
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println("关闭ws协议失败")
		}
	}(wsc.Conn)

	go StartReader(wsc)
	//go startWriter(wsc)
	select {}
}

func ChatListService(c *gin.Context, uid int64) {
	room_id := c.Query("room_id")
	if room_id == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "房间号不能为空",
		})
		return
	}
	//判断用户是否属于该房间
	if !model.CheckUserFromUserRoom(uid, room_id) {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "非法访问",
		})
		return
	}
	page, _ := strconv.Atoi(c.Query("page"))
	size, _ := strconv.Atoi(c.Query("size"))

	list, err := model.GetChatList(room_id, page, size)
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
