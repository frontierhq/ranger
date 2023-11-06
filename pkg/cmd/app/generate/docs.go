package generate

import (
	_ "embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/frontierdigital/ranger/pkg/core"
	rfile "github.com/frontierdigital/ranger/pkg/util/file"
	rtime "github.com/frontierdigital/ranger/pkg/util/time"
	igit "github.com/frontierdigital/utils/git"
	git "github.com/frontierdigital/utils/git/external_git"
	"github.com/frontierdigital/utils/output"
)

//go:embed tpl/workload.tpl
var workloadTemplate string

type Workload struct {
	Name    string
	Version string
	Build   string
}

func publish(ado *core.AzureDevOps, repoName string, repo interface{ igit.Git }) error {
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
		branchName := fmt.Sprintf("ranger/update/%s", us)

		err = repo.Checkout(branchName, true)
		if err != nil {
			return err
		}

		commitMessage := "Update docs."
		_, err = repo.Commit(commitMessage)
		if err != nil {
			return err
		}

		err = repo.Push(false)
		if err != nil {
			return err
		}

		prId, err := ado.CreatePullRequest(repoName, branchName, "main", "Update docs")
		if err != nil {
			return err
		}

		identityId, err := ado.GetIdentityId()
		if err != nil {
			return err
		}

		err = ado.SetPullRequestAutoComplete(repoName, prId, identityId)
		if err != nil {
			return err
		}
	}

	return nil
}

func createWorkLoadPages(workloads *[]core.Workload, localPath string) error {
	orderPath := filepath.Join(localPath, "workloads", ".order")

	err := rfile.Clear(orderPath)
	if err != nil {
		return errors.New("could not reset order file")
	}

	for _, w := range *workloads {
		fullPath := filepath.Join(localPath, "workloads", fmt.Sprintf("%s.md", w.Name))

		err := rfile.CreateOrUpdate(fullPath, w.Readme, false)
		if err != nil {
			return errors.New("could not create or update page")
		}

		err = rfile.CreateOrUpdate(orderPath, fmt.Sprintln(w.Name), true)
		if err != nil {
			return errors.New("could not create or update orderfile")
		}
	}

	return nil
}

func createWorkloadIndex(workloads *[]core.Workload, localPath string) error {
	wl := core.WorkloadIndex{
		Workloads: *workloads,
	}
	tmpl, err := template.New("workloadTemplate").Parse(workloadTemplate)
	if err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(localPath, "workloads.md"))
	if err != nil {
		return err
	}
	defer file.Close()

	err = tmpl.Execute(file, wl)
	if err != nil {
		return err
	}

	return nil
}

func GenerateDocs(config *core.Config, projectName string, organisationName string, wikiName string, feedName string) error {
	ado := &core.AzureDevOps{
		OrganisationName: organisationName,
		ProjectName:      projectName,
		PAT:              config.ADO.PAT,
		WorkloadFeedName: feedName,
	}

	err := ado.CreateWikiIfNotExists(wikiName, config.Git.UserName, config.Git.UserEmail)
	if err != nil {
		return err
	}

	output.PrintlnfInfo("Created wiki '%s' (%s)", wikiName, ado.WikiRemoteUrl)

	wikiRepo, err := git.NewClonedGit(ado.WikiRepoRemoteUrl, "x-oauth-basic", config.ADO.PAT, config.Git.UserEmail, config.Git.UserName)
	if err != nil {
		return err
	}
	defer os.RemoveAll(wikiRepo.GetRepositoryPath())

	workloads, err := ado.GetWorkloadInfo()
	if err != nil {
		return err
	}

	output.PrintlnfInfo("Fetched workload info from feed '%s' (https://dev.azure.com/%s/%s/_artifacts/feed/%s)", feedName, organisationName, projectName, feedName)

	err = createWorkLoadPages(workloads, wikiRepo.GetRepositoryPath())
	if err != nil {
		return err
	}

	err = createWorkloadIndex(workloads, wikiRepo.GetRepositoryPath())
	if err != nil {
		return err
	}

	output.PrintlnInfo("Generated workload index and pages")

	err = publish(ado, wikiName, wikiRepo)
	if err != nil {
		return errors.New("could not create or automerge PR")
	}

	output.PrintlnfInfo("Published Wiki '%s' (%s)", wikiName, ado.WikiRemoteUrl)

	return nil
}
