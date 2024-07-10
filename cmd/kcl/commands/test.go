// Copyright The KCL Authors. All rights reserved.

package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/acarl005/stripansi"
	"github.com/spf13/cobra"
	"kcl-lang.io/cli/pkg/options"
	kcl "kcl-lang.io/kcl-go"
	"kcl-lang.io/kcl-go/pkg/tools/testing"
)

const (
	testDesc = `This command automates testing the packages named by the import paths.

'KCL test' re-compiles each package along with any files with names matching
the file pattern "*_test.k". These additional files can contain test functions
that starts with "test_*".
`
	testExample = `  # Test whole current package recursively
  kcl test ./...

  # Test package named 'pkg'
  kcl test pkg

  # Test with the fail fast mode.
  kcl test ./... --fail-fast

  # Test with the regex expression filter 'test_func'
  kcl test ./... --run test_func`
)

// NewTestCmd returns the test command.
func NewTestCmd() *cobra.Command {
	o := new(kcl.TestOptions)
	runOpts := options.NewRunOptions()
	cmd := &cobra.Command{
		Use:     "test",
		Short:   "KCL test tool",
		Long:    testDesc,
		Example: testExample,
		RunE: func(_ *cobra.Command, args []string) error {
			if len(args) == 0 {
				args = append(args, ".")
			}
			o.PkgList = args
			return test(o, runOpts)
		},
		SilenceErrors: true,
		SilenceUsage:  true,
		Aliases:       []string{"t"},
	}

	flags := cmd.Flags()
	flags.BoolVar(&o.FailFast, "fail-fast", false,
		"Exist when meet the first fail test case in the test process.")
	flags.StringVar(&o.RunRegRxp, "run", "",
		"If specified, only run tests containing this string in their names.")
	appendRunnerFlags(runOpts, flags)

	return cmd
}

func test(o *kcl.TestOptions, runOpts *options.RunOptions) error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	depsOpt, err := options.LoadDepsFrom(pwd, runOpts.Quiet)
	if err != nil {
		return err
	}
	result, err := kcl.Test(
		o,
		*options.CompileOptionFromCli(runOpts).Option,
		*depsOpt,
	)
	if err != nil {
		if runOpts.NoStyle {
			err = errors.New(stripansi.Strip(err.Error()))
		}
		return err
	}
	if len(result.Info) == 0 {
		fmt.Println("no test files")
		return nil
	} else {
		reporter := testing.DefaultReporter(os.Stdout)
		err := reporter.Report(&result)
		if err != nil {
			return err
		}
		for _, info := range result.Info {
			if info.ErrMessage != "" && !info.Skip() {
				return errors.New("")
			}
		}
		return nil
	}
}
