package podman

import (
	"reflect"
	"testing"
)

func TestPodmanPodCreate(t *testing.T) {
	tests := []struct {
		name string
		opts []PodCreateOpt
		exp  []string
	}{
		{"no args", []PodCreateOpt{}, []string{"podman", "pod", "create", "--name", "my-pod"}},
		{"-p (--publish)", []PodCreateOpt{OptPublishPort("80:80", "81:81")}, []string{"podman", "pod", "create", "--name", "my-pod", "-p", "80:80", "-p", "81:81"}},
		{"-p (--publish) multiple times", []PodCreateOpt{OptPublishPort("80:80"), OptPublishPort("81:81")}, []string{"podman", "pod", "create", "--name", "my-pod", "-p", "80:80", "-p", "81:81"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pr := makeTestPodman()

			err := pr.PodCreate("my-pod", test.opts...)
			if err != nil {
				t.Fatalf("Podman.PodCreate failed: %v", err)
			}

			act := pr.executor.(*executorMock).lastCommand
			if !reflect.DeepEqual(act, test.exp) {
				t.Errorf("Incorrect command called: want %q, got %q", test.exp, act)
			}
		})
	}
}
