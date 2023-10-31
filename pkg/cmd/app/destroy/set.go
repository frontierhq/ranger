package destroy

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/frontierdigital/ranger/pkg/core"
	"github.com/frontierdigital/ranger/pkg/util/workload"
	"github.com/frontierdigital/utils/azuredevops"
	"github.com/frontierdigital/utils/output"
)

const (
	WaitForBuildAttempts uint = 240
	WaitForBuildInterval int  = 15
)

func DestroySet(config *core.Config, projectName string, organisationName string) error {
	azureDevOps := azuredevops.NewAzureDevOps(organisationName, config.ADO.PAT)

	manifestFilepath, _ := filepath.Abs("./manifest.yml")
	manifest, err := core.LoadManifest(manifestFilepath)
	if err != nil {
		return err
	}

	manifest.PrintHeader()

	manifest.PrintWorkloadsSummary()

	output.PrintfInfo("Action: Destroy (in reverse order)\n\n")

	var hasErrors bool
	for i := len(manifest.Workloads) - 1; i >= 0; i-- {
		workloadInstance := manifest.Workloads[i]

		workloadInstance.PrintHeader()

		result := DestroyWorkload(*azureDevOps, config, projectName, organisationName, manifest.Environment, manifest.Set, workloadInstance)
		if result.Error != nil {
			output.PrintfError(result.Error.Error())
			hasErrors = true
		}

		result.PrintResult()
	}

	if hasErrors {
		return fmt.Errorf("one or more errors occurred during set destroy")
	}

	return nil
}

func DestroyWorkload(azureDevOps azuredevops.AzureDevOps, config *core.Config, projectName string, organisationName string, environment string, set string, workloadInstance *core.WorkloadInstance) (result *core.WorkloadResult) {
	result = &core.WorkloadResult{
		Workload: workloadInstance,
	}

	if workloadInstance.PreventDestroy {
		output.PrintlnfInfo("Instance configuration prevents destroy; skipping")
		now := time.Now()
		result.Link = "N/A"
		result.QueueTime = &now
		result.FinishTime = &now
		return
	}

	typeProjectName, typeRepositoryName := workloadInstance.GetTypeProjectAndRepositoryNames()

	pipelineName := fmt.Sprintf("%s (destroy)", typeRepositoryName)
	buildDefinition, err := azureDevOps.GetBuildDefinitionByName(typeProjectName, pipelineName)
	if err != nil {
		result.Error = err
		return
	}

	output.PrintlnfInfo("Found destroy pipeline definition with Id '%d' for workload type '%s' (https://dev.azure.com/%s/%s/_build?definitionId=%d)",
		*buildDefinition.Id, workloadInstance.Type, organisationName, projectName, *buildDefinition.Id)

	repository, err := azureDevOps.GetRepository(typeProjectName, typeRepositoryName)
	if err != nil {
		result.Error = err
		return
	}

	output.PrintlnfInfo("Found repository with Id '%s' for workload type '%s' (%s)", repository.Id, workloadInstance.Type, *repository.WebUrl)

	configRef, err := workload.PrepareConfig(config, projectName, organisationName, environment, set, workloadInstance)
	if err != nil {
		result.Error = err
		return
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

	build, err := azureDevOps.QueueBuild(typeProjectName, buildDefinition.Id, sourceBranchName, templateParameters, tags)
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

	err = azureDevOps.WaitForBuild(typeProjectName, build.Id, WaitForBuildAttempts, WaitForBuildInterval)
	if err != nil {
		result.Error = err
		return
	}

	build, err = azureDevOps.GetBuild(typeProjectName, build.Id)
	if err != nil {
		result.Error = err
		return
	}

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
