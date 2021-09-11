package tests

import (
	"log"
	"strings"
	"testing"

	"github.com/avoronkov/composeman/lib/cli"
)

func Test02RunHelloworldFromDockerfile(t *testing.T) {
	pwd := chdir("./02_run_helloworld_from_dockerfile")
	log.Printf("pwd: %v", pwd)
	defer chdir(pwd)

	var stdout, stderr strings.Builder
	c := cli.New(&stdout, &stderr)
	rc := c.Run([]string{"run", "hello-world"})
	if rc != 0 {
		t.Fatalf("Command 'run hello-world' finished with non-zero exit code: %v\nStderr: %v", rc, stderr.String())
	}
	out := stdout.String()
	t.Logf("Stdout: %v\n", out)
	t.Logf("Stderr: %v\n", stderr.String())

	greeting := "Hello from Go's HelloWorld!"
	if !strings.Contains(out, greeting) {
		t.Errorf("String %q not found in output.", greeting)
	}
}
