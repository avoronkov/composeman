package podman

import (
	"reflect"
	"testing"
)

func TestRun(t *testing.T) {
	// podman run my-service
	pr := NewPodmanRun()
	pr.executor = &executorMock{}

	err := pr.Exec("my-service")
	if err != nil {
		t.Fatalf("PodmanRun.Exec failed: %v", err)
	}

	exp := []string{"podman", "run", "my-service"}
	act := pr.executor.(*executorMock).lastCommand
	if !reflect.DeepEqual(act, exp) {
		t.Errorf("Incorrect command called: want %q, got %q", exp, act)
	}
}

func TestRunRm(t *testing.T) {
	// podman run --rm my-service
	pr := NewPodmanRun()
	pr.executor = &executorMock{}

	err := pr.Exec("my-service", OptRm())
	if err != nil {
		t.Fatalf("PodmanRun.Exec failed: %v", err)
	}

	exp := []string{"podman", "run", "--rm", "my-service"}
	act := pr.executor.(*executorMock).lastCommand
	if !reflect.DeepEqual(act, exp) {
		t.Errorf("Incorrect command called: want %q, got %q", exp, act)
	}
}
