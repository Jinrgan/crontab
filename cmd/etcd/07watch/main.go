package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"time"
)

func main() {
	config := clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	}

	// 建立连接
	clt, err := clientv3.New(config)
	if err != nil {
		panic(err)
	}

	// KV
	kv := clientv3.NewKV(clt)

	// 模拟 etcd 中 KV 的变化
	ctx := context.Background()
	go func() {
		for {
			_, err := kv.Put(ctx, "/cron/jobs/job7", "i am job7")
			if err != nil {
				panic(err)
			}

			_, err = kv.Delete(ctx, "/cron/jobs/job7")
			if err != nil {
				panic(err)
			}

			time.Sleep(1 * time.Second)
		}
	}()

	// 先 GET 到当前的值，并监听后续变化
	getResp, err := kv.Get(ctx, "/cron/jobs/job7")
	if err != nil {
		panic(err)
	}

	// 现在 key 是存在的
	if len(getResp.Kvs) != 0 {
		fmt.Println("当前值：", string(getResp.Kvs[0].Value))
	}

	// 当前 etcd 集群事务 ID，单调递增的
	rev := getResp.Header.Revision + 1

	// 创建一个 watcher
	wat := clientv3.NewWatcher(clt)

	// 启动监听
	fmt.Println("从该版本向后监听", rev)

	tCTX, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	wCh := wat.Watch(tCTX, "/cron/jobs/job7", clientv3.WithRev(rev))

	// 处理 kv 变化事件
	for resp := range wCh {
		for _, event := range resp.Events {
			switch event.Type {
			case mvccpb.PUT:
				fmt.Println("修改为：", string(event.Kv.Value), "Revision: ", event.Kv.CreateRevision, event.Kv.ModRevision)
			case mvccpb.DELETE:
				fmt.Println("删除了 Revision: ", event.Kv.ModRevision)
			}
		}
	}
}
