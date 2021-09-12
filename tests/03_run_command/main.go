package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Printf("CLI args: %q\n", os.Args[1:])
}
