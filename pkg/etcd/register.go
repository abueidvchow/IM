package etcd

import (
	"IM/config"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

func RegisterETCDServer(config *config.ETCDConfig) (err error) {
	_, err = clientv3.New(clientv3.Config{
		Endpoints:   config.Endpoints,
		DialTimeout: time.Duration(config.Timeout) * time.Second,
	})
	if err != nil {
		return err
	}
	return
}
