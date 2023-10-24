// Copyright The KCL Authors. All rights reserved.
//
// #### Language & Tooling Commands
//
// ```
// kcl
//
//	run           compile kcl package from a url or filepath
//	build         build the kcl package
//	check         check the current package, but don't build target files
//	doc           documentation tool
//	fmt           format tool
//	lint          lint tool
//	test          unit/integration/benchmark test tool
//	lsp           language server tool
//	clean         remove object files and cached files
//
// ```
//
// #### Package & Registry Related Commands (mod and registry workspace)
//
// ```
// kcl
//
//	mod init         initialize new module in current directory
//	mod search       search a command from regisry
//	mod add          add new dependency
//	mod remove       remove dependency
//	mod update       update dependency
//	mod pkg          package a kcl package into tar
//	mod metadata     output the resolved dependencies of a package
//	mod push         push kcl package to OCI registry.
//	mod pull         pull kcl package from OCI registry.
//	registry login   login to a registry
//	registry logout  logout from a registry
//
// ```
//
// #### Integration Commands
//
// ```
// kcl
//
//	import     migration other data and schema to kcl e.g., openapi, jsonschema, json, yaml
//	export     convert kcl schema to other schema e.g., openapi
//
// ```
//
// #### Plugin Commands (plugin workspace)
//
// ```
// kcl
//
//	plugin install     install one or more kcl command plugins
//	plugin list        list installed command plugins
//	plugin uninstall   uninstall one or more command plugins
//	plugin update      update one or more command plugins
//
// ```
//
// #### Version and Help Commands
//
// ```
// kcl
//
//	help, h   Shows a list of commands or help for one command
//	version Shows the command version
//
// ```
// #### Alias
//
// ```
// alias kcl="kcl run"
// alias kpm="kcl mod"
// ```
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"kcl-lang.io/cli/pkg/plugin"
	"kcl-lang.io/cli/pkg/version"
)

const (
	cmdName  = "kcl"
	rootDesc = `The KCL Command Line Interface (CLI).

KCL is an open-source, constraint-based record and functional language that
enhances the writing of complex configurations, including those for cloud-native
scenarios. The KCL website: https://kcl-lang.io
`
)

// New creates a new cobra client
func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:          cmdName,
		Short:        "The KCL Command Line Interface (CLI).",
		Long:         rootDesc,
		SilenceUsage: true,
		Version:      version.GetVersionString(),
	}
	cmd.AddCommand(NewVersionCmd())
	cmd.AddCommand(NewRunCmd())
	cmd.AddCommand(NewModCmd())
	cmd.AddCommand(NewRegistryCmd())

	bootstrapCmdPlugin(cmd, plugin.NewDefaultPluginHandler([]string{cmdName}))

	return cmd
}

func bootstrapCmdPlugin(cmd *cobra.Command, pluginHandler plugin.PluginHandler) {
	if pluginHandler == nil {
		return
	}
	if len(os.Args) > 1 {
		cmdPathPieces := os.Args[1:]

		// only look for suitable extension executables if
		// the specified command does not already exist
		if foundCmd, _, err := cmd.Find(cmdPathPieces); err != nil {
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
					// Run the run command for the root command. alias kcl="kcl run"
					cmd := NewRunCmd()
					cmd.SetArgs(cmdPathPieces)
					if err := cmd.Execute(); err != nil {
						os.Exit(1)
					}
					os.Exit(0)
				}
			}
		}
	}
}
