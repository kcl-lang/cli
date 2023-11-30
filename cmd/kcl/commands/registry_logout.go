package cmd

import (
	"github.com/spf13/cobra"
	"kcl-lang.io/kpm/pkg/client"
	"kcl-lang.io/kpm/pkg/reporter"
)

const (
	registryLogoutDesc = `This command can be used to logout from the current registry.
`
	registryLogoutExample = `  # Logout the registry
  kcl registry logout docker.io
`
)

// NewRegistryLogoutCmd returns the registry logout command.
func NewRegistryLogoutCmd(cli *client.KpmClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "logout",
		Short:   "logout from a registry",
		Long:    registryLogoutDesc,
		Example: registryLogoutExample,
		RunE: func(_ *cobra.Command, args []string) error {
			registry := args[0]
			err := cli.LogoutOci(registry)
			if err != nil {
				return err
			}
			reporter.ReportMsgTo("Logout Succeeded", cli.GetLogWriter())
			return nil
		},
		// One positional argument that is the registry name.
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
	}

	return cmd
}
