package main

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

type result struct {
	err    error
	output []byte
}

func main() {
	// 执行 1 个 cmd，让它在一个协程里去执行，让它执行 2 秒：sleep 2; echo hello;
	// 1 秒的时候，我们杀死 cmd

	resCh := make(chan *result)

	// context: chan byte
	// cancel: close(chan byte)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		res := &result{}
		cmd := exec.CommandContext(ctx, "bash", "-c", "sleep 2; echo hello;")
		output, err := cmd.CombinedOutput()
		if err != nil {
			res.err = err
			resCh <- res
			return
		}

		res.output = output
		resCh <- res
	}()

	// 继续往下走
	time.Sleep(time.Second)

	// 取消上下文
	cancel()

	// 在 main 协程里，等待子协程的退出，并打印任务执行结果
	res := <-resCh

	// 打印任务执行结果
	fmt.Println(res.err, string(res.output))
}
