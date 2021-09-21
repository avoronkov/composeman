package tests

import (
	"bufio"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/avoronkov/composeman/lib/cli"
	"github.com/avoronkov/composeman/lib/podman"
)

const podName05 = "05_up_services"

func Test05UpDownServices(t *testing.T) {
	podName := podName05
	pwd := chdir(podName)
	defer chdir(pwd)
	defer removePod(podName)

	// compile binaries
	compile("-o", "server.exe", "./cmd/server")
	compile("-o", "client.exe", "./cmd/client")

	rp, wp := io.Pipe()
	c := cli.New(wp, os.Stderr)

	// up
	go func() {
		c.Run([]string{"up"})
	}()

	needle := "Success: got response from demo server"
	found := false
	sc := bufio.NewScanner(rp)
	for sc.Scan() {
		line := sc.Text()
		t.Log(line)
		if strings.Contains(line, needle) {
			found = true
			break
		}
	}
	if sc.Err() != nil {
		t.Errorf("Reading from pipe failed: %v", sc.Err())
	}
	if !found {
		t.Errorf("String %q not found in command 'up' output.", needle)
	}

	c2 := cli.New(os.Stdout, os.Stderr)
	if rc := c2.Run([]string{"down"}); rc != 0 {
		t.Errorf("down failed with exit-code: %v", rc)
	}

	// Check that pod is removed
	_, err := podman.InspectPod(podName05)
	if exp := podman.NotFoundError; err != exp {
		t.Errorf("Incorrect error: expected %v, found %v", exp, err)
	}
}

// up -d app-server
func Test05UpServiceDashD(t *testing.T) {
	podName := podName05
	pwd := chdir(podName)
	defer chdir(pwd)
	defer removePod(podName)

	// compile binaries
	compile("-o", "server.exe", "./cmd/server")

	c := cli.New(os.Stdout, os.Stderr)

	if rc := c.Run([]string{"up", "-d", "app-server"}); rc != 0 {
		t.Fatalf("Command 'up -d app-server' finished with non-zero exit code: %v", rc)
	}

	podInfo, err := podman.InspectPod(podName)
	if err != nil {
		t.Fatal(err)
	}
	if l := len(podInfo.Containers); l != 2 {
		t.Fatalf("Incorrect number of containers in pod: expected 2, found %v", l)
	}
	for _, cnt := range podInfo.Containers {
		if cnt.Id == podInfo.InfraContainerID {
			continue
		}
		if exp := "running"; cnt.State != exp {
			t.Errorf("Incorrect app containter state: expected %v, found %v", exp, cnt.State)
		}
	}
}
