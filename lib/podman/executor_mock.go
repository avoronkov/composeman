package podman

import "io"

type executorMock struct {
	lastCommand []string
}

var _ Executor = (*executorMock)(nil)

func (e *executorMock) SetStdout(out io.Writer) {}
func (e *executorMock) SetStderr(err io.Writer) {}
func (e *executorMock) SetStdin(in io.Reader)   {}
func (e *executorMock) Exec(cmd string, args ...string) error {
	e.lastCommand = append([]string{cmd}, args...)
	return nil
}
