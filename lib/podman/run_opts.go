package podman

import (
	"fmt"

	"github.com/avoronkov/composeman/lib/utils"
)

type RunOpt interface {
	SetRunOpt(*RunOpts)
}

type RunOpts struct {
	Rm      bool
	Pod     string
	Detach  bool
	Volumes []string
	EnvFile []string
	Env     []string
	Hosts   []string
	Cmd     []string

	// variable to store errors
	Err error
}

// --rm
func OptRm(rm bool) RunOpt {
	return &optRm{rm}
}

type optRm struct{ rm bool }

func (o *optRm) SetRunOpt(opts *RunOpts) {
	opts.Rm = o.rm
}

// --pod
func OptPod(pod string) RunOpt {
	return &optPod{pod}
}

type optPod struct {
	pod string
}

func (o *optPod) SetRunOpt(opts *RunOpts) {
	opts.Pod = o.pod
}

// --detach (-d)
func OptDetach(d bool) RunOpt { return &optDetach{d} }

type optDetach struct{ detach bool }

func (o *optDetach) SetRunOpt(opts *RunOpts) {
	opts.Detach = o.detach
}

// --volume (-v)
func OptVolume(volumes ...string) RunOpt { return &optVolume{volumes} }

type optVolume struct{ volumes []string }

func (o *optVolume) SetRunOpt(opts *RunOpts) {
	opts.Volumes = append(opts.Volumes, o.volumes...)
}

// --env-file
func OptEnvFile(envFile ...string) RunOpt { return &optEnvFile{envFile} }

type optEnvFile struct{ envFile []string }

func (o *optEnvFile) SetRunOpt(opts *RunOpts) {
	opts.EnvFile = o.envFile
}

// --env (-e)
func OptEnv(envs ...string) RunOpt { return &optEnv{envs} }

type optEnv struct{ envs []string }

func (o *optEnv) SetRunOpt(opts *RunOpts) {
	opts.Env = o.envs
}

// --add-host %v:127.0.0.1
func OptLocalHost(hosts ...string) RunOpt { return &optLocalHost{hosts} }

type optLocalHost struct{ hosts []string }

func (o *optLocalHost) SetRunOpt(opts *RunOpts) {
	for _, h := range o.hosts {
		opts.Hosts = append(opts.Hosts, fmt.Sprintf("%v:127.0.0.1", h))
	}
}

// cmd (string)
func OptCmdString(cmd string) RunOpt { return &optCmdString{cmd} }

type optCmdString struct{ cmd string }

func (o *optCmdString) SetRunOpt(opts *RunOpts) {
	shell, err := utils.ShellCmdFromString(o.cmd)
	if err != nil && opts.Err == nil {
		opts.Err = err
		return
	}
	opts.Cmd = shell.Split()
}

// cmd (list)
func OptCmdList(cmd ...string) RunOpt { return &optCmdList{cmd} }

type optCmdList struct{ cmd []string }

func (o *optCmdList) SetRunOpt(opts *RunOpts) {
	opts.Cmd = o.cmd
}
