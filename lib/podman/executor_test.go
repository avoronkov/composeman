package podman

import (
	"os/exec"
	"strings"
	"testing"
)

func TestExecutorTrue(t *testing.T) {
	var e Executor
	e = NewRealExecutor()

	err := e.Exec("/bin/true")
	if err != nil {
		t.Errorf("Real executor failed to run /bin/true: %v", err)
	}
}

func TestExecutorFalse(t *testing.T) {
	e := NewRealExecutor()
	err := e.Exec("/bin/false")
	errExit, ok := err.(*exec.ExitError)
	if !ok {
		t.Fatalf("Incorrect error returned: want *exec.ExitError, got %v", err)
	}
	if act, exp := errExit.ProcessState.ExitCode(), 1; act != exp {
		t.Errorf("Incorrect exit code returned: want %v, got %v", exp, act)
	}
}

func TestExecutorStdout(t *testing.T) {
	e := NewRealExecutor()
	var output strings.Builder
	e.SetStdout(&output)
	err := e.Exec("echo", "foobar")
	if err != nil {
		t.Fatalf("Command failed: %v", err)
	}
	if act, exp := output.String(), "foobar\n"; act != exp {
		t.Errorf("Incorrect output from command: want %q, got %q", exp, act)
	}
}

func TestExecutorStderr(t *testing.T) {
	e := NewRealExecutor()
	var stderr strings.Builder
	e.SetStderr(&stderr)
	if err := e.Exec("/bin/sh", "-c", "echo one-two >&2"); err != nil {
		t.Fatalf("Command failed; %v", err)
	}
	if act, exp := stderr.String(), "one-two\n"; act != exp {
		t.Errorf("Incorrect stderr from command: want %q, got %q", exp, act)
	}
}

func TestExecutorStdin(t *testing.T) {
	e := NewRealExecutor()
	input := strings.NewReader("hello-world")
	var stdout strings.Builder
	e.SetStdin(input)
	e.SetStdout(&stdout)
	if err := e.Exec("cat"); err != nil {
		t.Fatalf("Command failed: %v", err)
	}
	if act, exp := stdout.String(), "hello-world"; act != exp {
		t.Errorf("Program is not processed stdin correctly: want output %q, got %q", exp, act)
	}
}
