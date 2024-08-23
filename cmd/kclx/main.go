// Copyright The KCL Authors. All rights reserved.
//go:build cgo
// +build cgo

package main

import (
	"fmt"
	"os"
	"strings"

	cmd "kcl-lang.io/cli/cmd/kcl/commands"
	_ "kcl-lang.io/kcl-plugin/hello" // Import the hello plugin
	_ "kcl-lang.io/kcl-plugin/http"  // Import the http plugin
)

func main() {
	if err := cmd.NewWithName("kclx").Execute(); err != nil {
		fmt.Fprintln(os.Stderr, strings.TrimLeft(err.Error(), "\n"))
		os.Exit(1)
	}
}
