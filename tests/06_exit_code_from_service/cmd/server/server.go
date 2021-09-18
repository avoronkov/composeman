package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	message := os.Getenv("EXIT_CODE_MESSAGE")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s", message)
	})
	log.Printf("Demo server is running on :8080")
	http.ListenAndServe(":8080", nil)
}
