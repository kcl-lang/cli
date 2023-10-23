package options

import (
	"io"
	"os"

	"kcl-lang.io/kcl-go/pkg/kcl"
	"kcl-lang.io/kpm/pkg/api"
	"kcl-lang.io/kpm/pkg/client"
	"kcl-lang.io/kpm/pkg/opt"
	"kcl-lang.io/kpm/pkg/reporter"
	"kcl-lang.io/kpm/pkg/runner"
)

type FormatType string

const (
	Json FormatType = "json"
	Yaml FormatType = "yaml"
)

type RunOptions struct {
	Entries          []string   // List of entry points
	Output           string     // Output path
	Settings         []string   // List of settings
	Arguments        []string   // List of arguments
	Overrides        []string   // List of overrides
	PathSelectors    []string   // List of path selectors
	ExternalPackages []string   // List of external packages
	NoStyle          bool       // Disable style check
	Vendor           bool       // Use vendor directory
	SortKeys         bool       // Sort keys
	DisableNone      bool       // Disable empty options
	StrictRangeCheck bool       // Strict range check
	Tag              string     // Tag
	CheckOnly        bool       // Check only mode
	Format           FormatType // Formatting type
	Writer           io.Writer  // Writer
}

func NewRunOptions() *RunOptions {
	return &RunOptions{
		Writer: os.Stdout,
		Format: Yaml,
	}
}

func (o *RunOptions) Run() error {
	opts := CompileOptionFromCli(o)
	cli, err := client.NewKpmClient()
	if err != nil {
		reporter.Fatal(err)
	}
	entry, errEvent := runner.FindRunEntryFrom(opts.Entries())
	if errEvent != nil {
		return errEvent
	}
	var result *kcl.KCLResultList
	// kcl compiles the current package under '$pwd'.
	if entry.IsEmpty() {
		pwd, err := os.Getwd()
		opts.SetPkgPath(pwd)
		if err != nil {
			return reporter.NewErrorEvent(
				reporter.Bug, err, "internal bugs, please contact us to fix it.",
			)
		}
		result, err = cli.CompileWithOpts(opts)
		if err != nil {
			return err
		}
	} else {
		var err error
		// kcl compiles the package from the local file system.
		if entry.IsLocalFile() || entry.IsLocalFileWithKclMod() {
			if entry.IsLocalFile() {
				// If there is only kcl file without kcl package,
				result, err = api.RunWithOpt(opts)
			} else {
				// Else compile the kcl package.
				result, err = cli.CompileWithOpts(opts)
			}
		} else if entry.IsTar() {
			// kcl compiles the package from the kcl package tar.
			result, err = cli.CompileTarPkg(entry.PackageSource(), opts)
		} else {
			// kcl compiles the package from the OCI reference or url.
			result, err = cli.CompileOciPkg(entry.PackageSource(), o.Tag, opts)
		}
		if err != nil {
			return err
		}
	}
	_, err = opts.LogWriter().Write([]byte(result.GetRawYamlResult() + "\n"))
	if err != nil {
		return err
	}
	return nil
}

func (o *RunOptions) Validate() error {
	return nil
}

// CompileOptionFromCli will parse the kcl options from the cli context.
func CompileOptionFromCli(o *RunOptions) *opt.CompileOptions {
	opts := opt.DefaultCompileOptions()

	// --input
	opts.ExtendEntries(o.Entries)

	// --vendor
	opts.SetVendor(o.Vendor)

	// --setting, -Y
	if len(o.Settings) != 0 {
		for _, sPath := range o.Settings {
			opts.Merge(kcl.WithSettings(sPath))
		}
		opts.SetHasSettingsYaml(true)
	}

	// --argument, -D
	if len(o.Arguments) != 0 {
		for _, arg := range o.Arguments {
			opts.Merge(kcl.WithOptions(arg))
		}
	}

	// --overrides, -O
	if len(o.Overrides) != 0 {
		for _, override := range o.Overrides {
			opts.Merge(kcl.WithOverrides(override))
		}
	}

	// --disable_none, -n
	opts.Merge(kcl.WithDisableNone(o.DisableNone))

	// --sort_keys, -k
	opts.Merge(kcl.WithSortKeys(o.SortKeys))

	return opts
}
