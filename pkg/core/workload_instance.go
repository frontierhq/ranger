package core

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/gofrontier-com/go-utils/azuredevops"
	git "github.com/gofrontier-com/go-utils/git/external_git"
	"github.com/gofrontier-com/go-utils/output"
	"github.com/otiai10/copy"
	"github.com/segmentio/ksuid"
	"gopkg.in/yaml.v2"
)

const (
	WaitForBuildAttempts uint = 240
	WaitForBuildInterval int  = 15
)

func (w *WorkloadInstance) Deploy(azureDevOps azuredevops.AzureDevOps, config *Config, projectName string, organisationName string, environment string, set string, workloadConfigPath string, workloadSecretsPath string) (result *WorkloadResult) {
	defer func() {
		if e, ok := recover().(error); ok {
			result.Status = WorkloadResultStatusValuesType.Failed
			output.PrintfError(e.Error())
		}
	}()

	result = &WorkloadResult{
		Workload: w,
	}

	typeProjectName, typeRepositoryName := w.GetTypeProjectAndRepositoryNames()

	repository, err := azureDevOps.GetRepository(typeProjectName, typeRepositoryName)
	if err != nil {
		panic(err)
	}

	output.PrintlnfInfo("Found repository with Id '%s' for workload type '%s' (%s)", repository.Id, w.Type, *repository.WebUrl)

	pipelineName := fmt.Sprintf("%s~%s (deploy)", typeProjectName, typeRepositoryName)
	buildDefinition, err := azureDevOps.GetOrCreateBuildDefinition(projectName, pipelineName, repository.Id.String(), "Ranger/Workloads", "azure-pipelines.deploy.yml")
	if err != nil {
		panic(err)
	}

	output.PrintlnfInfo("Found or created pipeline definition with Id '%d' for workload type '%s' (https://dev.azure.com/%s/%s/_build?definitionId=%d)",
		*buildDefinition.Id, w.Type, organisationName, projectName, *buildDefinition.Id)

	configRef, err := w.PrepareConfig(config, projectName, organisationName, environment, workloadConfigPath, workloadSecretsPath)
	if err != nil {
		panic(err)
	}

	tags := []string{environment, set}

	templateParameters := map[string]string{
		"environment":  environment,
		"set":          set,
		"workloadName": w.Name,
	}
	if configRef != nil {
		templateParameters["configRef"] = *configRef
	}
	for _, v := range w.ExtraParameters {
		templateParameters[v.Name] = v.Value
	}

	sourceBranchName := fmt.Sprintf("refs/tags/%s", w.Version)

	build, err := azureDevOps.QueueBuild(projectName, *buildDefinition.Id, sourceBranchName, templateParameters, tags)
	if err != nil {
		panic(err)
	}
	buildLinks := build.Links.(map[string]interface{})
	buildWebLinks := buildLinks["web"].(map[string]interface{})
	buildWebLink := buildWebLinks["href"].(string)

	result.Link = buildWebLink
	result.QueueTime = &build.QueueTime.Time

	output.PrintlnfInfo("Queued build '%s' (%s)", *build.BuildNumber, buildWebLink)

	time.Sleep(time.Duration(WaitForBuildInterval) * time.Second)

	output.PrintlnfInfo("Waiting for build '%s' (%s)", *build.BuildNumber, buildWebLink)

	err = azureDevOps.WaitForBuild(projectName, *build.Id, WaitForBuildAttempts, WaitForBuildInterval)
	if err != nil {
		panic(err)
	}

	build, err = azureDevOps.GetBuild(projectName, *build.Id)
	if err != nil {
		panic(err)
	}

	if build.FinishTime != nil {
		result.FinishTime = &build.FinishTime.Time
	}

	if err != nil {
		output.PrintlnfInfo("Build '%s' has not completed (%s)", *build.BuildNumber, buildWebLink)
		panic(err)
	}

	output.PrintlnfInfo("Build '%s' has completed with result '%s' (%s)", *build.BuildNumber, *build.Result, buildWebLink)

	if *build.Result != "succeeded" {
		panic(fmt.Errorf("build '%d' has result '%s'", *build.Id, *build.Result))
	}

	result.Status = WorkloadResultStatusValuesType.Succeeded

	return
}

