package cmd

import (
	"github.com/spf13/cobra"
	"kcl-lang.io/kpm/pkg/client"
)

const (
	modPullDesc = `
  This command pulls kcl modules from the registry.
  `
	modPullExample = `  # Pull the the module named "k8s" to the local path from the registry
	kcl mod pull k8s`
)

// NewModPullCmd returns the mod pull command.
func NewModPullCmd(cli *client.KpmClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pull",
		Short:   "pull kcl package from the registry",
		Long:    modPullDesc,
		Example: modPullExample,
		RunE: func(_ *cobra.Command, args []string) error {
			source := argsGet(args, 0)
			localPath := argsGet(args, 1)
			return cli.PullFromOci(localPath, source, tag)
		},
		SilenceUsage: true,
	}

	cmd.Flags().StringVar(&tag, "tag", "", "git repository tag")

	return cmd
}
