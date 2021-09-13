package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	serverUrl := os.Getenv("DEMO_SERVER_URL")
	if serverUrl == "" {
		log.Fatal("DEMO_SERVER_URL is not specified")
	}

	tries := 5
	for try := 1; try <= tries; try++ {
		resp, err := http.Get(serverUrl)
		if err == nil {
			defer resp.Body.Close()
			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Got %v", string(data))
			return
		}
		log.Printf("%v", err)
		if try == tries {
			break
		}
		time.Sleep(6 * time.Second)
	}
	log.Fatal("Could not get a response from server.")
}
