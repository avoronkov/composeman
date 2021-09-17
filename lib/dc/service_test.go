package dc

import (
	"reflect"
	"sort"
	"testing"

	"gopkg.in/yaml.v3"
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

func TestEnv(t *testing.T) {
	tests := []struct {
		name string
		in   string
		out  []string
		err  error
	}{
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
				t.Fatalf("Incorrect error returned: want %v, got%v", test.err, err)
			}
			sort.Strings(act)
			sort.Strings(test.out)
			if !reflect.DeepEqual(act, test.out) {
				t.Errorf("Incorrect Env(): want %v, got %v", test.out, act)
			}
		})
	}
}
