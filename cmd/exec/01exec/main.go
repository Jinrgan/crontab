package main

import (
	"os/exec"
)

func main() {
	cmd := exec.Command("bash", "-c", "echo 1; echo 2")

	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}
