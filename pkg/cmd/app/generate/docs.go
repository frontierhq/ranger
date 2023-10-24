package generate

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/frontierdigital/ranger/pkg/util/config"
	rfile "github.com/frontierdigital/ranger/pkg/util/file"
	rtime "github.com/frontierdigital/ranger/pkg/util/time"
	"github.com/frontierdigital/utils/azuredevops"
	egit "github.com/frontierdigital/utils/git/external_git"
	"github.com/frontierdigital/utils/output"
)

func publish(a *azuredevops.AzureDevOps, projectName string, repositoryName string, repoPath string) error {
	repo := egit.NewGit(repoPath)
	hasChanges, err := repo.HasChanges()
	if err != nil {
		return err
	}

	if hasChanges {
		err = repo.AddAll()
		if err != nil {
			return err
		}

		us := rtime.GetUnixTimestamp()
		branchName := "generate-docs-" + us

		err = repo.Checkout(branchName, true)
		if err != nil {
			return err
		}

		commitMessage := "Initial Commit"
		_, err = repo.Commit(commitMessage)
		if err != nil {
			return err
		}

		err = repo.Push(false)
		if err != nil {
			return err
		}

		output.PrintlnfInfo("Pushed.")

		pr, err := a.CreatePullRequest(projectName, repositoryName, fmt.Sprintf("refs/heads/%s", branchName), "refs/heads/main", "Update docs "+us)
		if err != nil {
			return err
		}

		identityId, err := a.GetIdentityId()
		if err != nil {
			return err
		}

		err = a.SetPullRequestAutoComplete(projectName, repositoryName, *pr.PullRequestId, identityId)
	}

	return nil
}

func GenerateDocs(config *config.Config, projectName string, organisationName string, repoName string, feedName string) error {
	azureDevOps := azuredevops.NewAzureDevOps(organisationName, config.ADO.PAT)

	localPath, err := azureDevOps.CreateWikiIfNotExists(projectName, repoName, config.Git.UserEmail, config.Git.UserName, config.ADO.PAT)
	if err != nil {
		return err
	}

	packages, err := azureDevOps.GetPackageVersion(projectName, feedName)
	if err != nil {
		return err
	}

	orderPath := filepath.Join(*localPath, "workloads", ".order")
	err = rfile.Clear(orderPath)
	if err != nil {
		return errors.New("Could not create or update page")
	}

	for _, p := range *packages {
		if len(*p.Versions) > 0 {
			c, _ := azureDevOps.GetFileContent(projectName, *p.Name, *(*p.Versions)[0].Version)
			fullPath := filepath.Join(*localPath, "workloads", fmt.Sprintf("%s.md", *p.Name))
			// orderPath := filepath.Join(*localPath, "workloads", ".order")
			err := rfile.CreateOrUpdate(fullPath, *c.Content, false)
			if err != nil {
				return errors.New("Could not create or update page")
			}
			err = rfile.CreateOrUpdate(orderPath, fmt.Sprintln(*p.Name), true)
			if err != nil {
				return errors.New("Could not create or update orderfile")
			}
		}
	}

	err = publish(azureDevOps, projectName, repoName, *localPath)
	if err != nil {
		return errors.New("Could not create or automerge PR")
	}

	return nil
}
