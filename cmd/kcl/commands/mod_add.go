package cmd

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"

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
	modAddExample = `  # Add the module dependency named "k8s" from the default OCI registry
  kcl mod add k8s

  # Add the module dependency named "k8s" with the version "1.28" from the default OCI registry
  kcl mod add k8s:1.28

  # Add the module dependency from the GitHub by git url
  kcl mod add git://github.com/kcl-lang/konfig --tag v0.4.0

  # Add the module dependency from the OCI Registry by oci url
  kcl mod add oci://ghcr.io/kcl-lang/helloworld --tag 0.1.0

  # Add the module dependency from the local file system by file url
  kcl mod add /path/to/another_module

  # Add the module dependency from the GitHub by the tag flag
  kcl mod add --git https://github.com/kcl-lang/konfig --tag v0.4.0

  # Add the sub module dependency named "helloworld" from the Git repo by the tag flag
  kcl mod add helloworld --git https://github.com/kcl-lang/modules --tag v0.1.0

  # Add the module dependency from the OCI registry named "" by the tag flag
  kcl mod add --oci https://ghcr.io/kcl-lang/helloworld --tag 0.1.0`
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
	cmd.Flags().StringVar(&oci, "oci", "", "oci repository url")
	cmd.Flags().StringVar(&tag, "tag", "", "git or oci repository tag")
	cmd.Flags().StringVar(&commit, "commit", "", "git repository commit")
	cmd.Flags().StringVar(&branch, "branch", "", "git repository branch")
	cmd.Flags().StringVar(&path, "path", "", "filesystem path to local dependency to add")
	cmd.Flags().StringVar(&rename, "rename", "", "rename the dependency")
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
	// parse the CLI command with the following style
	// kcl mod add <package>
	// kcl mod add <package>:<version>
	// kcl mod add /path/to/xxx
	// kcl mod add https://xxx/xxx --tag 0.0.1
	// kcl mod add oci://xxx/xxx --tag 0.0.1
	//
	// kcl mod add --git https://xxx/xxx --tag 0.0.1
	// kcl mod add <sub_package> --git https://xxx/xxx --tag 0.0.1
	// kcl mod add --oci https://xxx/xxx --tag 0.0.1
	// kcl mod add <sub_package> --oci https://xxx/xxx --tag 0.0.1
	// kcl mod add --path /path/to/xxx
	// kcl mod add <sub_package> --path /path/to/xxx
	if len(git) != 0 {
		gitUrl, err := url.Parse(git)
		if err != nil {
			return nil, err
		}
		gitOpts := opt.NewGitOptionsFromUrl(gitUrl)
		if gitOpts == nil {
			return nil, fmt.Errorf("invalid git url '%s'", git)
		}
		gitOpts.Tag = tag
		gitOpts.Commit = commit
		gitOpts.Branch = branch
		// Git sub package.
		if len(args) > 0 {
			gitOpts.Package = args[len(args)-1]
		}
		return &opt.AddOptions{
			LocalPath:    localPath,
			RegistryOpts: opt.RegistryOptions{Git: gitOpts},
			NoSumCheck:   noSumCheck,
			NewPkgName:   rename,
		}, nil
	} else if len(oci) != 0 {
		ociUrl, err := url.Parse(oci)
		if err != nil {
			return nil, err
		}
		ociOpts := opt.NewOciOptionsFromUrl(ociUrl)
		if ociOpts == nil {
			return nil, fmt.Errorf("invalid oci url '%s'", oci)
		}
		ociOpts.Tag = tag
		// OCI sub package
		if len(args) > 0 {
			ociOpts.Package = args[len(args)-1]
		}
		return &opt.AddOptions{
			LocalPath:    localPath,
			RegistryOpts: opt.RegistryOptions{Oci: ociOpts},
			NoSumCheck:   noSumCheck,
			NewPkgName:   rename,
		}, nil
	} else if len(path) != 0 {
		pathUrl, err := url.Parse(path)
		if err != nil {
			return nil, err
		}
		pathOpts, err := opt.NewLocalOptionsFromUrl(pathUrl)
		if err != (*reporter.KpmEvent)(nil) {
			return nil, err
		}
		// Local path sub package
		if len(args) > 0 {
			pathOpts.Package = args[len(args)-1]
		}
		return &opt.AddOptions{
			LocalPath:    localPath,
			RegistryOpts: opt.RegistryOptions{Local: pathOpts},
			NoSumCheck:   noSumCheck,
			NewPkgName:   rename,
		}, nil
	} else {
		localPkg, err := parseLocalPathOptions(args)
		pkgSource := argsGet(args, 0)
		if err != (*reporter.KpmEvent)(nil) {
			// parse url and ref
			regOpt, err := opt.NewRegistryOptionsFrom(pkgSource, cli.GetSettings())
			if err != nil {
				return nil, err
			}

			if regOpt.Git != nil {
				regOpt.Git.Tag = tag
				regOpt.Git.Commit = commit
				regOpt.Git.Branch = branch
			} else if regOpt.Oci != nil && len(tag) != 0 {
				regOpt.Oci.Tag = tag
			} else if regOpt.Registry != nil {
				var invalidFlag string
				if len(tag) != 0 {
					invalidFlag = tag
				} else if len(commit) != 0 {
					invalidFlag = commit
				} else if len(branch) != 0 {
					invalidFlag = branch
				}
				var diagMsg string
				if len(invalidFlag) != 0 {
					diagMsg = fmt.Sprintf("invalid flag '%s' for oci registry ref", invalidFlag)
				}
				if len(tag) != 0 {
					diagMsg = fmt.Sprintf("%s, try 'kcl mod add %s:%s'", diagMsg, pkgSource, tag)
				}
				return nil, fmt.Errorf(diagMsg)
			}

			return &opt.AddOptions{
				LocalPath:    localPath,
				RegistryOpts: *regOpt,
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
