// Copyright The KCL Authors. All rights reserved.

package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"kcl-lang.io/cli/pkg/fs"
	"kcl-lang.io/kcl-go/pkg/utils"
	"kcl-lang.io/kpm/pkg/env"
)

const (
	cleanDesc = `This command cleans the kcl build and module cache.
`
	cleanExample = `  # Clean the build and module cache
  kcl clean`
)

// NewCleanCmd returns the clean command.
func NewCleanCmd() *cobra.Command {
	var assumeYes bool
	cmd := &cobra.Command{
		Use:     "clean",
		Short:   "KCL clean tool",
		Long:    cleanDesc,
		Example: cleanExample,
		RunE: func(_ *cobra.Command, args []string) error {
			if len(args) == 0 {
				args = append(args, ".")
			}
			if ok := cmdBox("Are you sure you want to clean the build cache? [y/N]", assumeYes); ok {
				if err := cleanBuildCache(args[0]); err != nil {
					return err
				}
			}
			if ok := cmdBox("Are you sure you want to clean the module cache? [y/N]", assumeYes); ok {
				if err := cleanModCache(); err != nil {
					return err
				}
			}
			return nil
		},
		SilenceUsage: true,
	}

	cmd.Flags().BoolVarP(&assumeYes, "yes", "y", false, "Automatically say yes to prompts")

	return cmd
}

func cmdBox(msg string, assumeYes bool) bool {
	if !assumeYes {
		fmt.Println(msg)
		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Failed to read input:", err)
			return false
		}
		if strings.TrimSpace(strings.ToLower(response)) != "y" {
			fmt.Println("Aborted.")
			return false
		}
	}
	return true
}

func cleanBuildCache(pwd string) error {
	cachePaths := []string{
		filepath.Join(pwd, ".kclvm/cache"),
		filepath.Join(pwd, "__main__/.kclvm/cache"),
	}
	pkgroot, err := utils.FindPkgRoot(pwd)
	if err == nil {
		cachePaths = append(cachePaths, filepath.Join(pkgroot, ".kclvm/cache"), filepath.Join(pkgroot, "__main__/.kclvm/cache"))
	}
	for _, cachePath := range cachePaths {
		if fs.IsDir(cachePath) {
			if err := os.RemoveAll(cachePath); err == nil {
				fmt.Printf("%s removed\n", cachePath)
			} else {
				fmt.Printf("remove %s failed\n", cachePath)
				return err
			}
		}
	}
	return nil
}

func cleanModCache() error {
	modulePath, err := env.GetAbsPkgPath()
	if err != nil {
		return err
	}
	if fs.IsDir(modulePath) {
		if err := os.RemoveAll(modulePath); err == nil {
			fmt.Printf("%s removed\n", modulePath)
		} else {
			fmt.Printf("remove %s failed\n", modulePath)
			return err
		}
	}
	return nil
}
