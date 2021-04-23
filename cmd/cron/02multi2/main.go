package main

import (
	"fmt"
	"github.com/gorhill/cronexpr"
	"time"
)

func main() {
	// 1. 定义 2 个 cronjob
	wCh1 := createWorker(cronexpr.MustParse("*/5 * * * * * *"))
	wCh2 := createWorker(cronexpr.MustParse("*/8 * * * * * *"))

	for {
		select {
		case res := <-wCh1:
			fmt.Println(res)
		case res := <-wCh2:
			fmt.Println(res)
		}
	}
}

type worker struct {
	expr *cronexpr.Expression
	ch   chan string
}

func createWorker(expr *cronexpr.Expression) chan string {
	ch := make(chan string)
	w := &worker{
		expr: expr,
		ch:   ch,
	}
	go w.work()

	return ch
}

func (w *worker) work() func() {
	for {
		// 当前时间
		now := time.Now()
		// 下次调度的时间
		nextTime := w.expr.Next(now)
		// 等待这个定时器超时
		tCh := time.After(nextTime.Sub(now))

		<-tCh
		w.ch <- fmt.Sprintf("被调度了：%v", nextTime)
	}
}
