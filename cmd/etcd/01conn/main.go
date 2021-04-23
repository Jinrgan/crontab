package main

import (
	"github.com/coreos/etcd/clientv3"
	"time"
)

func main() {
	// 客户端配置
	config := clientv3.Config{
		Endpoints:   []string{""},
		DialTimeout: 5 * time.Second,
	}

	// 建立连接
	_, err := clientv3.New(config)
	if err != nil {
		panic(err)
	}
}
