package cli

import (
	"flag"

	"github.com/avoronkov/composeman/lib/proc"
)

type Down struct {
	Proc *proc.Proc
}

func NewDown() *Down {
	return &Down{}
}

func (d *Down) Init(p *proc.Proc) {
	d.Proc = p
}

// Arguments: [-v]
func (d *Down) Run(args []string) error {
	// Parse arguments
	flags := flag.NewFlagSet("composeman down", flag.ContinueOnError)
	removeVolumes := false
	flags.BoolVar(&removeVolumes, "v", false, "Remove anonymous volumes")
	flags.BoolVar(&removeVolumes, "volumes", false, "Remove anonymous volumes")
	removeOrphans := false
	flags.BoolVar(&removeOrphans, "remove-orphans", false, "(ignored)")
	if err := flags.Parse(args); err != nil {
		return err
	}

	// Perform actions
	return d.Proc.RemovePod(removeVolumes)
}
