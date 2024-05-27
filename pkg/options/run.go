// Copyright The KCL Authors. All rights reserved.

package options

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/acarl005/stripansi"
	"github.com/pkg/errors"
	"kcl-lang.io/cli/pkg/fs"
	"kcl-lang.io/kcl-go/pkg/kcl"
	"kcl-lang.io/kpm/pkg/api"
	"kcl-lang.io/kpm/pkg/client"
	"kcl-lang.io/kpm/pkg/git"
	"kcl-lang.io/kpm/pkg/opt"
	pkg "kcl-lang.io/kpm/pkg/package"
	"kcl-lang.io/kpm/pkg/runner"
)

const (
	DefaultSettingsFile = "kcl.yaml"
)

// RunOptions is a struct that holds the options for the run command.
type RunOptions struct {
	// Entries is the list of the kcl code entry including filepath, folder, OCI package, etc.
	Entries []string
	// Output is the result output filepath. Default is os.Stdout.
	Output string
	// Settings is the list of kcl setting files including all of the CLI config.
	Settings []string
	// Arguments is the list of top level dynamic arguments for the kcl option function, e.g., env="prod"
	Arguments []string
	// Overrides is the list of override paths and values, e.g., app.image="v2"
	Overrides []string
	// PathSelectors is the list of path selectors to select output result, e.g., a.b.c
	PathSelectors []string
	// ExternalPackages denotes the list of external packages, e.g., k8s=./vendor/k8s
	ExternalPackages []string
	// NoStyle denotes disabling the output information style and color.
	NoStyle bool
	// Quiet denotes disabling all the output information.
	Quiet bool
	// Vendor denotes running kcl in the vendor mode.
	Vendor bool
	// SortKeys denotes sorting the output result keys, e.g., `{b = 1, a = 2} => {a = 2, b = 1}`.
	SortKeys bool
	// ShowHidden denotes output the hidden attribute in the result.
	ShowHidden bool
	// DisableNone denotes running kcl and disable dumping None values.
	DisableNone bool
	// Debug denotes running kcl in debug mode.
	Debug bool
	// StrictRangeCheck performs the 32-bit strict numeric range checks on numbers.
	StrictRangeCheck bool
	// Tag is the package tag of the OCI or Git artifact.
	Tag string
	// CompileOnly is used to check a local package and all of its dependencies for errors.
	CompileOnly bool
	// Format is the output type, e.g., Json, Yaml, etc. Default is Yaml.
	Format string
	// Writer is used to output the run result. Default is os.Stdout.
	Writer io.Writer
}

// NewRunOptions returns a new instance of RunOptions with default values.
func NewRunOptions() *RunOptions {
	return &RunOptions{
		Writer: os.Stdout,
		Format: Yaml,
	}
}

// Run runs the kcl run command with options.
func (o *RunOptions) Run() error {
	var result *kcl.KCLResultList
	var err error
	cli, err := client.NewKpmClient()
	if err != nil {
		return err
	}
	if o.Quiet {
		cli.SetLogWriter(nil)
	}
	// acquire the lock of the package cache.
	err = cli.AcquirePackageCacheLock()
	if err != nil {
		return err
	}
	defer func() {
		// release the lock of the package cache after the function returns.
		releaseErr := cli.ReleasePackageCacheLock()
		if releaseErr != nil && err == nil {
			err = releaseErr
		}
	}()
	opts := CompileOptionFromCli(o)
	if err != nil {
		return err
	}
	entry, errEvent := runner.FindRunEntryFrom(opts.Entries())
	if errEvent != nil {
		return errEvent
	}
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	if entry.IsEmpty() {
		// kcl compiles the current package under '$pwd'.
		if _, e := api.GetKclPackage(pwd); e == nil {
			opts.SetPkgPath(pwd)
			result, err = cli.CompileWithOpts(opts)
		} else {
			// TODO: refactor the entry search logic.
			depsOpt, depErr := LoadDepsFrom(pwd, o.Quiet)
			if depErr != nil {
				return err
			}
			opts.Merge(*depsOpt)
			result, err = api.RunWithOpt(opts)
		}
	} else {
		// kcl compiles the package from the local file system, tar and OCI package, etc.
		if entry.IsLocalFile() {
			depsOpt, depErr := LoadDepsFrom(pwd, o.Quiet)
			if depErr != nil {
				return err
			}
			opts.Merge(*depsOpt)
			result, err = api.RunWithOpt(opts)
		} else if entry.IsLocalFileWithKclMod() {
			// Else compile the kcl package (kcl.mod)
			var transformedEntries []string
			for _, entry := range opts.Entries() {
				if !filepath.IsAbs(entry) {
					entry, err = filepath.Abs(entry)
					if err != nil {
						return err
					}
				}
				transformedEntries = append(transformedEntries, entry)
			}
			// Maybe a single KCL module folder, use the kcl.mod entry profile to run.
			if pkg, e := api.GetKclPackage(transformedEntries[0]); e == nil && fs.IsDir(transformedEntries[0]) {
				entries := pkg.GetPkgProfile().GetEntries()
				if len(entries) > 0 {
					opts.SetEntries(entries)
					opts.SetPkgPath(transformedEntries[0])
				} else {
					// Multiple entries with the kcl.mod file and deps.
					opts.SetEntries(transformedEntries)
					opts.SetPkgPath(entry.PackageSource())
				}
			} else {
				// Multiple entries with the kcl.mod file and deps.
				opts.SetEntries(transformedEntries)
				opts.SetPkgPath(entry.PackageSource())
			}
			result, err = cli.CompileWithOpts(opts)
		} else if entry.IsTar() {
			// compiles the package from the kcl package tar.
			opts.SetEntries([]string{})
			result, err = cli.CompileTarPkg(entry.PackageSource(), opts)
		} else if entry.IsGit() {
			opts.SetEntries([]string{})
			gitOpts := git.NewCloneOptions(entry.PackageSource(), "", o.Tag, "", "", nil)
			// compiles the package from the git url
			result, err = cli.CompileGitPkg(gitOpts, opts)
		} else if entry.IsUrl() {
			// compiles the package from the OCI reference or url.
			opts.SetEntries([]string{})
			result, err = cli.CompileOciPkg(entry.PackageSource(), o.Tag, opts)
		} else {
			// If there is only kcl file without kcl package (kcl.mod)
			result, err = api.RunWithOpt(opts)
		}
	}
	if err != nil {
		if o.NoStyle {
			err = errors.New(stripansi.Strip(err.Error()))
		}
		return err
	}
	return o.writeResult(result)
}

