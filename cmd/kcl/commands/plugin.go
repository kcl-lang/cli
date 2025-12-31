// Copyright The KCL Authors. All rights reserved.
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"kcl-lang.io/cli/pkg/plugin"
)

// isFileOrDir checks if the given path is an existing file or directory
func isFileOrDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir() || info.Mode().IsRegular()
}

// executeRunCmd the run command for the root command.
func executeRunCmd(args []string) {
	cmd := NewRunCmd()
	cmd.SetArgs(args)
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}

func isHelpOrVersionFlag(flag string) bool {
	return flag == "-h" || flag == "--help" || flag == "-v" || flag == "--version"
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
		if strings.HasPrefix(cmdPathPieces[0], "-") && !isHelpOrVersionFlag(cmdPathPieces[0]) {
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

			builtinSubCmdExist := false
			for _, cmd := range foundCmd.Commands() {
				if cmd.Name() == cmdName {
					builtinSubCmdExist = true
					break
				}
			}
			switch cmdName {
			// Don't search for a plugin
			case "help", "completion", cobra.ShellCompRequestCmd, cobra.ShellCompNoDescRequestCmd:
			default:
				if !builtinSubCmdExist {
					// Check if the first argument might be a file/directory before trying plugins
					// If it looks like a KCL file or existing path, execute run command directly
					if len(cmdPathPieces) > 0 && !strings.HasPrefix(cmdPathPieces[0], "-") {
						firstArg := cmdPathPieces[0]
						// If it looks like a KCL file or is an existing file/directory, try to run it
						if strings.HasSuffix(firstArg, ".k") || isFileOrDir(firstArg) {
							executeRunCmd(cmdPathPieces)
							return
						}
					}

					// Try to find and execute a plugin
					if err := plugin.HandlePluginCommand(pluginHandler, cmdPathPieces, false); err != nil {
						fmt.Fprintf(os.Stderr, "Error: %v\n", err)
						os.Exit(1)
					}

					// No plugin found and the argument doesn't look like a file
					// Show helpful error for unknown command
					if len(cmdPathPieces) > 0 && !strings.HasPrefix(cmdPathPieces[0], "-") {
						fmt.Fprintf(os.Stderr, "Error: unknown command \"%s\" for \"%s\"\n", cmdName, cmd.Name())
						fmt.Fprintf(os.Stderr, "Run '%s --help' for available commands.\n", cmd.Name())
						os.Exit(1)
					}
				}
			}
		}
	}
}
