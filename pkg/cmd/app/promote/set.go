package promote

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gofrontier-com/go-utils/azuredevops"
	git "github.com/gofrontier-com/go-utils/git/external_git"
	"github.com/gofrontier-com/go-utils/output"
	"github.com/gofrontier-com/ranger/pkg/core"
)

func PromoteSet(config *core.Config, projectName string, organisationName string) error {
	azureDevOps := azuredevops.NewAzureDevOps(organisationName, config.ADO.PAT)

	sourceManifestFilepath, _ := filepath.Abs("./manifest.yml")
	sourceManifest, err := core.LoadManifest(sourceManifestFilepath)
	if err != nil {
		return err
	}

	sourceManifest.PrintHeader()

	sourceManifest.PrintWorkloadsSummary()

	if sourceManifest.NextEnvironment == "" {
		output.PrintfInfo("Action: None (no next environment specified)\n\n")
		return nil
	}

	output.PrintfInfo("Action: Promote to %s\n\n", sourceManifest.NextEnvironment)

	nextEnvironmentSetRepoName := fmt.Sprintf("%s-%s-set", sourceManifest.NextEnvironment, sourceManifest.Set)
	nextEnvironmentSetRepoUrl := fmt.Sprintf("https://dev.azure.com/%s/%s/_git/%s", organisationName, projectName, nextEnvironmentSetRepoName)

	nextEnvironmentSetRepo, err := git.NewClonedGit(nextEnvironmentSetRepoUrl, "x-oauth-basic", config.ADO.PAT, config.Git.UserEmail, config.Git.UserName)
	if err != nil {
		return err
	}
	defer os.RemoveAll(nextEnvironmentSetRepo.GetRepositoryPath())

	output.PrintfInfo("Cloned target environment set repository '%s' (https://dev.azure.com/%s/%s/_git/%s)", nextEnvironmentSetRepoName, organisationName, projectName, nextEnvironmentSetRepoName)

	promoteBranchName := fmt.Sprintf("ranger/promote/%s", sourceManifest.Environment)
	err = nextEnvironmentSetRepo.Checkout(promoteBranchName, true)
	if err != nil {
		return err
	}

	targetManifestFilePath := nextEnvironmentSetRepo.GetFilePath("manifest.yml")
	targetManifest, err := core.LoadManifest(targetManifestFilePath)
	if err != nil {
		return err
	}

	targetManifest.Version = sourceManifest.Version
	targetManifest.Workloads = sourceManifest.Workloads

	targetManifest.Save()

	err = nextEnvironmentSetRepo.AddAll()
	if err != nil {
		return err
	}

	commitMessage := fmt.Sprintf("Promote set version %d from %s", sourceManifest.Version, sourceManifest.Environment)
	_, err = nextEnvironmentSetRepo.Commit(commitMessage)
	if err != nil {
		return err
	}

	err = nextEnvironmentSetRepo.Push(true)
	if err != nil {
		return err
	}

	output.PrintlnfInfo("Pushed branch '%s' (https://dev.azure.com/%s/%s/_git/%s?version=GB%s)", promoteBranchName, organisationName, projectName, nextEnvironmentSetRepoName, promoteBranchName)

	existingPullRequest, err := azureDevOps.FindPullRequest(projectName, nextEnvironmentSetRepoName, fmt.Sprintf("refs/heads/%s", promoteBranchName), "refs/heads/main")
	if err != nil {
		return err
	}

	if existingPullRequest != nil {
		_, err = azureDevOps.AbandonPullRequest(projectName, nextEnvironmentSetRepoName, *existingPullRequest.PullRequestId)
		if err != nil {
			return err
		}

		output.PrintlnfInfo("Abandoned existing pull request with Id '%d' (https://dev.azure.com/%s/%s/_git/%s/pullrequest/%d)", *existingPullRequest.PullRequestId, organisationName, projectName, nextEnvironmentSetRepoName, *existingPullRequest.PullRequestId)
	}

	pullRequestTitle := commitMessage
	pullRequest, err := azureDevOps.CreatePullRequest(projectName, nextEnvironmentSetRepoName, fmt.Sprintf("refs/heads/%s", promoteBranchName), "refs/heads/main", pullRequestTitle)
	if err != nil {
		return err
	}

	output.PrintfInfo("Created pull request with Id '%d' (https://dev.azure.com/%s/%s/_git/%s/pullrequest/%d)", *pullRequest.PullRequestId, organisationName, projectName, nextEnvironmentSetRepoName, *pullRequest.PullRequestId)

	return nil
}
