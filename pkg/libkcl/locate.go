// Copyright The KCL Authors. All rights reserved.

package libkcl

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// Config controls a FreeBSD runtime lookup. Tests set Install and
// SystemDirs; Prepare fills them from the process environment.
type Config struct {
	GOOS   string
	GOARCH string
	Home   string
	Env    map[string]string
	// SystemDirs are searched after the user cache. Nil selects the
	// FreeBSD defaults (/usr/local/lib, /usr/lib). An empty slice
	// searches no system directories.
	SystemDirs []string
	// Install is the upstream runtime installer (kcl-lang.io/lib
	// InstallKcl). Resolve calls it before returning a directory
	// because the native loader runs that installer and panics when
	// it fails. Nil skips the call. This repository does not
	// implement installLib.
	Install func(dir string) error
}

// Resolve returns the directory the native loader should use as
// KCL_LIB_HOME. The directory contains libkcl.so exporting
// kcl_service_new.
func Resolve(cfg Config) (string, error) {
	cfg = cfg.normalized()
	rep := &loadReport{goos: cfg.GOOS, goarch: cfg.GOARCH}

	if explicit := cfg.Env[EnvLibHome]; explicit != "" {
		if rep.consider(explicit) {
			return finish(cfg, explicit, rep)
		}
		if err := installTo(cfg, explicit); err != nil {
			return "", rep.err(err)
		}
		if libraryUsable(filepath.Join(explicit, LibraryFile)) {
			return explicit, nil
		}
		return "", rep.err(missingLibrary(explicit))
	}

	cache, err := cfg.cacheDir()
	if err != nil {
		return "", rep.err(err)
	}
	if rep.consider(cache) {
		return finish(cfg, cache, rep)
	}
	for _, dir := range cfg.systemDirs() {
		if !rep.consider(dir) {
			continue
		}
		staged, err := stage(cache, filepath.Join(dir, LibraryFile))
		if err != nil {
			return "", rep.err(err)
		}
		return finish(cfg, staged, rep)
	}
	if err := installTo(cfg, cache); err != nil {
		return "", rep.err(err)
	}
	if libraryUsable(filepath.Join(cache, LibraryFile)) {
		return cache, nil
	}
	return "", rep.err(missingLibrary(cache))
}

func missingLibrary(dir string) error {
	path := filepath.Join(dir, LibraryFile)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("libkcl runtime %s is missing from %s", LibraryFile, dir)
	}
	return fmt.Errorf("%s does not export %s", path, ServiceSymbol)
}

func finish(cfg Config, dir string, rep *loadReport) (string, error) {
	if err := installTo(cfg, dir); err != nil {
		return "", rep.err(err)
	}
	path := filepath.Join(dir, LibraryFile)
	if !libraryUsable(path) {
		return "", rep.err(fmt.Errorf("%s does not export %s", path, ServiceSymbol))
	}
	return dir, nil
}

func installTo(cfg Config, dir string) error {
	if cfg.Install == nil {
		return nil
	}
	return cfg.Install(dir)
}

func (cfg Config) normalized() Config {
	if cfg.GOOS == "" {
		cfg.GOOS = "freebsd"
	}
	if cfg.GOARCH == "" {
		cfg.GOARCH = runtime.GOARCH
	}
	if cfg.Env == nil {
		cfg.Env = map[string]string{}
	}
	return cfg
}

// cacheDir matches kcl-lang.io/lib: $XDG_CACHE_HOME/kcl, or ~/.cache/kcl.
func (cfg Config) cacheDir() (string, error) {
	if xdg := cfg.Env["XDG_CACHE_HOME"]; xdg != "" && filepath.IsAbs(xdg) {
		return filepath.Join(xdg, "kcl"), nil
	}
	if cfg.Home == "" {
		return "", fmt.Errorf("home directory is unknown; set %s to the directory that contains %s", EnvLibHome, LibraryFile)
	}
	return filepath.Join(cfg.Home, ".cache", "kcl"), nil
}

