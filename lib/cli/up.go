package cli

import (
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
	// TODO handle -d
	if len(args) < 1 {
		return fmt.Errorf("Incorrect usage: up [-d] <service>")
	}

	service := args[0]

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
	err = u.Proc.RunServiceInPod(pod, srv.Volumes, srv.Environment, image, srv.Command)
	if err != nil {
		return err
	}

	return nil
}