func (w *WorkloadInstance) Destroy(azureDevOps azuredevops.AzureDevOps, config *Config, projectName string, organisationName string, environment string, set string, workloadConfigPath string, workloadSecretsPath string) (result *WorkloadResult) {
	defer func() {
		if e, ok := recover().(error); ok {
			if _, ok := e.(*WorkloadDestroyPreventedError); ok {
				result.Status = WorkloadResultStatusValuesType.Skipped
			} else {
				result.Status = WorkloadResultStatusValuesType.Failed
			}
			output.PrintfError(e.Error())
		}
	}()

	result = &WorkloadResult{
		Workload: w,
	}

	if w.PreventDestroy {
		now := time.Now()
		result.Link = "N/A"
		result.QueueTime = &now
		result.FinishTime = &now
		panic(&WorkloadDestroyPreventedError{
			Workload: w,
		})
	}

	typeProjectName, typeRepositoryName := w.GetTypeProjectAndRepositoryNames()

	repository, err := azureDevOps.GetRepository(typeProjectName, typeRepositoryName)
	if err != nil {
		panic(err)
	}

	output.PrintlnfInfo("Found repository with Id '%s' for workload type '%s' (%s)", repository.Id, w.Type, *repository.WebUrl)

	pipelineName := fmt.Sprintf("%s~%s (destroy)", typeProjectName, typeRepositoryName)
	buildDefinition, err := azureDevOps.GetOrCreateBuildDefinition(projectName, pipelineName, repository.Id.String(), "Ranger/Workloads", "azure-pipelines.destroy.yml")
	if err != nil {
		panic(err)
	}

	output.PrintlnfInfo("Found or created pipeline definition with Id '%d' for workload type '%s' (https://dev.azure.com/%s/%s/_build?definitionId=%d)",
		*buildDefinition.Id, w.Type, organisationName, projectName, *buildDefinition.Id)

	configRef, err := w.PrepareConfig(config, projectName, organisationName, environment, workloadConfigPath, workloadSecretsPath)
	if err != nil {
		panic(err)
	}

	tags := []string{environment, set}

	templateParameters := map[string]string{
		"environment":  environment,
		"set":          set,
		"workloadName": w.Name,
	}
	if configRef != nil {
		templateParameters["configRef"] = *configRef
	}
	for _, v := range w.ExtraParameters {
		templateParameters[v.Name] = v.Value
	}

	sourceBranchName := fmt.Sprintf("refs/tags/%s", w.Version)

	build, err := azureDevOps.QueueBuild(projectName, *buildDefinition.Id, sourceBranchName, templateParameters, tags)
	if err != nil {
		panic(err)
	}
	buildLinks := build.Links.(map[string]interface{})
	buildWebLinks := buildLinks["web"].(map[string]interface{})
	buildWebLink := buildWebLinks["href"].(string)

	result.Link = buildWebLink
	result.QueueTime = &build.QueueTime.Time

	output.PrintlnfInfo("Queued build '%s' (%s)", *build.BuildNumber, buildWebLink)

	time.Sleep(time.Duration(WaitForBuildInterval) * time.Second)

	output.PrintlnfInfo("Waiting for build '%s' (%s)", *build.BuildNumber, buildWebLink)

	err = azureDevOps.WaitForBuild(projectName, *build.Id, WaitForBuildAttempts, WaitForBuildInterval)
	if err != nil {
		panic(err)
	}

	build, err = azureDevOps.GetBuild(projectName, *build.Id)
	if err != nil {
		panic(err)
	}

	if build.FinishTime != nil {
		result.FinishTime = &build.FinishTime.Time
	}

	if err != nil {
		output.PrintlnfInfo("Build '%s' has not completed (%s)", *build.BuildNumber, buildWebLink)
		panic(err)
	}

	output.PrintlnfInfo("Build '%s' has completed with result '%s' (%s)", *build.BuildNumber, *build.Result, buildWebLink)

	if *build.Result != "succeeded" {
		panic(fmt.Errorf("build '%d' has result '%s'", *build.Id, *build.Result))
	}

	result.Status = WorkloadResultStatusValuesType.Succeeded

	return
}

func (w *WorkloadInstance) GetTypeProjectAndRepositoryNames() (string, string) {
	typeParts := strings.Split(w.Type, "/")
	return typeParts[0], typeParts[1]
}

func (w *WorkloadInstance) PrintHeader() {
	builder := &strings.Builder{}
	builder.WriteString(fmt.Sprintf("%s\n", strings.Repeat("=", 78)))
	builder.WriteString(fmt.Sprintf("Name     | %s\n", w.Name))
	builder.WriteString(fmt.Sprintf("Type     | %s\n", w.Type))
	builder.WriteString(fmt.Sprintf("Version  | %s\n", w.Version))
	builder.WriteString(strings.Repeat("-", 78))
	output.Println(builder.String())
}

func LoadWorkloadInstance(filePath string) (WorkloadInstance, error) {
	workloadInstance := WorkloadInstance{}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return workloadInstance, err
	}

	err = yaml.Unmarshal(data, &workloadInstance)
	if err != nil {
		return workloadInstance, err
	}

	workloadInstance.FilePath = filePath

	return workloadInstance, nil
}

func (w *WorkloadInstance) PrepareConfig(config *Config, projectName string, organisationName string, environment string, workloadConfigPath string, workloadSecretsPath string) (*string, error) {
	var workloadConfigExists bool
	_, err := os.Stat(workloadConfigPath)
	if err != nil {
		workloadConfigExists = !os.IsNotExist(err)
	} else {
		workloadConfigExists = true
	}

	var workloadSecretsExists bool
	_, err = os.Stat(workloadSecretsPath)
	if err != nil {
		workloadSecretsExists = !os.IsNotExist(err)
	} else {
		workloadSecretsExists = true
	}

	configRepoName := "generated-set-config"

	if workloadConfigExists || workloadSecretsExists {
		configRepoUrl := fmt.Sprintf("https://dev.azure.com/%s/%s/_git/%s", organisationName, projectName, configRepoName)

		configRepo, err := git.NewClonedGit(configRepoUrl, "x-oauth-basic", config.ADO.PAT, config.Git.UserEmail, config.Git.UserName)
		if err != nil {
			return nil, err
		}
		defer os.RemoveAll(configRepo.GetRepositoryPath())

		configBranchName := fmt.Sprintf("%s_%s_%s", w.Name, environment, ksuid.New().String())
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

		commitMessage := fmt.Sprintf("Generate config for workload '%s' in environment '%s'", w.Name, environment)
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
