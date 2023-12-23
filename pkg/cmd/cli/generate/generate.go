package generate

import (
	"github.com/gofrontier-com/ranger/pkg/core"
	"github.com/spf13/cobra"
)

// NewCmdGenerate creates a command to generate documentation
func NewCmdGenerate(config *core.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate content",
	}

	cmd.AddCommand(NewCmdGenerateDocs(config))

	return cmd
}
