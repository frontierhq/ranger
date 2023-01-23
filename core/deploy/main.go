package deploy

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/frontierdigital/ranger/core/structs"
	"github.com/frontierdigital/ranger/core/util"
	"github.com/microsoft/azure-devops-go-api/azuredevops"
	"github.com/microsoft/azure-devops-go-api/azuredevops/core"
)

func DeployManifest(configuration *structs.Configuration, projectName string, organisationName string) error {
	organisationUrl := fmt.Sprintf("https://dev.azure.com/%s", organisationName)
	connection := azuredevops.NewPatConnection(organisationUrl, configuration.ADO.PAT)

	ctx := context.Background()

	coreClient, err := core.NewClient(ctx, connection)
	if err != nil {
		return err
	}

	responseValue, err := coreClient.GetProjects(ctx, core.GetProjectsArgs{})
	if err != nil {
		return err
	}

	manifestFilepath, _ := filepath.Abs("./manifest.yml")
	manifest, err := util.LoadManifest(manifestFilepath)
	if err != nil {
		return err
	}

	util.PrintlnInfo(responseValue)
	util.PrintlnInfo(manifest)

	return nil
}
