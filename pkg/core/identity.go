package core

import (
	"github.com/frontierdigital/utils/azuredevops"
	"github.com/google/uuid"
)

func (ado *AzureDevOps) GetIdentityId() (*uuid.UUID, error) {
	azureDevOps := azuredevops.NewAzureDevOps(ado.OrganisationName, ado.PAT)
	identityId, err := azureDevOps.GetIdentityId()
	if err != nil {
		return nil, err
	}
	return identityId, nil
}
