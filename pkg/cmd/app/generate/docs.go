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

func processTemplateFile(src string, tgt string, localPath string, wikiContent *WikiContent) error {
	tmpl, err := template.New(path.Base(src)).ParseFS(wikiTemplates, src)
	if err != nil {
		return err
	}
	var f *os.File
	f, err = os.Create(filepath.Join(localPath, tgt))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = tmpl.Execute(f, wikiContent)
	if err != nil {
		panic(err)
	}
	return nil
}

func processTemplateFileWorkload(src string, tgt string, localPath string, workload *core.Workload) error {
	tmpl, err := template.New(path.Base(src)).ParseFS(wikiTemplates, src)
	if err != nil {
		return err
	}
	var f *os.File
	f, err = os.Create(filepath.Join(localPath, tgt))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = tmpl.Execute(f, workload)
	if err != nil {
		panic(err)
	}
	return nil
}

func writeWiki(wikiContent *WikiContent, localPath string) error {
	fs.WalkDir(wikiTemplates, "tpl", func(src string, d fs.DirEntry, err error) error {
		if src == "tpl/workloads/workload.tpl" {
			for _, w := range wikiContent.Workloads {
				tgt := "workloads/" + w.Name + ".md"
				terr := processTemplateFileWorkload(src, tgt, localPath, &w)
				if terr != nil {
					return nil
				}
			}
		}
		if !d.IsDir() && src != "tpl/workloads/workload.tpl" {
			tgt := strings.Replace(src, "tpl"+string(os.PathSeparator), "", 1)
			tgt = strings.Replace(tgt, ".tpl", ".md", 1)
			if strings.HasSuffix(tgt, "order.md") {
				tgt = strings.ReplaceAll(tgt, "order.md", ".order")
			}
			terr := processTemplateFile(src, tgt, localPath, wikiContent)
			if terr != nil {
				return terr
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

	err = publish(ado)
	if err != nil {
		return errors.New("could not create or automerge pr")
	}

	return nil
}
