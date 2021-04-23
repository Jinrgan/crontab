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

	ctx := context.Background()
	// 删除 KV
	resp, err := kv.Delete(ctx, "/cron/jobs/job2", clientv3.WithPrevKV())
	if err != nil {
		panic(err)
	}

	// 被删除之前的 value 是什么
	for _, p := range resp.PrevKvs {
		fmt.Println("删除了：", string(p.Key), string(p.Value))
	}
}
