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

	output.PrintfInfo("Action: Destroy from %s (in reverse order)\n\n", manifest.Environment)

	var hasErrors bool
	for i := len(manifest.Workloads) - 1; i >= 0; i-- {
		workloadInstance := manifest.Workloads[i]

		workloadInstance.PrintHeader()

		result := DestroyWorkload(*azureDevOps, config, projectName, organisationName, manifest.Environment, manifest.Set, workloadInstance)
		if result.Status == core.WorkloadResultStatusValuesType.Failed {
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
	defer func() {
		if e, ok := recover().(error); ok {
			if _, ok := e.(*core.WorkloadDestroyPreventedError); ok {
				result.Status = core.WorkloadResultStatusValuesType.Skipped
			} else {
				result.Status = core.WorkloadResultStatusValuesType.Failed
			}
			output.PrintfError(e.Error())
		}
	}()

	result = &core.WorkloadResult{
		Workload: workloadInstance,
	}

	if workloadInstance.PreventDestroy {
		now := time.Now()
		result.Link = "N/A"
		result.QueueTime = &now
		result.FinishTime = &now
		panic(&core.WorkloadDestroyPreventedError{
			Workload: workloadInstance,
		})
	}

	typeProjectName, typeRepositoryName := workloadInstance.GetTypeProjectAndRepositoryNames()

	repository, err := azureDevOps.GetRepository(typeProjectName, typeRepositoryName)
	if err != nil {
		panic(err)
	}

	output.PrintlnfInfo("Found repository with Id '%s' for workload type '%s' (%s)", repository.Id, workloadInstance.Type, *repository.WebUrl)

	var buildDefinitionId *int
	pipelineName := fmt.Sprintf("%s~%s (destroy)", typeProjectName, typeRepositoryName)
	buildDefinitionRef, err := azureDevOps.GetBuildDefinitionByName(projectName, pipelineName)
	if err != nil {
		if _, ok := err.(*azuredevops.BuildNotFoundError); ok {
			buildDefinition, err := azureDevOps.CreateBuildDefinition(projectName, repository.Id.String(), "Ranger/Workloads", pipelineName, "azure-pipelines.destroy.yml")
			if err != nil {
				panic(err)
			}
			buildDefinitionId = buildDefinition.Id
			output.PrintlnfInfo("Created destroy pipeline definition with Id '%d' for workload type '%s' (https://dev.azure.com/%s/%s/_build?definitionId=%d)",
				*buildDefinitionId, workloadInstance.Type, organisationName, projectName, *buildDefinitionId)
		} else {
			panic(err)
		}
	} else {
		buildDefinitionId = buildDefinitionRef.Id
		output.PrintlnfInfo("Found destroy pipeline definition with Id '%d' for workload type '%s' (https://dev.azure.com/%s/%s/_build?definitionId=%d)",
			*buildDefinitionId, workloadInstance.Type, organisationName, projectName, *buildDefinitionId)
	}

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

	build, err := azureDevOps.QueueBuild(projectName, *buildDefinitionId, sourceBranchName, templateParameters, tags)
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
