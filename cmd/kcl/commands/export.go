// Copyright The KCL Authors. All rights reserved.

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"kcl-lang.io/kcl-go/pkg/kcl"
)

const (
	exportDesc = `This command exports KCL schemas to source code in another programming language.

For each schema in the input KCL file(s), it generates a corresponding
struct/class with one field per schema attribute, plus a validation method
that re-implements the schema's check: block in the target language — so
values produced outside of KCL can still be checked against the same
constraints the KCL schema declares.
`
	exportExample = `  # Generate Rust structs and print them to stdout
  kcl export --target rust path/to/schema.k

  # Generate Python dataclasses and write them to a file
  kcl export --target python path/to/schema.k --output schema.py

  # Generate Go structs with a specific package name
  kcl export --target go path/to/schema.k --package dto --output dto.go

  # Generate Java classes with a specific package
  kcl export --target java path/to/schema.k --package com.example.dto --output Schema.java

  # Generate TypeScript classes
  kcl export --target typescript path/to/schema.k --output schema.ts

  # Export schemas from multiple KCL files at once
  kcl export --target go path/to/base.k path/to/schema.k --output dto.go`
)

// exportTargets lists the supported --target values, in the order they are
// shown in --help and validated against.
var exportTargets = []string{"rust", "python", "go", "java", "typescript"}

// ExportOptions holds the options for the export command.
type ExportOptions struct {
	// Target is the output language: one of exportTargets.
	Target string
	// Output, if set, writes the generated source to this file instead of
	// stdout.
	Output string
	// Package is the package/namespace name for the generated file.
	// Only applies to the go and java targets: go defaults to "generated"
	// when empty (Go requires a package clause); java omits the package
	// clause entirely when empty (Java's default package).
	Package string
}

// NewExportCmd returns the export command.
func NewExportCmd() *cobra.Command {
	o := ExportOptions{}
	cmd := &cobra.Command{
		Use:     "export",
		Short:   "Export KCL schemas to source code",
		Long:    exportDesc,
		Example: exportExample,
		RunE: func(_ *cobra.Command, args []string) error {
			return doExport(args, &o)
		},
		Args:         cobra.MinimumNArgs(1),
		SilenceUsage: true,
	}

	cmd.Flags().StringVar(&o.Target, "target", "",
		fmt.Sprintf("Target language (currently supported: %s)", joinExportTargets()))
	_ = cmd.MarkFlagRequired("target")
	cmd.Flags().StringVarP(&o.Output, "output", "o", "",
		"Write generated source to this file")
	cmd.Flags().StringVar(&o.Package, "package", "",
		`Go/Java package name for the generated file (target=go/java only; go defaults to "generated", java omits the package clause when unset)`)

	return cmd
}

func joinExportTargets() string {
	result := ""
	for i, target := range exportTargets {
		if i > 0 {
			result += ", "
		}
		result += target
	}
	return result
}

// doExport dispatches to the kcl-go export wrapper for the requested
// --target and writes the result to --output, or stdout when unset.
func doExport(paths []string, o *ExportOptions) error {
	var source string
	var err error

	switch o.Target {
	case "rust":
		source, err = kcl.ExportRust(paths)
	case "python":
		source, err = kcl.ExportPython(paths)
	case "go":
		source, err = kcl.ExportGo(paths, o.Package)
	case "java":
		source, err = kcl.ExportJava(paths, o.Package)
	case "typescript":
		source, err = kcl.ExportTypeScript(paths)
	default:
		return fmt.Errorf("export target %q is not supported; available targets: %s", o.Target, joinExportTargets())
	}
	if err != nil {
		return err
	}

	if o.Output != "" {
		return os.WriteFile(o.Output, []byte(source), 0644)
	}
	fmt.Print(source)
	return nil
}
