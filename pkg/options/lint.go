// Copyright The KCL Authors. All rights reserved.

package options

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
)

// HasLintPattern reports whether any of the args is a `...` pattern
// (e.g. `./...` or `pkg/...`) that selects all the KCL packages under
// a directory, like `go build ./...`.
func HasLintPattern(args []string) bool {
	for _, arg := range args {
		if isLintPattern(arg) {
			return true
		}
	}
	return false
}

// isLintPattern reports whether the path is a `...` package pattern.
func isLintPattern(path string) bool {
	return filepath.Base(path) == "..."
}

// LintAllPackages lints every KCL package selected by args. It supports
// `...` patterns (e.g. `./...` or `pkg/...`): every directory holding KCL
// files under the pattern root is linted as one package, so symbols of
// distinct packages never collide. Plain files and directories passed
// alongside a pattern join the package of their directory.
func (o *RunOptions) LintAllPackages(args []string) error {
	packages, err := lintPackages(args)
	if err != nil {
		return err
	}
	var errs []error
	for _, files := range packages {
		// Lint every package as an independent program.
		po := *o
		po.Entries = files
		po.CompileOnly = true
		if err := po.Run(); err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", filepath.Dir(files[0]), err))
		}
	}
	return errors.Join(errs...)
}

// lintPackages groups the lint targets of args into one file list per
// directory, in deterministic order:
//   - `...` patterns walk their root directory recursively;
//   - plain directories contribute their top-level KCL files, i.e. their
//     root package only;
//   - plain files join the list of their parent directory.
func lintPackages(args []string) ([][]string, error) {
	filesByDir := map[string][]string{}
	addFile := func(file string) {
		file = filepath.Clean(file)
		dir := filepath.Dir(file)
		if !slices.Contains(filesByDir[dir], file) {
			filesByDir[dir] = append(filesByDir[dir], file)
		}
	}
	addDir := func(dir string, recursively bool) error {
		return filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				if !recursively && path != dir {
					return filepath.SkipDir
				}
				return nil
			}
			if strings.HasSuffix(path, ".k") {
				addFile(path)
			}
			return nil
		})
	}
	for _, arg := range args {
		if isLintPattern(arg) {
			if err := addDir(filepath.Dir(arg), true); err != nil {
				return nil, err
			}
			continue
		}
		if info, err := os.Stat(arg); err == nil && info.IsDir() {
			if err := addDir(filepath.Clean(arg), false); err != nil {
				return nil, err
			}
			continue
		}
		addFile(arg)
	}
	if len(filesByDir) == 0 {
		return nil, fmt.Errorf("no KCL files found in %s", strings.Join(args, " "))
	}
	dirs := make([]string, 0, len(filesByDir))
	for dir := range filesByDir {
		dirs = append(dirs, dir)
	}
	sort.Strings(dirs)
	packages := make([][]string, 0, len(dirs))
	for _, dir := range dirs {
		files := filesByDir[dir]
		sort.Strings(files)
		packages = append(packages, files)
	}
	return packages, nil
}
