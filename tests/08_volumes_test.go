package tests

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/avoronkov/composeman/lib/cli"
)

const podName08 = "08_volumes"

func Test08VolumesRo(t *testing.T) {
	podName := podName08
	pwd := chdir(podName)
	defer chdir(pwd)
	defer removePod(podName)

	compile("-o", "prog.exe")
	var stdout strings.Builder
	c := cli.New(&stdout, os.Stderr)
	command := []string{"run", "--rm", "prog-ro"}
	rc := c.Run(command)
	if rc != 0 {
		t.Fatalf("Command %v finished with non-zero exit code: %v", strings.Join(command, " "), rc)
	}

	out := stdout.String()
	t.Logf("Stdout: %v\n", out)

	exp := "Test file content: Mounted file test content."
	if !strings.Contains(out, exp) {
		t.Errorf("String %q not found in output.", exp)
	}
}

func Test08VolumesRw(t *testing.T) {
	podName := podName08
	pwd := chdir(podName)
	defer chdir(pwd)

	for _, existing := range []bool{true, false} {
		t.Run(fmt.Sprintf("existing mount-dir: %v", existing), func(t *testing.T) {
			defer removePod(podName)

			testDir := "volume-rw"
			testFile := filepath.Join(testDir, "testfile.out")

			defer func() {
				if err := os.RemoveAll(testDir); err != nil {
					panic(err)
				}
			}()

			if existing {
				if err := os.Mkdir(testDir, 0755); err != nil {
					panic(err)
				}
			}

			compile("-o", "prog.exe")
			var stdout strings.Builder
			c := cli.New(&stdout, os.Stderr)
			command := []string{"run", "--rm", "prog-rw"}
			rc := c.Run(command)
			if rc != 0 {
				t.Fatalf("Command %v finished with non-zero exit code: %v", strings.Join(command, " "), rc)
			}

			out := stdout.String()
			t.Logf("Stdout: %v\n", out)

			exp := "Data written to test file: Generated content."
			if !strings.Contains(out, exp) {
				t.Errorf("String %q not found in output.", exp)
			}

			// check that testFile is written
			data, err := ioutil.ReadFile(testFile)
			if err != nil {
				t.Fatalf("Cannot read test file %v: %v", testFile, err)
			}
			if act, exp := string(data), "Generated content."; act != exp {
				t.Errorf("Incorrect data in test file: want %q, got %q", exp, act)
			}
		})
	}
}
