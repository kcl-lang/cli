package cmd

import (
	"github.com/spf13/pflag"
	"kcl-lang.io/cli/pkg/options"
)

func appendLangFlags(o *options.RunOptions, flags *pflag.FlagSet) {
	flags.StringSliceVarP(&o.Arguments, "argument", "D", []string{},
		"Specify the top-level argument")
	flags.StringSliceVarP(&o.Settings, "setting", "Y", []string{},
		"Specify the command line setting files")
	flags.StringSliceVarP(&o.Overrides, "overrides", "O", []string{},
		"Specify the configuration override path and value")
	flags.StringSliceVarP(&o.PathSelectors, "path_selector", "S", []string{},
		"Specify the path selectors")
	flags.StringSliceVarP(&o.ExternalPackages, "external", "E", []string{},
		"Specify the mapping of package name and path where the package is located")
	flags.StringVarP(&o.Output, "output", "o", "",
		"Specify the YAML/JSON output file path")
	flags.StringVarP(&o.Tag, "tag", "t", "",
		"Specify the tag for the OCI or Git artifact")
	flags.StringVar(&o.Format, "format", "yaml",
		"Specify the output format")
	flags.BoolVarP(&o.DisableNone, "disable_none", "n", false,
		"Disable dumping None values")
	flags.BoolVarP(&o.StrictRangeCheck, "strict_range_check", "r", false,
		"Do perform strict numeric range checks")
	flags.BoolVarP(&o.Debug, "debug", "d", false,
		"Run in debug mode")
	flags.BoolVarP(&o.SortKeys, "sort_keys", "k", false,
		"Sort output result keys")
	flags.BoolVarP(&o.Vendor, "vendor", "V", false,
		"Run in vendor mode")
	flags.BoolVar(&o.NoStyle, "no_style", false,
		"Set to prohibit output of command line waiting styles, including colors, etc.")
	flags.BoolVarP(&o.Quiet, "quiet", "q", false,
		"Set the quiet mode (no output)")
}
