package promote

import (
	"github.com/frontierdigital/ranger/pkg/util/config"

	"github.com/spf13/cobra"
)

// NewCmdPromote creates a command to promote an artifact
func NewCmdPromote(config *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "promote",
		Short: "Promote an artifact",
	}

	cmd.AddCommand(NewCmdPromoteSet(config))

	return cmd
}
