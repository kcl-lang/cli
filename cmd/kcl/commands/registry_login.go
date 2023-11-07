package cmd

import (
	"github.com/spf13/cobra"
	"kcl-lang.io/kpm/pkg/client"
	"kcl-lang.io/kpm/pkg/reporter"
	"kcl-lang.io/kpm/pkg/utils"
)

const (
	registryLoginDesc = `
  This command can be used to login to a registry.
  `
	registryLoginExample = `  # Login the docker hub
	kcl registry login docker.io

	# Login the GitHub container registry
	kcl registry login ghcr.io

	# Login a localhost registry
	kcl registry login https://localhost:5001
	`
)

// NewRegistryLoginCmd returns the registry login command.
func NewRegistryLoginCmd(cli *client.KpmClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "init",
		Short:   "login to a registry",
		Long:    registryLoginDesc,
		Example: registryLoginExample,
		RunE: func(_ *cobra.Command, args []string) error {
			registry := args[0]

			username, password, err := utils.GetUsernamePassword(username, password, false)
			if err != nil {
				return err
			}

			err = cli.LoginOci(registry, username, password)
			if err != nil {
				return err
			}
			reporter.ReportMsgTo("Login Succeeded", cli.GetLogWriter())
			return nil
		},
		// One positional argument that is the registry name.
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
	}

	cmd.Flags().StringVarP(&username, "username", "u", "", "registry username")
	cmd.Flags().StringVarP(&password, "password", "p", "", "registry password or identity token")

	return cmd
}
