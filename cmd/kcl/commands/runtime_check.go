// Copyright The KCL Authors. All rights reserved.

package cmd

import (
	"github.com/spf13/cobra"

	"kcl-lang.io/cli/pkg/libkcl"
)

// ensureRuntime checks the FreeBSD libkcl runtime before a command loads
// the native service. Commands that do not evaluate KCL skip the check so
// `kcl version` and `kcl --help` keep working.
func ensureRuntime(c *cobra.Command) error {
	if !commandNeedsLibKCL(c) {
		return nil
	}
	return libkcl.Prepare()
}

func commandNeedsLibKCL(c *cobra.Command) bool {
	if c == nil {
		return false
	}
	switch c.Name() {
	case "run", "lint", "fmt", "test", "vet", "server", "generate":
		return true
	default:
		return false
	}
}
