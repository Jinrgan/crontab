package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// startTime 小于某时间
// {"$lt": timestamp}
type TimeBeforeCond struct {
	Before int64 `bson:"$lt"`
}

// {"timePoint.startTime": {"$lt": timestamp}}
type DeleteCond struct {
	beforeCond TimeBeforeCond `bson:"timePoint.startTime"`
}

func main() {
	// mongodb 读取回来的是 bson，需要反序列为 LogRecord 对象
	// 1. 建立连接
	ctx := context.TODO()
	client, err := mongo.Connect(ctx, options.Client().
		ApplyURI("mongodb://localhost:27017/").
		SetConnectTimeout(5*time.Second))
	if err != nil {
		panic(err)
	}

	// 2. 选择数据库 my_db
	database := client.Database("cron")

	// 3. 选择表 my_collection
	col := database.Collection("log")

	// 4. 要删除开始时间早于当前时间的所有日志（$lt 是 less than）
	// delete({"timePoint.startTime": {"$lt": 当前时间}})
	delCond := &DeleteCond{beforeCond: TimeBeforeCond{Before: time.Now().Unix()}}

	// 执行删除
	delMany, err := col.DeleteMany(ctx, delCond)
	if err != nil {
		panic(err)
	}

	fmt.Println("删除的行数：", delMany.DeletedCount)
}
