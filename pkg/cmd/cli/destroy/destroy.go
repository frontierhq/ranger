package destroy

import (
	"github.com/frontierdigital/ranger/pkg/core"
	"github.com/spf13/cobra"
)

// NewCmdDestroy creates a command to destroy an artifact
func NewCmdDestroy(config *core.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "destroy",
		Short: "Destroy an artifact",
	}

	cmd.AddCommand(NewCmdDestroySet(config))

	return cmd
}
