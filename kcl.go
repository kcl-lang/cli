// Copyright The KCL Authors. All rights reserved.

package main

import (
	"os"

	"kcl-lang.io/cli/cmd"
)

func main() {
	if err := cmd.New().Execute(); err != nil {
		os.Exit(1)
	}
}
