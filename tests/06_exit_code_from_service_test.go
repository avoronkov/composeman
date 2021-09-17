package tests

import (
	"os"
	"testing"

	"github.com/avoronkov/composeman/lib/cli"
	"github.com/avoronkov/composeman/lib/podman"
)

const podName06 = "06_exit_code_from_service"

func Test06ExitCodeFromService(t *testing.T) {
	pwd := chdir(podName06)
	defer chdir(pwd)
	defer removePod(podName06)

	// compile binaries
	compile("-o", "server.exe", "./cmd/server")
	compile("-o", "client.exe", "./cmd/client")

	// up
	c := cli.New(os.Stdout, os.Stderr)

	rc := c.Run([]string{"up", "--build", "--exit-code-from", "app-client", "app-client"})
	if exp := 4; rc != exp {
		t.Errorf("Incorrect exit code returned: want %v, got %v", exp, rc)
	}

	// Check that pod contains running app-server and exited app-client
	podInfo, err := podman.InspectPod(podName06)
	if err != nil {
		t.Fatal(err)
	}
	if l := len(podInfo.Containers); l != 3 {
		t.Fatalf("Incorrect number of containers in pod: expected 3, found %v", l)
	}
	for _, cnt := range podInfo.Containers {
		if cnt.Id == podInfo.InfraContainerID {
			continue
		}
		switch cnt.Name {
		case "app-client":
			if exp := "exited"; cnt.State != exp {
				t.Errorf("Incorrect state of container %v: want %v, got %v", cnt.Name, exp, cnt.State)
			}
		case "app-server":
			if exp := "running"; cnt.State != exp {
				t.Errorf("Incorrect state of container %v: want %v, got %v", cnt.Name, exp, cnt.State)
			}
		}
	}
}
