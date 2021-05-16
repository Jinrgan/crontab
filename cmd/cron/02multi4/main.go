package main

import (
	"fmt"
	"time"

	"github.com/gorhill/cronexpr"
)

func main() {
	// 1. 定义 2 个 cronjob
	s := NewScheduler()
	s.createWorker("*/5 * * * * * *")
	s.createWorker("*/8 * * * * * *")
	s.createWorker("*/15 * * * * * *")

	s.Run()

	time.Sleep(30 * time.Second)
}

type scheduler struct {
	logCh chan string
}

func NewScheduler() *scheduler {
	return &scheduler{logCh: make(chan string)}
}

func (s *scheduler) createWorker(line string) {
	wCh := createWorker(line)

	go func() {
		for res := range wCh {
			s.logCh <- res
		}
	}()
}

func (s *scheduler) Run() {
	for log := range s.logCh {
		fmt.Println(log)
	}
}

type worker struct {
	expr  *cronexpr.Expression
	timer *time.Timer
	ch    chan string
}

func createWorker(line string) chan string {
	expr := cronexpr.MustParse(line)
	// 当前时间
	now := time.Now()
	// 下次调度的时间
	nextTime := expr.Next(now)
	// 等待这个定时器超时
	tm := time.NewTimer(nextTime.Sub(now))

	ch := make(chan string)
	w := &worker{
		expr:  expr,
		timer: tm,
		ch:    ch,
	}
	go w.work()

	return ch
}

func (w *worker) work() {
	for {
		// 当前时间
		now := time.Now()
		// 下次调度的时间
		nextTime := w.expr.Next(now)
		// 等待这个定时器超时
		w.timer.Reset(nextTime.Sub(now))

		t := <-w.timer.C
		w.ch <- fmt.Sprintf("被调度了：%v", t)
	}
}
