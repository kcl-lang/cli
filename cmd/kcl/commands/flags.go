package cmd

import (
	"github.com/spf13/pflag"
	"kcl-lang.io/cli/pkg/options"
)

func appendLangFlags(o *options.RunOptions, flags *pflag.FlagSet) {
	// Use StringArrayVar (not StringSlice) for every repeatable flag below:
	// StringSlice splits the value on commas via encoding/csv, which breaks
	// legitimate inputs that contain commas inside the value, e.g.
	// `-O person.ids=[1,2]`, `-O 'person.name="Alice,Bob"'`, and any path
	// selector or settings file path. Each repeat of the flag should become
	// one distinct element of the slice, regardless of its contents.
	// See https://github.com/kcl-lang/cli/issues/232.
	flags.StringArrayVarP(&o.PathSelectors, "path_selector", "S", []string{},
		"Specify the path selectors")
	flags.StringVarP(&o.Output, "output", "o", "",
		"Specify the YAML/JSON output file path")
	flags.StringVarP(&o.Git, "git", "", "",
		"Specify the KCL module git url")
	flags.StringVarP(&o.Oci, "oci", "", "",
		"Specify the KCL module oci url")
	flags.StringVarP(&o.Tag, "tag", "t", "",
		"Specify the tag for the OCI or Git artifact")
	flags.StringVarP(&o.Commit, "commit", "c", "",
		"Specify the commit for the Git artifact")
	flags.StringVarP(&o.Branch, "branch", "b", "",
		"Specify the branch for the Git artifact")
	flags.StringVar(&o.Format, "format", "yaml",
		"Specify the output format (yaml, json, toml, xml). When xml is selected, schema attributes decorated with `@info(type=\"attr\")` are rendered as `name=\"value\"` attributes on the parent element instead of as child elements.")
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
	flags.StringArrayVarP(&o.Settings, "setting", "Y", []string{},
		"Specify the command line setting files")
	flags.StringArrayVarP(&o.Overrides, "overrides", "O", []string{},
		"Specify the configuration override path and value")
	flags.StringArrayVarP(&o.ExternalPackages, "external", "E", []string{},
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
