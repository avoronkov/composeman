package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/avoronkov/composeman/lib/proc"
)

type subCommand interface {
	Run(args []string) error
}

type Cli struct {
	Proc     *proc.Proc
	commands map[string]subCommand
}

func New(p *proc.Proc) *Cli {
	return &Cli{
		Proc: p,
		commands: map[string]subCommand{
			"up":   NewUp(p),
			"down": NewDown(p),
		},
	}
}

func (c *Cli) Run(args []string) (rc int) {
	if len(args) == 0 {
		c.usage(os.Stderr)
		return 2
	}

	cmd, ok := c.commands[args[0]]
	if !ok {
		fmt.Fprintf(os.Stderr, "Unknown command specified: %v", args[0])
		return 2
	}

	if err := cmd.Run(args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err.Error())
		return 1
	}

	// OK
	return 0
}

func (c *Cli) usage(out io.Writer) {
	fmt.Fprintf(out, "No command specified.\nPossible commands are:\n")
	for name := range c.commands {
		fmt.Fprintf(out, "- %v\n", name)
	}
}
