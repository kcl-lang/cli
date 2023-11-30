package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"kcl-lang.io/kpm/pkg/client"
	"kcl-lang.io/kpm/pkg/env"
	"kcl-lang.io/kpm/pkg/opt"
	pkg "kcl-lang.io/kpm/pkg/package"
	"kcl-lang.io/kpm/pkg/reporter"
)

const (
	modInitDesc = `This command initializes new kcl module in current directory.
`
	modInitExample = `  kcl mod init <command> [arguments]...
  # Init one kcl module with the current folder name
  kcl mod init
  
  # Init one kcl module with the name
  kcl mod init package-name`
)

// NewModInitCmd returns the mod init command.
func NewModInitCmd(cli *client.KpmClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "init",
		Short:   "initialize new module in current directory",
		Long:    modInitDesc,
		Example: modInitExample,
		RunE: func(_ *cobra.Command, args []string) error {
			pwd, err := os.Getwd()

			if err != nil {
				reporter.Fatal("kpm: internal bugs, please contact us to fix it")
			}

			var pkgName string
			var pkgRootPath string
			// 1. If no package name is given, the current directory name is used as the package name.
			if len(args) == 0 {
				pkgName = filepath.Base(pwd)
				pkgRootPath = pwd
			} else {
				// 2. If the package name is given, create a new directory for the package.
				pkgName = argsGet(args, 0)
				pkgRootPath = filepath.Join(pwd, pkgName)
				err = os.MkdirAll(pkgRootPath, 0755)
				if err != nil {
					return err
				}
			}

			initOpts := opt.InitOptions{
				Name:     pkgName,
				InitPath: pkgRootPath,
			}

			err = initOpts.Validate()
			if err != nil {
				return err
			}

			kclPkg := pkg.NewKclPkg(&initOpts)

			globalPkgPath, err := env.GetAbsPkgPath()

			if err != nil {
				return err
			}

			err = kclPkg.ValidateKpmHome(globalPkgPath)

			if err != (*reporter.KpmEvent)(nil) {
				return err
			}

			err = cli.InitEmptyPkg(&kclPkg)
			if err != nil {
				return err
			}

			reporter.ReportMsgTo(fmt.Sprintf("kpm: package '%s' init finished", pkgName), cli.GetLogWriter())
			return nil
		},
		SilenceUsage: true,
	}

	return cmd
}
