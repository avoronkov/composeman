package cli

import "fmt"

type Strings struct {
	values []string
	def    bool
}

func StringsDefault(values []string) *Strings {
	return &Strings{
		values: values,
		def:    true,
	}
}

func (s *Strings) String() string {
	return fmt.Sprintf("%v", []string(s.values))
}

func (s *Strings) Values() []string {
	return s.values
}

func (s *Strings) Set(v string) error {
	if s.def {
		s.values = []string{v}
		return nil
	}
	s.values = append(s.values, v)
	return nil
}
