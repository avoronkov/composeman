package tests

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

func removePod(pod string) {
	log.Printf("Removing pod %v", pod)
	cmd := exec.Command("podman", "pod", "rm", "-f", pod)
	var stderr strings.Builder
	cmd.Stderr = &stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		if strings.Contains(stderr.String(), "no such pod") {
			return
		}
		log.Printf("Removing pod failed: %v", stderr.String())
		panic(err)
	}
}
