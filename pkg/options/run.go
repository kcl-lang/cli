// Copyright The KCL Authors. All rights reserved.

package options

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"

	"github.com/acarl005/stripansi"
	"github.com/goccy/go-yaml"
	"github.com/pkg/errors"
	"kcl-lang.io/cli/pkg/fs"
	"kcl-lang.io/kcl-go/pkg/3rdparty/toml"
	"kcl-lang.io/kcl-go/pkg/kcl"
	"kcl-lang.io/kcl-go/pkg/tools/gen"
	"kcl-lang.io/kpm/pkg/client"
	"kcl-lang.io/kpm/pkg/constants"
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
	// InsecureSkipTLSverify denotes skipping the TLS verification.
	InsecureSkipTLSverify bool
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
	// Git Url is the package url of the Git artifact.
	Git string
	// Oci Url is the package url of the OCI artifact.
	Oci string
	// Tag is the package tag of the OCI or Git artifact.
	Tag string
	// Commit is the package commit of the Git artifact.
	Commit string
	// Branch is the package branch of the Git artifact.
	Branch string
	// CompileOnly is used to check a local package and all of its dependencies for errors.
	CompileOnly bool
	// Format is the output type, e.g., Json, Yaml, Toml etc. Default is Yaml.
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
	cli.SetInsecureSkipTLSverify(o.InsecureSkipTLSverify)
	// Acquire the lock of the package cache.
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
	// Generate temp entries from os.Stdin
	tempEntries := []string{}
	for i, entry := range o.Entries {
		if entry == "-" {
			entry, err := fs.GenTempFileFromStdin()
			if err != nil {
				return err
			}
			tempEntries = append(tempEntries, entry)
			o.Entries[i] = entry
		}
	}
	result, err = cli.Run(
		client.WithRunSourceUrls(o.Entries),
		client.WithSettingFiles(o.Settings),
		client.WithArguments(o.Arguments),
		client.WithOverrides(o.Overrides, o.Debug),
		client.WithPathSelectors(o.PathSelectors),
		client.WithExternalPkgs(o.ExternalPackages),
		client.WithVendor(o.Vendor),
		client.WithSortKeys(o.SortKeys),
		client.WithShowHidden(o.ShowHidden),
		client.WithDisableNone(o.DisableNone),
		client.WithDebug(o.Debug),
		client.WithStrictRange(o.StrictRangeCheck),
		client.WithCompileOnly(o.CompileOnly),
		client.WithLogger(o.Writer),
	)

	if err != nil {
		if o.NoStyle {
			err = errors.New(stripansi.Strip(err.Error()))
		}
		return err
	}
	// Remove temp entries
	for _, entry := range tempEntries {
		_ = os.Remove(entry)
	}
	return o.writeResult(result)
}

// Complete completes the options based on the provided arguments.
func (o *RunOptions) Complete(args []string) error {
	if len(o.Git) != 0 {
		gitUrl, err := url.Parse(o.Git)
		if err != nil {
			return err
		}
		if gitUrl.Scheme == constants.HttpsScheme || gitUrl.Scheme == constants.HttpScheme {
			gitUrl.Scheme = constants.GitScheme
		}
		query := gitUrl.Query()
		if o.Tag != "" {
			query.Set("tag", o.Tag)
		}
		if o.Commit != "" {
			query.Set("commit", o.Commit)
		}
		if o.Branch != "" {
			query.Set("branch", o.Branch)
		}
		gitUrl.RawQuery = query.Encode()
		o.Entries = append(o.Entries, gitUrl.String())
	}

	if len(o.Oci) != 0 {
		ociUrl, err := url.Parse(o.Oci)
		if err != nil {
			return err
		}
		if ociUrl.Scheme == constants.HttpsScheme || ociUrl.Scheme == constants.HttpScheme {
			ociUrl.Scheme = constants.OciScheme
		}
		query := ociUrl.Query()
		if o.Tag != "" {
			query.Set("tag", o.Tag)
		}
		ociUrl.RawQuery = query.Encode()
		o.Entries = append(o.Entries, ociUrl.String())
	}

	for _, arg := range args {
		argUrl, err := url.Parse(arg)
		if err != nil {
			return err
		}
		query := argUrl.Query()
		if o.Tag != "" {
			query.Set("tag", o.Tag)
		}
		if o.Commit != "" {
			query.Set("commit", o.Commit)
		}
		if o.Branch != "" {
			query.Set("branch", o.Branch)
		}
		argUrl.RawQuery = query.Encode()
		o.Entries = append(o.Entries, argUrl.String())
	}
	return nil
}

// Validate validates the options.
func (o *RunOptions) Validate() error {

	if len(o.Tag) != 0 || len(o.Commit) != 0 || len(o.Branch) != 0 {
		// Tag, commit, and branch are only valid with a single module.
		if len(o.Entries) > 1 {
			return fmt.Errorf("cannot specify tag, commit, or branch with multiple modules %s", o.Entries)
		}

		// Tag, commit, and branch must be specified with a module.
		if len(o.Entries) == 0 {
			return fmt.Errorf("cannot specify tag, commit, or branch without modules")
		}

		// Check that only one of tag, commit, or branch is specified
		specCount := 0
		if len(o.Tag) != 0 {
			specCount++
		}
		if len(o.Commit) != 0 {
			specCount++
		}
		if len(o.Branch) != 0 {
			specCount++
		}
		if specCount > 1 {
			return fmt.Errorf("only one of tag, commit, or branch can be specified")
		}
	}

	if len(o.Git) != 0 || len(o.Oci) != 0 {
		if len(o.Entries) > 1 {
			return fmt.Errorf("cannot specify multiple KCL modules %s", o.Entries)
		}
	}

	if o.Format != "" && strings.ToLower(o.Format) != Json && strings.ToLower(o.Format) != Yaml && strings.ToLower(o.Format) != Toml {
		return fmt.Errorf("invalid output format, expected %v, got %v", []string{Json, Yaml, Toml}, o.Format)
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
	} else if strings.ToLower(o.Format) == Toml {
		var out []byte
		var err error
		if o.SortKeys {
			yamlData := make(map[string]any)
			if err := yaml.UnmarshalWithOptions([]byte(result.GetRawYamlResult()), &yamlData); err != nil {
				return err
			}
			out, err = toml.Marshal(&yamlData)
		} else {
			yamlData := &yaml.MapSlice{}
			if err := yaml.UnmarshalWithOptions([]byte(result.GetRawYamlResult()), yamlData, yaml.UseOrderedMap()); err != nil {
				return err
			}
			out, err = gen.MarshalTOML(yamlData)
		}
		if err != nil {
			return err
		}
		output = []byte(string(out) + "\n")
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
