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
	registryDesc = `
This command manages the kcl registry
`
	registryExample = `  # Login the registry
  kcl registry login docker.io

  # Logout the registry
  kcl registry logout
  `
)

// NewModCmd returns the mod command.
func NewRegistryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "registry",
		Short:   "KCL registry management",
		Long:    registryDesc,
		Example: registryExample,
		RunE: func(_ *cobra.Command, args []string) error {
			return RunWithKpmRegistry("registry", args)
		},
		SilenceUsage: true,
	}

	return cmd
}

func RunWithKpmRegistry(cmd string, args []string) error {
	reporter.InitReporter()
	kpmcli, err := client.NewKpmClient()
	if err != nil {
		return err
	}
	app := cli.NewApp()
	app.Usage = "registry related functions"
	app.Name = "kcl registry"
	app.Version = version.GetVersionString()
	app.UsageText = "kcl registry <command> [arguments]..."
	app.Commands = []*cli.Command{
		kpmcmd.NewLoginCmd(kpmcli),
		kpmcmd.NewLogoutCmd(kpmcli),
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
