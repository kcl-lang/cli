// Copyright The KCL Authors. All rights reserved.

package cmd

import (
	"github.com/spf13/cobra"
	"kcl-lang.io/cli/pkg/options"
)

const (
	runDesc = `This command runs the kcl code and displays the output. 'kcl run' takes multiple input for arguments.

For example, 'kcl run path/to/kcl.k' will run the file named path/to/kcl.k 
`
	runExample = `  # Run the current package
  kcl run

  # Run a single file and output YAML
  kcl run path/to/kcl.k

  # Run a single file and output JSON
  kcl run path/to/kcl.k --format json

  # Run a single file and output TOML
  kcl run path/to/kcl.k --format toml

  # Run multiple files
  kcl run path/to/kcl1.k path/to/kcl2.k

  # Run OCI packages
  kcl run oci://ghcr.io/kcl-lang/helloworld --tag 0.1.0

  # Run remote Git repo
  kcl run git://github.com/kcl-lang/flask-demo-kcl-manifests --commit ade147b

  # Run OCI packages by flag
  kcl run --oci https://ghcr.io/kcl-lang/helloworld --tag 0.1.0

  # Run remote module from Git with branch repo by flag
  kcl run --git https://github.com/kcl-lang/flask-demo-kcl-manifests --branch main`
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
		Aliases:      []string{"r"},
		SilenceUsage: true,
	}

	appendLangFlags(o, cmd.Flags())

	return cmd
}
