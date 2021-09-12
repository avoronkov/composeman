package podman

import (
	"encoding/json"
	"os"
	"os/exec"
	"strings"
)

func InspectPod(pod string) (*PodInspect, error) {
	args := []string{"pod", "inspect", pod}
	output := &strings.Builder{}
	stderr := &strings.Builder{}
	cmd := exec.Command("podman", args...)
	cmd.Stdout = output
	cmd.Stderr = stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		if strings.Contains(stderr.String(), "no such pod") {
			return nil, NotFoundError
		}
		return nil, err
	}

	pi := &PodInspect{}
	if err := json.Unmarshal([]byte(output.String()), pi); err != nil {
		return nil, err
	}

	return pi, nil
}
