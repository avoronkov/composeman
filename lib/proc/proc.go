package proc

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"

	"github.com/avoronkov/composeman/lib/dc"
	shellquote "github.com/kballard/go-shellquote"
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
	return p.runPodmanCommand(args...)
}

// Run specified services in the pod.
// Run all services if empty list is specified.
func (p *Proc) RunServicesInPod(pod string, services []string, detach bool) (err error) {
	var interupt chan os.Signal
	if !detach {
		interupt = make(chan os.Signal, 1)
		signal.Notify(interupt, os.Interrupt, os.Kill)
	}
	if pod == "" {
		pod, err = p.DetectPodName()
		if err != nil {
			return err
		}
	}
	if len(services) == 0 {
		for name := range p.compose.Services {
			services = append(services, name)
		}
	} else {
		services, err = p.findDependingServices(services)
		if err != nil {
			return err
		}
	}
	// find all ports mappings
	ports := []string{}
	for _, service := range services {
		srv, ok := p.compose.Services[service]
		if !ok {
			return fmt.Errorf("Unknown service: %v", service)
		}
		ports = append(ports, srv.Ports...)
	}

	// start pod
	if err := p.CreatePod(pod, ports); err != nil {
		return err
	}

	for _, service := range services {
		srv, ok := p.compose.Services[service]
		if !ok {
			return fmt.Errorf("Unknown service: %v", service)
		}
		img, err := p.prepareServiceImage(pod, service, &srv)
		if err != nil {
			return err
		}
		env, err := srv.Env()
		if err != nil {
			return err
		}
		err = p.runServiceInPod(pod, srv.Volumes, srv.EnvFile, env, img, srv.Command, detach, services)
		if err != nil {
			return err
		}
	}
	if !detach {
		sig := <-interupt
		log.Printf("Signal caught (%v). Interupting...", sig)
		// maybe withVolumes should be false?
		return p.RemovePod(pod, true)
	}
	return nil
}

func (p *Proc) prepareServiceImage(pod, service string, srv *dc.Service) (imageName string, err error) {
	if image := srv.Image; image != "" {
		return image, nil
	}
	if srv.Build == nil {
		return "", fmt.Errorf("'image' or 'build' should be specified for service %v", service)
	}
	return p.BuildImage(pod, service, srv.Build.Context, srv.Build.Target, srv.Build.Args)
}

func (p *Proc) runServiceInPod(pod string, volumes []string, envFile string, env []string, image, cmd string, detach bool, hosts []string) error {
	args := []string{"run", "-t", "--pod", pod}
	if detach {
		args = append(args, "-d")
	}
	if len(volumes) > 0 {
		args = append(args, "--security-opt", "label=disable")
		for _, volume := range volumes {
			args = append(args, "-v", volume)
		}
	}
	if envFile != "" {
		args = append(args, "--env-file", envFile)
	}
	for _, e := range env {
		args = append(args, "-e", e)
	}
	for _, h := range hosts {
		args = append(args, "--add-host", fmt.Sprintf("%v:127.0.0.1", h))
	}
	args = append(args, p.canonicalImageName(image))
	if cmd != "" {
		words, err := shellquote.Split(cmd)
		if err != nil {
			return err
		}
		args = append(args, words...)
	}
	if !detach {
		go p.runPodmanCommand(args...)
		return nil
	}
	return p.runPodmanCommand(args...)
}

func (p *Proc) RemovePod(pod string, withVolumes bool) (err error) {
	var podVolumes []string
	if withVolumes {
		podVolumes, err = p.getPodVolumes(pod)
		if err != nil {
			return err
		}
	}
	if err = p.runPodmanCommand("pod", "rm", "-f", pod); err != nil {
		return err
	}
	if withVolumes {
		err = p.removeVolumes(podVolumes)
	}
	return err
}

func (p *Proc) BuildImage(pod, serviceName, context, target string, buildArgs map[string]string) (imageName string, err error) {
	tag := fmt.Sprintf("img-%v-%v", pod, serviceName)
	if context == "" {
		context = "."
	}
	args := []string{"build", context, "--tag", tag}
	if target != "" {
		args = append(args, "--target", target)
	}
	for k, v := range buildArgs {
		args = append(args, fmt.Sprintf("%v=%v", k, v))
	}
	err = p.runPodmanCommand(args...)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("localhost/%v:latest", tag), nil
}

func (p *Proc) removeVolumes(volumes []string) error {
	if len(volumes) == 0 {
		return nil
	}
	args := []string{"volume", "rm", "--force"}
	args = append(args, volumes...)
	return p.runPodmanCommand(args...)
}

func (p *Proc) getPodVolumes(pod string) (volumes []string, err error) {
	podInfo, err := p.podInspect(pod)
	if err != nil {
		return nil, err
	}
	for _, cnt := range podInfo.Containers {
		if cnt.Id == podInfo.InfraContainerID {
			continue
		}
		vols, err := p.getContainerVolumes(cnt.Id)
		if err != nil {
			return nil, err
		}
		volumes = append(volumes, vols...)
	}
	return volumes, nil
}

func (p *Proc) getContainerVolumes(cntId string) (volumes []string, err error) {
	info, err := p.inspect(cntId)
	if err != nil {
		return nil, err
	}
	for _, mount := range info.Mounts {
		if mount.Type != "volume" {
			continue
		}
		volumes = append(volumes, mount.Name)
	}
	return volumes, nil
}

func (p *Proc) podInspect(pod string) (*PodInspect, error) {
	args := []string{"pod", "inspect", pod}
	output := &strings.Builder{}
	cmd := exec.Command("podman", args...)
	cmd.Stdout = output
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	pi := &PodInspect{}
	if err := json.Unmarshal([]byte(output.String()), pi); err != nil {
		return nil, err
	}

	return pi, nil
}

func (p *Proc) inspect(cntId string) (*Inspect, error) {
	args := []string{"inspect", cntId}
	output := &strings.Builder{}
	cmd := exec.Command("podman", args...)
	cmd.Stdout = output
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	ci := []Inspect{}
	if err := json.Unmarshal([]byte(output.String()), &ci); err != nil {
		return nil, err
	}

	return &ci[0], nil
}

func (p *Proc) canonicalImageName(image string) string {
	/*
		slashes := strings.Count(image, "/")
		if slashes == 0 {
			return "docker.io/library/" + image
		}
		if slashes == 1 {
			return "docker.io/" + image
		}
	*/
	return image
}

func (p *Proc) runPodmanCommand(args ...string) error {
	log.Printf("Running: podman %v", strings.Join(args, " "))
	cmd := exec.Command("podman", args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func (p *Proc) findDependingServices(services []string) ([]string, error) {
	result := make(map[string]bool)
	for _, service := range services {
		if err := p.addDependindServices(service, &result); err != nil {
			return nil, err
		}
	}
	list := make([]string, 0, len(result))
	for s := range result {
		list = append(list, s)
	}
	return list, nil
}

func (p *Proc) addDependindServices(service string, result *map[string]bool) error {
	if (*result)[service] == true {
		// already added
		return nil
	}
	srv, ok := p.compose.Services[service]
	if !ok {
		return fmt.Errorf("Unknown service: %v", service)
	}
	(*result)[service] = true
	for _, dep := range srv.DependsOn {
		if err := p.addDependindServices(dep, result); err != nil {
			return err
		}
	}
	return nil
}
