package cmd

import (
	"github.com/spf13/pflag"
	"kcl-lang.io/cli/pkg/options"
)

func appendLangFlags(o *options.RunOptions, flags *pflag.FlagSet) {
	flags.StringSliceVarP(&o.PathSelectors, "path_selector", "S", []string{},
		"Specify the path selectors")
	flags.StringVarP(&o.Output, "output", "o", "",
		"Specify the YAML/JSON output file path")
	flags.StringVarP(&o.Git, "git", "", "",
		"Specify the KCL module git url")
	flags.StringVarP(&o.Oci, "oci", "", "",
		"Specify the KCL module oci url")
	flags.StringVarP(&o.Path, "path", "", "",
		"Specify the KCL module local path")
	flags.StringVarP(&o.Tag, "tag", "t", "",
		"Specify the tag for the OCI or Git artifact")
	flags.StringVarP(&o.Commit, "commit", "c", "",
		"Specify the commit for the Git artifact")
	flags.StringVarP(&o.Branch, "branch", "b", "",
		"Specify the branch for the Git artifact")
	flags.StringVar(&o.Format, "format", "yaml",
		"Specify the output format")
	flags.BoolVarP(&o.DisableNone, "disable_none", "n", false,
		"Disable dumping None values")
	flags.BoolVarP(&o.Debug, "debug", "d", false,
		"Run in debug mode")
	flags.BoolVarP(&o.SortKeys, "sort_keys", "k", false,
		"Sort output result keys")
	flags.BoolVarP(&o.ShowHidden, "show_hidden", "H", false,
		"Display hidden attributes")
	appendRunnerFlags(o, flags)
}

func appendRunnerFlags(o *options.RunOptions, flags *pflag.FlagSet) {
	flags.StringArrayVarP(&o.Arguments, "argument", "D", []string{},
		"Specify the top-level argument")
	flags.StringSliceVarP(&o.Settings, "setting", "Y", []string{},
		"Specify the command line setting files")
	flags.StringSliceVarP(&o.Overrides, "overrides", "O", []string{},
		"Specify the configuration override path and value")
	flags.StringSliceVarP(&o.ExternalPackages, "external", "E", []string{},
		"Specify the mapping of package name and path where the package is located")
	flags.BoolVarP(&o.Vendor, "vendor", "V", false,
		"Run in vendor mode")
	flags.BoolVar(&o.NoStyle, "no_style", false,
		"Set to prohibit output of command line waiting styles, including colors, etc.")
	flags.BoolVarP(&o.Quiet, "quiet", "q", false,
		"Set the quiet mode (no output)")
	flags.BoolVarP(&o.StrictRangeCheck, "strict_range_check", "r", false,
		"Do perform strict numeric range checks")
}
