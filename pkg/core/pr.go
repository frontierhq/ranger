package core

import (
	"fmt"

	"github.com/frontierdigital/utils/azuredevops"
	"github.com/google/uuid"
)

// GetWorkloadInfo gets the workload info
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
