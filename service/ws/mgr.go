package ws

import (
	"IM/config"
	"IM/pkg/db"
	"fmt"
	"sync"
)

var OnlineUser []int64

type WebSocketConnMgr struct {
	Wscs      map[int64]*WebSocketConn
	Lock      sync.RWMutex
	TaskQueue []chan *Req // 工作池
}

var WSCMgr *WebSocketConnMgr = &WebSocketConnMgr{Wscs: make(map[int64]*WebSocketConn), TaskQueue: make([]chan *Req, config.Conf.WorkerPoolSize)}

func (this *WebSocketConnMgr) AddConn(uid int64, wsc *WebSocketConn) {
	this.Lock.Lock()
	defer this.Lock.Unlock()
	this.Wscs[uid] = wsc
}

func (this *WebSocketConnMgr) RemoveConn(uid int64) {
	this.Lock.Lock()
	defer this.Lock.Unlock()
	delete(this.Wscs, uid)
	return
}

func (this *WebSocketConnMgr) GetConn(uid int64) *WebSocketConn {
	this.Lock.RLock()
	defer this.Lock.RUnlock()
	return this.Wscs[uid]
}

func (this *WebSocketConnMgr) GetAllConn() map[int64]*WebSocketConn {
	this.Lock.Lock()
	defer this.Lock.Unlock()
	return this.Wscs
}

// SendMessageAll 进行本地推送
func (this *WebSocketConnMgr) SendMessageAll(userId2Msg map[int64][]byte) {
	var wg sync.WaitGroup
	ch := make(chan struct{}, 5) // 限制并发数
	for userId, data := range userId2Msg {
		ch <- struct{}{}
		wg.Add(1)
		go func(userId int64, data []byte) {
			defer func() {
				<-ch
				wg.Done()
			}()
			conn := this.GetConn(userId)
			if conn != nil {
				conn.SendMsg(data)
			}
		}(userId, data)
	}
	close(ch)
	wg.Wait()
}

// StartWorkerPool 启动 worker 工作池
func (this *WebSocketConnMgr) StartWorkerPool() {
	// 初始化并启动 worker 工作池
	for i := 0; i < len(this.TaskQueue); i++ {
		// 初始化
		this.TaskQueue[i] = make(chan *Req, config.Conf.MaxWorkerTask) // 初始化worker队列中，每个worker的队列长度
		// 启动worker
		go this.StartOneWorker(i, this.TaskQueue[i])
	}
}

// StartOneWorker 启动 worker 的工作流程
func (this *WebSocketConnMgr) StartOneWorker(workerID int, taskQueue chan *Req) {
	fmt.Println("Worker ID = ", workerID, " is started.")
	for {
		select {
		case req := <-taskQueue:
			req.f()
		}
	}
}

// SendMsgToTaskQueue 将消息交给 taskQueue，由 worker 调度处理
func (this *WebSocketConnMgr) SendMsgToTaskQueue(req *Req) {
	if len(this.TaskQueue) > 0 {
		// 根据ConnID来分配当前的连接应该由哪个worker负责处理，保证同一个连接的消息处理串行
		// 轮询的平均分配法则

		//得到需要处理此条连接的workerID
		workerID := req.conn.UserID % int64(len(this.TaskQueue))

		// 将消息发给对应的 taskQueue
		this.TaskQueue[workerID] <- req
	} else {
		// 可能导致消息乱序
		go req.f()
	}
}

func (this *WebSocketConnMgr) Stop() {
	fmt.Println("server stop ...")
	ch := make(chan struct{}, 1000) // 控制并发数，防止过多协程竞争资源
	var wg sync.WaitGroup           // 确保所有的协程运行完毕后，再继续执行主线程的逻辑
	connAll := this.GetAllConn()
	for _, conn := range connAll {
		ch <- struct{}{}
		wg.Add(1)
		c := conn
		go func() {
			//c.Close()
			//wg.Done()
			//<-ch
			// 如果是上面这样写，如果Close出现错误或导致 panic，wg.Done() 和 <-ch 可能不会执行，可能会导致程序卡住或等待结束不了。
			// wg.Done() 不被调用：主线程会永远卡在 wg.Wait()，因为 Goroutine 执行失败，未能减少 WaitGroup 的计数器。
			// <-ch 不被调用：通道 ch 的空间无法释放，新的 Goroutine 也可能无法启动。
			defer func() { // 如果Close出错了，可以确保wg和ch的执行；
				wg.Done()
				<-ch
			}()
			c.Close()
		}()
	}
	close(ch)

	for _, userID := range OnlineUser {
		err := db.DelUserOnline(userID)
		if err != nil {
			fmt.Println("Stop.db.DelUserOnline error:", err)
			return
		}
	}

	wg.Wait()
}
