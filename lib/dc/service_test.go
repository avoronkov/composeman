package dc

import (
	"reflect"
	"sort"
	"testing"

	yaml "gopkg.in/yaml.v3"
)

const inputEnvList = `
environment:
  - RACK_ENV=development
  - SHOW=true
  - SESSION_SECRET
`

const inputEnvMap = `
environment:
  RACK_ENV: development
  SHOW: 'true'
  SESSION_SECRET:
`

func TestServiceEnv(t *testing.T) {
	tests := []struct {
		name string
		in   string
		out  []string
		err  error
	}{
		{
			"no env",
			"",
			nil,
			nil,
		},
		{
			"env list",
			inputEnvList,
			[]string{"RACK_ENV=development", "SHOW=true", "SESSION_SECRET"},
			nil,
		},
		{
			"env map",
			inputEnvMap,
			[]string{"RACK_ENV=development", "SHOW=true", "SESSION_SECRET"},
			nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var s Service
			err := yaml.Unmarshal([]byte(test.in), &s)
			if err != nil {
				t.Fatal(err)
			}
			act, err := s.Env()
			if err != test.err {
				t.Fatalf("Incorrect error returned: want %v, got %v", test.err, err)
			}
			sort.Strings(act)
			sort.Strings(test.out)
			if !reflect.DeepEqual(act, test.out) {
				t.Errorf("Incorrect Env(): want %v, got %v", test.out, act)
			}
		})
	}
}

const inputEnvFile = `
env_file: .env
`

const inputEnvFileList = `
env_file:
- ./common.env
- ./apps/web.env
- /opt/runtime_opts.env
`

func TestServiceEnvFile(t *testing.T) {
	tests := []struct {
		name string
		in   string
		out  []string
		err  error
	}{
		{
			"no env_file",
			"",
			nil,
			nil,
		},
		{
			"single env_file",
			inputEnvFile,
			[]string{".env"},
			nil,
		},
		{
			"env_file list",
			inputEnvFileList,
			[]string{"./common.env", "./apps/web.env", "/opt/runtime_opts.env"},
			nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var s Service
			err := yaml.Unmarshal([]byte(test.in), &s)
			if err != nil {
				t.Fatal(err)
			}
			act, err := s.EnvFile()
			if err != test.err {
				t.Fatalf("Incorrect error returned: want %v, got %v", test.err, err)
			}
			if !reflect.DeepEqual(act, test.out) {
				t.Errorf("Incorrect EnvFile: want %v (%T), got %v (%T)", test.out, test.out, act, act)
			}
		})
	}
}
