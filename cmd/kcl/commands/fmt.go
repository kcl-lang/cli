// Copyright The KCL Authors. All rights reserved.

package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	kcl "kcl-lang.io/kcl-go"
)

const (
	fmtDesc = `This command formats all kcl files of the current crate.
`
	fmtExample = `  # Format the single file
  kcl fmt /path/to/file.k

  # Format all files in this folder recursively
  kcl fmt ./...`
)

// NewFmtCmd returns the fmt command.
func NewFmtCmd() *cobra.Command {
	o := new(kcl.FormatPathOptions)
	cmd := &cobra.Command{
		Use:     "fmt",
		Short:   "KCL format tool",
		Long:    fmtDesc,
		Example: fmtExample,
		RunE: func(_ *cobra.Command, args []string) error {
			var changedPaths []string
			if len(args) == 0 {
				args = append(args, ".")
			}
			for _, p := range args {
				paths, err := kcl.FormatPathWithOptions(p, kcl.FormatPathOptions{
					DryRun: o.DryRun,
				})
				if err != nil {
					return err
				}
				changedPaths = append(changedPaths, paths...)
			}
			if len(changedPaths) > 0 {
				fmt.Println(strings.Join(changedPaths, "\n"))
			}
			if o.DryRun && len(changedPaths) > 0 {
				return fmt.Errorf("%d KCL file(s) require formatting", len(changedPaths))
			}
			return nil
		},
		SilenceUsage: true,
	}

	flags := cmd.Flags()
	flags.BoolVar(&o.DryRun, "dry-run", false, "Report files requiring formatting without modifying them.")

	return cmd
}
