package tests

import (
	"os"
	"strings"
	"testing"

	"github.com/avoronkov/composeman/lib/cli"
)

const podName07 = "07_env"

func Test07Env(t *testing.T) {
	podName := podName07
	pwd := chdir(podName)
	defer chdir(pwd)

	compile("-o", "prog.exe")

	tests := []struct {
		name     string
		command  []string
		expValue string
	}{
		{
			"docker-compose.yml env",
			[]string{"run", "--rm", "prog-env"},
			"from-docker-compose-yml-env",
		},
		{
			"docker-compose.yml envfile",
			[]string{"run", "--rm", "prog-envfile"},
			"from-envfile",
		},
		{
			"cli -e",
			[]string{"run", "--rm", "-e", "MY_TEST_ENV_VAR=from-cli-env", "prog-envfile"},
			"from-cli-env",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var stdout strings.Builder
			c := cli.New(&stdout, os.Stderr)
			rc := c.Run(test.command)
			if rc != 0 {
				t.Fatalf("Command %v finished with non-zero exit code: %v", strings.Join(test.command, " "), rc)
			}

			out := stdout.String()
			t.Logf("Stdout: %v\n", out)

			if !strings.Contains(out, test.expValue) {
				t.Errorf("String %q not found in output.", test.expValue)
			}
		})
	}
}
