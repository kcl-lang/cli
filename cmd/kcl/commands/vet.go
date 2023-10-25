// Copyright The KCL Authors. All rights reserved.

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"kcl-lang.io/kcl-go/pkg/tools/validate"
)

const (
	vetDesc = `
This command validates the data file using the kcl code.
`
	vetExample = `  # Validate the JSON data using the kcl code
  kcl vet data.json code.k
`
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
			data, err := os.ReadFile(args[0])
			if err != nil {
				return err
			}
			code, err := os.ReadFile(args[1])
			if err != nil {
				return err
			}
			ok, err := validate.ValidateCode(string(data), string(code), &o)
			if err != nil {
				return err
			}
			if ok {
				fmt.Println("Validate success")
			}
			return nil
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
