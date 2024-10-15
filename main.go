package main

import (
	"IM/config"
	"IM/pkg/db"
	"IM/pkg/logger"
	"IM/pkg/mq"
	sf "IM/pkg/snowflake"
	"IM/router"
	"IM/service/ws"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// 启动服务（优雅关机）
func run(r *gin.Engine) {
	// 开启心跳超时检测
	checker := ws.NewHeartBeatChecker(time.Second*time.Duration(config.Conf.HeartbeatInterval), ws.WSCMgr)
	go checker.Start()

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%v", config.Conf.Port),
		Handler: r,
	}

	go func() {
		// 开启一个goroutine启动服务
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("listen: ", zap.Error(err))
		}
	}()

	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit                                               // 阻塞在此，当接收到上述两种信号时才会往下执行

	checker.Stop()
	zap.L().Info("Shutdown Server ...")
	// 创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown: ", zap.Error(err))
	}

	zap.L().Info("Server exiting")
}

func init() {
	if len(os.Args) < 2 {
		fmt.Println("需要配置文件路径")
		os.Exit(0)
	}
	// 读取配置文件
	err := config.Init(os.Args[1])
	if err != nil {
		fmt.Println("读取配置文件错误:", err)
		return
	}

	// 数据库初始化
	err = db.InitMySQL(config.Conf.MySQLConfig)
	if err != nil {
		fmt.Println("数据库初始化失败：", err)
		return
	}
	// 日志初始化
	err = logger.Init(config.Conf.LogConfig)
	if err != nil {
		fmt.Println("日志初始化失败：", err)
		return
	}

	// 初始化雪花算法
	if err := sf.Init(config.Conf.StartTime, config.Conf.MachineID); err != nil {
		fmt.Println("雪花算法初始化失败：", err)
		return
	}

	// 初始化消息队列
	if err := mq.InitRabbitMQ(config.Conf.RabbitMQConfig); err != nil {
		fmt.Println("消息队列初始化失败：", err)
		return
	}
}

func main() {

	// 路由注册
	r := router.SetUpRouter()
	run(r)

}
