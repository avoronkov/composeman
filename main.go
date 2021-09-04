package main

import (
	"log"
	"os"

	"github.com/avoronkov/composeman/lib/cli"
	"github.com/avoronkov/composeman/lib/dc"
	"github.com/avoronkov/composeman/lib/proc"
)

func main() {
	cfg, err := dc.NewDockerCompose("docker-compose.yml")
	if err != nil {
		log.Fatal(err)
	}

	pr := proc.New(cfg)

	cl := cli.New(pr)
	os.Exit(cl.Run(os.Args[1:]))
}
