// Copyright The KCL Authors. All rights reserved.

package cmd

import (
	"github.com/spf13/cobra"
	"kcl-lang.io/cli/pkg/options"
)

const (
	lintDesc = `This command lints the kcl code. 'kcl lint' takes multiple input for arguments.

For example, 'kcl lint path/to/kcl.k' will lint the file named path/to/kcl.k 
`
	lintExample = `  # Lint a single file and output YAML
  kcl lint path/to/kcl.k

  # Lint multiple files
  kcl lint path/to/kcl1.k path/to/kcl2.k

  # Lint OCI packages
  kcl lint oci://ghcr.io/kcl-lang/helloworld

  # Lint the current package
  kcl lint`
)

// NewLintCmd returns the lint command.
func NewLintCmd() *cobra.Command {
	o := options.NewRunOptions()
	cmd := &cobra.Command{
		Use:     "lint",
		Short:   "Lint KCL codes.",
		Long:    lintDesc,
		Example: lintExample,
		RunE: func(_ *cobra.Command, args []string) error {
			if err := o.Complete(args); err != nil {
				return err
			}
			if err := o.Validate(); err != nil {
				return err
			}
			o.CompileOnly = true
			return o.Run()
		},
		SilenceUsage: true,
	}

	appendLangFlags(o, cmd.Flags())

	return cmd
}
