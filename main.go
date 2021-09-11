package main

import (
	"os"

	"github.com/avoronkov/composeman/lib/cli"
)

func main() {
	cl := cli.New(os.Stdout, os.Stderr)
	os.Exit(cl.Run(os.Args[1:]))
}
