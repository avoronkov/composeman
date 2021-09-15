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
	pwd := chdir(podName05)
	defer chdir(pwd)
	defer removePod(podName05)

	rp, wp := io.Pipe()
	c := cli.New(wp, os.Stderr)

	// up
	go func() {
		c.Run([]string{"up"})
	}()

	sc := bufio.NewScanner(rp)
	for sc.Scan() {
		line := sc.Text()
		t.Log(line)
		if strings.Contains(line, "Success: got response from demo server") {
			break
		}
	}
	if sc.Err() != nil {
		t.Errorf("Reading from pipe failed: %v", sc.Err())
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
