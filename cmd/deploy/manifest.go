package deploy

import (
	"github.com/spf13/cobra"

	"github.com/frontierdigital/ranger/core/output"
)

var (
	pat         = ""
	projectName = ""
	orgName     = ""
)

// NewCmdDeploy creates a new deploy command
func NewCmdDeployManifest() *cobra.Command {
	c := &cobra.Command{
		Use:   "manifest",
		Short: "deploy a manifest",
		Long:  "deploy a manifest",
		RunE: func(_ *cobra.Command, _ []string) error {
			output.Println("deploy manifest")
			output.Printf("%s %s %s\n", pat, projectName, orgName)

			return nil
		},
	}

	c.Flags().StringVarP(&pat, "pat", "t", "", "Personal Access Token for ADO")
	c.Flags().StringVarP(&projectName, "proj", "p", "", "ADO Project Name")
	c.Flags().StringVarP(&orgName, "org", "o", "", "ADO Organisation")
	return c
}
