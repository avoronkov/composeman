package cli

import (
	"flag"
	"fmt"

	"github.com/avoronkov/composeman/lib/proc"
)

type Run struct {
	Proc *proc.Proc
}

func NewRun() *Run {
	return &Run{}
}

func (r *Run) Init(p *proc.Proc) {
	r.Proc = p
}

func (r *Run) Run(args []string) error {
	// Parse arguments
	flags := flag.NewFlagSet("composeman run", flag.ContinueOnError)
	// Not used at the momemt
	rm := false
	flags.BoolVar(&rm, "rm", false, "(ignored)")

	// Not used at the momemt
	user := ""
	flags.StringVar(&user, "user", "", "(ignored)")

	env := &Strings{}
	flags.Var(env, "e", "set environment variable")

	if err := flags.Parse(args); err != nil {
		return err
	}

	if flags.NArg() == 0 {
		return fmt.Errorf("Service to run is not specified")
	}
	service := flags.Arg(0)

	// Perform actions

	return r.Proc.RunService("", service, flags.Args()[1:], env.Values())
}
