package main

import (
	"context"
	"fmt"
	"time"

	"github.com/coreos/etcd/clientv3"
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

	// KV
	kv := clientv3.NewKV(client)

	// 创建 OP
	putOP := clientv3.OpPut("/cron/jobs/job8", "")

	// 执行 OP
	ctx := context.Background()
	opResp, err := kv.Do(ctx, putOP)
	if err != nil {
		panic(err)
	}

	fmt.Println("写入 Revision：", opResp.Put().Header.Revision)

	// 创建 OP
	getOP := clientv3.OpGet("/cron/jobs/job8")

	// 执行 OP
	opResp, err = kv.Do(ctx, getOP)
	if err != nil {
		panic(err)
	}

	// 打印
	fmt.Println("数据 Revision：", opResp.Get().Kvs[0].ModRevision)
	fmt.Println("数据 value：", string(opResp.Get().Kvs[0].Value))
}
