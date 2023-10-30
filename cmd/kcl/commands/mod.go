package cmd

import (
	"github.com/spf13/cobra"
	cli "github.com/urfave/cli/v2"
	"kcl-lang.io/cli/pkg/version"
	"kcl-lang.io/kpm/pkg/client"
	kpmcmd "kcl-lang.io/kpm/pkg/cmd"
	"kcl-lang.io/kpm/pkg/reporter"
)

const (
	modDesc = `
This command manages the kcl module
`
	modExample = `  # Init one kcl module
  kcl mod init

  # Add dependencies for the current module
  kcl mod add k8s

  # Pull external packages to local
  kcl mod pull k8s

  # Push the module
  kcl mod push`
)

// NewModCmd returns the mod command.
func NewModCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "mod",
		Short:   "KCL module management",
		Long:    modDesc,
		Example: modExample,
		RunE: func(_ *cobra.Command, args []string) error {
			return RunWithKpmMod("mod", args)
		},
		SilenceUsage: true,
	}

	return cmd
}

func RunWithKpmMod(cmd string, args []string) error {
	reporter.InitReporter()
	kpmcli, err := client.NewKpmClient()
	if err != nil {
		return err
	}
	app := cli.NewApp()
	app.Usage = "module related functions"
	app.Name = "kcl mod"
	app.Version = version.GetVersionString()
	app.UsageText = "kcl mod <command> [arguments]..."
	app.Commands = []*cli.Command{
		kpmcmd.NewInitCmd(kpmcli),
		kpmcmd.NewAddCmd(kpmcli),
		kpmcmd.NewPkgCmd(kpmcli),
		kpmcmd.NewMetadataCmd(kpmcli),
		kpmcmd.NewRunCmd(kpmcli),
		kpmcmd.NewLoginCmd(kpmcli),
		kpmcmd.NewLogoutCmd(kpmcli),
		kpmcmd.NewPushCmd(kpmcli),
		kpmcmd.NewPullCmd(kpmcli),
	}
	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:  kpmcmd.FLAG_QUIET,
			Usage: "push in vendor mode",
		},
	}
	app.Before = func(c *cli.Context) error {
		if c.Bool(kpmcmd.FLAG_QUIET) {
			kpmcli.SetLogWriter(nil)
		}
		return nil
	}
	argsWithCmd := []string{cmd}
	argsWithCmd = append(argsWithCmd, args...)
	return app.Run(argsWithCmd)
}
