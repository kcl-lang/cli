package version

import (
	"fmt"
	"runtime"
	"testing"
)

func TestGetVersionString(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "test get version in string",
			want: VersionTypeLatest.String(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetVersionString(); got != tt.want {
				t.Errorf(" GetVersionString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuildInfoVersion(t *testing.T) {
	// Save and restore the package-level version variable so we can exercise the
	// build info path without leaking state across tests.
	prevVersion := version
	defer func() { version = prevVersion }()

	// Force the package-level fallback so buildInfoVersion is consulted.
	version = ""

	// When running under `go test`, ReadBuildInfo returns "(devel)" for the main
	// module, which buildInfoVersion must treat as "no version available".
	if got := buildInfoVersion(); got != "" {
		t.Errorf("buildInfoVersion() under `go test` = %q, want empty string", got)
	}
}

// TestGetVersionString_AlwaysHasSuffix is a regression test for kcl-lang/cli#386:
// every code path must produce a string of the form `{version}-{goos}-{goarch}`,
// regardless of which version source fires. The previous implementation only
// formatted the `-{goos}-{goarch}` suffix on the `VersionTypeLatest` fallback,
// so callers that resolved the version via the `-ldflags` injection or via
// `debug.ReadBuildInfo()` printed a bare `0.12.10` instead.
func TestGetVersionString_AlwaysHasSuffix(t *testing.T) {
	wantSuffix := fmt.Sprintf("-%s-%s", runtime.GOOS, runtime.GOARCH)

	cases := []struct {
		name    string
		setup   func()
		cleanup func()
	}{
		{
			name: "fallback to VersionTypeLatest",
		},
		{
			name: "via buildInfoVersion (forced by clearing version)",
			setup: func() {
				version = ""
			},
			cleanup: func() {},
		},
		{
			name: "via ldflags-injected version",
			setup: func() {
				version = "1.2.3-test"
			},
			cleanup: func() {
				version = ""
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setup != nil {
				tc.setup()
			}
			defer func() {
				if tc.cleanup != nil {
					tc.cleanup()
				}
			}()

			got := GetVersionString()
			if got == "" {
				t.Fatalf("GetVersionString() = empty string, want non-empty")
			}
			// Must end with the platform suffix.
			if !hasSuffix(got, wantSuffix) {
				t.Errorf("GetVersionString() = %q, want suffix %q", got, wantSuffix)
			}
		})
	}
}

func hasSuffix(s, suffix string) bool {
	if len(suffix) > len(s) {
		return false
	}
	return s[len(s)-len(suffix):] == suffix
}
