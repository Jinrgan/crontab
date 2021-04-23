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

	// 申请一个 lease （租约）
	lease := clientv3.NewLease(client)

	// 申请一个 5s 的租约
	ctx := context.Background()
	gResp, err := lease.Grant(ctx, 5)
	if err != nil {
		panic(err)
	}

	// 租约的 ID
	leaseID := gResp.ID

	// 自动续租
	// 处理续约应答的协程
	go func() {
		// 续租了 10 秒，停止了续租，5 秒的生命期 = 15 秒的生命期

		// 10 秒后会取消自动续租
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		AliCh, err := lease.KeepAlive(ctx, leaseID)
		if err != nil {
			panic(err)
		}

		for AliResp := range AliCh {
			fmt.Println("收到自动续租应答：", AliResp.ID) // 每秒会续租一次，所以就会收到一次应答
		}

		fmt.Println("租约已经失效了")
	}()

	// 获得 kv API 子集
	kv := clientv3.NewKV(client)

	// Put 一个 KV，让它与租约关联起来，从而实现 5 秒后自动过期
	_, err = kv.Put(ctx, "/cron/lock/job1", "", clientv3.WithLease(leaseID))
	if err != nil {
		panic(err)
	}

	// 定期地看一下 key 过期了没
	for {
		getResp, err := kv.Get(ctx, "/cron/lock/job1")
		if err != nil {
			panic(err)
		}
		if getResp.Count == 0 {
			fmt.Println("KV 过期了")
			break
		} else {
			fmt.Println("还没过期：", getResp.Kvs)
			time.Sleep(2 * time.Second)
		}
	}
}
