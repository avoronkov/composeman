package dc

type dockerComposeFile struct {
	Services map[string]Service `yaml:"services"`
}
