// Copyright The KCL Authors. All rights reserved.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"kcl-lang.io/kcl-go/pkg/tools/gen"
)

const (
	docDesc = `This command shows documentation for KCL modules or symbols.
`
	docExample = `  # Generate document for current package
  kcl doc generate`

	docGenDesc = `This command generates documents for KCL modules.
`
	docGenExample = `  # Generate Markdown document for current package
  kcl doc generate

  # Generate Markdown document for current package and escape the HTML symbols | to \|, \n to <br>, etc.
  kcl doc generate --escape-html

  # Generate Html document for current package
  kcl doc generate --format html

  # Generate Markdown document for specific package
  kcl doc generate --file-path <package path>

  # Generate Markdown document for specific package to a <target directory>
  kcl doc generate --file-path <package path> --target <target directory>`
)

// NewDocCmd returns the doc command.
func NewDocCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "doc",
		Short:        "KCL document tool",
		Long:         docDesc,
		Example:      docExample,
		SilenceUsage: true,
		Aliases:      []string{"d"},
	}

	cmd.AddCommand(NewDocGenerateCmd())

	return cmd
}

// NewDocGenerateCmd returns the doc generate command.
func NewDocGenerateCmd() *cobra.Command {
	o := gen.GenOpts{}
	cmd := &cobra.Command{
		Use:     "generate",
		Short:   "Generates documents from code and examples",
		Long:    docGenDesc,
		Example: docGenExample,
		RunE: func(*cobra.Command, []string) error {
			genContext, err := o.ValidateComplete()
			if err != nil {
				fmt.Println(fmt.Errorf("generate failed: %s", err))
				return err
			}

			err = genContext.GenDoc()
			if err != nil {
				fmt.Println(fmt.Errorf("generate failed: %s", err))
				return err
			} else {
				fmt.Printf("Generate Complete! Check generated docs in %s\n", genContext.Target)
				return nil
			}
		},
		Aliases:      []string{"gen", "g"},
		SilenceUsage: true,
	}

	cmd.Flags().StringVar(&o.Path, "file-path", "",
		`Relative or absolute path to the KCL package root when running kcl-doc command from
outside of the KCL package root directory.
If not specified, the current work directory will be used as the KCL package root.`)
	cmd.Flags().StringVar(&o.Format, "format", string(gen.Markdown),
		`The document format to generate. Supported values: markdown, html, openapi.`)
	cmd.Flags().StringVar(&o.Target, "target", "",
		`If not specified, the current work directory will be used. A docs/ folder will be created under the target directory.`)
	cmd.Flags().StringVar(&o.TemplateDir, "template", "",
		`The template directory based on the KCL package root. If not specified, the built-in templates will be used.`)
	cmd.Flags().BoolVar(&o.IgnoreDeprecated, "ignore-deprecated", false,
		`Do not generate documentation for deprecated schemas.`)
	cmd.Flags().BoolVar(&o.EscapeHtml, "escape-html", false,
		`Whether to escape html symbols when the output format is markdown. Always scape when the output format is html. Default to false.`)

	return cmd
}
