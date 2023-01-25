package deploy

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/frontierdigital/ranger/pkg/cmd/app"
	"github.com/frontierdigital/ranger/pkg/util/manifest"
	"github.com/frontierdigital/ranger/pkg/util/output"
	"github.com/microsoft/azure-devops-go-api/azuredevops"
	"github.com/microsoft/azure-devops-go-api/azuredevops/build"
	"github.com/microsoft/azure-devops-go-api/azuredevops/git"
	"golang.org/x/exp/slices"
)

func DeployManifest(config *app.Config, projectName string, organisationName string) error {
	organisationUrl := fmt.Sprintf("https://dev.azure.com/%s", organisationName)
	connection := azuredevops.NewPatConnection(organisationUrl, config.ADO.PAT)

	ctx := context.Background()

	buildClient, err := build.NewClient(ctx, connection)
	if err != nil {
		return err
	}

	gitClient, err := git.NewClient(ctx, connection)
	if err != nil {
		return err
	}

	manifestFilepath, _ := filepath.Abs("/Users/fraserdavidson/Code/Frontier/GitHub/ranger/manifest.yml")
	manifest, err := manifest.LoadManifest(manifestFilepath)
	if err != nil {
		return err
	}

	manifestName := fmt.Sprintf("%s-%s", manifest.Environment, manifest.Layer)

	manifest.PrintHeader(manifestName, manifest.Layer, manifest.Environment, manifest.Version)

	for _, workload := range manifest.Workloads {
		sourceProjectName, sourceRepositoryName := workload.GetSourceProjectAndRepositoryNames()

		pipelineName := fmt.Sprintf("%s (deploy)", sourceRepositoryName)
		getDefinitionsArgs := build.GetDefinitionsArgs{
			Name:    &pipelineName,
			Project: &sourceProjectName,
		}
		pipelines, err := buildClient.GetDefinitions(ctx, getDefinitionsArgs)
		if err != nil {
			return err
		}
		if len(pipelines.Value) == 0 {
			return fmt.Errorf("pipeline with name '%s' not found", pipelineName)
		}
		if len(pipelines.Value) > 1 {
			return fmt.Errorf("multiple pipeline with name '%s' found", pipelineName)
		}
		pipeline := pipelines.Value[0]

		output.PrintlnfInfo("Found deploy pipeline definition with Id '%d' for workload '%s' (https://dev.azure.com/%s/%s/_build?definitionId=%d)",
			*pipeline.Id, workload.Source, organisationName, projectName, *pipeline.Id)

		getRepositoriesArgs := git.GetRepositoriesArgs{
			Project: &sourceProjectName,
		}
		repositories, err := gitClient.GetRepositories(ctx, getRepositoriesArgs)
		if err != nil {
			return err
		}
		findRepositoryFunc := func(r git.GitRepository) bool { return r.Name == &sourceRepositoryName }
		repositoryIdx := slices.IndexFunc(*repositories, findRepositoryFunc)

		_ = repositories[repositoryIdx]
	}

	return nil
}
