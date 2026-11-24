// Copyright The KCL Authors. All rights reserved.

// Package libkcl locates the native KCL runtime (libkcl.so) before the
// FreeBSD loader panics.
//
// kcl-lang.io/lib binds the current entry point kcl_service_new. Older
// FreeBSD builds panicked in purego while looking up kclvm_service_new.
// This package searches the FreeBSD library locations, asks the upstream
// installer to unpack a runtime when one is not already present, and
// returns an error that names libkcl and the platform when the runtime
// is still missing. It does not implement installLib; that stays in
// kcl-lang.io/lib.
package libkcl

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"kcl-lang.io/lib/go/install"
)

const (
	// LibraryFile is the shared library the native loader opens.
	LibraryFile = "libkcl.so"
	// LegacyLibraryFile is the previous FreeBSD artifact name. The loader
	// does not open it.
	LegacyLibraryFile = "libkclvm_cli_cdylib.so"
	// ServiceSymbol is the entry point current libkcl exports.
	ServiceSymbol = "kcl_service_new"
	// LegacySymbol is the entry point named by the original FreeBSD panic.
	LegacySymbol = "kclvm_service_new"
	// EnvLibHome is the directory the native loader searches first.
	EnvLibHome = "KCL_LIB_HOME"

	freebsdArtifactsIssue = "https://github.com/kcl-lang/kcl/issues/1029"

	defaultLocalLibDir = "/usr/local/lib"
	defaultBaseLibDir  = "/usr/lib"
)

// LoadError is returned when the libkcl runtime cannot be loaded.
type LoadError struct {
	GOOS    string
	GOARCH  string
	Checked []string
	Notes   []string
	Detail  string
}

func (e *LoadError) Error() string {
	if e == nil {
		return "kcl: libkcl runtime is unavailable"
	}
	goos, goarch := e.GOOS, e.GOARCH
	if goos == "" {
		goos = runtime.GOOS
	}
	if goarch == "" {
		goarch = runtime.GOARCH
	}
	var b strings.Builder
	fmt.Fprintf(&b, "kcl: libkcl runtime is unavailable for %s/%s\n", goos, goarch)
	fmt.Fprintf(&b, "The native loader requires %s exporting %s.\n", LibraryFile, ServiceSymbol)
	if len(e.Checked) > 0 {
		b.WriteString("Looked for libkcl.so in:\n")
		for _, path := range e.Checked {
			fmt.Fprintf(&b, "  - %s\n", path)
		}
	}
	for _, note := range e.Notes {
		b.WriteString(note)
		if !strings.HasSuffix(note, "\n") {
			b.WriteByte('\n')
		}
	}
	if e.Detail != "" {
		fmt.Fprintf(&b, "Detail: %s\n", e.Detail)
	}
	fmt.Fprintf(&b, "Set %s to the directory that contains %s.\n", EnvLibHome, LibraryFile)
	if goos == "freebsd" {
		fmt.Fprintf(&b, "FreeBSD runtime artifacts: %s\n", freebsdArtifactsIssue)
	}
	return strings.TrimRight(b.String(), "\n")
}

// Prepare points the native loader at a FreeBSD libkcl.so that exports
// kcl_service_new. On other operating systems it does nothing, because
// those builds unpack a vendored runtime on first use.
//
// A missing or unusable runtime returns *LoadError. The process does not
// panic with an undefined-symbol message.
func Prepare() error {
	if runtime.GOOS != "freebsd" {
		return nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		home = ""
	}
	dir, err := Resolve(Config{
		GOOS:   runtime.GOOS,
		GOARCH: runtime.GOARCH,
		Home:   home,
		Env: map[string]string{
			EnvLibHome:       os.Getenv(EnvLibHome),
			"XDG_CACHE_HOME": os.Getenv("XDG_CACHE_HOME"),
		},
		Install: install.InstallKcl,
	})
	if err != nil {
		return err
	}
	if dir == "" {
		return nil
	}
	return os.Setenv(EnvLibHome, dir)
}

// ErrorFromPanic converts a native-loader panic into *LoadError.
// Unrelated panics return nil so callers can crash as usual.
func ErrorFromPanic(recovered any) error {
	if recovered == nil {
		return nil
	}
	msg := fmt.Sprint(recovered)
	if !isLoaderPanic(msg) {
		return nil
	}
	err := &LoadError{
		GOOS:   runtime.GOOS,
		GOARCH: runtime.GOARCH,
		Detail: msg,
	}
	if strings.Contains(msg, LegacySymbol) {
		err.Notes = append(err.Notes, fmt.Sprintf("The loader entry point is %s in %s.", ServiceSymbol, LibraryFile))
	}
	if runtime.GOOS == "freebsd" {
		home, homeErr := os.UserHomeDir()
		if homeErr != nil {
			home = ""
		}
		err.Checked = candidatePaths(Config{
			GOOS:   runtime.GOOS,
			GOARCH: runtime.GOARCH,
			Home:   home,
			Env: map[string]string{
				EnvLibHome:       os.Getenv(EnvLibHome),
				"XDG_CACHE_HOME": os.Getenv("XDG_CACHE_HOME"),
			},
		})
	}
	return err
}

func isLoaderPanic(msg string) bool {
	needles := []string{
		"Undefined symbol",
		ServiceSymbol,
		LegacySymbol,
		"installLib",
		"libkcl",
		"prebuilt libkcl",
		EnvLibHome,
	}
	for _, needle := range needles {
		if strings.Contains(msg, needle) {
			return true
		}
	}
	return false
}
