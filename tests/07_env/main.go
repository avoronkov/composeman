package main

import (
	"fmt"
	"os"
)

func main() {
	const varName = "MY_TEST_ENV_VAR"
	fmt.Printf("%v: %v\n", varName, os.Getenv(varName))
}
