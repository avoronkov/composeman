package podman

import (
	"io"
	"log"
	"os/exec"
	"strings"
)

type Executor interface {
	SetStdout(out io.Writer)
	SetStderr(err io.Writer)
	SetStdin(in io.Reader)
	Exec(cmd string, args ...string) error
}

type RealExecutor struct {
	stdout io.Writer
	stderr io.Writer
	stdin  io.Reader
}

var _ Executor = (*RealExecutor)(nil)

func NewRealExecutor() *RealExecutor {
	return &RealExecutor{}
}

func (e *RealExecutor) Exec(cmd string, args ...string) error {
	log.Printf("Running: %v %v", cmd, strings.Join(args, " "))
	com := exec.Command(cmd, args...)
	com.Stdout = e.stdout
	com.Stderr = e.stderr
	com.Stdin = e.stdin
	return com.Run()
}

func (e *RealExecutor) SetStdout(out io.Writer) {
	e.stdout = out
}

func (e *RealExecutor) SetStderr(out io.Writer) {
	e.stderr = out
}
func (e *RealExecutor) SetStdin(input io.Reader) {
	e.stdin = input
}
