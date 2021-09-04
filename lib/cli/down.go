package cli

import "github.com/avoronkov/composeman/lib/proc"

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
	pod, err := d.Proc.DetectPodName()
	if err != nil {
		return err
	}

	return d.Proc.RemovePod(pod)
}
