package tests

import (
	"os"
	"strings"
	"testing"

	"github.com/avoronkov/composeman/lib/cli"
)

const podName03 = "03_run_command"

func Test03RunCommand(t *testing.T) {
	pwd := chdir(podName03)
	defer chdir(pwd)
	defer removePod(podName03)

	compile("-o", "prog.exe")

	var stdout strings.Builder
	c := cli.New(&stdout, os.Stderr)
	rc := c.Run([]string{"run", "prog"})
	if rc != 0 {
		t.Fatalf("Command 'run prog' finished with non-zero exit code: %v", rc)
	}
	out := stdout.String()
	t.Logf("Stdout: %v\n", out)

	needle := `CLI args: ["foo" "--bar" "baz"]`
	if !strings.Contains(out, needle) {
		t.Errorf("String %v not found in output.", needle)
	}
}

func Test03RunCommandCli(t *testing.T) {
	pwd := chdir(podName03)
	defer chdir(pwd)
	defer removePod(podName03)

	compile("-o", "prog.exe")

	var stdout strings.Builder
	c := cli.New(&stdout, os.Stderr)
	rc := c.Run([]string{"run", "prog", "one", "-two", "--three"})
	if rc != 0 {
		t.Fatalf("Command 'run prog' finished with non-zero exit code: %v", rc)
	}
	out := stdout.String()
	t.Logf("Stdout: %v\n", out)

	needle := `CLI args: ["one" "-two" "--three"]`
	if !strings.Contains(out, needle) {
		t.Errorf("String %v not found in output.", needle)
	}
}
