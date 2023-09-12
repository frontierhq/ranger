package promote

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/frontierdigital/ranger/pkg/util/config"
	"github.com/frontierdigital/ranger/pkg/util/manifest"
	"github.com/frontierdigital/utils/azuredevops"
	git "github.com/frontierdigital/utils/git/external_git"
	"github.com/frontierdigital/utils/output"
)

func PromoteSet(config *config.Config, projectName string, organisationName string, nextEnvironment string) error {
	azureDevOps := azuredevops.NewAzureDevOps(organisationName, config.ADO.PAT)

	sourceManifestFilepath, _ := filepath.Abs("./manifest.yml")
	sourceManifest, err := manifest.LoadManifest(sourceManifestFilepath)
	if err != nil {
		return err
	}

	sourceManifest.PrintHeader()

	sourceManifest.PrintWorkloadsSummary()

	nextEnvironmentSetRepoName := fmt.Sprintf("%s-%s-set", nextEnvironment, sourceManifest.Set)
	nextEnvironmentSetRepoUrl := fmt.Sprintf("https://dev.azure.com/%s/%s/_git/%s", organisationName, projectName, nextEnvironmentSetRepoName)

	nextEnvironmentSetRepoPath, err := os.MkdirTemp("", "")
	if err != nil {
		return err
	}
	nextEnvironmentSetRepo := git.NewGit(nextEnvironmentSetRepoPath)
	err = nextEnvironmentSetRepo.CloneOverHttp(nextEnvironmentSetRepoUrl, config.ADO.PAT, "x-oauth-basic")
	if err != nil {
		return err
	}
	err = nextEnvironmentSetRepo.SetConfig("user.email", config.Git.UserEmail)
	if err != nil {
		return err
	}
	err = nextEnvironmentSetRepo.SetConfig("user.name", config.Git.UserName)
	if err != nil {
		return err
	}

	output.PrintfInfo("Cloned target environment set repository '%s' (https://dev.azure.com/%s/%s/_git/%s)", nextEnvironmentSetRepoName, organisationName, projectName, nextEnvironmentSetRepoName)

	promoteBranchName := fmt.Sprintf("ranger/promote/%s", sourceManifest.Environment)
	err = nextEnvironmentSetRepo.Checkout(promoteBranchName, true)
	if err != nil {
		return err
	}

	targetManifestFilePath := nextEnvironmentSetRepo.GetFilePath("manifest.yml")
	targetManifest, err := manifest.LoadManifest(targetManifestFilePath)
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
