package cli

import (
	"flag"
	"fmt"

	"github.com/avoronkov/composeman/lib/proc"
)

type Up struct {
	Proc *proc.Proc
}

func NewUp(p *proc.Proc) *Up {
	return &Up{
		Proc: p,
	}
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

	// TODO handle -d
	if len(flags.Args()) < 1 {
		return fmt.Errorf("Incorrect usage: up [-d] <service>")
	}

	service := flags.Arg(0)

	srv, ok := u.Proc.FindService(service)
	if !ok {
		return fmt.Errorf("Unknown service: %v", service)
	}

	// start pod
	pod, err := u.Proc.DetectPodName()
	if err != nil {
		return err
	}

	err = u.Proc.CreatePod(pod, srv.Ports)
	if err != nil {
		return err
	}

	image := srv.Image
	if image == "" {
		if srv.Build == nil {
			return fmt.Errorf("'image' or 'build' should be specified for service %v", service)
		}
		builtImage, err := u.Proc.BuildImage(pod, service, srv.Build.Context, srv.Build.Target, srv.Build.Args)
		if err != nil {
			return err
		}
		image = builtImage
	}

	// Run service
	err = u.Proc.RunServiceInPod(pod, srv.Volumes, srv.Environment, image, srv.Command, detach)
	if err != nil {
		return err
	}

	return nil
}
