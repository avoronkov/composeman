package proc

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"github.com/avoronkov/composeman/lib/dc"
	"github.com/avoronkov/composeman/lib/utils"
)

type Proc struct {
	compose *dc.DockerCompose
	pod     string
	stdout  io.Writer
	stderr  io.Writer
}

func New(compose *dc.DockerCompose, pod string, stdout, stderr io.Writer) *Proc {
	if pod == "" {
		panic("pod name should not be empty")
	}
	return &Proc{
		compose: compose,
		pod:     pod,
		stdout:  stdout,
		stderr:  stderr,
	}
}

func (p *Proc) CreatePod(ports []string) error {
	args := []string{"pod", "create", "--name", p.pod}
	for _, p := range ports {
		args = append(args, "-p", p)
	}
	return p.runPodmanCommand(args...)
}

// Run specified services in the pod.
// Run all services if empty list is specified.
func (p *Proc) RunServicesInPod(services []string, detach bool) (err error) {
	var interupt chan os.Signal
	if !detach {
		interupt = make(chan os.Signal, 1)
		signal.Notify(interupt, os.Interrupt, os.Kill, syscall.SIGTERM)
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
	if err := p.CreatePod(ports); err != nil {
		return err
	}

	for _, service := range services {
		srv, ok := p.compose.Services[service]
		if !ok {
			return fmt.Errorf("Unknown service: %v", service)
		}
		img, err := p.prepareServiceImage(service, &srv)
		if err != nil {
			return err
		}
		env, err := srv.Env()
		if err != nil {
			return err
		}
		cmd, err := utils.ShellCmdFromString(srv.Command)
		if err != nil {
			return err
		}
		err = p.runServiceInPod(srv.Volumes, srv.EnvFile, env, img, cmd, detach, services, false)
		if err != nil {
			return err
		}
	}
	if !detach {
		sig := <-interupt
		log.Printf("Signal caught (%v). Interupting...", sig)
		// maybe withVolumes should be false?
		return p.RemovePod(true)
	}
	return nil
}

// Implementing "run" command
func (p *Proc) RunService(service string, cmd []string, cliEnv []string) (err error) {
	services, err := p.findDependingServices([]string{service})
	if err != nil {
		return err
	}
	// find all ports mappings
	// TODO deduplicate
	ports := []string{}
	for _, s := range services {
		srv, ok := p.compose.Services[s]
		if !ok {
			return fmt.Errorf("Unknown service: %v", s)
		}
		ports = append(ports, srv.Ports...)
	}

	// start pod
	if err := p.CreatePod(ports); err != nil {
		return err
	}

	for _, s := range services {
		if s == service {
			// "main" service will be started later
			continue
		}
		srv, ok := p.compose.Services[s]
		if !ok {
			return fmt.Errorf("Unknown service: %v", s)
		}
		img, err := p.prepareServiceImage(s, &srv)
		if err != nil {
			return err
		}
		env, err := srv.Env()
		if err != nil {
			return err
		}
		cmd, err := utils.ShellCmdFromString(srv.Command)
		if err != nil {
			return err
		}
		err = p.runServiceInPod(srv.Volumes, srv.EnvFile, env, img, cmd, true, services, false)
		if err != nil {
			return err
		}
	}
	srv, ok := p.compose.Services[service]
	if !ok {
		return fmt.Errorf("Unknown service: %v", service)
	}
	img, err := p.prepareServiceImage(service, &srv)
	if err != nil {
		return err
	}
	env, err := srv.Env()
	if err != nil {
		return err
	}
	env = mergeEnvs(env, cliEnv)
	var command *utils.ShellCmd
	if len(cmd) > 0 {
		command = utils.ShellCmdFromArgs(cmd)
	} else {
		var err error
		command, err = utils.ShellCmdFromString(srv.Command)
		if err != nil {
			return err
		}
	}
	err = p.runServiceInPod(srv.Volumes, srv.EnvFile, env, img, command, false, services, true)
	var errRm error
	if len(services) == 1 {
		errRm = p.RemovePod(true)
	}
	if err != nil {
		return err
	}
	return errRm
}

func (p *Proc) prepareServiceImage(service string, srv *dc.Service) (imageName string, err error) {
	if image := srv.Image; image != "" {
		return image, nil
	}
	if srv.Build == nil {
		return "", fmt.Errorf("'image' or 'build' should be specified for service %v", service)
	}
	return p.BuildImage(service, srv.Build.Context, srv.Build.Target, srv.Build.Args)
}

func (p *Proc) runServiceInPod(volumes []string, envFile string, env []string, image string, cmd *utils.ShellCmd, detach bool, hosts []string, sync bool) error {
	args := []string{"run", "-t", "--pod", p.pod}
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
	if !cmd.Empty() {
		args = append(args, cmd.Split()...)
	}
	if !detach && !sync {
		go p.runPodmanCommand(args...)
		return nil
	}
	return p.runPodmanCommand(args...)
}

func (p *Proc) RemovePod(withVolumes bool) (err error) {
	var podVolumes []string
	if withVolumes {
		podVolumes, err = p.getPodVolumes()
		if err != nil {
			return err
		}
	}
	if err = p.runPodmanCommand("pod", "rm", "-f", p.pod); err != nil {
		return err
	}
	if withVolumes {
		err = p.removeVolumes(podVolumes)
	}
	return err
}

func (p *Proc) BuildImage(serviceName, context, target string, buildArgs map[string]string) (imageName string, err error) {
	tag := fmt.Sprintf("img-%v-%v", p.pod, serviceName)
	if context == "" {
		context = "."
	}
	args := []string{"build", context, "--tag", tag}
	if target != "" {
		args = append(args, "--target", target)
	}
	for k, v := range buildArgs {
		args = append(args, "--build-arg", fmt.Sprintf("%v=%v", k, v))
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

func (p *Proc) getPodVolumes() (volumes []string, err error) {
	podInfo, err := p.podInspect()
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

func (p *Proc) podInspect() (*PodInspect, error) {
	args := []string{"pod", "inspect", p.pod}
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
	return image
}

func (p *Proc) runPodmanCommand(args ...string) error {
	log.Printf("Running: podman %v", strings.Join(args, " "))
	cmd := exec.Command("podman", args...)
	cmd.Stdout = p.stdout
	cmd.Stderr = p.stderr
	// TODO handle stdin
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

func mergeEnvs(env1, env2 []string) []string {
	// TODO implement precise merring of env variables
	return append(env1, env2...)
}
