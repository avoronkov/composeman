package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	name := os.Getenv("TEST_FILE")
	if name == "" {
		log.Fatal("TEST_FILE is not set")
	}
	action := os.Getenv("PROG_ACTION")
	switch action {
	case "read":
		data, err := ioutil.ReadFile(name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Test file content: %s\n", data)
	case "write":
		content := []byte("Generated content.")
		err := ioutil.WriteFile(name, content, 0644)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Data written to test file: %s\n", content)
	default:
		log.Fatalf("Unknown PROG_ACTION specified: %v", action)
	}
}
