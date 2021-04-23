package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"time"
)

func main() {
	config := clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	}

	// 建立一个客户端
	client, err := clientv3.New(config)
	if err != nil {
		panic(err)
	}

	// 用于读写 etcd 的键值对
	kv := clientv3.NewKV(client)

	// 写入另外一个 job
	ctx := context.Background()
	_, err = kv.Put(ctx, "/cron/jobs/job2", "{...}")
	if err != nil {
		panic(err)
	}

	// 读取 /cron/jobs/ 为前缀的所有 key
	resp, err := kv.Get(ctx, "/cron/jobs/", clientv3.WithPrefix())
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(resp.Kvs)
}
