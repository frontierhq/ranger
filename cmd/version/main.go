package version

import (
	"github.com/spf13/cobra"
	goVersion "go.hein.dev/go-version"

	"github.com/frontierdigital/ranger/core/output"
)

var (
	outputFmt = "json"
	shortened = false
)

// NewCmdVersion creates a command to output the current version of SoloDeploy
func NewCmdVersion(version string, commit string, date string) *cobra.Command {
	c := &cobra.Command{
		Use:   "version",
		Short: "Version will output the current build information",
		Long:  "Prints the version, Git commit ID and commit date in JSON or YAML format using the go.hein.dev/go-version package.",
		Run: func(_ *cobra.Command, _ []string) {
			resp := goVersion.FuncWithOutput(shortened, version, commit, date, outputFmt)
			output.PrintfInfo(resp)
		},
	}

	c.Flags().BoolVarP(&shortened, "short", "s", false, "Print just the version number.")
	c.Flags().StringVarP(&outputFmt, "output", "o", "json", "Output format. One of 'yaml' or 'json'.")

	return c
}
