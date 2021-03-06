package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "response from demo server\n")
	})
	log.Printf("Demo server is running on :8080")
	http.ListenAndServe(":8080", nil)
}
