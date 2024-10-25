package etcd

import (
	"IM/common"
	"IM/config"
	"fmt"
	"time"
)

var (
	DiscoverySer *Discovery
)

func InitETCD(config *config.AppConfig) {
	fmt.Println("ETCD InitETCD")
	host := fmt.Sprintf("%v:%v", config.IP, config.RPCPort)
	// 注册服务并设置 k-v 租约
	err := RegisterETCDServer(common.EtcdServerList+host, host, 5)
	if err != nil {
		fmt.Println("InitETCD.RegisterETCDServer Error:", err)
		panic(err)
	}
	time.Sleep(100 * time.Millisecond)

	// 服务发现
	DiscoverySer, err = NewDiscovery()
	if err != nil {
		fmt.Println("InitETCD.NewDiscovery Error:", err)
		panic(err)
	}
	fmt.Println("服务发现完成")
	// 阻塞监听
	err = DiscoverySer.WatchService(common.EtcdServerList)
	if err != nil {
		fmt.Println("InitETCD.WatchService Error:", err)
		panic(err)
	}
}
