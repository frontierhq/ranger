package core

import (
	"time"
)

type ADOConfig struct {
	PAT string `mapstructure:"PAT"`
}

type AzureDevOps struct {
	OrganisationName  string
	ProjectName       string
	PAT               string
	WikiRemoteUrl     string
	WikiRepoRemoteUrl string
	WorkloadFeedName  string
}

type Config struct {
	ADO ADOConfig `mapstructure:"ADO"`
	Git GitConfig `mapstructure:"Git"`
}

type ExtraParameter struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type GitConfig struct {
	UserEmail string `mapstructure:"UserEmail"`
	UserName  string `mapstructure:"UserName"`
}

type Manifest struct {
	Version         int64               `yaml:"version"`
	Environment     string              `yaml:"environment"`
	FilePath        string              `yaml:"-"`
	NextEnvironment string              `yaml:"nextEnvironment"`
	Set             string              `yaml:"set"`
	Workloads       []*WorkloadInstance `yaml:"workloads"`
}

type Set struct {
	Name        string
	Environment string
	Next        *Set
	Previous    *Set
	Manifest    *Manifest
}

type SetCollection struct {
	Name  string
	Entry *Set
	Sets  []*Set
}

type Workload struct {
	Name      string
	Version   string
	Build     string
	Readme    string
	Instances []*WorkloadInstance
}

type WorkloadIndex struct {
	Workloads []Workload
}

type WorkloadInstance struct {
	ExtraParameters []ExtraParameter `yaml:"extraParameters"`
	FilePath        string           `yaml:"-"`
	Name            string           `yaml:"name"`
	PreventDestroy  bool             `yaml:"preventDestroy"`
	Type            string           `yaml:"type"`
	Version         string           `yaml:"version"`
}

type WorkloadResult struct {
	FinishTime *time.Time
	Link       string
	QueueTime  *time.Time
	Status     WorkloadResultStatus
	Workload   *WorkloadInstance
}

type WorkloadResultStatus string

type workloadResultStatusValuesType struct {
	Succeeded WorkloadResultStatus
	Failed    WorkloadResultStatus
	Skipped   WorkloadResultStatus
}

var WorkloadResultStatusValuesType = workloadResultStatusValuesType{
	Succeeded: "Succeeded",
	Failed:    "Failed",
	Skipped:   "Skipped",
}

type WorkloadDestroyPreventedError struct {
	Workload *WorkloadInstance
}
