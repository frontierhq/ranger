package deploy

import (
	"fmt"
	"path"
	"path/filepath"

	"github.com/gofrontier-com/go-utils/azuredevops"
	"github.com/gofrontier-com/go-utils/output"
	"github.com/gofrontier-com/ranger/pkg/core"
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
		workloadConfigPath := path.Join("config", "workloads", workloadInstance.Name)
		workloadSecretsPath := path.Join("secrets", "workloads", workloadInstance.Name)

		workloadInstance.PrintHeader()

		result := workloadInstance.Deploy(*azureDevOps, config, projectName, organisationName, manifest.Environment, manifest.Set, workloadConfigPath, workloadSecretsPath)
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
