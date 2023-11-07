package core

import (
	"fmt"

	"github.com/frontierdigital/utils/azuredevops"
	"github.com/google/uuid"
)

func (ado *AzureDevOps) CreatePullRequest(repositoryName string, sourceBranchName string, targetBranchName string, title string) (*int, error) {
	azureDevOps := azuredevops.NewAzureDevOps(ado.OrganisationName, ado.PAT)

	pr, err := azureDevOps.CreatePullRequest(ado.ProjectName, repositoryName, fmt.Sprintf("refs/heads/%s", sourceBranchName), fmt.Sprintf("refs/heads/%s", targetBranchName), title)
	if err != nil {
		return nil, err
	}

	return pr.PullRequestId, nil
}

func (ado *AzureDevOps) SetPullRequestAutoComplete(repositoryName string, pullRequestId *int, identityId *uuid.UUID) error {
	azureDevOps := azuredevops.NewAzureDevOps(ado.OrganisationName, ado.PAT)
	return azureDevOps.SetPullRequestAutoComplete(ado.ProjectName, repositoryName, *pullRequestId, identityId)
}

func (ado *AzureDevOps) CreateWikiIfNotExists(wikiName string, gitUserName string, gitUserEmail string) error {
	azureDevOps := azuredevops.NewAzureDevOps(ado.OrganisationName, ado.PAT)

	wiki, repo, err := azureDevOps.CreateWikiIfNotExists(ado.ProjectName, wikiName, gitUserEmail, gitUserName)
	if err != nil {
		return err
	}
	ado.WikiRemoteUrl = *wiki.RemoteUrl
	ado.WikiRepoRemoteUrl = *repo.RemoteUrl

	return nil
}

func (ado *AzureDevOps) GetIdentityId() (*uuid.UUID, error) {
	azureDevOps := azuredevops.NewAzureDevOps(ado.OrganisationName, ado.PAT)

	identityId, err := azureDevOps.GetIdentityId()
	if err != nil {
		return nil, err
	}

	return identityId, nil
}

func (ado *AzureDevOps) GetWorkloadInfo() (*[]Workload, error) {
	azureDevOps := azuredevops.NewAzureDevOps(ado.OrganisationName, ado.PAT)

	var workloads []Workload

	packages, err := azureDevOps.GetPackageVersion(ado.ProjectName, ado.WorkloadFeedName)
	if err != nil {
		return nil, err
	}

	for _, p := range *packages {
		if len(*p.Versions) > 0 {
			c, _ := azureDevOps.GetFileContent(ado.ProjectName, *p.Name, *(*p.Versions)[0].Version, "README.md")
			workloads = append(workloads, Workload{
				Name:    *p.Name,
				Version: *(*p.Versions)[0].Version,
				Build:   "N/A",
				Readme:  *c.Content,
				// Instances: [],
			})
		}
	}

	return &workloads, nil
}
