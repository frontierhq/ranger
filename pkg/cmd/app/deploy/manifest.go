package deploy

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/frontierdigital/ranger/pkg/cmd/app"
	"github.com/frontierdigital/ranger/pkg/util/manifest"
	"github.com/microsoft/azure-devops-go-api/azuredevops"
	"github.com/microsoft/azure-devops-go-api/azuredevops/core"
)

func DeployManifest(config *app.Config, projectName string, organisationName string) error {
	organisationUrl := fmt.Sprintf("https://dev.azure.com/%s", organisationName)
	connection := azuredevops.NewPatConnection(organisationUrl, config.ADO.PAT)

	ctx := context.Background()

	coreClient, err := core.NewClient(ctx, connection)
	if err != nil {
		return err
	}
	_ = coreClient

	manifestFilepath, _ := filepath.Abs("./manifest.yml")
	manifest, err := manifest.LoadManifest(manifestFilepath)
	if err != nil {
		return err
	}

	manifestName := fmt.Sprintf("%s-%s", manifest.Environment, manifest.Layer)

	manifest.PrintHeader(manifestName, manifest.Layer, manifest.Environment, manifest.Version)

	// responseValue, err := coreClient.GetProjects(ctx, core.GetProjectsArgs{})
	// if err != nil {
	// 	return err
	// }

	return nil
}
