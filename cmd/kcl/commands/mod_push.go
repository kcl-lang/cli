package cmd

import (
	"fmt"
	"net/url"
	"os"

	"github.com/spf13/cobra"
	"kcl-lang.io/kpm/pkg/client"
	"kcl-lang.io/kpm/pkg/errors"
	kpmoci "kcl-lang.io/kpm/pkg/oci"
	"kcl-lang.io/kpm/pkg/opt"
	pkg "kcl-lang.io/kpm/pkg/package"
	"kcl-lang.io/kpm/pkg/reporter"
	"kcl-lang.io/kpm/pkg/utils"
)

const (
	modPushDesc = `This command pushes kcl modules to the registry.
`
	modPushExample = `  # Push the current module to the registry
  kcl mod push
  
  # Push the current module to the registry in the vendor mode
  kcl mod push --vendor`
)

// NewModPushCmd returns the mod push command.
func NewModPushCmd(cli *client.KpmClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "push",
		Short:   "push kcl package to the registry",
		Long:    modPushDesc,
		Example: modPushExample,
		RunE: func(_ *cobra.Command, args []string) error {
			return ModPush(cli, args)
		},
		SilenceUsage: true,
	}

	cmd.Flags().BoolVar(&vendor, "vendor", false, "run in vendor mode (default: false)")
	cmd.Flags().StringVar(&target, "tar_path", "", "packaged target path that will be pushed")

	return cmd
}

func ModPush(cli *client.KpmClient, args []string) error {
	localTarPath := target
	ociUrl := argsGet(args, 0)

	var err error

	if len(localTarPath) == 0 {
		// If the tar package to be pushed is not specified,
		// the current kcl package is packaged into tar and pushed.
		err = pushCurrentPackage(ociUrl, vendor, cli)
	} else {
		// Else push the tar package specified.
		err = pushTarPackage(ociUrl, localTarPath, vendor, cli)
	}

	if err != nil {
		return err
	}

	return nil
}

// genDefaultOciUrlForKclPkg will generate the default oci url from the current package.
func genDefaultOciUrlForKclPkg(pkg *pkg.KclPkg, cli *client.KpmClient) (string, error) {

	urlPath := utils.JoinPath(cli.GetSettings().DefaultOciRepo(), pkg.GetPkgName())

	u := &url.URL{
		Scheme: kpmoci.OCI_SCHEME,
		Host:   cli.GetSettings().DefaultOciRegistry(),
		Path:   urlPath,
	}

	return u.String(), nil
}

// pushCurrentPackage will push the current package to the oci registry.
func pushCurrentPackage(ociUrl string, vendorMode bool, kpmcli *client.KpmClient) error {
	pwd, err := os.Getwd()

	if err != nil {
		reporter.ReportEventToStderr(reporter.NewEvent(reporter.Bug, "internal bug: failed to load working directory"))
		return err
	}
	// 1. Load the current kcl packege.
	kclPkg, err := pkg.LoadKclPkg(pwd)

	if err != nil {
		reporter.ReportEventToStderr(reporter.NewEvent(reporter.FailedLoadKclMod, fmt.Sprintf("failed to load package in '%s'", pwd)))
		return err
	}

	// 2. push the package
	return pushPackage(ociUrl, kclPkg, vendorMode, kpmcli)
}

// pushTarPackage will push the kcl package in tarPath to the oci registry.
// If the tar in 'tarPath' is not a kcl package tar, pushTarPackage will return an error.
func pushTarPackage(ociUrl, localTarPath string, vendorMode bool, kpmcli *client.KpmClient) error {
	var kclPkg *pkg.KclPkg
	var err error

	// clean the temp dir used to untar kcl package tar file.
	defer func() {
		if kclPkg != nil && utils.DirExists(kclPkg.HomePath) {
			err = os.RemoveAll(kclPkg.HomePath)
			if err != nil {
				err = errors.InternalBug
			}
		}
	}()

	// 1. load the kcl package from the tar path.
	kclPkg, err = pkg.LoadKclPkgFromTar(localTarPath)
	if err != nil {
		return err
	}

	// 2. push the package
	return pushPackage(ociUrl, kclPkg, vendorMode, kpmcli)
}

// pushPackage will push the kcl package to the oci registry.
// 1. pushPackage will package the current kcl package into default tar path.
// 2. If the oci url is not specified, generate the default oci url from the current package.
// 3. Generate the OCI options from oci url and the version of current kcl package.
// 4. Push the package to the oci registry.
func pushPackage(ociUrl string, kclPkg *pkg.KclPkg, vendorMode bool, cli *client.KpmClient) error {

	tarPath, err := cli.PackagePkg(kclPkg, vendorMode)
	if err != nil {
		return err
	}

	// clean the tar path.
	defer func() {
		if kclPkg != nil && utils.DirExists(tarPath) {
			err = os.RemoveAll(tarPath)
			if err != nil {
				err = errors.InternalBug
			}
		}
	}()

	// 2. If the oci url is not specified, generate the default oci url from the current package.
	if len(ociUrl) == 0 {
		ociUrl, err = genDefaultOciUrlForKclPkg(kclPkg, cli)
		if err != nil || len(ociUrl) == 0 {
			return reporter.NewErrorEvent(
				reporter.InvalidCmd,
				fmt.Errorf("failed to generate default oci url for current package"),
				"run 'kpm push help' for more information",
			)
		}
	}

	// 3. Generate the OCI options from oci url and the version of current kcl package.
	ociOpts, err := opt.ParseOciOptionFromOciUrl(ociUrl, kclPkg.GetPkgTag())
	if err != (*reporter.KpmEvent)(nil) {
		return reporter.NewErrorEvent(
			reporter.UnsupportOciUrlScheme,
			errors.InvalidOciUrl,
			"only support url scheme 'oci://'.",
		)
	}
	ociOpts.Annotations, err = kpmoci.GenOciManifestFromPkg(kclPkg)
	if err != nil {
		return err
	}

	reporter.ReportMsgTo(fmt.Sprintf("package '%s' will be pushed", kclPkg.GetPkgName()), cli.GetLogWriter())
	// 4. Push it.
	err = cli.PushToOci(tarPath, ociOpts)
	if err != (*reporter.KpmEvent)(nil) {
		return err
	}
	return nil
}
