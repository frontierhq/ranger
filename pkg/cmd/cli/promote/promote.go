package promote

import (
	"github.com/gofrontier-com/ranger/pkg/core"
	"github.com/spf13/cobra"
)

// NewCmdPromote creates a command to promote an artifact
func NewCmdPromote(config *core.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "promote",
		Short: "Promote an artifact",
	}

	cmd.AddCommand(NewCmdPromoteSet(config))

	return cmd
}
