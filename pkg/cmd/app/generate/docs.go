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
	egit "github.com/frontierdigital/utils/git/external_git"
	"github.com/frontierdigital/utils/output"
)

//go:embed tpl/workload.tpl
var workloadTemplate string

type Workload struct {
	Name    string
	Version string
	Build   string
}

func publish(ado *core.AzureDevOps) error {
	repo := egit.NewGit(ado.WikiRepo.LocalPath)
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

		prId, err := ado.CreatePullRequest(branchName, "Update docs "+us)
		if err != nil {
			return err
		}

		identityId, err := ado.GetIdentityId()
		if err != nil {
			return err
		}

		err = ado.SetPullRequestAutoComplete(prId, identityId)
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
		return errors.New("Could not reset order file")
	}

	for _, w := range *workloads {
		fullPath := filepath.Join(localPath, "workloads", fmt.Sprintf("%s.md", w.Name))

		err := rfile.CreateOrUpdate(fullPath, w.Readme, false)
		if err != nil {
			return errors.New("Could not create or update page")
		}

		err = rfile.CreateOrUpdate(orderPath, fmt.Sprintln(w.Name), true)
		if err != nil {
			return errors.New("Could not create or update orderfile")
		}
	}

	return nil
}

type WorkloadIndex struct {
	Workloads []core.Workload
}

func createWorkloadIndex(workloads *[]core.Workload, localPath string) error {
	wl := WorkloadIndex{
		Workloads: *workloads,
	}
	tmpl, err := template.New("workloadTemplate").Parse(workloadTemplate)
	if err != nil {
		return err
	}
	var f *os.File
	f, err = os.Create(filepath.Join(localPath, "workloads.md"))
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(f, wl)
	if err != nil {
		panic(err)
	}
	err = f.Close()
	if err != nil {
		panic(err)
	}

	return nil
}

func (*Workload) GetTemplate() string {
	return workloadTemplate
}

func GenerateDocs(config *core.Config, projectName string, organisationName string, repoName string, feedName string) error {
	ado := &core.AzureDevOps{
		OrganisationName: organisationName,
		ProjectName:      projectName,
		PAT:              config.ADO.PAT,
		WorkloadFeedName: feedName,
		WikiRepoName:     repoName,
	}

	workloads, err := ado.GetWorkloadInfo()
	if err != nil {
		return err
	}

	err = ado.CreateWikiIfNotExists(config.Git.UserName, config.Git.UserEmail)
	if err != nil {
		return err
	}

	err = createWorkLoadPages(workloads, ado.WikiRepo.LocalPath)
	if err != nil {
		return err
	}

	err = createWorkloadIndex(workloads, ado.WikiRepo.LocalPath)
	if err != nil {
		return err
	}

	output.Println(ado.WikiRepo.LocalPath)

	err = publish(ado)
	if err != nil {
		return errors.New("Could not create or automerge PR")
	}

	return nil
}
