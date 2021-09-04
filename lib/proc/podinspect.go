package proc

type PodInspect struct {
	InfraContainerID string `json:"InfraContainerID"`
	Containers       []struct {
		Id    string
		Name  string
		State string
	} `json:"Containers"`
}
