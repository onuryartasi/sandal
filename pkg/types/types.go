package types


type Project struct{
	Name string
	Image string
	Containers []string
	//ContainerOptions types.ContainerCreateConfig ## this config later
	Min	int
	Max int
	//CpuLimit float32
}
