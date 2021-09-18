package podman

// "podman run ..."
func (p *Podman) Run(service string, opts ...RunOpt) error {
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
	for _, f := range ro.EnvFile {
		args = append(args, "--env-file", f)
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
