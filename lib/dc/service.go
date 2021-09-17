package dc

import "fmt"

type Service struct {
	Image       string      `yaml:"image"`
	Environment interface{} `yaml:"environment"`
	EnvFile     string      `yaml:"env_file"`
	Ports       []string    `yaml:"ports"`
	Volumes     []string    `yaml:"volumes"`
	Command     string      `yaml:"command"`
	DependsOn   []string    `yaml:"depends_on"`
	Build       *struct {
		Context    string            `yaml:"context"`
		Dockerfile string            `yaml:"dockerfile"`
		Target     string            `yaml:"target"`
		Args       map[string]string `yaml:"args"`
	} `yaml:"build"`
}

func (s *Service) String() string {
	return fmt.Sprintf("Image: %v\nEnvironment: %v\nPorts: %v\nVolumes: %v\n", s.Image, s.Environment, s.Ports, s.Volumes)
}

func (s *Service) Env() (envs []string, err error) {
	if s.Environment == nil {
		return nil, nil
	}
	if list, ok := s.Environment.([]interface{}); ok {
		for _, v := range list {
			envs = append(envs, v.(string))
		}
		return envs, nil
	}
	if mp, ok := s.Environment.(map[string]interface{}); ok {
		for k, v := range mp {
			if v != nil {
				envs = append(envs, fmt.Sprintf("%v=%v", k, v))
			} else {
				envs = append(envs, k)
			}
		}
		return envs, nil
	}
	return nil, fmt.Errorf("Unexpected type in environment: %v (%T)", s.Environment, s.Environment)
}
