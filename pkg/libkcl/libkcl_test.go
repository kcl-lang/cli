// Copyright The KCL Authors. All rights reserved.

package libkcl

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func freebsdConfig(home string, system []string) Config {
	if system == nil {
		system = []string{}
	}
	return Config{
		GOOS:       "freebsd",
		GOARCH:     "amd64",
		Home:       home,
		Env:        map[string]string{},
		SystemDirs: system,
	}
}

func TestExportsCurrentSymbol(t *testing.T) {
	dir := t.TempDir()
	path := writeTestLibrary(t, dir, ServiceSymbol)
	ok, err := exportsSymbol(path, ServiceSymbol)
	if err != nil || !ok {
		t.Fatalf("exports %s: ok=%v err=%v", ServiceSymbol, ok, err)
	}
	ok, err = exportsSymbol(path, LegacySymbol)
	if err != nil || ok {
		t.Fatalf("legacy symbol present: ok=%v err=%v", ok, err)
	}
}

func TestResolveUsesExplicitLibHome(t *testing.T) {
	dir := t.TempDir()
	writeTestLibrary(t, dir, ServiceSymbol)
	cfg := freebsdConfig(t.TempDir(), nil)
	cfg.Env[EnvLibHome] = dir

	got, err := Resolve(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if got != dir {
		t.Fatalf("dir = %s, want %s", got, dir)
	}
}

func TestResolveStagesSystemLibraryIntoCache(t *testing.T) {
	home := t.TempDir()
	system := t.TempDir()
	src := writeTestLibrary(t, system, ServiceSymbol)
	cfg := freebsdConfig(home, []string{system})

	got, err := Resolve(cfg)
	if err != nil {
		t.Fatal(err)
	}
	cache := filepath.Join(home, ".cache", "kcl")
	if got != cache {
		t.Fatalf("dir = %s, want %s", got, cache)
	}
	link := filepath.Join(cache, LibraryFile)
	target, err := os.Readlink(link)
	if err != nil {
		t.Fatal(err)
	}
	if target != src {
		t.Fatalf("symlink = %s, want %s", target, src)
	}
	ok, err := exportsSymbol(link, ServiceSymbol)
	if err != nil || !ok {
		t.Fatalf("staged library exports %s: ok=%v err=%v", ServiceSymbol, ok, err)
	}
}

func TestResolveMissingRuntimeNamesLibkclAndPlatform(t *testing.T) {
	home := t.TempDir()
	system := t.TempDir()
	legacy := filepath.Join(system, LegacyLibraryFile)
	if err := os.WriteFile(legacy, []byte("old"), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg := freebsdConfig(home, []string{system})
	cfg.Install = func(string) error {
		return errors.New("the prebuilt libkcl binary is not vendored in this build")
	}

	_, err := Resolve(cfg)
	if err != nil {
		assertLoadError(t, err)
		return
	}
	t.Fatal("expected missing runtime error")
}

func TestResolveRejectsLegacySymbol(t *testing.T) {
	dir := t.TempDir()
	writeTestLibrary(t, dir, LegacySymbol)
	cfg := freebsdConfig(t.TempDir(), nil)
	cfg.Env[EnvLibHome] = dir

	_, err := Resolve(cfg)
	if err == nil {
		t.Fatal("expected error")
	}
	assertLoadError(t, err)
	if !strings.Contains(err.Error(), LegacySymbol) {
		t.Fatalf("error should name %s:\n%s", LegacySymbol, err)
	}
}

func TestResolveInstallPopulatesCache(t *testing.T) {
	home := t.TempDir()
	cfg := freebsdConfig(home, nil)
	cfg.Install = func(dir string) error {
		writeTestLibrary(t, dir, ServiceSymbol)
		return nil
	}

	got, err := Resolve(cfg)
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(home, ".cache", "kcl")
	if got != want {
		t.Fatalf("dir = %s, want %s", got, want)
	}
}

func TestResolveInstallFailureIsClearError(t *testing.T) {
	cfg := freebsdConfig(t.TempDir(), nil)
	cfg.Install = func(string) error {
		return errors.New("undefined: installLib")
	}
	_, err := Resolve(cfg)
	if err == nil {
		t.Fatal("expected error")
	}
	assertLoadError(t, err)
	if !strings.Contains(err.Error(), "installLib") {
		t.Fatalf("missing installer detail:\n%s", err)
	}
}

func TestResolveHonorsXDGCacheHome(t *testing.T) {
	xdg := t.TempDir()
	cache := filepath.Join(xdg, "kcl")
	writeTestLibrary(t, cache, ServiceSymbol)
	cfg := freebsdConfig(t.TempDir(), nil)
	cfg.Env["XDG_CACHE_HOME"] = xdg

	got, err := Resolve(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if got != cache {
		t.Fatalf("dir = %s, want %s", got, cache)
	}
}

func TestErrorFromPanicUndefinedLegacySymbol(t *testing.T) {
	err := ErrorFromPanic(`Undefined symbol "kclvm_service_new"`)
	if err == nil {
		t.Fatal("expected error")
	}
	msg := err.Error()
	for _, want := range []string{LibraryFile, ServiceSymbol, LegacySymbol, "libkcl", runtime.GOOS, runtime.GOARCH} {
		if !strings.Contains(msg, want) {
			t.Errorf("panic translation missing %q:\n%s", want, msg)
		}
	}
	if strings.Contains(msg, "panic:") {
		t.Fatalf("loader failure still looks like a panic:\n%s", msg)
	}
}

func TestErrorFromPanicIgnoresUnrelated(t *testing.T) {
	if err := ErrorFromPanic("index out of range"); err != nil {
		t.Fatal(err)
	}
	if err := ErrorFromPanic(nil); err != nil {
		t.Fatal(err)
	}
}

func TestPrepareSkipsNonFreeBSD(t *testing.T) {
	if runtime.GOOS == "freebsd" {
		t.Skip("Prepare searches for libkcl on FreeBSD")
	}
	t.Setenv(EnvLibHome, "")
	if err := Prepare(); err != nil {
		t.Fatal(err)
	}
	if got := os.Getenv(EnvLibHome); got != "" {
		t.Fatalf("%s = %q", EnvLibHome, got)
	}
}

func assertLoadError(t *testing.T, err error) {
	t.Helper()
	msg := err.Error()
	for _, want := range []string{LibraryFile, ServiceSymbol, "freebsd", "amd64", "libkcl"} {
		if !strings.Contains(msg, want) {
			t.Errorf("error missing %q:\n%s", want, msg)
		}
	}
	var load *LoadError
	if !errors.As(err, &load) {
		t.Fatalf("error type = %T, want *LoadError", err)
	}
}
