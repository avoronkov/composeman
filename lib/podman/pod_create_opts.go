package podman

type PodCreateOpt interface {
	SetPodCreateOpt(*PodCreateOpts)
}

type PodCreateOpts struct {
	Ports []string
}

// -p (--publish)
func OptPublishPort(ports ...string) PodCreateOpt {
	return &optPublishPort{
		ports,
	}
}

type optPublishPort struct{ ports []string }

func (o *optPublishPort) SetPodCreateOpt(opts *PodCreateOpts) {
	opts.Ports = o.ports
}
