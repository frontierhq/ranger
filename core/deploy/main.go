package deploy

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/frontierdigital/ranger/core"
	"github.com/frontierdigital/ranger/core/print"
	"github.com/frontierdigital/ranger/core/util"
	"github.com/microsoft/azure-devops-go-api/azuredevops"
	azuredevopscore "github.com/microsoft/azure-devops-go-api/azuredevops/core"
)

func DeployManifest(configuration *core.Configuration, projectName string, organisationName string) error {
	organisationUrl := fmt.Sprintf("https://dev.azure.com/%s", organisationName)
	connection := azuredevops.NewPatConnection(organisationUrl, configuration.ADO.PAT)

	ctx := context.Background()

	coreClient, err := azuredevopscore.NewClient(ctx, connection)
	if err != nil {
		return err
	}
	_ = coreClient

	// responseValue, err := coreClient.GetProjects(ctx, core.GetProjectsArgs{})
	// if err != nil {
	// 	return err
	// }

	manifestFilepath, _ := filepath.Abs("./manifest.yml")
	manifest, err := util.LoadManifest(manifestFilepath)
	if err != nil {
		return err
	}

	manifestName := fmt.Sprintf("%s-%s", manifest.Environment, manifest.Layer)

	print.PrintManifestHeader(manifestName, manifest.Layer, manifest.Environment, manifest.Version)

	return nil
}
