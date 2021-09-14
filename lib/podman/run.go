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
	if ro.Err != nil {
		return ro.Err
	}
	args := []string{"run"}
	if ro.Rm {
		args = append(args, "--rm")
	}
	if ro.Pod != "" {
		args = append(args, "--pod", ro.Pod)
	}
	if ro.Detach {
		args = append(args, "-d")
	}
	if len(ro.Volumes) > 0 {
		args = append(args, "--security-opt", "label=disable")
		for _, v := range ro.Volumes {
			args = append(args, "-v", v)
		}
	}
	if ro.EnvFile != "" {
		args = append(args, "--env-file", ro.EnvFile)
	}
	for _, e := range ro.Env {
		args = append(args, "-e", e)
	}
	for _, h := range ro.Hosts {
		args = append(args, "--add-host", h)
	}
	args = append(args, service)
	args = append(args, ro.Cmd...)
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
