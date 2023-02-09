package deploy

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/frontierdigital/ranger/pkg/cmd/app"
	"github.com/frontierdigital/ranger/pkg/util/azure_devops"
	"github.com/frontierdigital/ranger/pkg/util/git"
	"github.com/frontierdigital/ranger/pkg/util/manifest"
	"github.com/frontierdigital/ranger/pkg/util/output"

	"github.com/microsoft/azure-devops-go-api/azuredevops"
	"github.com/otiai10/copy"
	"github.com/segmentio/ksuid"
)

func DeployManifest(config *app.Config, projectName string, organisationName string) error {
	organisationUrl := fmt.Sprintf("https://dev.azure.com/%s", organisationName)
	connection := azuredevops.NewPatConnection(organisationUrl, config.ADO.PAT)

	ctx := context.Background()

	manifestFilepath, _ := filepath.Abs("./manifest.yml")
	manifest, err := manifest.LoadManifest(manifestFilepath)
	if err != nil {
		return err
	}

	manifestName := fmt.Sprintf("%s-%s", manifest.Environment, manifest.Layer)

	manifest.PrintHeader(manifestName, manifest.Layer, manifest.Environment, manifest.Version)

	for _, workload := range manifest.Workloads {
		sourceProjectName, sourceRepositoryName := workload.GetSourceProjectAndRepositoryNames()

		pipelineName := fmt.Sprintf("%s (deploy)", sourceRepositoryName)
		pipeline, err := azure_devops.GetPipelineByName(ctx, connection, sourceProjectName, pipelineName)
		if err != nil {
			return err
		}

		output.PrintlnfInfo("Found deploy pipeline definition with Id '%d' for workload '%s' (https://dev.azure.com/%s/%s/_build?definitionId=%d)",
			*pipeline.Id, workload.Source, organisationName, projectName, *pipeline.Id)

		repository, err := azure_devops.GetRepositoryByName(ctx, connection, sourceProjectName, sourceRepositoryName)
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
			configRepoUrl := fmt.Sprintf("https://frontierdigital@dev.azure.com/%s/%s/_git/%s", organisationName, projectName, configRepoName)

			configRepo, configRepoPath, err := git.CloneOverHttp(configRepoUrl, config.ADO.PAT, "x-oauth-basic")
			if err != nil {
				return err
			}

			configBranchName := fmt.Sprintf("%s_%s_%s_%s", workload.Name, manifest.Environment, manifest.Layer, ksuid.New().String())
			branch, err := git.CheckoutBranch(configRepo, configBranchName, true)
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

			commitMessage := fmt.Sprintf("Generate config for workload '%s', environment '%s' and layer '%s'", workload.Name, manifest.Environment, manifest.Layer)
			_, err = git.Commit(configRepo, branch, config.Git.UserEmail, config.Git.UserName, commitMessage)
			if err != nil {
				return err
			}

			err = git.Push(configRepo, branch, config.ADO.PAT, "x-oauth-basic")
			if err != nil {
				return err
			}

			output.PrintlnfInfo("Pushed config (https://dev.azure.com/%s/%s/_git/%s?version=GB%s)", organisationName, projectName, configRepoName, configBranchName)

			defer os.RemoveAll(configRepoPath)
		}
	}

	return nil
}
