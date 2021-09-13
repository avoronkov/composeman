package tests

import (
	"os"
	"strings"
	"testing"

	"github.com/avoronkov/composeman/lib/cli"
)

const podName04 = "04_dependent_services"

func Test04RunDependentServices(t *testing.T) {
	pwd := chdir(podName04)
	defer chdir(pwd)
	defer removePod(podName04)

	var stdout strings.Builder
	c := cli.New(&stdout, os.Stderr)
	rc := c.Run([]string{"run", "app-client"})
	if rc != 0 {
		t.Fatalf("Command 'run app-client' finished with non-zero exit code: %v", rc)
	}
	out := stdout.String()
	t.Logf("Stdout: %v\n", out)

	needle := "Got response from demo server"
	if !strings.Contains(out, needle) {
		t.Errorf("String %q not found in output.", needle)
	}
}
