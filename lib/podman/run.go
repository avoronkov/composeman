package podman

import "io"

type PodmanRun struct {
	executor Executor
}

func NewPodmanRun() *PodmanRun {
	return &PodmanRun{
		executor: NewRealExecutor(),
	}
}

func (p *PodmanRun) Exec(service string, opts ...RunOpt) error {
	ro := &RunOpts{}
	for _, opt := range opts {
		opt.SetRunOpt(ro)
	}
	args := []string{"run"}
	if ro.Rm {
		args = append(args, "--rm")
	}
	args = append(args, service)
	return p.executor.Exec("podman", args...)
}

func (p *PodmanRun) SetStdout(out io.Writer) {
	p.executor.SetStdout(out)
}

func (p *PodmanRun) SetStderr(err io.Writer) {
	p.executor.SetStderr(err)
}

func (p *PodmanRun) SetStdin(in io.Reader) {
	p.executor.SetStdin(in)
}
