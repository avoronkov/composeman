package dc

import "fmt"

type Service struct {
	Image       string   `yaml:"image"`
	Environment []string `yaml:"environment"`
	Ports       []string `yaml:"ports"`
	Volumes     []string `yaml:"volumes"`
}

func (s *Service) String() string {
	return fmt.Sprintf("Image: %v\nEnvironment: %v\nPorts: %v\nVolumes: %v\n", s.Image, s.Environment, s.Ports, s.Volumes)
}
