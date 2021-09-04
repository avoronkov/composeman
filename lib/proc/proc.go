package proc

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/avoronkov/composeman/lib/dc"
)

type Proc struct {
	compose *dc.DockerCompose
}

func New(compose *dc.DockerCompose) *Proc {
	return &Proc{
		compose: compose,
	}
}

func (p *Proc) DetectPodName() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", nil
	}
	return filepath.Base(dir), nil
}

func (p *Proc) FindService(service string) (dc.Service, bool) {
	srv, ok := p.compose.Services[service]
	return srv, ok
}

func (p *Proc) CreatePod(pod string, ports []string) error {
	args := []string{"pod", "create", "--name", pod}
	for _, p := range ports {
		args = append(args, "-p", p)
	}
	return p.runPodmanCommand(args)
}

func (p *Proc) RunServiceInPod(pod string, volumes []string, env []string, image string) error {
	args := []string{"run", "-dt", "--pod", pod}
	if len(volumes) > 0 {
		args = append(args, "--security-opt", "label=disable")
		for _, volume := range volumes {
			args = append(args, "-v", volume)
		}
	}
	for _, e := range env {
		args = append(args, "-e", e)
	}

	args = append(args, p.canonicalImageName(image))
	return p.runPodmanCommand(args)
}

func (p *Proc) canonicalImageName(image string) string {
	slashes := strings.Count(image, "/")
	if slashes == 0 {
		return "docker.io/library/" + image
	}
	if slashes == 1 {
		return "docker.io/" + image
	}
	return image
}

func (p *Proc) runPodmanCommand(args []string) error {
	log.Printf("Running: podman %v", strings.Join(args, " "))
	cmd := exec.Command("podman", args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
