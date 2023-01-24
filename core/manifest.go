package core

type Manifest struct {
	Environment string      `yaml:"environment"`
	Layer       string      `yaml:"layer"`
	Version     int64       `yaml:"version"`
	Workloads   []*Workload `yaml:"workloads"`
}
