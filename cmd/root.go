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
package cmd

import (
	"github.com/spf13/cobra"
)

const rootCmdShortUsage = "The KCL Command Line Interface (CLI)."
const rootCmdLongUsage = `The KCL Command Line Interface (CLI).

KCL is an open-source, constraint-based record and functional language that
enhances the writing of complex configurations, including those for cloud-native
scenarios. The KCL website: https://kcl-lang.io
`

// New creates a new cobra client
func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "kcl",
		Short:        rootCmdShortUsage,
		Long:         rootCmdLongUsage,
		SilenceUsage: true,
	}
	cmd.AddCommand(NewVersionCmd())
	cmd.AddCommand(NewRunCmd())
	cmd.SetHelpCommand(&cobra.Command{}) // Disable the help command
	return cmd
}
