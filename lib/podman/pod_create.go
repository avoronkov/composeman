package podman

// "podman pod create"
func (p *Podman) PodCreate(pod string, opts ...PodCreateOpt) error {
	pco := &PodCreateOpts{}
	for _, opt := range opts {
		opt.SetPodCreateOpt(pco)
	}
	args := []string{"pod", "create", "--name", pod}
	for _, p := range pco.Ports {
		args = append(args, "-p", p)
	}
	return p.executor.Exec("podman", args...)
}
