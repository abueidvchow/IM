package main

import (
	"IM/common"
	"IM/pkg/protocol/pb"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const (
	httpAddr       = "http://127.0.0.1:8080"
	websocketAddr  = "ws://127.0.0.1:8080"
	ResendCountMax = 3 // 超时重传最大次数
)

type WebSocketClient struct {
	conn                 *websocket.Conn
	token                string
	userId               int64
	clientId             int64
	clientId2Cancel      map[int64]context.CancelFunc // clientId 到 context 的映射
	clientId2CancelMutex sync.Mutex
	seq                  int64 // 本地消息最大同步序列号
}

// websocket 客户端
func main() {
	// http 登录，获取 token
	client := Login()

	// 连接 websocket 服务
	client.Start()
}

func (wsc *WebSocketClient) Start() {
	// 连接 websocket
	header := http.Header{}
	header.Set("Authorization", "Bearer "+wsc.token)
	conn, _, err := websocket.DefaultDialer.Dial(websocketAddr+"/api/ws", header)
	if err != nil {
		panic(err)
	}
	wsc.conn = conn

	fmt.Println("与 websocket 建立连接")
	// 向 websocket 发送登录请求
	//wsc.Login()

	// 心跳
	go wsc.Heartbeat()

	time.Sleep(time.Millisecond * 100)

	// 离线消息同步
	go wsc.Sync()

	// 收取消息
	go wsc.Receive()

	time.Sleep(time.Millisecond * 100)

	wsc.ReadLine()
}

// ReadLine 读取用户消息并发送
func (wsc *WebSocketClient) ReadLine() {
	var (
		msg         string
		receiverId  int64
		sessionType int8
	)

	readLine := func(hint string, v interface{}) {
		fmt.Println(hint)
		_, err := fmt.Scanln(v)
		if err != nil {
			panic(err)
		}
	}
	for {
		readLine("发送消息", &msg)
		readLine("接收人id(user_id/group_id)：", &receiverId)
		readLine("发送消息类型(1-单聊、2-群聊)：", &sessionType)
		message := &pb.Message{
			SessionType: pb.SessionType(sessionType),
			ReceiverId:  receiverId,
			SenderId:    wsc.userId,
			MessageType: pb.MessageType_MT_TEXT,
			Content:     []byte(msg),
			SendTime:    time.Now().UnixMilli(),
		}
		UpMsg := &pb.UpMsg{
			Msg:      message,
			ClientId: wsc.GetClientId(),
		}

		wsc.SendMsg(pb.CmdType_CT_MESSAGE, UpMsg)

		// 启动超时重传
		ctx, cancel := context.WithCancel(context.Background())

		go func(ctx context.Context) {
			maxRetry := ResendCountMax // 最大重试次数
			retryCount := 0
			retryInterval := time.Millisecond * 100 // 重试间隔
			ticker := time.NewTicker(retryInterval)
			defer ticker.Stop()
			for {
				select {
				case <-ctx.Done():
					//fmt.Println("收到 ACK，不再重试")
					return
				case <-ticker.C:
					if retryCount >= maxRetry {
						fmt.Println("达到最大超时次数，不再重试")
						// TODO 进行消息发送失败处理
						return
					}
					fmt.Println("消息超时 msg:", msg, "，第", retryCount+1, "次重试")
					wsc.SendMsg(pb.CmdType_CT_MESSAGE, UpMsg)
					retryCount++
				}
			}
		}(ctx)

		// 多协程里有访问到map都要加锁
		wsc.clientId2CancelMutex.Lock()
		wsc.clientId2Cancel[UpMsg.ClientId] = cancel
		wsc.clientId2CancelMutex.Unlock()

		time.Sleep(time.Second)
	}
}

func (wsc *WebSocketClient) Heartbeat() {
	//  2min 一次
	ticker := time.NewTicker(time.Second * 2 * 60)
	for range ticker.C {
		wsc.SendMsg(pb.CmdType_CT_HEARTBEAT, &pb.HeartbeatMsg{})
		//fmt.Println("发送心跳", time.Now().Format("2006-01-02 15:04:05"))
	}
}

func (wsc *WebSocketClient) Sync() {
	wsc.SendMsg(pb.CmdType_CT_SYNC, &pb.SyncInputMsg{Seq: wsc.seq})
}

func (wsc *WebSocketClient) Receive() {
	for {
		_, data, err := wsc.conn.ReadMessage()
		if err != nil {
			panic(err)
		}
		wsc.HandlerMessage(data)
	}
}

// HandlerMessage 客户端消息处理
func (wsc *WebSocketClient) HandlerMessage(bytes []byte) {
	outputBatchMsg := new(pb.OutputBatch)
	err := proto.Unmarshal(bytes, outputBatchMsg)
	if err != nil {
		panic(err)
	}
	for _, output := range outputBatchMsg.Outputs {
		msg := new(pb.Output)
		err := proto.Unmarshal(output, msg)
		if err != nil {
			panic(err)
		}

		fmt.Println("收到顶层 OutPut 消息：", msg)

		switch msg.Type {
		case pb.CmdType_CT_SYNC:
			syncMsg := new(pb.SyncOutputMsg)
			err = proto.Unmarshal(msg.Data, syncMsg)
			if err != nil {
				panic(err)
			}

			seq := wsc.seq
			for _, message := range syncMsg.Messages {
				fmt.Println("收到离线消息：", message)
				if seq < message.Seq {
					seq = message.Seq
				}
			}
			wsc.seq = seq
			// 如果还有，继续拉取
			if syncMsg.HasMore {
				wsc.Sync()
			}
		case pb.CmdType_CT_MESSAGE:
			pushMsg := new(pb.PushMsg)
			err = proto.Unmarshal(msg.Data, pushMsg)
			if err != nil {
				panic(err)
			}
			fmt.Printf("收到消息：%s, 发送人id：%d, 会话类型：%s, 接收时间:%s\n", pushMsg.Msg.GetContent(), pushMsg.Msg.GetSenderId(), pushMsg.Msg.SessionType, time.Now().Format("2006-01-02 15:04:05"))
			// 更新 seq
			seq := pushMsg.Msg.Seq
			if wsc.seq < seq {
				wsc.seq = seq
				fmt.Println("更新 seq:", wsc.seq)
			}

			// 需要向服务端回复 ACKType：AT_PUSH的消息
			ackMsg := new(pb.ACKMsg)
			ackMsg.Type = pb.ACKType_AT_PUSH
			ackMsg.ClientId = wsc.clientId
			ackMsg.Seq = seq

			wsc.SendMsg(pb.CmdType_CT_ACK, ackMsg)

		case pb.CmdType_CT_ACK: // 收到 ACK
			ackMsg := new(pb.ACKMsg)
			err = proto.Unmarshal(msg.Data, ackMsg)
			if err != nil {
				panic(err)
			}

			switch ackMsg.Type {
			case pb.ACKType_AT_LOGIN:
				fmt.Println("登录成功")
			case pb.ACKType_AT_UP: // 收到上行消息的 ACK
				// 取消超时重传
				clientId := ackMsg.ClientId // 更新clientId
				wsc.clientId2CancelMutex.Lock()
				if cancel, ok := wsc.clientId2Cancel[clientId]; ok {
					// 取消超时重传
					cancel()
					delete(wsc.clientId2Cancel, clientId)
					//fmt.Println("取消超时重传，clientId:", clientId)
				}
				wsc.clientId2CancelMutex.Unlock()
				// 更新客户端本地维护的 seq
				seq := ackMsg.Seq
				if wsc.seq < seq {
					wsc.seq = seq
				}
			}
		default:
			fmt.Println("未知消息类型")
		}
	}
}

// Login websocket 登录
func (wsc *WebSocketClient) Login() {
	fmt.Println("websocket login...")
	// 组装底层数据
	loginMsg := &pb.LoginMsg{
		Token: []byte(wsc.token),
	}
	wsc.SendMsg(pb.CmdType_CT_LOGIN, loginMsg)
}

// SendMsg 客户端向服务端发送上行消息
func (wsc *WebSocketClient) SendMsg(cmdType pb.CmdType, msg proto.Message) {
	// 组装顶层数据
	cmdMsg := &pb.Input{
		Type: cmdType,
		Data: nil,
	}
	if msg != nil {
		data, err := proto.Marshal(msg)
		if err != nil {
			panic(err)
		}
		cmdMsg.Data = data
	}

	bytes, err := proto.Marshal(cmdMsg)
	if err != nil {
		panic(err)
	}

	// 发送
	err = wsc.conn.WriteMessage(websocket.BinaryMessage, bytes)
	if err != nil {
		panic(err)
	}
}

func (wsc *WebSocketClient) GetClientId() int64 {
	wsc.clientId++
	return wsc.clientId
}

// Login 用户http登录获取 token
func Login() *WebSocketClient {
	// 读取 phone_number 和 password 参数
	var username, password string
	fmt.Print("Enter username: ")
	fmt.Scanln(&username)
	fmt.Print("Enter password: ")
	fmt.Scanln(&password)

	// 设置请求头
	headers := map[string]string{
		"Content-Type": "application/json",
	}

	// 准备请求体的原始数据
	body := fmt.Sprintf(`{"username":"%s", "password":"%s"}`, username, password)
	data := bytes.NewBufferString(body)

	// 创建HTTP的 POST 请求
	req, err := http.NewRequest("POST", httpAddr+"/login", data)
	if err != nil {
		fmt.Println("Error sending HTTP request:", err)
		panic(err)
	}
	// 添加请求头
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	// 创建 HTTP 客户端
	httpClient := &http.Client{}

	// 发送请求
	resp, err := httpClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// 读取响应体
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	//fmt.Println("responseBody:", string(responseBody))

	// 解析服务端返回的登录数据
	var respData struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    []struct {
			Token  string `json:"token"`
			UserId string `json:"user_id"`
		} `json:"data"`
	}
	err = json.Unmarshal(responseBody, &respData)
	if err != nil {
		panic(err)
	}

	// 读取返回数据
	//bodyData, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	panic(err)
	//}
	// 检查响应状态码
	if respData.Code != common.SUCCESS_LOGIN {
		panic(fmt.Sprintf("登录失败, %s\n", respData))
	}
	// 获取客户端 seq 序列号
	var seq int64
	fmt.Print("Enter seq: ")
	fmt.Scanln(&seq)

	client := &WebSocketClient{
		clientId2Cancel: make(map[int64]context.CancelFunc),
		seq:             seq,
	}

	client.token = respData.Data[0].Token
	clientStr := respData.Data[0].UserId
	client.userId, err = strconv.ParseInt(clientStr, 10, 64)
	if err != nil {
		panic(err)
	}

	fmt.Println("client code:", respData.Code, " ", respData.Message)
	fmt.Println("token:", client.token, "userId", client.userId)
	return client
}
