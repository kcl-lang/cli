// Copyright The KCL Authors. All rights reserved.

package options

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"kcl-lang.io/cli/pkg/fs"
	"kcl-lang.io/kcl-go/pkg/logger"
	"kcl-lang.io/kcl-go/pkg/tools/gen"
	crdGen "kcl-lang.io/kcl-openapi/pkg/kube_resource/generator"
	"kcl-lang.io/kcl-openapi/pkg/swagger/generator"
)

type ImportOptions struct {
	Mode           string
	Files          []string
	Output         string
	Force          bool
	SkipValidation bool
	ModelPackage   string
	Recursive      bool
}

// NewImportOptions returns a new instance of ImportOptions with default values.
func NewImportOptions() *ImportOptions {
	return &ImportOptions{
		Mode: Auto,
	}
}

// Run runs the kcl import command with options.
func (o *ImportOptions) Run() error {
	opts := &gen.GenKclOptions{}
	mode := strings.ToLower(o.Mode)
	files, err := fs.ExpandInputFiles(o.Files, o.Recursive)
	if err != nil {
		return err
	}
	switch mode {
	case Json:
		opts.Mode = gen.ModeJson
	case Yaml:
		opts.Mode = gen.ModeYaml
	case Toml:
		opts.Mode = gen.ModeToml
	case GoStruct:
		opts.Mode = gen.ModeGoStruct
	case JsonSchema:
		opts.Mode = gen.ModeJsonSchema
	case TerraformSchema:
		opts.Mode = gen.ModeTerraformSchema
	case Auto:
		opts.Mode = gen.ModeAuto
	case Crd, OpenAPI:
		for _, p := range files {
			opts := new(generator.GenOpts)
			// cli opts to generator.GenOpts
			opts.Spec = p
			if o.Output != "" {
				opts.Target = o.Output
			} else {
				opts.Target = "."
			}
			opts.ValidateSpec = !o.SkipValidation
			opts.ModelPackage = o.ModelPackage
			// set default configurations
			if err := opts.EnsureDefaults(); err != nil {
				return err
			}
			var specs []string
			// when the spec is a crd, get openapi spec file from it
			if mode == Crd {
				specs, err = crdGen.GetSpecs(&crdGen.GenOpts{
					Spec: opts.Spec,
				})
				if err != nil {
					logger.GetLogger().Error(err)
				}
				// do not run validate spec on spec file generated from crd
				opts.ValidateSpec = false
			} else {
				specs = []string{opts.Spec}
			}
			for _, spec := range specs {
				opts.Spec = spec
				if err := generator.Generate(opts); err != nil {
					logger.GetLogger().Error(err)
				}
			}
		}
		return nil
	default:
		return fmt.Errorf("invalid mode: %s", o.Mode)
	}

	if o.Output == "-" {
		for _, p := range files {
			err := gen.GenKcl(os.Stdout, p, nil, opts)
			if err != nil {
				return err
			}
		}
	} else {
		for _, p := range files {
			outputFile := o.Output
			if outputFile == "" {
				filenameWithExtension := filepath.Base(p)
				filename := strings.TrimSuffix(filenameWithExtension, filepath.Ext(filenameWithExtension))
				outputFile = fmt.Sprintf("%s.k", filename)
			}
			if _, err := os.Stat(outputFile); err == nil && !o.Force {
				return fmt.Errorf("output file already exist, use --force to overwrite: %s", outputFile)
			}
			outputWriter, err := os.Create(outputFile)
			if err != nil {
				return fmt.Errorf("failed to create output file: %s", outputFile)
			}
			err = gen.GenKcl(outputWriter, p, nil, opts)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
