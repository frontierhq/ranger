package workload

import (
	"fmt"
	"os"
	"path"

	"github.com/frontierdigital/ranger/pkg/core"

	git "github.com/frontierdigital/utils/git/external_git"
	"github.com/frontierdigital/utils/output"
	"github.com/otiai10/copy"
	"github.com/segmentio/ksuid"
)

func PrepareConfig(config *core.Config, projectName string, organisationName string, environment string, set string, workload *core.WorkloadInstance) (*string, error) {
	workloadConfigPath := path.Join("config", "workloads", workload.Name)
	workloadConfigExists := true
	_, err := os.Stat(workloadConfigPath)
	if err != nil {
		workloadConfigExists = !os.IsNotExist(err)
	}

	workloadSecretsPath := path.Join("secrets", "workloads", workload.Name)
	workloadSecretsExists := true
	_, err = os.Stat(workloadSecretsPath)
	if err != nil {
		workloadSecretsExists = !os.IsNotExist(err)
	}

	configRepoName := "generated-set-config"

	if workloadConfigExists || workloadSecretsExists {
		configRepoUrl := fmt.Sprintf("https://dev.azure.com/%s/%s/_git/%s", organisationName, projectName, configRepoName)

		configRepo, err := git.NewClonedGit(configRepoUrl, "x-oauth-basic", config.ADO.PAT, config.Git.UserEmail, config.Git.UserName)
		if err != nil {
			return nil, err
		}
		defer os.RemoveAll(configRepo.GetRepositoryPath())

		configBranchName := fmt.Sprintf("%s_%s_%s_%s", workload.Name, environment, set, ksuid.New().String())
		err = configRepo.Checkout(configBranchName, true)
		if err != nil {
			return nil, err
		}

		if workloadConfigExists {
			err = copy.Copy(workloadConfigPath, path.Join(configRepo.GetRepositoryPath(), ".config"))
			if err != nil {
				return nil, err
			}
		}

		if workloadSecretsExists {
			err = copy.Copy(workloadSecretsPath, path.Join(configRepo.GetRepositoryPath(), ".secrets"))
			if err != nil {
				return nil, err
			}
		}

		err = configRepo.AddAll()
		if err != nil {
			return nil, err
		}

		commitMessage := fmt.Sprintf("Generate config for workload '%s', environment '%s' and set '%s'", workload.Name, environment, set)
		commitSha, err := configRepo.Commit(commitMessage)
		if err != nil {
			return nil, err
		}

		err = configRepo.Push(false)
		if err != nil {
			return nil, err
		}

		output.PrintlnfInfo("Pushed config (https://dev.azure.com/%s/%s/_git/%s?version=GB%s)", organisationName, projectName, configRepoName, configBranchName)

		configRef := fmt.Sprintf("%s/%s@%s", projectName, configRepoName, commitSha)

		return &configRef, nil
	}

	return nil, nil
}
