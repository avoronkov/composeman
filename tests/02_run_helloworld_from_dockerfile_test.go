package tests

import (
	"os"
	"strings"
	"testing"

	"github.com/avoronkov/composeman/lib/cli"
)

func Test02RunHelloworldFromDockerfile(t *testing.T) {
	pwd := chdir("./02_run_helloworld_from_dockerfile")
	defer chdir(pwd)

	var stdout strings.Builder
	c := cli.New(&stdout, os.Stderr)
	rc := c.Run([]string{"run", "hello-world"})
	if rc != 0 {
		t.Fatalf("Command 'run hello-world' finished with non-zero exit code: %v", rc)
	}
	out := stdout.String()
	t.Logf("Stdout: %v\n", out)

	greeting := "Hello from Go's HelloWorld!"
	if !strings.Contains(out, greeting) {
		t.Errorf("String %q not found in output.", greeting)
	}
}
