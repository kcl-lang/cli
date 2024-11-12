package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/mod/module"
	"kcl-lang.io/kpm/pkg/client"
	"kcl-lang.io/kpm/pkg/env"
	pkg "kcl-lang.io/kpm/pkg/package"
	"kcl-lang.io/kpm/pkg/reporter"
)

const (
	modGraphDesc = `This command prints the module dependency graph.
Each module is identified as a string of the form path@version.
`
	modGraphExample = `  # Print the current module dependency graph.
  kcl mod graph`
)

// NewModGraphCmd returns the mod graph command.
func NewModGraphCmd(cli *client.KpmClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "graph",
		Short:   "prints dependencies",
		Long:    modGraphDesc,
		Example: modGraphExample,
		RunE: func(_ *cobra.Command, args []string) error {
			return ModGraph(cli, args)
		},
		SilenceUsage: true,
	}

	cmd.Flags().BoolVar(&insecureSkipTLSverify, "insecure-skip-tls-verify", false, "skip tls certificate checks for the KCL module download")
	return cmd
}

func ModGraph(cli *client.KpmClient, args []string) error {
	cli.SetInsecureSkipTLSverify(insecureSkipTLSverify)
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

	kclPkg, err := pkg.LoadKclPkgWithOpts(
		pkg.WithPath(pwd),
		pkg.WithSettings(cli.GetSettings()),
	)
	if err != nil {
		return err
	}

	err = kclPkg.ValidateKpmHome(globalPkgPath)
	if err != (*reporter.KpmEvent)(nil) {
		return err
	}

	gra, err := cli.Graph(
		client.WithGraphMod(kclPkg),
	)

	graStr, err := gra.DisplayGraphFromVertex(
		module.Version{Path: kclPkg.GetPkgName(), Version: kclPkg.GetPkgVersion()},
	)

	if err != nil {
		return err
	}

	reporter.ReportMsgTo(graStr, cli.GetLogWriter())

	return nil
}

func format(m module.Version) string {
	formattedMsg := m.Path
	if m.Version != "" {
		formattedMsg += "@" + m.Version
	}
	return formattedMsg
}
