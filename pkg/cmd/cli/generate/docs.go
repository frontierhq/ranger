package generate

import (
	"github.com/frontierdigital/ranger/pkg/cmd/app/generate"
	"github.com/frontierdigital/ranger/pkg/util/config"

	"github.com/spf13/cobra"
)

var (
	projectName = ""
	orgName     = ""
)

// NewCmdGenerateDocs creates a command to deploy a set
func NewCmdGenerateDocs(config *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "docs",
		Short: "documentation",
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := generate.GenerateDocs(config, projectName, orgName); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&projectName, "project-name", "p", "", "Project name")
	cmd.Flags().StringVarP(&orgName, "organisation-name", "o", "", "Organisation name")

	return cmd
}
