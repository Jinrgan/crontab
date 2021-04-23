package main

import (
	"fmt"
	"github.com/gorhill/cronexpr"
	"time"
)

func main() {
	// linux crontab
	// 秒粒度，年配置（2018-2099）
	// 哪一分钟（0-59），哪小时（0-23），哪天（1-31），哪月（1-12），星期几（0-6）

	// 每隔 5 分钟执行一次
	expr := cronexpr.MustParse("*/2 * * * * * *")

	go work(expr)()

	time.Sleep(5 * time.Second)
}

func work(expr *cronexpr.Expression) func() {
	return func() {
		for {
			// 当前时间
			now := time.Now()
			// 下次调度的时间
			nextTime := expr.Next(now)
			// 等待这个定时器超时
			tCh := time.After(nextTime.Sub(now))

			<-tCh
			nextTime = expr.Next(now)
			fmt.Println("被调度了：", nextTime)
		}
	}
}
