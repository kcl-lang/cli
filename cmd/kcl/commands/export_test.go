// Copyright The KCL Authors. All rights reserved.

package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

// TestNewExportCmdExposesFlags is a structural test (mirroring
// TestNewVetCmdExposesExternalFlag) that ensures the `export` command
// registers its flags correctly without requiring the native service.
func TestNewExportCmdExposesFlags(t *testing.T) {
	cmd := NewExportCmd()

	if cmd.Use != "export" {
		t.Errorf("unexpected command use: %s", cmd.Use)
	}

	targetFlag := cmd.Flags().Lookup("target")
	if targetFlag == nil {
		t.Fatal("expected `export` to expose the `--target` flag")
	}
	required := targetFlag.Annotations[cobra.BashCompOneRequiredFlag]
	if len(required) == 0 || required[0] != "true" {
		t.Fatal("expected `--target` to be marked required")
	}

	outputFlag := cmd.Flags().Lookup("output")
	if outputFlag == nil {
		t.Fatal("expected `export` to expose the `--output` flag")
	}
	if outputFlag.Shorthand != "o" {
		t.Fatalf("expected `--output` shorthand to be `o`, got %q", outputFlag.Shorthand)
	}

	if cmd.Flags().Lookup("package") == nil {
		t.Fatal("expected `export` to expose the `--package` flag")
	}

	if !cmd.SilenceUsage {
		t.Error("expected SilenceUsage to be true")
	}
}

// TestJoinExportTargets covers the pure-Go helper that renders the
// supported --target values for --help and error messages.
func TestJoinExportTargets(t *testing.T) {
	want := "rust, python, go, java, typescript"
	if got := joinExportTargets(); got != want {
		t.Fatalf("joinExportTargets() = %q, want %q", got, want)
	}
}

// TestDoExportUnsupportedTarget exercises doExport's validation branch,
// which returns before ever touching the native service, so it can run
// without a native libkcl binary.
func TestDoExportUnsupportedTarget(t *testing.T) {
	o := &ExportOptions{Target: "undefined"}
	err := doExport([]string{"schema.k"}, o)
	if err == nil {
		t.Fatal("expected an error for an unsupported target")
	}
	want := `export target "undefined" is not supported; available targets: rust, python, go, java, typescript`
	if err.Error() != want {
		t.Fatalf("doExport() error = %q, want %q", err.Error(), want)
	}
}
