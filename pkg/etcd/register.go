package etcd

import (
	"IM/config"
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

// 服务 注册
type Register struct {
	client        *clientv3.Client                        // etcd client
	leaseID       clientv3.LeaseID                        // 租约ID
	keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse // 租约KeepAlive 响应channel
	key           string                                  // key
	val           string                                  // val
}

// 新建注册服务
func RegisterETCDServer(key, value string, lease int64) (err error) {
	fmt.Println("config.Conf.ETCDConfig.Endpoints:", config.Conf.ETCDConfig.Endpoints)
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   config.Conf.ETCDConfig.Endpoints,
		DialTimeout: time.Duration(config.Conf.ETCDConfig.Timeout) * time.Second,
	})
	if err != nil {
		fmt.Println("create etcd server failed, err:", err)
		return err
	}

	ser := &Register{
		client: client,
		key:    key,
		val:    value,
	}

	// 申请租约设置时间keepalive
	if err = ser.putKeyWithLease(lease); err != nil {
		return err
	}
	// 监听续租相应的channel
	go ser.ListenLeaseRespChannel()

	return
}

// 设置key和对应的租约
func (r *Register) putKeyWithLease(timeNum int64) error {
	// 设置租约时间

	// 如果etcd服务没有启动的话，这步会直接没法执行了，但是程序也不报错会接着执行下去
	resp, err := r.client.Grant(context.TODO(), timeNum)
	if err != nil {
		return err
	}

	// 注册服务并绑定租约
	_, err = r.client.Put(context.TODO(), r.key, r.val, clientv3.WithLease(resp.ID))
	if err != nil {
		return err
	}
	// 不设置续租的话，到达租约时间就会消失
	// 设置续租，定期发送需求请求
	leaseRespChan, err := r.client.KeepAlive(context.TODO(), resp.ID)
	if err != nil {
		return err
	}
	r.leaseID = resp.ID
	r.keepAliveChan = leaseRespChan
	return nil
}

// ListenLeaseRespChan 监听 续租情况
func (r *Register) ListenLeaseRespChannel() {
	// 如果续租失败要关闭
	defer func(r *Register) {
		err := r.Close()
		if err != nil {
			fmt.Printf("续约失败，leaseID:%d, Put key:%s,val:%s\n", r.leaseID, r.key, r.val)
		}
	}(r)

	for range r.keepAliveChan {
	}

}

// Close 撤销租约
func (r *Register) Close() (err error) {
	if _, err = r.client.Revoke(context.TODO(), r.leaseID); err != nil {
		return err
	}
	fmt.Printf("撤销租约成功, leaseID:%d, Put key:%s,val:%s\n", r.leaseID, r.key, r.val)
	return r.client.Close()
}
