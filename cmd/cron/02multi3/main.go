package main

import (
	"fmt"
	"github.com/gorhill/cronexpr"
	"time"
)

func main() {
	// 1. 定义 2 个 cronjob
	s := &scheduler{}
	s.createWork(cronexpr.MustParse("*/5 * * * * * *"))
	s.createWork(cronexpr.MustParse("*/8 * * * * * *"))
	s.createWork(cronexpr.MustParse("*/15 * * * * * *"))

	time.Sleep(30 * time.Second)
}

type scheduler struct {
	wChans []chan string
}

func (s *scheduler) createWork(expr *cronexpr.Expression) {
	wCh := createWorker(expr)

	go func() {
		for {
			res := <-wCh
			fmt.Println(res)
		}
	}()
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
