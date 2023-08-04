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

func PromoteSet(config *config.Config, projectName string, organisationName string, targetEnvironment string) error {
	azureDevOps := azuredevops.NewAzureDevOps(organisationName, config.ADO.PAT)

	sourceManifestFilepath, _ := filepath.Abs("./manifest.yml")
	sourceManifest, err := manifest.LoadManifest(sourceManifestFilepath)
	if err != nil {
		return err
	}

	sourceManifest.PrintHeader()

	sourceManifest.PrintWorkloadsSummary()

	targetEnvironmentSetRepoName := fmt.Sprintf("%s-%s-set", targetEnvironment, sourceManifest.Set)
	targetEnvironmentSetRepoUrl := fmt.Sprintf("https://dev.azure.com/%s/%s/_git/%s", organisationName, projectName, targetEnvironmentSetRepoName)

	targetEnvironmentSetRepoPath, err := os.MkdirTemp("", "")
	if err != nil {
		return err
	}
	targetEnvironmentSetRepo := git.NewGit(targetEnvironmentSetRepoPath)
	err = targetEnvironmentSetRepo.CloneOverHttp(targetEnvironmentSetRepoUrl, config.ADO.PAT, "x-oauth-basic")
	if err != nil {
		return err
	}
	err = targetEnvironmentSetRepo.SetConfig("user.email", config.Git.UserEmail)
	if err != nil {
		return err
	}
	err = targetEnvironmentSetRepo.SetConfig("user.name", config.Git.UserName)
	if err != nil {
		return err
	}

	output.PrintfInfo("Cloned target environment set repository '%s' (https://dev.azure.com/%s/%s/_git/%s)", targetEnvironmentSetRepoName, organisationName, projectName, targetEnvironmentSetRepoName)

	promoteBranchName := fmt.Sprintf("ranger/promote/%s", sourceManifest.Environment)
	err = targetEnvironmentSetRepo.Checkout(promoteBranchName, true)
	if err != nil {
		return err
	}

	targetManifestFilePath := targetEnvironmentSetRepo.GetFilePath("manifest.yml")
	targetManifest, err := manifest.LoadManifest(targetManifestFilePath)
	if err != nil {
		return err
	}

	targetManifest.Version = sourceManifest.Version
	targetManifest.Workloads = sourceManifest.Workloads

	targetManifest.Save()

	err = targetEnvironmentSetRepo.AddAll()
	if err != nil {
		return err
	}

	commitMessage := fmt.Sprintf("Promote set version %d from %s", sourceManifest.Version, sourceManifest.Environment)
	_, err = targetEnvironmentSetRepo.Commit(commitMessage)
	if err != nil {
		return err
	}

	err = targetEnvironmentSetRepo.Push(true)
	if err != nil {
		return err
	}

	output.PrintlnfInfo("Pushed branch '%s' (https://dev.azure.com/%s/%s/_git/%s?version=GB%s)", promoteBranchName, organisationName, projectName, targetEnvironmentSetRepoName, promoteBranchName)

	existingPullRequest, err := azureDevOps.FindPullRequest(projectName, targetEnvironmentSetRepoName, fmt.Sprintf("refs/heads/%s", promoteBranchName), "refs/heads/main")
	if err != nil {
		return err
	}

	if existingPullRequest != nil {
		_, err = azureDevOps.AbandonPullRequest(projectName, targetEnvironmentSetRepoName, *existingPullRequest.PullRequestId)
		if err != nil {
			return err
		}

		output.PrintlnfInfo("Abandoned existing pull request with Id '%d' (https://dev.azure.com/%s/%s/_git/%s/pullrequest/%d)", *existingPullRequest.PullRequestId, organisationName, projectName, targetEnvironmentSetRepoName, *existingPullRequest.PullRequestId)
	}

	pullRequestTitle := commitMessage
	pullRequest, err := azureDevOps.CreatePullRequest(projectName, targetEnvironmentSetRepoName, fmt.Sprintf("refs/heads/%s", promoteBranchName), "refs/heads/main", pullRequestTitle)
	if err != nil {
		return err
	}

	output.PrintfInfo("Created pull request with Id '%d' (https://dev.azure.com/%s/%s/_git/%s/pullrequest/%d)", *pullRequest.PullRequestId, organisationName, projectName, targetEnvironmentSetRepoName, *pullRequest.PullRequestId)

	return nil
}
