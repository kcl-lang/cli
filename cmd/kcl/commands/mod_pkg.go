package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"kcl-lang.io/kpm/pkg/client"
	"kcl-lang.io/kpm/pkg/errors"
	pkg "kcl-lang.io/kpm/pkg/package"
	"kcl-lang.io/kpm/pkg/reporter"
	"kcl-lang.io/kpm/pkg/utils"
)

const (
	modPkgDesc = `This command converts a kcl module into tar
`
	modPkgExample = `  # Package the current module
  kcl mod pkg
  
  # Package the current module in the vendor mode
  kcl mod pkg --vendor
	
  # Package the current module into the target directory
  kcl mod pkg --target /path/to/target_dir`
)

// NewModPkgCmd returns the mod pkg command.
func NewModPkgCmd(cli *client.KpmClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pkg",
		Short:   "package a kcl package into tar",
		Long:    modPkgDesc,
		Example: modPkgExample,
		RunE: func(*cobra.Command, []string) error {
			if len(target) == 0 {
				return reporter.NewErrorEvent(
					reporter.InvalidCmd,
					fmt.Errorf("the directory where the tar is generated is required"),
					"run 'kpm pkg help' for more information",
				)
			}

			pwd, err := os.Getwd()

			if err != nil {
				reporter.ExitWithReport("internal bug: failed to load working directory")
			}

			kclPkg, err := pkg.LoadKclPkg(pwd)

			if err != nil {
				reporter.ExitWithReport("failed to load package in " + pwd + ".")
				return err
			}

			// If the file path used to save the package tar file does not exist, create this file path.
			if !utils.DirExists(target) {
				err := os.MkdirAll(target, os.ModePerm)
				if err != nil {
					return errors.InternalBug
				}
			}

			return cli.Package(kclPkg, filepath.Join(target, kclPkg.GetPkgTarName()), vendor)
		},
		SilenceUsage: true,
	}

	cmd.Flags().StringVar(&target, "target", "", "packaged target path")
	cmd.Flags().BoolVar(&vendor, "vendor", false, "package in vendor mode (default: false)")

	return cmd
}
