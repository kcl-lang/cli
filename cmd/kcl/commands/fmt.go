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
				paths, err := kcl.FormatPath(p)
				if err != nil {
					return err
				}
				changedPaths = append(changedPaths, paths...)
			}
			fmt.Println(strings.Join(changedPaths, "\n"))
			return nil
		},
		SilenceUsage: true,
	}

	return cmd
}
