package podman

import "io"

type Podman struct {
	executor Executor
}

func NewPodman() *Podman {
	return &Podman{
		executor: NewRealExecutor(),
	}
}

func (p *Podman) SetStdout(out io.Writer) {
	p.executor.SetStdout(out)
}

func (p *Podman) SetStderr(err io.Writer) {
	p.executor.SetStderr(err)
}

func (p *Podman) SetStdin(in io.Reader) {
	p.executor.SetStdin(in)
}
