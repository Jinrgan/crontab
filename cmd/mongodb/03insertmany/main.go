package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
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

func main() {
	// 1. 建立连接
	client, err := mongo.Connect(context.TODO(), options.Client().
		ApplyURI("mongodb://localhost:27017/").
		SetConnectTimeout(5*time.Second))
	if err != nil {
		panic(err)
	}

	// 2. 选择数据库 my_db
	database := client.Database("cron")

	// 3. 选择表 my_collection
	col := database.Collection("log")

	// 4. 插入记录（bson）
	record := &LogRecord{
		JobName: "job10",
		Command: "echo hello",
		Err:     "",
		Content: "hello",
		TimePoint: &TimePoint{
			StartTime: time.Now().Unix(),
			EndTime:   time.Now().Unix() + 10,
		},
	}

	// 5. 批量插入多条 document
	manyResult, err := col.InsertMany(context.Background(), []interface{}{
		record,
		record,
		record})
	if err != nil {
		log.Fatal(err)
	}

	// 推特很早的时候开源的，tweet 的 ID
	// snowflake：毫秒/微秒的当前时间 + 机器的 ID + 当前毫秒/微秒内的自增 ID（每当毫秒变化了，会重置成 0，继续自增）
	fmt.Println(manyResult)
}
