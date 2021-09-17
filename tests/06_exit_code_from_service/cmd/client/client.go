package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	serverUrl := os.Getenv("DEMO_SERVER_URL")
	if serverUrl == "" {
		log.Fatal("DEMO_SERVER_URL is not specified")
	}

	for {
		resp, err := http.Get(serverUrl)
		if err == nil {
			defer resp.Body.Close()
			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Success: got %v\n", string(data))
			rc, err := strconv.Atoi(string(data))
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Exit with %v\n", rc)
			os.Exit(rc)
		}
		log.Printf("%v", err)
		time.Sleep(6 * time.Second)
	}
}
