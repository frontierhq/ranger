package core

import (
	"github.com/frontierdigital/utils/azuredevops"
)

// GetWorkloadInfo gets the workload info
func (ado *AzureDevOps) CreateWikiIfNotExists(gitUserName string, gitUserEmail string) error {
	azureDevOps := azuredevops.NewAzureDevOps(ado.OrganisationName, ado.PAT)
	localPath, err := azureDevOps.CreateWikiIfNotExists(ado.ProjectName, ado.ProjectName, gitUserEmail, gitUserName, ado.PAT)
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
