// Copyright The KCL Authors. All rights reserved.
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"kcl-lang.io/cli/pkg/plugin"
)

// executeRunCmd the run command for the root command.
func executeRunCmd(args []string) {
	cmd := NewRunCmd()
	cmd.SetArgs(args)
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}

func bootstrapCmdPlugin(cmd *cobra.Command, pluginHandler plugin.PluginHandler) {
	if pluginHandler == nil {
		return
	}
	if len(os.Args) > 1 {
		cmdPathPieces := os.Args[1:]

		// only look for suitable extension executables if
		// the specified command does not already exist
		// flags cannot be placed before plugin name
		if strings.HasPrefix(cmdPathPieces[0], "-") {
			executeRunCmd(cmdPathPieces)
		} else if foundCmd, _, err := cmd.Find(cmdPathPieces); err != nil {
			// Also check the commands that will be added by Cobra.
			// These commands are only added once rootCmd.Execute() is called, so we
			// need to check them explicitly here.
			var cmdName string // first "non-flag" arguments
			for _, arg := range cmdPathPieces {
				if !strings.HasPrefix(arg, "-") {
					cmdName = arg
					break
				}
			}

			builtinSubcmdExist := false
			for _, subcmd := range foundCmd.Commands() {
				if subcmd.Name() == cmdName {
					builtinSubcmdExist = true
					break
				}
			}
			switch cmdName {
			// Don't search for a plugin
			case "help", cobra.ShellCompRequestCmd, cobra.ShellCompNoDescRequestCmd:
			default:
				if !builtinSubcmdExist {
					if err := plugin.HandlePluginCommand(pluginHandler, cmdPathPieces, false); err != nil {
						fmt.Fprintf(os.Stderr, "Error: %v\n", err)
						os.Exit(1)
					}
					executeRunCmd(cmdPathPieces)
				}
			}
		}
	}
}