// Complete completes the options based on the provided arguments.
func (o *RunOptions) Complete(args []string) error {
	o.Entries = args
	return nil
}

// Validate validates the options.
func (o *RunOptions) Validate() error {
	if o.Format != "" && strings.ToLower(o.Format) != Json && strings.ToLower(o.Format) != Yaml {
		return fmt.Errorf("invalid output format, expected %v, got %v", []string{Json, Yaml}, o.Format)
	}
	for _, setting := range o.Settings {
		if _, err := os.Stat(setting); err != nil {
			return fmt.Errorf("failed to load '%s', no such file or directory", setting)
		}
	}
	return nil
}

func (o *RunOptions) writeResult(result *kcl.KCLResultList) error {
	if result == nil {
		return nil
	}
	var output []byte
	if strings.ToLower(o.Format) == Json {
		var out bytes.Buffer
		err := json.Indent(&out, []byte(result.GetRawJsonResult()), "", "    ")
		if err != nil {
			return err
		}
		output = []byte(out.String() + "\n")
	} else {
		// Both considering the raw YAML format and the YAML stream format that contains the `---` separator.
		output = []byte(result.GetRawYamlResult() + "\n")
	}

	if o.Output == "" {
		_, err := o.Writer.Write(output)
		if err != nil {
			return err
		}
	} else {
		file, err := os.OpenFile(o.Output, os.O_CREATE|os.O_RDWR, 0744)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = file.Write(output)
		if err != nil {
			return err
		}
	}
	return nil
}

// CompileOptionFromCli will parse the kcl options from the cli options.
func CompileOptionFromCli(o *RunOptions) *opt.CompileOptions {
	opts := opt.DefaultCompileOptions()

	// <input>
	opts.ExtendEntries(o.Entries)

	// --setting, -Y
	if len(o.Settings) != 0 {
		for _, sPath := range o.Settings {
			opts.Merge(kcl.WithSettings(sPath))
		}
		opts.SetHasSettingsYaml(true)
	} else if fs.FileExists(DefaultSettingsFile) {
		// If exists default kcl.yaml, load it.
		opts.Merge(kcl.WithSettings(DefaultSettingsFile))
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
		opts.Merge(kcl.WithOverrides(o.Overrides...))
		if o.Debug {
			opts.PrintOverrideAst = true
		}
	}

	// --path_selector, -S
	if len(o.PathSelectors) != 0 {
		opts.Merge(kcl.WithSelectors(o.PathSelectors...))
	}

	// --debug, -d
	if o.Debug {
		opts.Debug = 1
	}

	// --disable_none, -n
	opts.Merge(kcl.WithDisableNone(o.DisableNone))

	// --external, -E
	opts.Merge(kcl.WithExternalPkgs(o.ExternalPackages...))

	// --sort_keys, -k
	opts.Merge(kcl.WithSortKeys(o.SortKeys))

	// --show_hidden, -H
	opts.Merge(kcl.WithShowHidden(o.ShowHidden))

	// --strict_range_check, -r
	opts.StrictRangeCheck = o.StrictRangeCheck

	opts.CompileOnly = o.CompileOnly

	// --vendor
	opts.SetVendor(o.Vendor)

	// Set logger to stdout to show the kcl values of the print function.
	opts.Merge(kcl.WithLogger(os.Stdout))

	return opts
}

// LoadDepsFrom parses the kcl external package option from a path.
// It will find `kcl.mod` recursively from the path, resolve deps
// in the `kcl.mod` and return the option. If not found, return the
// empty option.
func LoadDepsFrom(path string, quiet bool) (*kcl.Option, error) {
	o := kcl.NewOption()
	entry, errEvent := runner.FindRunEntryFrom([]string{path})
	if errEvent != nil {
		return o, errEvent
	}
	if entry.IsLocalFileWithKclMod() {
		cli, err := client.NewKpmClient()
		if err != nil {
			return o, err
		}
		if quiet {
			cli.SetLogWriter(nil)
		}
		pkg, err := pkg.LoadKclPkg(entry.PackageSource())
		if err != nil {
			return o, err
		}
		depsMap, err := cli.ResolveDepsIntoMap(pkg)
		if err != nil {
			return o, err
		}
		for depName, depPath := range depsMap {
			o.Merge(kcl.WithExternalPkgs(fmt.Sprintf("%s=%s", depName, depPath)))
		}
	}
	return o, nil
}
