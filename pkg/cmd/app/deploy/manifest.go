package deploy

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/frontierdigital/ranger/pkg/util/config"
	"github.com/frontierdigital/ranger/pkg/util/deploy"
	"github.com/frontierdigital/ranger/pkg/util/manifest"
	"github.com/frontierdigital/ranger/pkg/util/workload"
	"github.com/frontierdigital/utils/azuredevops"
	git "github.com/frontierdigital/utils/git/external_git"
	"github.com/frontierdigital/utils/output"

	"github.com/otiai10/copy"
	"github.com/segmentio/ksuid"
)

const (
	WaitForBuildAttempts uint = 240
	WaitForBuildInterval int  = 15
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

	var hasErrors bool
	for _, workload := range manifest.Workloads {
		result := DeployWorkload(*azureDevOps, config, organisationName, projectName, manifest.Environment, manifest.Set, workload)

		if result.Error != nil {
			output.PrintfError(result.Error.Error())
			hasErrors = true
		}
		result.PrintResult()
	}

	if hasErrors {
		return fmt.Errorf("one or more errors occurred during manifest deploy")
	}

	return nil
}

func DeployWorkload(azureDevOps azuredevops.AzureDevOps, config *config.Config, organisationName string, projectName string, environment string, set string, workload *workload.Workload) (result *deploy.DeployWorkloadResult) {
	result = &deploy.DeployWorkloadResult{
		Workload: workload,
	}

	workload.PrintHeader()

	typeProjectName, typeRepositoryName := workload.GetTypeProjectAndRepositoryNames()

	pipelineName := fmt.Sprintf("%s (deploy)", typeRepositoryName)
	buildDefinition, err := azureDevOps.GetBuildDefinitionByName(typeProjectName, pipelineName)
	if err != nil {
		result.Error = err
		return
	}

	output.PrintlnfInfo("Found deploy pipeline definition with Id '%d' for workload type '%s' (https://dev.azure.com/%s/%s/_build?definitionId=%d)",
		*buildDefinition.Id, workload.Type, organisationName, projectName, *buildDefinition.Id)

	repository, err := azureDevOps.GetRepositoryByName(typeProjectName, typeRepositoryName)
	if err != nil {
		result.Error = err
		return
	}

	output.PrintlnfInfo("Found repository with Id '%s' for workload type '%s' (%s)", repository.Id, workload.Type, *repository.WebUrl)

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

	configRepoName := "generated-manifest-config"
	var configRef string

	if workloadConfigExists || workloadSecretsExists {
		configRepoUrl := fmt.Sprintf("https://dev.azure.com/%s/%s/_git/%s", organisationName, projectName, configRepoName)

		configRepoPath, err := os.MkdirTemp("", "")
		if err != nil {
			result.Error = err
			return
		}
		configRepo := git.NewGit(configRepoPath)
		err = configRepo.CloneOverHttp(configRepoUrl, config.ADO.PAT, "x-oauth-basic")
		if err != nil {
			result.Error = err
			return
		}
		err = configRepo.SetConfig("user.email", config.Git.UserEmail)
		if err != nil {
			result.Error = err
			return
		}
		err = configRepo.SetConfig("user.name", config.Git.UserName)
		if err != nil {
			result.Error = err
			return
		}

		configBranchName := fmt.Sprintf("%s_%s_%s_%s", workload.Name, environment, set, ksuid.New().String())
		err = configRepo.Checkout(configBranchName, true)
		if err != nil {
			result.Error = err
			return
		}

		if workloadConfigExists {
			err = copy.Copy(workloadConfigPath, path.Join(configRepoPath, ".config"))
			if err != nil {
				result.Error = err
				return
			}
		}

		if workloadSecretsExists {
			err = copy.Copy(workloadSecretsPath, path.Join(configRepoPath, ".secrets"))
			if err != nil {
				result.Error = err
				return
			}
		}

		err = configRepo.AddAll()
		if err != nil {
			result.Error = err
			return
		}

		commitMessage := fmt.Sprintf("Generate config for workload '%s', environment '%s' and set '%s'", workload.Name, environment, set)
		commitSha, err := configRepo.Commit(commitMessage)
		if err != nil {
			result.Error = err
			return
		}
		_ = commitSha

		err = configRepo.Push(false)
		if err != nil {
			result.Error = err
			return
		}

		output.PrintlnfInfo("Pushed config (https://dev.azure.com/%s/%s/_git/%s?version=GB%s)", organisationName, projectName, configRepoName, configBranchName)

		configRef = fmt.Sprintf("%s/%s@%s", projectName, configRepoName, commitSha)

		defer os.RemoveAll(configRepoPath)
	}

	tags := []string{environment, set}

	templateParameters := map[string]string{
		"environment": environment,
		"name":        workload.Name,
		"set":         set,
		"version":     workload.Version,
	}
	if configRef != "" {
		templateParameters["configRef"] = configRef
	}

	sourceBranchName := fmt.Sprintf("refs/tags/%s", workload.Version)

	build, err := azureDevOps.QueueBuild(projectName, buildDefinition.Id, sourceBranchName, templateParameters, tags)
	if err != nil {
		result.Error = err
		return
	}
	buildLinks := build.Links.(map[string]interface{})
	buildWebLinks := buildLinks["web"].(map[string]interface{})
	buildWebLink := buildWebLinks["href"].(string)

	result.Link = buildWebLink
	result.QueueTime = &build.QueueTime.Time

	output.PrintlnfInfo("Queued build '%s' (%s)", *build.BuildNumber, buildWebLink)

	time.Sleep(time.Duration(WaitForBuildInterval) * time.Second)

	output.PrintlnfInfo("Waiting for build '%s' (%s)", *build.BuildNumber, buildWebLink)

	build, err = azureDevOps.WaitForBuild(projectName, build.Id, WaitForBuildAttempts, WaitForBuildInterval)

	if build.FinishTime != nil {
		result.FinishTime = &build.FinishTime.Time
	}

	if err != nil {
		result.Error = err
		output.PrintlnfInfo("Build '%s' has not completed (%s)", *build.BuildNumber, buildWebLink)
		return
	}

	output.PrintlnfInfo("Build '%s' has completed with result '%s' (%s)", *build.BuildNumber, *build.Result, buildWebLink)

	if *build.Result != "succeeded" {
		result.Error = fmt.Errorf("build '%d' has result '%s'", *build.Id, *build.Result)
	}

	return
}
