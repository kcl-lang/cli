package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"kcl-lang.io/cli/pkg/version"
)

// NewVersionCmd returns the version command.
func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show version of the KCL CLI",
		Run: func(*cobra.Command, []string) {
			fmt.Println(version.VersionTypeLatest)
		},
		SilenceUsage: true,
	}
}
