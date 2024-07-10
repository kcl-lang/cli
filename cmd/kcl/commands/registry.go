package cmd

import (
	"github.com/spf13/cobra"
	"kcl-lang.io/kpm/pkg/client"
	"kcl-lang.io/kpm/pkg/reporter"
)

const (
	registryDesc = `This command manages the kcl registry
`
	registryExample = `  # Login the registry
  kcl registry login docker.io

  # Logout the registry
  kcl registry logout`
)

var (
	username string
	password string
)

// NewModCmd returns the mod command.
func NewRegistryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "registry",
		Short:        "KCL registry management",
		Long:         registryDesc,
		Example:      registryExample,
		SilenceUsage: true,
		Aliases:      []string{"reg", "r"},
	}

	cmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "Set the quiet mode (no output)")

	reporter.InitReporter()
	cli, err := client.NewKpmClient()
	if err != nil {
		panic(err)
	}
	if quiet {
		cli.SetLogWriter(nil)
	}

	cmd.AddCommand(NewRegistryLoginCmd(cli))
	cmd.AddCommand(NewRegistryLogoutCmd(cli))

	return cmd
}
