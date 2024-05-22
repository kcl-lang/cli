package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"kcl-lang.io/kpm/pkg/client"
	"kcl-lang.io/kpm/pkg/env"
	"kcl-lang.io/kpm/pkg/errors"
	"kcl-lang.io/kpm/pkg/reporter"
)

const (
	modUpdateDesc = `This command updates dependencies listed in kcl.mod.lock based on kcl.mod.
`
	modUpdateExample = `  # Update the current module
  kcl mod update
  
  # Update the module with the specified path
  kcl mod update path/to/package`
)

// NewModUpdateCmd returns the mod update command.
func NewModUpdateCmd(cli *client.KpmClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update",
		Short:   "updates dependencies",
		Long:    modUpdateDesc,
		Example: modUpdateExample,
		RunE: func(_ *cobra.Command, args []string) error {
			return ModUpdate(cli, args)
		},
		SilenceUsage: true,
	}

	cmd.Flags().BoolVar(&noSumCheck, "no_sum_check", false, "do not check the checksum of the package and update kcl.mod.lock")

	return cmd
}

func ModUpdate(cli *client.KpmClient, args []string) error {
	pkgPath := argsGet(args, 0)

	if len(pkgPath) == 0 {
		pwd, err := os.Getwd()
		if err != nil {
			return errors.InternalBug
		}
		pkgPath = pwd
	}
	cli.SetNoSumCheck(noSumCheck)
	kclPkg, err := cli.LoadPkgFromPath(pkgPath)
	if err != nil {
		return err
	}

	globalPkgPath, err := env.GetAbsPkgPath()
	if err != nil {
		return err
	}

	err = kclPkg.ValidateKpmHome(globalPkgPath)
	if err != (*reporter.KpmEvent)(nil) {
		return err
	}

	_, err = cli.ResolveDepsMetadataInJsonStr(kclPkg, true)
	if err != nil {
		return err
	}

	err = kclPkg.UpdateModAndLockFile()
	if err != nil {
		return err
	}
	return nil
}
