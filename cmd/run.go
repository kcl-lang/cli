// Copyright The KCL Authors. All rights reserved.

package cmd

import (
	"github.com/spf13/cobra"
	"kcl-lang.io/cli/pkg/options"
)

const (
	runDesc = `
This command runs the kcl code and displays the output. 'kcl run' takes multiple input for arguments.

For example, 'kcl run path/to/kcl.k' will run the file named path/to/kcl.k 
`
	runExample = `  # Run a single file and output YAML
  kcl run path/to/kcl.k

  # Run a single file and output JSON
  kcl run path/to/kcl.k --format json

  # Run multiple files
  kcl run path/to/kcl1.k path/to/kcl2.k
  
  # Run OCI packages
  kcl run oci://ghcr.io/kcl-lang/hello-world
  
  # Run the current package
  kcl run
  `
)

// NewRunCmd returns the run command.
func NewRunCmd() *cobra.Command {
	o := options.NewRunOptions()
	cmd := &cobra.Command{
		Use:     "run",
		Short:   "Run KCL codes.",
		Long:    runDesc,
		Example: runExample,
		RunE: func(_ *cobra.Command, args []string) error {
			if err := o.Complete(args); err != nil {
				return err
			}
			if err := o.Validate(); err != nil {
				return err
			}
			return o.Run()
		},
		SilenceUsage: true,
	}

	cmd.Flags().StringSliceVarP(&o.Arguments, "argument", "D", []string{},
		"Specify the top-level argument")
	cmd.Flags().StringSliceVarP(&o.Settings, "setting", "Y", []string{},
		"Specify the command line setting files")
	cmd.Flags().StringSliceVarP(&o.Overrides, "overrides", "O", []string{},
		"Specify the configuration override path and value")
	cmd.Flags().StringSliceVarP(&o.PathSelectors, "path_selectors", "S", []string{},
		"Specify the path selectors")
	cmd.Flags().StringSliceVarP(&o.ExternalPackages, "external", "E", []string{},
		" Mapping of package name and path where the package is located")
	cmd.Flags().StringVarP(&o.Output, "output", "o", "",
		"Specify the YAML/JSON output file path")
	cmd.Flags().StringVarP(&o.Tag, "tag", "t", "",
		"Specify the tag for the OCI or Git artifact")
	cmd.Flags().StringVar(&o.Format, "format", "yaml",
		"Specify the output format")
	cmd.Flags().BoolVarP(&o.DisableNone, "disable_none", "n", false,
		"Disable dumping None values")
	cmd.Flags().BoolVarP(&o.StrictRangeCheck, "strict_range_check", "r", false,
		"Do perform strict numeric range checks")
	cmd.Flags().BoolVarP(&o.Debug, "debug", "d", false,
		"Run in debug mode")
	cmd.Flags().BoolVarP(&o.SortKeys, "sort_keys", "k", false,
		"Sort output result keys")
	cmd.Flags().BoolVarP(&o.Vendor, "vendor", "V", false,
		"Sort output result keys")
	cmd.Flags().BoolVar(&o.NoStyle, "no_style", false,
		"Sort output result keys")

	return cmd
}
