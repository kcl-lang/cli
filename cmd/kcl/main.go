// Copyright The KCL Authors. All rights reserved.

package main

import (
	"fmt"
	"os"
	"strings"

	cmd "kcl-lang.io/cli/cmd/kcl/commands"
)

func main() {
	if err := cmd.New().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, strings.TrimLeft(err.Error(), "\n"))
		os.Exit(1)
	}
}
