package tests

import (
	"os"
	"strings"
	"testing"

	"github.com/avoronkov/composeman/lib/cli"
	"github.com/avoronkov/composeman/lib/podman"
)

const podName02 = "02_run_helloworld_from_dockerfile"

func Test02RunHelloworldFromDockerfile(t *testing.T) {
	pwd := chdir(podName02)
	defer chdir(pwd)
	defer removePod(podName02)

	compile("-o", "helloworld.exe")

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

	// Check that pod contains stopped container
	podInfo, err := podman.InspectPod(podName02)
	if err != nil {
		t.Fatal(err)
	}
	if l := len(podInfo.Containers); l != 2 {
		t.Fatalf("Incorrect number of containers in pod: expected 2, found %v", l)
	}
	for _, cnt := range podInfo.Containers {
		if cnt.Id == podInfo.InfraContainerID {
			continue
		}

		if exp := []string{"exited", "stopped"}; !stringsContain(exp, cnt.State) {
			t.Errorf("Incorrect app containter state: expected one of %v, found %v", exp, cnt.State)
		}
	}
}

func stringsContain(list []string, s string) bool {
	for _, x := range list {
		if x == s {
			return true
		}
	}
	return false
}

// run --rm removes container after finishing.
func Test02RunHelloworldFromDockerfileRm(t *testing.T) {
	pwd := chdir(podName02)
	defer chdir(pwd)
	defer removePod(podName02)

	compile("-o", "helloworld.exe")

	var stdout strings.Builder
	c := cli.New(&stdout, os.Stderr)
	rc := c.Run([]string{"run", "--rm", "hello-world"})
	if rc != 0 {
		t.Fatalf("Command 'run hello-world' finished with non-zero exit code: %v", rc)
	}
	out := stdout.String()
	t.Logf("Stdout: %v\n", out)

	greeting := "Hello from Go's HelloWorld!"
	if !strings.Contains(out, greeting) {
		t.Errorf("String %q not found in output.", greeting)
	}

	// Check that pod is removed
	_, err := podman.InspectPod(podName02)
	if exp := podman.NotFoundError; err != exp {
		t.Errorf("Incorrect error: expected %v, found %v", exp, err)
	}
}
