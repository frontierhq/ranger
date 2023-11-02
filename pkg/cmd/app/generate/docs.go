package generate

import (
	"embed"
	"errors"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/frontierdigital/ranger/pkg/core"
	rtime "github.com/frontierdigital/ranger/pkg/util/time"
	egit "github.com/frontierdigital/utils/git/external_git"
	"github.com/frontierdigital/utils/output"
)

//go:embed tpl/*
var wikiTemplates embed.FS

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

type WikiContent struct {
	Sets      []core.Set
	Workloads []core.Workload
}

func writeWiki(wikiContent *WikiContent, localPath string) error {
	fs.WalkDir(wikiTemplates, "tpl", func(p string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			output.Println(p)
			target := strings.Replace(p, "tpl"+string(os.PathSeparator), "", 1)
			target = strings.Replace(target, ".tpl", ".md", 1)
			tmpl, err := template.New(path.Base(p)).ParseFS(wikiTemplates, p)
			if err != nil {
				return err
			}
			var f *os.File
			f, err = os.Create(filepath.Join(localPath, target))
			if err != nil {
				panic(err)
			}
			defer f.Close()
			err = tmpl.Execute(f, wikiContent)
			if err != nil {
				panic(err)
			}
		}
		return nil
	})
	return nil
}

func GenerateDocs(config *core.Config, projectName string, organisationName string, repoName string, feedName string) error {
	ado := &core.AzureDevOps{
		OrganisationName: organisationName,
		ProjectName:      projectName,
		PAT:              config.ADO.PAT,
		WorkloadFeedName: feedName,
		WikiRepoName:     repoName,
	}
	sets, err := ado.GetSets()
	if err != nil {
		return err
	}
	output.Println(sets)

	workloads, err := ado.GetWorkloadInfo()
	if err != nil {
		return err
	}

	wc := WikiContent{
		Workloads: *workloads,
		Sets:      *sets,
	}

	err = ado.CreateWikiIfNotExists(config.Git.UserName, config.Git.UserEmail)
	if err != nil {
		return err
	}

	err = writeWiki(&wc, ado.WikiRepo.LocalPath)
	if err != nil {
		return err
	}
	output.Println(ado.WikiRepo.LocalPath)

	err = publish(ado)
	if err != nil {
		return errors.New("could not create or automerge pr")
	}

	return nil
}
