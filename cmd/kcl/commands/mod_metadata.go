package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"kcl-lang.io/kpm/pkg/client"
	"kcl-lang.io/kpm/pkg/env"
	"kcl-lang.io/kpm/pkg/errors"
	"kcl-lang.io/kpm/pkg/reporter"
)

const (
	modMetadataDesc = `This command outputs the resolved dependencies of a package
`
	modMetadataExample = `  # Output the resolved dependencies the current module
  kcl mod metadata
  
  # Output the resolved dependencies the current module in the vendor mode
  kcl mod metadata --vendor
	
  # Output the resolved dependencies the current module with the update check
  kcl mod metadata --update`
)

// NewModMetadataCmd returns the mod metadata command.
func NewModMetadataCmd(cli *client.KpmClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "metadata",
		Short:   "output the resolved dependencies of a package",
		Long:    modMetadataDesc,
		Example: modMetadataExample,
		RunE: func(*cobra.Command, []string) error {
			pwd, err := os.Getwd()
			if err != nil {
				return errors.InternalBug
			}

			cli.SetLogWriter(nil)
			kclPkg, err := cli.LoadPkgFromPath(pwd)
			if err != nil {
				return err
			}

			globalPkgPath, err := env.GetAbsPkgPath()
			if err != nil {
				return err
			}

			kclPkg.SetVendorMode(vendor)

			err = kclPkg.ValidateKpmHome(globalPkgPath)
			if err != (*reporter.KpmEvent)(nil) {
				return err
			}

			jsonStr, err := cli.ResolveDepsMetadataInJsonStr(kclPkg, update)
			if err != nil {
				return err
			}

			if update {
				err = kclPkg.UpdateModAndLockFile()
				if err != nil {
					return err
				}
			}

			fmt.Println(jsonStr)

			return nil
		},
		SilenceUsage: true,
	}

	cmd.Flags().BoolVar(&vendor, "vendor", false, "run in vendor mode (default: false)")
	cmd.Flags().BoolVar(&update, "update", false, "check the local package and update and download the local package. (default: false)")

	return cmd
}