func (cfg Config) systemDirs() []string {
	if cfg.SystemDirs != nil {
		return cfg.SystemDirs
	}
	return []string{defaultLocalLibDir, defaultBaseLibDir}
}

func candidatePaths(cfg Config) []string {
	cfg = cfg.normalized()
	var dirs []string
	if explicit := cfg.Env[EnvLibHome]; explicit != "" {
		dirs = append(dirs, explicit)
	} else {
		if cache, err := cfg.cacheDir(); err == nil {
			dirs = append(dirs, cache)
		}
		dirs = append(dirs, cfg.systemDirs()...)
	}
	paths := make([]string, 0, len(dirs))
	for _, dir := range dirs {
		paths = append(paths, filepath.Join(dir, LibraryFile))
	}
	return paths
}

type loadReport struct {
	goos    string
	goarch  string
	checked []string
	notes   []string
}

func (r *loadReport) consider(dir string) bool {
	r.checked = append(r.checked, filepath.Join(dir, LibraryFile))
	legacy := filepath.Join(dir, LegacyLibraryFile)
	if st, err := os.Stat(legacy); err == nil && !st.IsDir() {
		r.notes = append(r.notes, fmt.Sprintf("Found legacy library %s. The loader entry point is %s in %s.", legacy, ServiceSymbol, LibraryFile))
	}
	path := filepath.Join(dir, LibraryFile)
	switch kind, detail := classify(path); kind {
	case kindUsable:
		r.notes = append(r.notes, fmt.Sprintf("Found %s exporting %s.", path, ServiceSymbol))
		return true
	case kindNoSymbol:
		note := fmt.Sprintf("Found %s without %s.", path, ServiceSymbol)
		if ok, err := exportsSymbol(path, LegacySymbol); err == nil && ok {
			note = fmt.Sprintf("Found %s exporting %s. The loader entry point is %s.", path, LegacySymbol, ServiceSymbol)
		}
		r.notes = append(r.notes, note)
		return false
	case kindBad:
		r.notes = append(r.notes, fmt.Sprintf("Found %s but it is not a readable libkcl shared library: %s.", path, detail))
		return false
	default:
		return false
	}
}

func (r *loadReport) err(detail error) error {
	e := &LoadError{
		GOOS:    r.goos,
		GOARCH:  r.goarch,
		Checked: append([]string(nil), r.checked...),
		Notes:   append([]string(nil), r.notes...),
	}
	if detail != nil {
		e.Detail = detail.Error()
	}
	return e
}

type libKind int

const (
	kindMissing libKind = iota
	kindUsable
	kindNoSymbol
	kindBad
)

func classify(path string) (libKind, string) {
	st, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return kindMissing, ""
		}
		return kindBad, err.Error()
	}
	if st.IsDir() {
		return kindBad, "is a directory"
	}
	ok, err := exportsSymbol(path, ServiceSymbol)
	if err != nil {
		return kindBad, err.Error()
	}
	if !ok {
		return kindNoSymbol, ""
	}
	return kindUsable, ""
}

func libraryUsable(path string) bool {
	kind, _ := classify(path)
	return kind == kindUsable
}

func stage(cacheDir, src string) (string, error) {
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		return "", err
	}
	dst := filepath.Join(cacheDir, LibraryFile)
	if same, err := sameFile(dst, src); err == nil && same {
		return cacheDir, nil
	}
	if err := os.Remove(dst); err != nil && !os.IsNotExist(err) {
		return "", err
	}
	absSrc, err := filepath.Abs(src)
	if err != nil {
		absSrc = src
	}
	if err := os.Symlink(absSrc, dst); err != nil {
		return "", err
	}
	return cacheDir, nil
}

func sameFile(a, b string) (bool, error) {
	sa, err := os.Stat(a)
	if err != nil {
		return false, err
	}
	sb, err := os.Stat(b)
	if err != nil {
		return false, err
	}
	return os.SameFile(sa, sb), nil
}
