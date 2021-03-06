package podman

import (
	"reflect"
	"strings"
	"testing"
)

func makeTestPodman() *Podman {
	pr := NewPodman()
	pr.executor = &executorMock{}
	return pr
}

func TestPodmanRun(t *testing.T) {
	tests := []struct {
		name string
		opts []RunOpt
		exp  []string
	}{
		{"no args", []RunOpt{}, []string{"podman", "run", "my-image"}},
		{"--rm", []RunOpt{OptRm(true)}, []string{"podman", "run", "--rm", "my-image"}},
		{"--pod", []RunOpt{OptPod("my-pod")}, []string{"podman", "run", "--pod", "my-pod", "my-image"}},
		{"-d (--detach)", []RunOpt{OptDetach(true)}, []string{"podman", "run", "-d", "my-image"}},
		{
			"-v (--volume)",
			[]RunOpt{OptVolume("./apps/common:/app/apps/common:ro", "./apps/frontend:/app/apps/frontend:ro")},
			[]string{"podman", "run", "--security-opt", "label=disable", "-v", "./apps/common:/app/apps/common:ro", "-v", "./apps/frontend:/app/apps/frontend:ro", "my-image"},
		},
		{
			"-v (--volume) multiple times",
			[]RunOpt{OptVolume("./apps/common:/app/apps/common:ro"), OptVolume("./apps/frontend:/app/apps/frontend:ro")},
			[]string{"podman", "run", "--security-opt", "label=disable", "-v", "./apps/common:/app/apps/common:ro", "-v", "./apps/frontend:/app/apps/frontend:ro", "my-image"},
		},
		{"--env-file", []RunOpt{OptEnvFile(".env.local")}, []string{"podman", "run", "--env-file", ".env.local", "my-image"}},
		{"--env-file (2)", []RunOpt{OptEnvFile(".env.global", ".env.local")}, []string{"podman", "run", "--env-file", ".env.global", "--env-file", ".env.local", "my-image"}},
		{"--env-file multiple times", []RunOpt{OptEnvFile(".env.global"), OptEnvFile(".env.local")}, []string{"podman", "run", "--env-file", ".env.global", "--env-file", ".env.local", "my-image"}},
		{"--env-file", []RunOpt{OptEnvFile(".env.global", ".env.local")}, []string{"podman", "run", "--env-file", ".env.global", "--env-file", ".env.local", "my-image"}},
		{"-e (--env)", []RunOpt{OptEnv("FOO=bar", "HELLO=world")}, []string{"podman", "run", "-e", "FOO=bar", "-e", "HELLO=world", "my-image"}},
		{"-e (--env) multiple times", []RunOpt{OptEnv("FOO=bar"), OptEnv("HELLO=world")}, []string{"podman", "run", "-e", "FOO=bar", "-e", "HELLO=world", "my-image"}},
		{
			"--add-host",
			[]RunOpt{OptLocalHost("my-service", "db")},
			[]string{"podman", "run", "--add-host", "my-service:127.0.0.1", "--add-host", "db:127.0.0.1", "my-image"},
		},
		{
			"--add-host multiple times",
			[]RunOpt{OptLocalHost("my-service"), OptLocalHost("db")},
			[]string{"podman", "run", "--add-host", "my-service:127.0.0.1", "--add-host", "db:127.0.0.1", "my-image"},
		},
		{
			"command (string)",
			[]RunOpt{OptCmdString("sh -c \"echo hello-world\"")},
			[]string{"podman", "run", "my-image", "sh", "-c", "echo hello-world"},
		},
		{
			"command (list)",
			[]RunOpt{OptCmdList("sh", "-c", "echo hello-world")},
			[]string{"podman", "run", "my-image", "sh", "-c", "echo hello-world"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pr := makeTestPodman()

			err := pr.Run("my-image", test.opts...)
			if err != nil {
				t.Fatalf("Podman.Run failed: %v", err)
			}

			act := pr.executor.(*executorMock).lastCommand
			if !reflect.DeepEqual(act, test.exp) {
				t.Errorf("Incorrect command called: want %q, got %q", test.exp, act)
			}
		})
	}
}

// Run
func TestRunOptCmdStringIncorrect(t *testing.T) {
	pr := makeTestPodman()

	err := pr.Run("my-image", OptCmdString("sh -c \"hello"))
	if err == nil || !strings.Contains(err.Error(), "Unterminated double-quoted string") {
		t.Errorf("Incorrect error returned on bad shell string command: %v", err)
	}
}
