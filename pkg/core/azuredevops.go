package core

import (
	"fmt"
	"strings"

	"github.com/gofrontier-com/go-utils/azuredevops"
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
			c, _ := azureDevOps.GetFileContent(ado.ProjectName, *p.Name, *(*p.Versions)[0].Version, "README.md", "tag")
			workloads = append(workloads, Workload{
				Name:    strings.ReplaceAll(*p.Name, "-workload", ""),
				Version: *(*p.Versions)[0].Version,
				Build:   "N/A",
				Readme:  *c.Content,
				// Instances: [],
			})
		}
	}

	return &workloads, nil
}

func getManifestContent(azureDevOps *azuredevops.AzureDevOps, projectName *string, repoName *string) (*Manifest, error) {
	m := "main"
	c, err := azureDevOps.GetFileContent(*projectName, *repoName, m, "manifest.yml", "branch")
	if err != nil {
		return nil, err
	}

	man, err := LoadManifestFromString(*c.Content)
	if err != nil {
		return nil, err
	}

	return &man, nil
}

func updateSetCollections(sets *[]SetCollection, sc *SetCollection) {
	for i, s := range *sets {
		if s.Name == sc.Name {
			(*sets)[i] = *sc
		}
	}
}

func (ado *AzureDevOps) GetSets() (*[]SetCollection, error) {
	azureDevOps := azuredevops.NewAzureDevOps(ado.OrganisationName, ado.PAT)
	var sets []SetCollection

	repos, err := azureDevOps.GetRepositories(ado.ProjectName)
	if err != nil {
		return nil, err
	}

	var sc *SetCollection
	for _, r := range *repos {
		if strings.HasSuffix(*r.Name, "-set") {
			n := getSetNameFromRepoName(r.Name)

			m, err := getManifestContent(azureDevOps, &ado.ProjectName, r.Name)
			if err != nil {
				return nil, err
			}

			sc = getSetCollectionByName(&sets, *n)
			if sc == nil {
				sets = newSetCollection(&sets, *n)
				sc = getSetCollectionByName(&sets, *n)
			}
			sc.addSet(&Set{
				Name:        *n,
				Manifest:    m,
				Environment: m.Environment,
				Next:        nil,
				Previous:    nil,
			})
			updateSetCollections(&sets, sc)
		}
	}

	for _, s := range sets {
		s.Order()
		updateSetCollections(&sets, &s)
	}
	return &sets, nil
}
