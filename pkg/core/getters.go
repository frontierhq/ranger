package core

import (
	"github.com/frontierdigital/utils/azuredevops"
)

// GetWorkloadInfo gets the workload info
func (ado *AzureDevOps) GetWorkloadInfo() (*[]Workload, error) {
	a := azuredevops.NewAzureDevOps(ado.OrganisationName, ado.PAT)

	var workloads []Workload

	packages, err := a.GetPackageVersion(ado.ProjectName, ado.WorkloadFeedName)
	if err != nil {
		return nil, err
	}

	for _, p := range *packages {
		if len(*p.Versions) > 0 {
			c, _ := a.GetFileContent(ado.ProjectName, *p.Name, *(*p.Versions)[0].Version)
			workloads = append(workloads, Workload{
				Name:      *p.Name,
				Version:   *(*p.Versions)[0].Version,
				Build:     "N/A",
				Readme:    *c.Content,
				Instances: 1,
			})
		}
	}

	return &workloads, nil
}
