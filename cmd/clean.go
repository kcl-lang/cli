// Copyright The KCL Authors. All rights reserved.

package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"kcl-lang.io/kcl-go/pkg/utils"
)

const (
	cleanDesc = `
This command cleans the kcl build cache.
`
	cleanExample = `  # Clean the build cache
  kcl clean
`
)

// NewCleanCmd returns the clean command.
func NewCleanCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "clean",
		Short:   "KCL clean tool",
		Long:    cleanDesc,
		Example: cleanExample,
		RunE: func(_ *cobra.Command, args []string) error {
			if len(args) == 0 {
				args = append(args, ".")
			}
			pkgroot, err := utils.FindPkgRoot(args[0])
			if err != nil {
				fmt.Println("no cache found")
				return err
			}
			cachePaths := []string{
				filepath.Join(pkgroot, ".kclvm/cache"),
				filepath.Join(pkgroot, "__main__/.kclvm/cache"),
				filepath.Join(args[0], ".kclvm/cache"),
				filepath.Join(args[0], "__main__/.kclvm/cache"),
			}
			for _, cachePath := range cachePaths {
				if isDir(cachePath) {
					if err := os.RemoveAll(cachePath); err == nil {
						fmt.Printf("%s removed\n", cachePath)
					} else {
						fmt.Printf("remove %s failed\n", cachePath)
						return err
					}
				}
			}
			return nil
		},
		SilenceUsage: true,
	}

	return cmd
}

func isDir(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}
