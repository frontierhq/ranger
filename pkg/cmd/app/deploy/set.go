package deploy

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/gofrontier-com/go-utils/azuredevops"
	"github.com/gofrontier-com/go-utils/output"
	"github.com/gofrontier-com/ranger/pkg/core"
	"github.com/gofrontier-com/ranger/pkg/util/workload"
)

const (
	WaitForBuildAttempts uint = 240
	WaitForBuildInterval int  = 15
)

func DeploySet(config *core.Config, projectName string, organisationName string) error {
	azureDevOps := azuredevops.NewAzureDevOps(organisationName, config.ADO.PAT)

	manifestFilepath, _ := filepath.Abs("./manifest.yml")
	manifest, err := core.LoadManifest(manifestFilepath)
	if err != nil {
		return err
	}

	manifest.PrintHeader()

	manifest.PrintWorkloadsSummary()

	output.PrintfInfo("Action: Deploy to %s\n\n", manifest.Environment)

	var hasErrors bool
	for _, workloadInstance := range manifest.Workloads {
		workloadInstance.PrintHeader()

		result := DeployWorkload(*azureDevOps, config, projectName, organisationName, manifest.Environment, manifest.Set, workloadInstance)
		if result.Status == core.WorkloadResultStatusValuesType.Failed {
			hasErrors = true
		}

		result.PrintResult()
	}

	if hasErrors {
		return fmt.Errorf("one or more errors occurred during set deploy")
	}

	return nil
}

func DeployWorkload(azureDevOps azuredevops.AzureDevOps, config *core.Config, projectName string, organisationName string, environment string, set string, workloadInstance *core.WorkloadInstance) (result *core.WorkloadResult) {
	defer func() {
		if e, ok := recover().(error); ok {
			result.Status = core.WorkloadResultStatusValuesType.Failed
			output.PrintfError(e.Error())
		}
	}()

	result = &core.WorkloadResult{
		Workload: workloadInstance,
	}

	typeProjectName, typeRepositoryName := workloadInstance.GetTypeProjectAndRepositoryNames()

	repository, err := azureDevOps.GetRepository(typeProjectName, typeRepositoryName)
	if err != nil {
		panic(err)
	}

	output.PrintlnfInfo("Found repository with Id '%s' for workload type '%s' (%s)", repository.Id, workloadInstance.Type, *repository.WebUrl)

	pipelineName := fmt.Sprintf("%s~%s (deploy)", typeProjectName, typeRepositoryName)
	buildDefinition, err := azureDevOps.GetOrCreateBuildDefinition(projectName, pipelineName, repository.Id.String(), "Ranger/Workloads", "azure-pipelines.deploy.yml")
	if err != nil {
		panic(err)
	}

	output.PrintlnfInfo("Found or created pipeline definition with Id '%d' for workload type '%s' (https://dev.azure.com/%s/%s/_build?definitionId=%d)",
		*buildDefinition.Id, workloadInstance.Type, organisationName, projectName, *buildDefinition.Id)

	configRef, err := workload.PrepareConfig(config, projectName, organisationName, environment, set, workloadInstance)
	if err != nil {
		panic(err)
	}

	tags := []string{environment, set}

	templateParameters := map[string]string{
		"environment":  environment,
		"set":          set,
		"workloadName": workloadInstance.Name,
	}
	if configRef != nil {
		templateParameters["configRef"] = *configRef
	}
	for _, v := range workloadInstance.ExtraParameters {
		templateParameters[v.Name] = v.Value
	}

	sourceBranchName := fmt.Sprintf("refs/tags/%s", workloadInstance.Version)

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

	result.Status = core.WorkloadResultStatusValuesType.Succeeded

	return
}
