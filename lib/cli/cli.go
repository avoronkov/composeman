package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/avoronkov/composeman/lib/dc"
	"github.com/avoronkov/composeman/lib/proc"
)

type subCommand interface {
	Init(p *proc.Proc)
	Run(args []string) error
}

type Cli struct {
	commands map[string]subCommand
}

func New() *Cli {
	return &Cli{
		commands: map[string]subCommand{
			"up":    NewUp(),
			"down":  NewDown(),
			"build": NewBuild(),
			"run":   NewRun(),
			"rm":    NewRm(),
		},
	}
}

func (c *Cli) Run(args []string) (rc int) {
	// Parse command line arguments
	flags := flag.NewFlagSet("composeman", flag.ContinueOnError)
	composeFiles := c.defaultComposeFiles()
	flags.Var(composeFiles, "f", "Specify an alternate compose file")
	// ignored
	project := ""
	flags.StringVar(&project, "p", os.Getenv("COMPOSE_PROJECT_NAME"), "project name (ignored)")
	noAnsi := false
	flags.BoolVar(&noAnsi, "no-ansi", false, "ignored")
	if err := flags.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return 2
	}

	if project == "" {
		var err error
		project, err = c.detectProjectName()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			return 1
		}
	}

	// Init Proc and DockerCompose
	cfg, err := dc.NewDockerCompose(composeFiles.Values()...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return 1
	}
	pr := proc.New(cfg, project)

	if flags.NArg() == 0 {
		c.usage(os.Stderr)
		return 2
	}

	cmd, ok := c.commands[flags.Arg(0)]
	if !ok {
		fmt.Fprintf(os.Stderr, "Unknown command specified: %v\n", flags.Arg(0))
		return 2
	}

	cmd.Init(pr)

	if err := cmd.Run(flags.Args()[1:]); err != nil {
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

func (c *Cli) defaultComposeFiles() *Strings {
	composeFile := os.Getenv("COMPOSE_FILE")
	if composeFile == "" {
		return &Strings{}
	}
	sep := os.Getenv("COMPOSE_PATH_SEPARATOR")
	if sep == "" {
		// TODO: OS-dependent
		sep = ":"
	}
	files := strings.Split(composeFile, sep)
	return StringsDefault(files)
}

func (c *Cli) detectProjectName() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", nil
	}
	return filepath.Base(dir), nil
}
