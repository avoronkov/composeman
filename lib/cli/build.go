package cli

import (
	"log"

	"github.com/avoronkov/composeman/lib/proc"
)

type Build struct{}

func NewBuild() *Build {
	return &Build{}
}

func (b *Build) Init(p *proc.Proc) {}

func (b *Build) Run(args []string) error {
	log.Printf("[Warning] 'build' command does nothing at the moment.")
	return nil
}
