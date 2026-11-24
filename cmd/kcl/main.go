// Copyright The KCL Authors. All rights reserved.

package main

import (
	"fmt"
	"os"
	"runtime/debug"
	"strings"

	cmd "kcl-lang.io/cli/cmd/kcl/commands"
	"kcl-lang.io/cli/pkg/libkcl"
)

func main() {
	defer func() {
		recovered := recover()
		if recovered == nil {
			return
		}
		if err := libkcl.ErrorFromPanic(recovered); err != nil {
			fmt.Fprintln(os.Stderr, strings.TrimLeft(err.Error(), "\n"))
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "panic: %v\n\n%s", recovered, debug.Stack())
		os.Exit(2)
	}()
	if err := cmd.New().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, strings.TrimLeft(err.Error(), "\n"))
		os.Exit(1)
	}
}
