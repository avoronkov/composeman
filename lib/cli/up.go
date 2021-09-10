package cli

import (
	"flag"

	"github.com/avoronkov/composeman/lib/proc"
)

type Up struct {
	Proc *proc.Proc
}

func NewUp() *Up {
	return &Up{}
}

func (u *Up) Init(p *proc.Proc) {
	u.Proc = p
}

// Arguments: [-d] <service>
func (u *Up) Run(args []string) error {
	// Parse arguments
	flags := flag.NewFlagSet("composeman up", flag.ContinueOnError)
	detach := false
	flags.BoolVar(&detach, "d", false, "Run containers in the background")
	if err := flags.Parse(args); err != nil {
		return err
	}

	return u.Proc.RunServicesInPod(flags.Args(), detach)
}
