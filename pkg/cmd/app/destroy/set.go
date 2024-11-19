package destroy

import (
	"fmt"
	"path"
	"path/filepath"

	"github.com/gofrontier-com/go-utils/azuredevops"
	"github.com/gofrontier-com/go-utils/output"
	"github.com/gofrontier-com/ranger/pkg/core"
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

		workloadConfigPath := path.Join("config", "workloads", workloadInstance.Name)
		workloadSecretsPath := path.Join("secrets", "workloads", workloadInstance.Name)

		workloadInstance.PrintHeader()

		result := workloadInstance.Destroy(*azureDevOps, config, projectName, organisationName, manifest.Environment, manifest.Set, workloadConfigPath, workloadSecretsPath)
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
