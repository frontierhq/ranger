package generate

import (
	"github.com/frontierdigital/ranger/pkg/cmd/app/generate"
	"github.com/frontierdigital/ranger/pkg/util/config"

	"github.com/spf13/cobra"
)

var (
	projectName = ""
	orgName     = ""
	wikiName    = ""
	feedName    = ""
)

// NewCmdGenerateDocs creates a command to deploy a set
func NewCmdGenerateDocs(config *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "docs",
		Short: "documentation",
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := generate.GenerateDocs(config, projectName, orgName, wikiName, feedName); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&projectName, "project-name", "p", "", "Project name")
	cmd.Flags().StringVarP(&orgName, "organisation-name", "o", "", "Organisation name")
	cmd.Flags().StringVarP(&wikiName, "wiki-name", "w", "", "Repository name the stores the wiki docs")
	cmd.Flags().StringVarP(&feedName, "feed-name", "f", "", "ADO artifact feed name")

	return cmd
}
