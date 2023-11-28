// Copyright The KCL Authors. All rights reserved.
//
// #### Language & Tooling Commands
//
// ```
// kcl
//
//		run           compile kcl package from a url or filepath
//		build         build the kcl package (Not yet implemented)
//		check         check the current package, but don't build target files (Not yet implemented)
//		doc           documentation tool
//		fmt           format tool
//		lint          lint tool
//		vet           vet tool
//		test          unit/integration/benchmark test tool
//	 	deps          dependency analysis, providing dependency diagrams for KCL modules and packages (Not yet implemented)
//	 	server        run a KCL server to provider REST APIs for other applications
//		clean         remove object files and cached files
//		play          open the playground
//
// ```
//
// #### Package & Registry Related Commands (mod and registry workspace)
//
// ```
// kcl
//
//	mod init         initialize new module in current directory
//	mod search       search a command from registry
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
//	export     convert kcl schema to other schema e.g., openapi (Not yet implemented)
//
// ```
//
// #### Plugin Commands (Not yet implemented)
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
//	version   Shows the command version
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
	// Language commands
	cmd.AddCommand(NewRunCmd())
	// Tool commands
	cmd.AddCommand(NewLintCmd())
	cmd.AddCommand(NewDocCmd())
	cmd.AddCommand(NewFmtCmd())
	cmd.AddCommand(NewTestCmd())
	cmd.AddCommand(NewVetCmd())
	cmd.AddCommand(NewCleanCmd())
	cmd.AddCommand(NewImportCmd())
	cmd.AddCommand(NewPlayCmd())
	// Module & Registry commands
	cmd.AddCommand(NewModCmd())
	cmd.AddCommand(NewRegistryCmd())
	// Server commands
	cmd.AddCommand(NewServerCmd())
	// Version & Help commands
	cmd.AddCommand(NewVersionCmd())
	cmd.SetHelpCommand(&cobra.Command{}) // Disable the help command
	// Plugin commands
	bootstrapCmdPlugin(cmd, plugin.NewDefaultPluginHandler([]string{cmdName}))

	return cmd
}
