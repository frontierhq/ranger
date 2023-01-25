package app

import "strings"

type ADOConfig struct {
	PAT string `mapstructure:"PAT"`
}

type Config struct {
	ADO ADOConfig `mapstructure:"ADO"`
}

type Manifest struct {
	Environment string      `yaml:"environment"`
	Layer       string      `yaml:"layer"`
	Version     int64       `yaml:"version"`
	Workloads   []*Workload `yaml:"workloads"`
}

func (m Manifest) PrintHeader(name string, layer string, environment string, version int64) {

}

type Workload struct {
	Name    string `yaml:"name"`
	Source  string `yaml:"source"`
	Version string `yaml:"version"`
}

func (w Workload) GetSourceProjectAndRepositoryNames() (string, string) {
	sourceParts := strings.Split(w.Source, "/")
	return sourceParts[0], sourceParts[1]
}
