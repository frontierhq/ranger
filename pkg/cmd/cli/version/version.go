package version

import (
	"github.com/frontierdigital/utils/output"

	"github.com/spf13/cobra"
	goVersion "go.hein.dev/go-version"
)

var (
	outputFmt = "json"
	shortened = false
)

// NewCmdVersion creates a command to output the current version of Ranger
func NewCmdVersion(version string, commit string, date string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Version will output the current build information",
		Long:  "Prints the version, Git commit ID and commit date in JSON or YAML format using the go.hein.dev/go-version package.",
		RunE: func(_ *cobra.Command, _ []string) error {
			resp := goVersion.FuncWithOutput(shortened, version, commit, date, outputFmt)
			output.PrintfInfo(resp)

			return nil
		},
	}

	cmd.Flags().BoolVarP(&shortened, "short", "s", false, "Print just the version number.")
	cmd.Flags().StringVarP(&outputFmt, "output", "o", "json", "Output format. One of 'yaml' or 'json'.")

	return cmd
}
