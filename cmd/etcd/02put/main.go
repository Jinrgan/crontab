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

	resp, err := kv.Put(context.Background(), "/cron/jobs/job1", "bye", clientv3.WithPrevKV())
	if err != nil {
		panic(err)
	}

	fmt.Println("Revision", resp.Header.Revision)
	if resp.PrevKv != nil { // 打印 hello
		fmt.Println("PrevValue", string(resp.PrevKv.Value))
	}
}
