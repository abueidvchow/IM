package etcd

import "IM/config"

func InitETCD(config *config.ETCDConfig) (err error) {

	// 注册服务并设置 k-v 租约
	err = RegisterETCDServer(config)
	if err != nil {
		return err
	}

	return nil
}
