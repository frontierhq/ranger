package generate

import (
	"github.com/frontierdigital/ranger/pkg/util/config"

	"github.com/spf13/cobra"
)

// NewCmdGenerate creates a command to generate documentation
func NewCmdGenerate(config *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate content",
	}

	cmd.AddCommand(NewCmdGenerateDocs(config))

	return cmd
}
