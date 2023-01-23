package structs

type Manifest struct {
	Environment string      `yaml:"environment"`
	Layer       string      `yaml:"layer"`
	Version     string      `yaml:"version"`
	Workloads   []*Workload `yaml:"workloads"`
}
