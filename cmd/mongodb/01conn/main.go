package main

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func main() {
	// 1. 建立连接
	client, err := mongo.Connect(context.TODO(), options.Client().
		ApplyURI("mongodb://localhost:27017/").
		SetConnectTimeout(5*time.Second))
	if err != nil {
		panic(err)
	}

	// 2. 选择数据库 my_db
	database := client.Database("my_db")

	// 3. 选择表 my_collection
	_ = database.Collection("my_collection")
}
