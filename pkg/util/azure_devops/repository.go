package azure_devops

import (
	"context"
	"fmt"

	"github.com/microsoft/azure-devops-go-api/azuredevops"
	"github.com/microsoft/azure-devops-go-api/azuredevops/git"
	"golang.org/x/exp/slices"
)

func GetRepositories(ctx context.Context, connection *azuredevops.Connection, projectName string) (*[]git.GitRepository, error) {
	gitClient, err := git.NewClient(ctx, connection)
	if err != nil {
		return nil, err
	}

	getRepositoriesArgs := git.GetRepositoriesArgs{
		Project: &projectName,
	}
	return gitClient.GetRepositories(ctx, getRepositoriesArgs)
}

func GetRepositoryByName(ctx context.Context, connection *azuredevops.Connection, projectName string, name string) (*git.GitRepository, error) {
	repositories, err := GetRepositories(ctx, connection, projectName)
	if err != nil {
		return nil, err
	}

	findRepositoryFunc := func(r git.GitRepository) bool { return *r.Name == name }
	repositoryIdx := slices.IndexFunc(*repositories, findRepositoryFunc)

	if repositoryIdx == -1 {
		return nil, fmt.Errorf("repository with name '%s' not found in project '%s'", name, projectName)
	}

	return &(*repositories)[repositoryIdx], nil
}
