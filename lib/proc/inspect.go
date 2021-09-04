package proc

type Inspect struct {
	Mounts []struct {
		Type string
		Name string
	} `json:"Mounts"`
}
