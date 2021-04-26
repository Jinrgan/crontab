package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
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
	collection := database.Collection("log")

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

	result, err := collection.InsertOne(context.Background(), record)
	if err != nil {
		panic(err)
	}

	// _id: 默认生成一个全局唯一 ID，ObjectID：12 字节的二进制
	docID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		panic("no such type")
	}

	fmt.Println("自增 ID：", docID)
}
