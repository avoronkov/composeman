package tests

import (
	"fmt"
	"os"
	"os/exec"
)

func compile(args ...string) {
	cmd := exec.Command("go", append([]string{"build"}, args...)...)
	cmd.Stderr = os.Stderr
	cmd.Env = []string{
		"CGO_ENABLED=0",
		fmt.Sprintf("GOPATH=%v", os.Getenv("GOPATH")),
		fmt.Sprintf("HOME=%v", os.Getenv("HOME")),
	}
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}
