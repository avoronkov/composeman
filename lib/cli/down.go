package cli

import (
	"flag"

	"github.com/avoronkov/composeman/lib/proc"
)

type Down struct {
	Proc *proc.Proc
}

func NewDown(p *proc.Proc) *Down {
	return &Down{
		Proc: p,
	}
}

// Arguments: [-v]
func (d *Down) Run(args []string) error {
	// Parse arguments
	flags := flag.NewFlagSet("composeman down", flag.ContinueOnError)
	removeVolumes := false
	flags.BoolVar(&removeVolumes, "v", false, "Remove anonymous volumes")
	if err := flags.Parse(args); err != nil {
		return err
	}

	// Perform actions
	pod, err := d.Proc.DetectPodName()
	if err != nil {
		return err
	}

	return d.Proc.RemovePod(pod, removeVolumes)
}
