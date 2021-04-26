package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//TimePoint 任务的执行时间点
type TimePoint struct {
	StartTime int64 `bson:"start_time"`
	EndTime   int64 `bson:"end_time"`
}

//LogRecord 一条日志
type LogRecord struct {
	JobName   string     `bson:"job_name"`   // 任务名
	Command   string     `bson:"command"`    // shell 命令
	Err       string     `bson:"err"`        // 脚本错误
	Content   string     `bson:"content"`    // 脚本输出
	TimePoint *TimePoint `bson:"time_point"` // 执行时间点
}

//FindByJobName jobName 过滤条件
type FindByJobName struct {
	JobName string `bson:"job_name"` // JobName 赋值为 job10
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

	// 4. 按照 jobName 字段过滤，想找出 jobName = job10，找出 5 条
	cond := &FindByJobName{JobName: "job10"} // {"job_name": "job10"}

	// 5. 查询（过滤 + 翻页参数）
	cursor, err := col.Find(ctx, cond, options.Find().SetSkip(0).SetLimit(2))
	if err != nil {
		panic(err)
	}
	defer cursor.Close(ctx) // 延迟释放游标

	// 6. 遍历结果集
	for cursor.Next(ctx) {
		// 定义一个日志数据结构
		var record LogRecord

		// 反序列化 bson 到数据结构
		err := cursor.Decode(&record)
		if err != nil {
			panic(err)
		}

		// 把日志行打印出来
		fmt.Println(record)
	}
}
