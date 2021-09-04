package dc

import (
	"fmt"
	"os"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

type DockerCompose struct {
	Services map[string]Service
}

func NewDockerCompose(files ...string) (*DockerCompose, error) {
	c := &DockerCompose{
		Services: map[string]Service{},
	}

	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			return nil, err
		}
		dec := yaml.NewDecoder(f)

		dcf := &dockerComposeFile{}
		if err := dec.Decode(dcf); err != nil {
			return nil, err
		}

		// Handle services
		for name, srv := range dcf.Services {
			c.Services[name] = srv
		}
	}

	return c, nil
}

func (c *DockerCompose) String() string {
	buf := &strings.Builder{}
	fmt.Fprintf(buf, "Services:\n")
	for name, srv := range c.Services {
		fmt.Fprintf(buf, "%v:\n%v\n", name, srv.String())
	}
	return buf.String()
}
