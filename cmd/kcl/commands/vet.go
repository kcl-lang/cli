// Copyright The KCL Authors. All rights reserved.

package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"kcl-lang.io/cli/pkg/fs"
	"kcl-lang.io/kcl-go/pkg/kcl"
	"kcl-lang.io/kcl-go/pkg/spec/gpyrpc"
	"kcl-lang.io/kcl-go/pkg/tools/validate"
)

const (
	vetDesc = `This command validates the data file using the kcl code.
`
	vetExample = `  # Validate the JSON data using the kcl code
  kcl vet data.json code.k

  # Validate the YAML data using the kcl code
  kcl vet data.yaml code.k --format yaml

  # Validate the JSON data using the kcl code with the schema name
  kcl vet data.json code.k -s Schema`
)

// NewVetCmd returns the vet command.
func NewVetCmd() *cobra.Command {
	o := validate.ValidateOptions{}
	cmd := &cobra.Command{
		Use:     "vet",
		Short:   "KCL validation tool",
		Long:    vetDesc,
		Example: vetExample,
		RunE: func(_ *cobra.Command, args []string) error {
			dataFile := args[0]
			codeFile := args[1]
			return doValidate(dataFile, codeFile, &o)
		},
		SilenceUsage: true,
	}

	// Two positional arguments <data_file> <kcl_file>
	cmd.Args = cobra.ExactArgs(2)
	cmd.Flags().StringVarP(&o.Schema, "schema", "s", "",
		"Specify the validate schema.")
	cmd.Flags().StringVarP(&o.Schema, "attribute_name", "a", "",
		"Specify the validate config attribute name.")
	cmd.Flags().StringVar(&o.Format, "format", "",
		"Specify the validate data format. e.g., yaml, json. Default is json")

	return cmd
}

func doValidate(dataFile, codeFile string, o *validate.ValidateOptions) error {
	var ok bool
	if dataFile == "-" {
		// Read data from stdin
		input, err := io.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
		code, err := os.ReadFile(codeFile)
		if err != nil {
			return err
		}
		ok, err = validate.ValidateCode(string(input), string(code), o)
		if err != nil {
			return err
		}
	} else {
		// Read data from files
		dataFiles, err := fs.ExpandInputFiles([]string{dataFile}, false)
		if err != nil {
			return err
		}
		for _, dataFile := range dataFiles {
			ok, err = validateFile(dataFile, codeFile, o)
			if err != nil {
				return err
			}
		}
	}
	if ok {
		fmt.Println("Validate success!")
	}
	return nil
}

func validateFile(dataFile, codeFile string, opts *validate.ValidateOptions) (ok bool, err error) {
	if opts == nil {
		opts = &validate.ValidateOptions{}
	}
	svc := kcl.Service()
	resp, err := svc.ValidateCode(&gpyrpc.ValidateCodeArgs{
		Datafile:      dataFile,
		File:          codeFile,
		Schema:        opts.Schema,
		AttributeName: opts.AttributeName,
		Format:        opts.Format,
	})
	if err != nil {
		return false, err
	}
	var e error = nil
	if resp.ErrMessage != "" {
		e = errors.New(resp.ErrMessage)
	}
	return resp.Success, e
}
