package tests

import (
	"os"
	"strings"
	"testing"

	"github.com/avoronkov/composeman/lib/cli"
)

const podName01 = "01_run_helloworld"

func Test01RunHelloworld(t *testing.T) {
	pwd := chdir(podName01)
	defer chdir(pwd)
	defer removePod(podName01)

	var stdout strings.Builder
	c := cli.New(&stdout, os.Stderr)
	rc := c.Run([]string{"run", "hello-world"})
	if rc != 0 {
		t.Fatalf("Command 'run hello-world' finished with non-zero exit code: %v", rc)
	}
	out := stdout.String()
	t.Logf("Stdout: %v\n", out)

	greeting := "Hello from Docker!"
	if !strings.Contains(out, greeting) {
		t.Errorf("String %q not found in output.", greeting)
	}
}
