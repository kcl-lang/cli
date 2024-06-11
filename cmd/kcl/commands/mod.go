package cmd

import (
	"github.com/spf13/cobra"
	"kcl-lang.io/kpm/pkg/client"
	"kcl-lang.io/kpm/pkg/reporter"
)

const (
	modDesc = `This command manages the kcl module
`
	modExample = `  kcl mod <command> [arguments]...
  # Init one kcl module
  kcl mod init

  # Add dependencies for the current module
  kcl mod add k8s

  # Pull external packages to local
  kcl mod pull k8s

  # Push the module
  kcl mod push
  
  # Print the current module dependency graph.
  kcl mod graph`
)

var (
	quiet      bool
	vendor     bool
	update     bool
	git        string
	oci        string
	path       string
	tag        string
	commit     string
	branch     string
	target     string
	rename     string
	noSumCheck bool
)

// NewModCmd returns the mod command.
func NewModCmd() *cobra.Command {
	reporter.InitReporter()
	cli, err := client.NewKpmClient()
	if err != nil {
		panic(err)
	}
	cmd := &cobra.Command{
		Use:          "mod",
		Short:        "KCL module management",
		Long:         modDesc,
		Example:      modExample,
		SilenceUsage: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if quiet {
				cli.SetLogWriter(nil)
			}
			return nil
		},
	}

	cmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "Set the quiet mode (no output)")

	cmd.AddCommand(NewModInitCmd(cli))
	cmd.AddCommand(NewModAddCmd(cli))
	cmd.AddCommand(NewModPkgCmd(cli))
	cmd.AddCommand(NewModMetadataCmd(cli))
	cmd.AddCommand(NewModPushCmd(cli))
	cmd.AddCommand(NewModPullCmd(cli))
	cmd.AddCommand(NewModUpdateCmd(cli))
	cmd.AddCommand(NewModGraphCmd(cli))

	return cmd
}
