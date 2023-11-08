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
	Version     int64               `yaml:"version"`
	Environment string              `yaml:"environment"`
	FilePath    string              `yaml:"-"`
	Set         string              `yaml:"set"`
	Workloads   []*WorkloadInstance `yaml:"workloads"`
}

type Set struct {
	Name        string
	Environment string
	Next        *Set
	Previous    *Set
	Manifest    *Manifest
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
	Name            string           `yaml:"name"`
	PreventDestroy  bool             `yaml:"preventDestroy"`
	Type            string           `yaml:"type"`
	Version         string           `yaml:"version"`
}

type WorkloadResult struct {
	Error      error
	FinishTime *time.Time
	Link       string
	QueueTime  *time.Time
	Workload   *WorkloadInstance
}
