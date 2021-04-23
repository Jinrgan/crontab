package main

import (
	"fmt"
	"github.com/gorhill/cronexpr"
	"time"
)

// 代表一个任务
type CronJob struct {
	expr     *cronexpr.Expression
	nextTime time.Time
}

func main() {
	// 需要有 1 个调度协程，它定时检查所有的 Cron 任务，谁过期了就执行谁

	// key: 任务的名字
	scheduleTable := make(map[string]*CronJob)

	now := time.Now()

	// 1. 定义 2 个 cronjob
	expr := cronexpr.MustParse("*/5 * * * * * *")
	job := &CronJob{
		expr:     expr,
		nextTime: expr.Next(now),
	}
	// 任务注册调度表
	scheduleTable["job1"] = job

	expr = cronexpr.MustParse("*/8 * * * * * *")
	job = &CronJob{
		expr:     expr,
		nextTime: expr.Next(now),
	}
	// 任务注册调度表
	scheduleTable["job2"] = job

	// 启动一个调度协程
	// 定时检查一下任务调度
	for {
		now = time.Now()

		for jobName, cronJob := range scheduleTable {
			// 判断是否过期
			if cronJob.nextTime.Before(now) || cronJob.nextTime.Equal(now) {
				// 启动一个协程，执行这个任务
				go func(jobName string) {
					fmt.Println("执行：", jobName)
				}(jobName)

				// 计算下一次执行时间
				cronJob.nextTime = cronJob.expr.Next(now)
				fmt.Println("下次执行时间：", cronJob.nextTime)
			}
		}

		// 睡眠 100 毫秒
		select {
		case <-time.NewTimer(100 * time.Millisecond).C: // 将在 100 毫秒可读，返回

		}
	}
}
