package tests

import (
	"fmt"
	"log"
	"os"
)

func chdir(to string) string {
	log.Printf("chdir: %v", to)
	from, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	if err := os.Chdir(to); err != nil {
		panic(fmt.Errorf("Cannot change directory from %v: %v", from, err))
	}
	return from
}
