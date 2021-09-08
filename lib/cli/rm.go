package cli

import (
	"log"

	"github.com/avoronkov/composeman/lib/proc"
)

type Rm struct{}

func NewRm() *Rm {
	return &Rm{}
}

func (r *Rm) Init(p *proc.Proc) {}

func (r *Rm) Run(args []string) error {
	log.Printf("[Warning] 'rm' command does nothing at the moment")
	return nil
}
