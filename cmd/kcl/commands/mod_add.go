package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"kcl-lang.io/kpm/pkg/client"
	"kcl-lang.io/kpm/pkg/env"
	"kcl-lang.io/kpm/pkg/errors"
	"kcl-lang.io/kpm/pkg/opt"
	pkg "kcl-lang.io/kpm/pkg/package"
	"kcl-lang.io/kpm/pkg/reporter"
)

const (
	modAddDesc = `This command adds new dependency
`
	modAddExample = `  # Add the module dependency named "k8s"
  kcl mod add k8s

  # Add the module dependency named "k8s" with the version "1.28"
  kcl mod add k8s:1.28

  # Add the module dependency from the GitHub
  kcl mod add --git https://github.com/kcl-lang/konfig --tag v0.4.0

  # Add a local dependency
  kcl mod add /path/to/another_module`
)

// NewModAddCmd returns the mod add command.
func NewModAddCmd(cli *client.KpmClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "add",
		Short:   "add new dependency",
		Long:    modAddDesc,
		Example: modAddExample,
		RunE: func(_ *cobra.Command, args []string) error {
			return ModAdd(cli, args)
		},
		SilenceUsage: true,
	}

	cmd.Flags().StringVar(&git, "git", "", "git repository url")
	cmd.Flags().StringVar(&tag, "tag", "", "git repository tag")
	cmd.Flags().StringVar(&commit, "commit", "", "git repository commit")
	cmd.Flags().StringVar(&branch, "branch", "", "git repository branch")
	cmd.Flags().StringVar(&rename, "rename", "", "rename the dependency")
	cmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "quiet (no output)")
	cmd.Flags().BoolVar(&noSumCheck, "no_sum_check", false, "do not check the checksum of the package and update kcl.mod.lock")

	return cmd
}

func ModAdd(cli *client.KpmClient, args []string) error {
	// acquire the lock of the package cache.
	err := cli.AcquirePackageCacheLock()
	if err != nil {
		return err
	}

	defer func() {
		// release the lock of the package cache after the function returns.
		releaseErr := cli.ReleasePackageCacheLock()
		if releaseErr != nil && err == nil {
			err = releaseErr
		}
	}()

	pwd, err := os.Getwd()

	if err != nil {
		return reporter.NewErrorEvent(reporter.Bug, err, "internal bugs, please contact us to fix it.")
	}

	globalPkgPath, err := env.GetAbsPkgPath()
	if err != nil {
		return err
	}

	kclPkg, err := pkg.LoadKclPkg(pwd)
	if err != nil {
		return err
	}

	err = kclPkg.ValidateKpmHome(globalPkgPath)
	if err != (*reporter.KpmEvent)(nil) {
		return err
	}

	addOpts, err := parseAddOptions(cli, globalPkgPath, args)
	if err != nil {
		return err
	}

	if addOpts.RegistryOpts.Local != nil {
		absAddPath, err := filepath.Abs(addOpts.RegistryOpts.Local.Path)
		if err != nil {
			return reporter.NewErrorEvent(reporter.Bug, err, "internal bugs, please contact us to fix it.")
		}
		if absAddPath == kclPkg.HomePath {
			return reporter.NewErrorEvent(
				reporter.AddItselfAsDep,
				fmt.Errorf("cannot add '%s' as a dependency to itself", kclPkg.GetPkgName()),
			)
		}
	}

	err = addOpts.Validate()
	if err != nil {
		return err
	}

	_, err = cli.AddDepWithOpts(kclPkg, addOpts)
	if err != nil {
		return err
	}
	return nil
}

// parseAddOptions will parse the user cli inputs.
func parseAddOptions(cli *client.KpmClient, localPath string, args []string) (*opt.AddOptions, error) {
	if len(args) == 0 {
		return &opt.AddOptions{
			LocalPath: localPath,
			RegistryOpts: opt.RegistryOptions{
				Git: &opt.GitOptions{
					Url:    git,
					Tag:    tag,
					Commit: commit,
					Branch: branch,
				},
			},
			NoSumCheck: noSumCheck,
			NewPkgName: rename,
		}, nil
	} else {
		localPkg, err := parseLocalPathOptions(args)
		if err != (*reporter.KpmEvent)(nil) {
			// parse from 'kpm add xxx:0.0.1'.
			ociReg, err := parseOciRegistryOptions(cli, args)
			if err != nil {
				return nil, err
			}
			return &opt.AddOptions{
				LocalPath:    localPath,
				RegistryOpts: *ociReg,
				NoSumCheck:   noSumCheck,
				NewPkgName:   rename,
			}, nil
		} else {
			return &opt.AddOptions{
				LocalPath:    localPath,
				RegistryOpts: *localPkg,
				NoSumCheck:   noSumCheck,
				NewPkgName:   rename,
			}, nil
		}
	}
}

// parseOciRegistryOptions will parse the oci registry information from user cli inputs.
func parseOciRegistryOptions(cli *client.KpmClient, args []string) (*opt.RegistryOptions, error) {
	ociPkgRef := argsGet(args, 0)
	name, version, err := parseOciPkgNameAndVersion(ociPkgRef)
	if err != nil {
		return nil, err
	}

	return &opt.RegistryOptions{
		Oci: &opt.OciOptions{
			Reg:     cli.GetSettings().DefaultOciRegistry(),
			Repo:    cli.GetSettings().DefaultOciRepo(),
			PkgName: name,
			Tag:     version,
		},
	}, nil
}

// parseLocalPathOptions will parse the local path information from user cli inputs.
func parseLocalPathOptions(args []string) (*opt.RegistryOptions, *reporter.KpmEvent) {
	localPath := argsGet(args, 0)
	if localPath == "" {
		return nil, reporter.NewErrorEvent(reporter.PathIsEmpty, errors.PathIsEmpty)
	}
	// check if the local path exists.
	if _, err := os.Stat(localPath); os.IsNotExist(err) {
		return nil, reporter.NewErrorEvent(reporter.LocalPathNotExist, err)
	} else {
		return &opt.RegistryOptions{
			Local: &opt.LocalOptions{
				Path: localPath,
			},
		}, nil
	}
}

// parseOciPkgNameAndVersion will parse package name and version
// from string "<pkg_name>:<pkg_version>".
func parseOciPkgNameAndVersion(s string) (string, string, error) {
	parts := strings.Split(s, ":")
	if len(parts) == 1 {
		return parts[0], "", nil
	}

	if len(parts) > 2 {
		return "", "", reporter.NewErrorEvent(reporter.InvalidPkgRef, fmt.Errorf("invalid oci package reference '%s'", s))
	}

	if parts[1] == "" {
		return "", "", reporter.NewErrorEvent(reporter.InvalidPkgRef, fmt.Errorf("invalid oci package reference '%s'", s))
	}

	return parts[0], parts[1], nil
}
