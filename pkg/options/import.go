// Copyright The KCL Authors. All rights reserved.

package options

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"kcl-lang.io/cli/pkg/fs"
	"kcl-lang.io/cli/pkg/import/crd"
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
	// DataIdentifier, when non-empty, assigns the imported data to this
	// top-level KCL name (e.g. `myData = { ... }`). Only meaningful for
	// data-emitting modes (json, yaml, toml, auto-detected data files);
	// schema modes ignore it because they emit `schema Foo: ...`
	// declarations rather than a free-standing dict literal.
	DataIdentifier string
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

	// Process input files, fetching URLs to temp files if needed
	processedFiles := make([]string, 0, len(o.Files))
	tempFiles := []string{}

	for _, f := range o.Files {
		if fs.IsURL(f) {
			// Fetch URL content to temp file
			tempPath, err := fs.GenTempFileFromURL(f)
			if err != nil {
				for _, tf := range tempFiles {
					os.Remove(tf)
				}
				return err
			}
			tempFiles = append(tempFiles, tempPath)
			processedFiles = append(processedFiles, tempPath)
		} else {
			processedFiles = append(processedFiles, f)
		}
	}

	// Ensure temp files are cleaned up when function returns
	defer func() {
		for _, tf := range tempFiles {
			os.Remove(tf)
		}
	}()

	files, err := fs.ExpandInputFiles(processedFiles, o.Recursive)
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
			// Convert CLI options to generator.GenOpts
			opts.Spec = p
			if o.Output != "" {
				opts.Target = o.Output
			} else {
				opts.Target = "."
			}
			opts.ValidateSpec = !o.SkipValidation
			opts.ModelPackage = o.ModelPackage
			// Set default configurations
			if err := opts.EnsureDefaults(); err != nil {
				return err
			}
			var specs []string
			// When the spec is a crd, get OpenAPI spec file from it
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
			// Generate specs to KCL files
			for _, spec := range specs {
				opts.Spec = spec
				if err := generator.Generate(opts); err != nil {
					logger.GetLogger().Error(err)
				}
			}
			// Group by the api version
			if mode == Crd {
				err := crd.GroupByKclFiles(opts.Target, opts.ModelPackage)
				if err != nil {
					return err
				}
			}
		}
		return nil
	default:
		return fmt.Errorf("invalid mode: %s", o.Mode)
	}

	if o.Output == "-" {
		for _, p := range files {
			if err := o.writeGenKcl(os.Stdout, p, opts); err != nil {
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
			if err := o.writeGenKcl(outputWriter, p, opts); err != nil {
				return err
			}
		}
	}
	return nil
}

// writeGenKcl generates KCL source for the given input file and writes it to
// w. When DataIdentifier is set and the resolved mode emits a free-standing
// data literal (JSON / YAML / TOML), the literal is wrapped in
// `<DataIdentifier> = { ... }` so the result can be referenced from other
// KCL files. Schema-producing modes (jsonschema, gostruct, ...) bypass the
// wrapping because their output contains `schema ...:` declarations that
// cannot be assigned to a name.
func (o *ImportOptions) writeGenKcl(w io.Writer, input string, opts *gen.GenKclOptions) error {
	wrap := o.shouldWrapData()
	if !wrap {
		return gen.GenKcl(w, input, nil, opts)
	}
	var buf bytes.Buffer
	if err := gen.GenKcl(&buf, input, nil, opts); err != nil {
		return err
	}
	wrapped, err := wrapDataLiteral(buf.String(), o.DataIdentifier)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(wrapped))
	return err
}

// shouldWrapData reports whether the resolved mode produces a data literal
// that can be assigned to a name. Schema modes (jsonschema / gostruct /
// terraformschema / crd / openapi) return false because their output
// contains named `schema` declarations rather than a free-standing dict.
func (o *ImportOptions) shouldWrapData() bool {
	if o.DataIdentifier == "" {
		return false
	}
	switch strings.ToLower(o.Mode) {
	case Json, Yaml, Toml:
		return true
	case Auto:
		// Auto mode is resolved by the gen package at runtime; we cannot
		// know the resolved mode until generation happens. Inspect every
		// input file extension to decide.
		for _, f := range o.Files {
			switch strings.ToLower(filepath.Ext(f)) {
			case ".json", ".yaml", ".yml", ".toml":
				return true
			}
		}
		return false
	default:
		return false
	}
}

// wrapDataLiteral turns the free-standing dict literal that gen.GenKcl emits
// for data modes into `<name> = { ... }`. The literal can be preceded by a
// generated header comment block; that block is kept verbatim and the
// literal is re-indented by four spaces before being assigned to `name`.
func wrapDataLiteral(src, name string) (string, error) {
	openIdx := strings.Index(src, "{")
	if openIdx < 0 {
		return "", fmt.Errorf("kcl import: data literal not found in generated source")
	}
	closeIdx := strings.LastIndex(src, "}")
	if closeIdx < openIdx {
		return "", fmt.Errorf("kcl import: malformed data literal in generated source")
	}
	header := src[:openIdx]
	body := src[openIdx+1 : closeIdx]
	// Trim surrounding blank lines so the wrapped block does not start
	// with a blank line right after `{` or end with a trailing blank line
	// before `}`.
	body = strings.Trim(body, "\n")
	// Re-indent body by four spaces, preserving the existing newlines.
	var indented bytes.Buffer
	lines := strings.Split(body, "\n")
	for i, line := range lines {
		if i > 0 {
			indented.WriteString("\n")
		}
		if line == "" {
			continue
		}
		indented.WriteString("    ")
		indented.WriteString(line)
	}
	var out bytes.Buffer
	out.WriteString(header)
	out.WriteString(name)
	out.WriteString(" = {\n")
	out.WriteString(indented.String())
	out.WriteString("\n}\n")
	return out.String(), nil
}
