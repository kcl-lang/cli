// Copyright The KCL Authors. All rights reserved.

package cmd

import (
	"github.com/spf13/cobra"
	runtime "kcl-lang.io/kcl-go/pkg/kclvm_runtime"
	"kcl-lang.io/kcl-go/pkg/service"
)

const (
	serverDesc = `
This command runs a kcl server with multiple REST APIs. See https://kcl-lang.io/docs/reference/xlang-api/rest-api for more information.
`
	serverExample = `  # Run a kcl server
  kcl server

  # Ping the server
  curl -X POST http://127.0.0.1:2021/api:protorpc/BuiltinService.Ping --data '{}'

  # List the API list
  curl -X POST http://127.0.0.1:2021/api:protorpc/BuiltinService.ListMethod --data '{}'

  # Use the Run API
  curl -X POST http://127.0.0.1:2021/api:protorpc/KclvmService.ExecProgram -H  "accept: application/json" --data '{"k_filename_list": ["main.k"]}'
  `
)

var (
	http         string
	processCount int = 4
)

// NewServerCmd returns the server command.
func NewServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "server",
		Short:   "Run a KCL server",
		Long:    serverDesc,
		Example: serverExample,
		RunE: func(_ *cobra.Command, args []string) error {
			runtime.InitRuntime(processCount)
			return service.RunRestServer(http)
		},
		SilenceUsage: true,
	}

	cmd.Flags().StringVar(&http, "http", ":2021", "set listen address")
	cmd.Flags().IntVarP(&processCount, "proc", "n", 4, "set max process count")

	return cmd
}
