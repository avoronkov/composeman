package podman

type RunOpt interface {
	SetRunOpt(*RunOpts)
}

type RunOpts struct {
	Rm bool
}

// --rm
func OptRm() RunOpt {
	return &optRm{}
}

type optRm struct{}

func (o *optRm) SetRunOpt(opts *RunOpts) {
	opts.Rm = true
}
