// Copyright The KCL Authors. All rights reserved.

package cmd

import (
	"github.com/spf13/cobra"
	"kcl-lang.io/cli/pkg/options"
)

const (
	importDesc = `This command converts other formats to KCL file.

Supported conversion modes:
- json:            convert JSON data to KCL data
- yaml:            convert YAML data to KCL data
- toml:            convert TOML data to KCL data
- gostruct:        convert Go struct to KCL schema
- jsonschema:      convert JSON schema to KCL schema
- terraformschema: convert Terraform schema to KCL schema
- openapi:         convert OpenAPI spec to KCL schema
- crd:             convert Kubernetes CRD to KCL schema
- auto:            automatically detect the input format
`
	importExample = `  # Generate KCL models from OpenAPI spec
  kcl import -m openapi swagger.json

  # Generate KCL models from Kubernetes CRD
  kcl import -m crd crd.yaml

  # Generate KCL models from JSON
  kcl import data.json

  # Generate KCL models from YAML
  kcl import data.yaml

  # Generate KCL models from TOML
  kcl import data.toml

  # Generate KCL models from JSON Schema
  kcl import -m jsonschema schema.json

  # Generate KCL models from Terraform provider schema
  kcl import -m terraformschema schema.json

  # Generate KCL models from Go structs
  kcl import -m gostruct schema.go`
)

// NewImportCmd returns the import command.
func NewImportCmd() *cobra.Command {
	o := options.NewImportOptions()
	cmd := &cobra.Command{
		Use:     "import",
		Short:   "KCL import tool",
		Long:    importDesc,
		Example: importExample,
		RunE: func(_ *cobra.Command, args []string) error {
			o.Files = args
			return o.Run()
		},
		SilenceUsage: true,
	}

	cmd.Args = cobra.MinimumNArgs(1)
	cmd.Flags().StringVarP(&o.Mode, "mode", "m", "auto",
		"Specify the import mode. Default is mode")
	cmd.Flags().StringVarP(&o.Output, "output", "o", "",
		"Specify the output file path")
	cmd.Flags().BoolVarP(&o.Force, "force", "f", false,
		"Force overwrite output file")
	cmd.Flags().BoolVarP(&o.SkipValidation, "skip-validation", "s", false,
		"Skips validation of spec prior to generation")
	cmd.Flags().StringVarP(&o.ModelPackage, "package", "p", "models",
		"The package to save the models. Default is models")

	return cmd
}
