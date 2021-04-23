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

	// 建立连接
	client, err := clientv3.New(config)
	if err != nil {
		panic(err)
	}

	// lease 实现锁自动过期
	// op 操作
	// txn 事务：if else then

	// 1. 上锁（创建租约，自动续租，拿着租约去抢占一个 key）
	// 申请一个 lease （租约）
	lease := clientv3.NewLease(client)

	// 申请一个 5s 的租约
	grtResp, err := lease.Grant(context.Background(), 5)
	if err != nil {
		panic(err)
	}
	leaseID := grtResp.ID
	defer lease.Revoke(context.Background(), leaseID)

	go func() {
		// 准备一个用于取消自动续租的 context
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// 5 秒后会取消自动续租
		alvCh, err := lease.KeepAlive(ctx, leaseID)
		if err != nil {
			panic(err)
		}

		// 处理续约应答
		for aliResp := range alvCh {
			fmt.Println("收到自动续租应答：", aliResp.ID) // 每秒会续租一次，所以就会收到一次应答
		}

		fmt.Println("租约已经失效了")
	}()

	// if 不存在 key，then 设置它，else 抢锁失败
	kv := clientv3.NewKV(client)

	// 创建事务
	txn := kv.Txn(context.Background())
	// 定义事务
	// 如果 key 不存在
	txn.If(clientv3.Compare(clientv3.CreateRevision("/cron/lock/job9"), "=", 0)).
		Then(clientv3.OpPut("/cron/lock/job9", "xxx", clientv3.WithLease(leaseID))).
		Else(clientv3.OpGet("/cron/lock/job9")) // 否则抢锁失败

	// 提交事务
	txnResp, err := txn.Commit()
	if err != nil {
		panic(err)
	}

	// 判断是否抢到了锁
	if !txnResp.Succeeded {
		fmt.Println("锁被占用：", string(txnResp.Responses[0].GetResponseRange().Kvs[0].Value))
		return
	}

	// 2. 处理事务
	fmt.Println("处理事务")
	time.Sleep(5 * time.Second)

	// 在锁内，很安全

	// 3. 释放锁（取消自动续租，释放租约）
	// defer 会把租约释放掉，关联的 KV 就被删除了
}
