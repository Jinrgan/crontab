package main

import (
	"fmt"
	"os/exec"
)

func main() {
	// 生成 cmd
	cmd := exec.Command("bash", "-c", "sleep 5; ls -l")
	// 执行了命令，捕获了子进程的输出（pipe）
	output, err := cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}

	// 打印子进程的输出
	fmt.Println(string(output))
}
