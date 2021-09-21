package utils

import shellquote "github.com/kballard/go-shellquote"

type ShellCmd struct {
	args []string
}

func ShellCmdFromString(cmd string) (*ShellCmd, error) {
	if cmd == "" {
		return &ShellCmd{args: nil}, nil
	}
	args, err := shellquote.Split(cmd)
	if err != nil {
		return nil, err
	}
	return &ShellCmd{args: args}, nil
}

func (s *ShellCmd) Split() []string {
	return s.args
}
