package generate

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"text/template"

	igit "github.com/gofrontier-com/go-utils/git"
	git "github.com/gofrontier-com/go-utils/git/external_git"
	"github.com/gofrontier-com/go-utils/output"
	"github.com/gofrontier-com/ranger/pkg/core"
	rtime "github.com/gofrontier-com/ranger/pkg/util/time"
)

//go:embed tpl/*
var wikiTemplates embed.FS

type Workload struct { // TODO: Move
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
	} else {
		output.PrintlnfInfo("No changes for wiki")
	}

	return nil
}

type WikiContent struct { // TODO: Move
	Sets      []core.SetCollection
	Workloads []core.Workload
}

func processTemplateFile(src string, tgt string, localPath string, wikiContent interface{}) error {
	tmpl, err := template.New(path.Base(src)).ParseFS(wikiTemplates, src)
	if err != nil {
		return err
	}
	var f *os.File
	err = os.MkdirAll(filepath.Dir(filepath.Join(localPath, tgt)), 0700)
	if err != nil {
		panic(err)
	}
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

func writeWiki(wikiContent *WikiContent, localPath string) error {
	loopDirs := []string{"tpl/sets", "tpl/workloads"}
	fs.WalkDir(wikiTemplates, "tpl", func(src string, d fs.DirEntry, err error) error {
		if slices.Contains(loopDirs, filepath.Dir(src)) {
			if strings.Contains(src, "workloads") {
				for _, w := range wikiContent.Workloads {
					tgt := "workloads/" + w.Name + ".md"
					terr := processTemplateFile(src, tgt, localPath, &w)
					if terr != nil {
						return nil
					}
				}
			}
			if strings.Contains(src, "sets") {
				for _, s := range wikiContent.Sets {
					tgt := "sets/" + s.Name + ".md"
					terr := processTemplateFile(src, tgt, localPath, &s)
					if terr != nil {
						return nil
					}
				}
			}
		}

		if !d.IsDir() && !slices.Contains(loopDirs, filepath.Dir(src)) {
			tgt := strings.Replace(src, "tpl"+string(os.PathSeparator), "", 1)
			tgt = strings.ReplaceAll(tgt, ".tpl", "")
			terr := processTemplateFile(src, tgt, localPath, wikiContent)
			if terr != nil {
				return terr
			}
		}
		return nil
	})
	return nil
}

func GenerateDocs(config *core.Config, projectName string, organisationName string, wikiName string, feedName string) error {
	ado := &core.AzureDevOps{
		OrganisationName: organisationName,
		ProjectName:      projectName,
		PAT:              config.ADO.PAT,
		WorkloadFeedName: feedName,
	}
	sets, err := ado.GetSets()
	if err != nil {
		return err
	}
	output.PrintlnfInfo("Fetched set info from project '%s/%s'", organisationName, projectName)

	workloads, err := ado.GetWorkloadInfo()
	if err != nil {
		return err
	}
	output.PrintlnfInfo("Fetched workload info from feed '%s' (https://dev.azure.com/%s/%s/_artifacts/feed/%s)", feedName, organisationName, projectName, feedName)

	wc := WikiContent{
		Workloads: *workloads,
		Sets:      *sets,
	}

	err = ado.CreateWikiIfNotExists(wikiName, config.Git.UserName, config.Git.UserEmail)
	if err != nil {
		return err
	}

	output.PrintlnfInfo("Created wiki '%s' (%s)", wikiName, ado.WikiRemoteUrl)

	wikiRepo, err := git.NewClonedGit(ado.WikiRepoRemoteUrl, "x-oauth-basic", config.ADO.PAT, config.Git.UserEmail, config.Git.UserName)
	if err != nil {
		return err
	}
	defer os.RemoveAll(wikiRepo.GetRepositoryPath())

	err = writeWiki(&wc, wikiRepo.GetRepositoryPath())
	if err != nil {
		return err
	}

	err = publish(ado, wikiName, wikiRepo)
	if err != nil {
		return errors.New("could not create or automerge pr")
	}

	output.PrintlnfInfo("Published Wiki '%s' (%s)", wikiName, ado.WikiRemoteUrl)

	return nil
}
