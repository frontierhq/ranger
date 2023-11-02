package core

import (
	"fmt"
	"strings"

	"github.com/frontierdigital/utils/azuredevops"
	"github.com/google/uuid"
)

func (ado *AzureDevOps) CreatePullRequest(branchName string, message string) (*int, error) {
	azureDevOps := azuredevops.NewAzureDevOps(ado.OrganisationName, ado.PAT)

	pr, err := azureDevOps.CreatePullRequest(ado.ProjectName, ado.WikiRepoName, fmt.Sprintf("refs/heads/%s", branchName), "refs/heads/main", message)
	if err != nil {
		return nil, err
	}

	return pr.PullRequestId, nil
}

func (ado *AzureDevOps) SetPullRequestAutoComplete(pullRequestId *int, identityId *uuid.UUID) error {
	azureDevOps := azuredevops.NewAzureDevOps(ado.OrganisationName, ado.PAT)
	return azureDevOps.SetPullRequestAutoComplete(ado.ProjectName, ado.WikiRepoName, *pullRequestId, identityId)
}

func (ado *AzureDevOps) CreateWikiIfNotExists(gitUserName string, gitUserEmail string) error {
	azureDevOps := azuredevops.NewAzureDevOps(ado.OrganisationName, ado.PAT)
	localPath, err := azureDevOps.CreateWikiIfNotExists(ado.ProjectName, ado.WikiRepoName, gitUserEmail, gitUserName, ado.PAT)
	if err != nil {
		return err
	}
	ado.WikiRepo = &GitRepository{
		LocalPath: *localPath,
		UserName:  gitUserName,
		UserEmail: gitUserEmail,
	}
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
			c, _ := azureDevOps.GetFileContent(ado.ProjectName, *p.Name, *(*p.Versions)[0].Version)
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

func (ado *AzureDevOps) GetSets() (*[]Set, error) {
	azureDevOps := azuredevops.NewAzureDevOps(ado.OrganisationName, ado.PAT)
	var sets []Set
	repos, err := azureDevOps.GetRepositories(ado.ProjectName)
	if err != nil {
		return nil, err
	}
	for _, r := range *repos {
		if strings.HasSuffix(*r.Name, "-set") {
			sets = append(sets, Set{
				Name: *r.Name,
			})
		}
	}
	return &sets, nil
}
