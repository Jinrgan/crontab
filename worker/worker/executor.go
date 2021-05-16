package worker

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/gorhill/cronexpr"
)

type Result struct {
	Output    string
	Err       error
	StartTime time.Time
	EndTime   time.Time
}

type Executor struct {
	command string
	expr    *cronexpr.Expression
	ResCh   chan *Result
	timer   *time.Timer
}

func NewExecutor(cmd, line string) (*Executor, error) {
	expr, err := cronexpr.Parse(line)
	if err != nil {
		return nil, fmt.Errorf("cannot parse crontab expression: %v", err)
	}
	now := time.Now()
	timer := time.NewTimer(expr.Next(now).Sub(now))

	return &Executor{
		command: cmd,
		expr:    expr,
		ResCh:   make(chan *Result),
		timer:   timer,
	}, nil
}

func (e *Executor) Execute(ctx context.Context) {
	for {
		// 当前时间
		now := time.Now()
		// 下次调度的时间
		nextTime := e.expr.Next(now)
		// 等待这个定时器超时
		e.timer.Reset(nextTime.Sub(now))

		t := <-e.timer.C
		res := &Result{
			StartTime: t,
		}
		cmd := exec.CommandContext(ctx, "/bin/bash", "-c", e.command)
		output, err := cmd.CombinedOutput()
		res.EndTime = time.Now()
		res.Output = string(output)
		res.Err = err

		e.ResCh <- res
	}
}
