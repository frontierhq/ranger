package deploy

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/frontierdigital/ranger/pkg/cmd/app/type/config"
	"github.com/frontierdigital/ranger/pkg/util/manifest"
	"github.com/frontierdigital/utils/azuredevops"
	"github.com/frontierdigital/utils/git"
	"github.com/frontierdigital/utils/output"

	"github.com/otiai10/copy"
	"github.com/segmentio/ksuid"
)

func DeployManifest(config *config.Config, projectName string, organisationName string) error {
	azureDevOps := azuredevops.NewAzureDevOps(organisationName, config.ADO.PAT)

	manifestFilepath, _ := filepath.Abs("./manifest.yml")
	manifest, err := manifest.LoadManifest(manifestFilepath)
	if err != nil {
		return err
	}

	manifest.PrintHeader()

	manifest.PrintWorkloadsSummary()

	for _, workload := range manifest.Workloads {
		workload.PrintHeader()

		sourceProjectName, sourceRepositoryName := workload.GetSourceProjectAndRepositoryNames()

		pipelineName := fmt.Sprintf("%s (deploy)", sourceRepositoryName)
		pipeline, err := azureDevOps.GetPipelineByName(sourceProjectName, pipelineName)
		if err != nil {
			return err
		}

		output.PrintlnfInfo("Found deploy pipeline definition with Id '%d' for workload '%s' (https://dev.azure.com/%s/%s/_build?definitionId=%d)",
			*pipeline.Id, workload.Source, organisationName, projectName, *pipeline.Id)

		repository, err := azureDevOps.GetRepositoryByName(sourceProjectName, sourceRepositoryName)
		if err != nil {
			return err
		}

		output.PrintlnfInfo("Found repository with Id '%s' for workload '%s' (%s)", repository.Id, workload.Source, *repository.WebUrl)

		workloadConfigPath := path.Join("config", "workloads", workload.Name)
		workloadConfigExists := true
		_, err = os.Stat(workloadConfigPath)
		if err != nil {
			workloadConfigExists = !os.IsNotExist(err)
		}

		workloadSecretsPath := path.Join("secrets", "workloads", workload.Name)
		workloadSecretsExists := true
		_, err = os.Stat(workloadSecretsPath)
		if err != nil {
			workloadSecretsExists = !os.IsNotExist(err)
		}

		if workloadConfigExists || workloadSecretsExists {
			configRepoName := "generated-manifest-config"
			configRepoUrl := fmt.Sprintf("https://dev.azure.com/%s/%s/_git/%s", organisationName, projectName, configRepoName)

			configRepoPath, err := os.MkdirTemp("", "")
			if err != nil {
				return err
			}
			configRepo := git.NewGit(configRepoPath)
			err = configRepo.CloneOverHttp(configRepoUrl, config.ADO.PAT, "x-oauth-basic")
			if err != nil {
				return err
			}
			err = configRepo.SetConfig("user.email", config.Git.UserEmail)
			if err != nil {
				return err
			}
			err = configRepo.SetConfig("user.name", config.Git.UserName)
			if err != nil {
				return err
			}

			configBranchName := fmt.Sprintf("%s_%s_%s_%s", workload.Name, manifest.Environment, manifest.Layer, ksuid.New().String())
			err = configRepo.Checkout(configBranchName, true)
			if err != nil {
				return err
			}

			if workloadConfigExists {
				err = copy.Copy(workloadConfigPath, path.Join(configRepoPath, ".config"))
				if err != nil {
					return err
				}
			}

			if workloadSecretsExists {
				err = copy.Copy(workloadSecretsPath, path.Join(configRepoPath, ".secrets"))
				if err != nil {
					return err
				}
			}

			err = configRepo.AddAll()
			if err != nil {
				return err
			}

			commitMessage := fmt.Sprintf("Generate config for workload '%s', environment '%s' and layer '%s'", workload.Name, manifest.Environment, manifest.Layer)
			err = configRepo.Commit(commitMessage)
			if err != nil {
				return err
			}

			err = configRepo.Push(false)
			if err != nil {
				return err
			}

			output.PrintlnfInfo("Pushed config (https://dev.azure.com/%s/%s/_git/%s?version=GB%s)", organisationName, projectName, configRepoName, configBranchName)

			workload.PrintFooter("", "", "", "")
			defer os.RemoveAll(configRepoPath)
		}
	}

	return nil
}
